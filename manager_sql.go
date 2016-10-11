package ladon

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/pkg/errors"
	"github.com/ory-am/common/compiler"
	"github.com/ory-am/common/pkg"
	"github.com/jmoiron/sqlx"
)

var sqlSchema = []string{
	`CREATE TABLE IF NOT EXISTS ladon_policy (
		id           varchar(255) NOT NULL PRIMARY KEY,
		description  text NOT NULL,
		effect       text NOT NULL CHECK (effect='allow' OR effect='deny'),
		conditions 	 text NOT NULL
	)`,
	`CREATE TABLE IF NOT EXISTS ladon_policy_subject (
    	compiled text NOT NULL,
    	template varchar(1023) NOT NULL,
    	policy   varchar(255) NOT NULL REFERENCES ladon_policy (id) ON DELETE CASCADE
	)`,
	`CREATE TABLE IF NOT EXISTS ladon_policy_permission (
    	compiled text NOT NULL,
    	template varchar(1023) NOT NULL,
    	policy   varchar(255) NOT NULL REFERENCES ladon_policy (id) ON DELETE CASCADE
	)`,
	`CREATE TABLE IF NOT EXISTS ladon_policy_resource (
    	compiled text NOT NULL,
    	template varchar(1023) NOT NULL,
    	policy   varchar(255) NOT NULL REFERENCES ladon_policy (id) ON DELETE CASCADE
	)`,
}

// SQLManager is a postgres implementation for Manager to store policies persistently.
type SQLManager struct {
	db     *sqlx.DB
	schema []string
}

// NewSQLManager initializes a new SQLManager for given db instance.
func NewSQLManager(db *sqlx.DB, schema []string) *SQLManager {
	return &SQLManager{
		db: db,
		schema: schema,
	}
}

// CreateSchemas creates ladon_policy tables
func (s *SQLManager) CreateSchemas() error {
	schema := s.schema
	if schema == nil {
		schema = sqlSchema
	}
	for _, query := range schema {
		if _, err := s.db.Exec(query); err != nil {
			log.Printf("Error creating schema %s", query)
			return err
		}
	}
	return nil
}

// Create inserts a new policy
func (s *SQLManager) Create(policy Policy) (err error) {
	conditions := []byte("[]")
	if policy.GetConditions() != nil {
		cs := policy.GetConditions()
		conditions, err = json.Marshal(&cs)
		if err != nil {
			return err
		}
	}

	if tx, err := s.db.Begin(); err != nil {
		return err
	} else if _, err = tx.Exec(s.db.Rebind("INSERT INTO ladon_policy (id, description, effect, conditions) VALUES (?, ?, ?, ?)"), policy.GetID(), policy.GetDescription(), policy.GetEffect(), conditions); err != nil {
		return err
	} else if err = createLinkSQL(s.db, tx, "ladon_policy_subject", policy, policy.GetSubjects()); err != nil {
		return err
	} else if err = createLinkSQL(s.db, tx, "ladon_policy_permission", policy, policy.GetActions()); err != nil {
		return err
	} else if err = createLinkSQL(s.db, tx, "ladon_policy_resource", policy, policy.GetResources()); err != nil {
		return err
	} else if err = tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return nil
}

// Get retrieves a policy.
func (s *SQLManager) Get(id string) (Policy, error) {
	var p DefaultPolicy
	var conditions []byte

	if err := s.db.QueryRow(s.db.Rebind("SELECT id, description, effect, conditions FROM ladon_policy WHERE id=?"), id).Scan(&p.ID, &p.Description, &p.Effect, &conditions); err == sql.ErrNoRows {
		return nil, pkg.ErrNotFound
	} else if err != nil {
		return nil, errors.Wrap(err, "")
	}

	p.Conditions = Conditions{}
	if err := json.Unmarshal(conditions, &p.Conditions); err != nil {
		return nil, errors.Wrap(err, "")
	}

	subjects, err := getLinkedSQL(s.db, "ladon_policy_subject", id)
	if err != nil {
		return nil, err
	}
	permissions, err := getLinkedSQL(s.db, "ladon_policy_permission", id)
	if err != nil {
		return nil, err
	}
	resources, err := getLinkedSQL(s.db, "ladon_policy_resource", id)
	if err != nil {
		return nil, err
	}

	p.Actions = permissions
	p.Subjects = subjects
	p.Resources = resources
	return &p, nil
}

// Delete removes a policy.
func (s *SQLManager) Delete(id string) error {
	_, err := s.db.Exec(s.db.Rebind("DELETE FROM ladon_policy WHERE id=?"), id)
	return err
}

// FindPoliciesForSubject returns Policies (an array of policy) for a given subject
func (s *SQLManager) FindPoliciesForSubject(subject string) (policies Policies, err error) {
	find := func(query string, args ...interface{}) (ids []string, err error) {
		rows, err := s.db.Query(query, args...)
		if err == sql.ErrNoRows {
			return nil, pkg.ErrNotFound
		} else if err != nil {
			return nil, errors.Wrap(err, "")
		}
		defer rows.Close()
		for rows.Next() {
			var urn string
			if err = rows.Scan(&urn); err != nil {
				return nil, errors.Wrap(err, "")
			}
			ids = append(ids, urn)
		}
		return ids, nil
	}

	var query string
	switch s.db.DriverName() {
	case "postgres", "pgx":
		query = "SELECT policy FROM ladon_policy_subject WHERE $1 ~ ('^' || compiled || '$')"
	case "mysql":
		query = "SELECT policy FROM ladon_policy_subject WHERE ? REGEXP BINARY CONCAT('^', compiled, '$') GROUP BY policy"
	}

	if query == "" {
		return nil, errors.Errorf("driver %s not supported", s.db.DriverName())
	}

	subjects, err := find(query, subject)
	if err != nil {
		return policies, err
	}

	for _, id := range subjects {
		p, err := s.Get(id)
		if err != nil {
			return nil, err
		}
		policies = append(policies, p)
	}
	return policies, nil
}

func getLinkedSQL(db *sqlx.DB, table, policy string) ([]string, error) {
	urns := []string{}
	rows, err := db.Query(db.Rebind(fmt.Sprintf("SELECT template FROM %s WHERE policy=?", table)), policy)
	if err == sql.ErrNoRows {
		return nil, pkg.ErrNotFound
	} else if err != nil {
		return nil, errors.Wrap(err, "")
	}

	defer rows.Close()
	for rows.Next() {
		var urn string
		if err = rows.Scan(&urn); err != nil {
			return []string{}, errors.Wrap(err, "")
		}
		urns = append(urns, urn)
	}
	return urns, nil
}

func createLinkSQL(db *sqlx.DB, tx *sql.Tx, table string, p Policy, templates []string) error {
	for _, template := range templates {
		reg, err := compiler.CompileRegex(template, p.GetStartDelimiter(), p.GetEndDelimiter())

		// Execute SQL statement
		query := db.Rebind(fmt.Sprintf("INSERT INTO %s (policy, template, compiled) VALUES (?, ?, ?)", table))
		if _, err = tx.Exec(query, p.GetID(), template, reg.String()); err != nil {
			if rb := tx.Rollback(); rb != nil {
				return rb
			}
			return err
		}
	}
	return nil
}