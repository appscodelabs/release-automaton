/*
Copyright AppsCode Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package lib

import (
	"strconv"
	"strings"

	"github.com/Masterminds/semver"
)

// SemverCollection is a collection of Version instances and implements the sort
// interface. See the sort package for more details.
// https://golang.org/pkg/sort/
type SemverCollection []*semver.Version

// Len returns the length of a collection. The number of Version instances
// on the slice.
func (c SemverCollection) Len() int {
	return len(c)
}

// Less is needed for the sort interface to compare two Version objects on the
// slice. If checks if one is less than the other.
func (c SemverCollection) Less(i, j int) bool {
	vi := c[i]
	mi, _ := vi.SetPrerelease("")
	vj := c[j]
	mj, _ := vj.SetPrerelease("")

	if mi.Equal(&mj) &&
		(vi.Prerelease() == "" || strings.HasPrefix(vi.Prerelease(), "v")) &&
		(vj.Prerelease() == "" || strings.HasPrefix(vj.Prerelease(), "v")) &&
		!(vi.Prerelease() == "" && vj.Prerelease() == "") {

		si := -1
		sj := -1
		if strings.HasPrefix(vi.Prerelease(), "v") {
			si, _ = strconv.Atoi(vi.Prerelease()[1:])
		}
		if strings.HasPrefix(vj.Prerelease(), "v") {
			sj, _ = strconv.Atoi(vj.Prerelease()[1:])
		}
		return si < sj
	}
	return c[i].LessThan(c[j])
}

// Swap is needed for the sort interface to replace the Version objects
// at two different positions in the slice.
func (c SemverCollection) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
