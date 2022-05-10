/*
 * Copyright © 2016-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package ladon

import (
	"fmt"
	"strings"

	"regexp"

	"github.com/dgraph-io/ristretto"
	"github.com/pkg/errors"

	"github.com/ory/ladon/compiler"
)

func NewRegexpMatcher(size int) *RegexpMatcher {
	if size <= 0 {
		size = 16 * 1024
	}

	// golang-lru only returns an error if the cache's size is 0. This, we can safely ignore this error.
	// cache, _ := lru.NewARC(size)
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 10e7,
		MaxCost:     1 << 30,
		BufferItems: 64,
	})
	if err != nil {
		panic(err)
	}
	return &RegexpMatcher{
		Cache: cache,
	}
}

type RegexpMatcher struct {
	Cache *ristretto.Cache

	C map[string]*regexp.Regexp
}

func (m *RegexpMatcher) get(pattern string) (bool, bool) {
	if val, ok := m.Cache.Get(pattern); !ok {
		return false, false
	} else if match, ok := val.(bool); !ok {
		return false, true
	} else {
		return match, true
	}
}

func (m *RegexpMatcher) set(pattern string, match bool) {
	m.Cache.Set(pattern, match, 0)
}

// Matches a needle with an array of regular expressions and returns true if a match was found.
func (m *RegexpMatcher) Matches(p Policy, haystack []string, needle string) (bool, error) {
	key := fmt.Sprintf("%s+%s", strings.Join(haystack, ","), needle)
	if matched, ok := m.get(key); ok {
		if matched {
			return true, nil
		} else {
			return false, nil
		}
	}

	//var reg *regexp.Regexp
	//var err error
	matched := false
	for _, h := range haystack {

		// This means that the current haystack item does not contain a regular expression
		if strings.Count(h, string(p.GetStartDelimiter())) == 0 {
			// If we have a simple string match, we've got a match!
			if h == needle {
				matched = true
				break
			}

			// Not string match, but also no regexp, continue with next haystack item
			continue
		}

		reg, err := compiler.CompileRegex(h, p.GetStartDelimiter(), p.GetEndDelimiter())
		if err != nil {
			return false, errors.WithStack(err)
		}

		matched, err = reg.MatchString(needle)
		if matched {
			break
		}
	}
	m.set(key, matched)
	return matched, nil
}
