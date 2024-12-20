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
	"hash/fnv"
	"strings"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/pkg/errors"

	"github.com/paullesiak/ladon/compiler"
)

func NewRegexpMatcher(size int) *RegexpMatcher {
	if size <= 0 {
		size = 16 * 1024
	}
	size *= 10

	cache, err := ristretto.NewCache[string, bool](
		&ristretto.Config[string, bool]{
			NumCounters: int64(size),
			MaxCost:     1 << 30,
			BufferItems: 64,
			Metrics:     true,
		})
	if err != nil {
		panic(err)
	}
	return &RegexpMatcher{
		Cache: cache,
	}
}

type RegexpMatcher struct {
	Cache *ristretto.Cache[string, bool]
}

func (m *RegexpMatcher) get(pattern string) (bool, bool) {
	val, ok := m.Cache.Get(pattern)
	if !ok {
		return false, false
	}
	return val, true
}

func (m *RegexpMatcher) set(pattern string, match bool) {
	m.Cache.Set(pattern, match, int64(len(pattern)))
}

// Matches a needle with an array of regular expressions and returns true if a match was found.
func (m *RegexpMatcher) Matches(p Policy, haystack []string, needle string) (bool, error) {
	keyHash := fnv.New64()
	_, _ = keyHash.Write([]byte(needle))
	for i := 0; i < len(haystack); i++ {
		_, _ = keyHash.Write([]byte(haystack[i]))
	}
	key := string(keyHash.Sum(nil))
	if matched, ok := m.get(key); ok {
		return matched, nil
	}

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

		matched = reg.MatchString(needle)
		if matched {
			break
		}
	}
	m.set(key, matched)
	return matched, nil
}

func (m *RegexpMatcher) CacheMetrics() *ristretto.Metrics {
	return m.Cache.Metrics
}
