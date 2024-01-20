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

package cmds

import (
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
)

func MustTime(t time.Time, e error) time.Time {
	if e != nil {
		panic(e)
	}
	return t
}

func computeTag(v, prerelease string) string {
	prerelease = strings.TrimPrefix(prerelease, "-")
	if prerelease == "" {
		return v
	}

	sm, err := semver.NewVersion(v)
	if err != nil {
		panic(err)
	}
	// patch versions can't have prerelease component
	if sm.Patch() > 0 {
		return v
	}
	return fmt.Sprintf("%s-%s", v, prerelease)
}

func TagP(v, prerelease string) *string {
	tag := computeTag(v, prerelease)
	return &tag
}
