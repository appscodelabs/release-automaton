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
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/appscodelabs/release-automaton/api"
	"github.com/appscodelabs/release-automaton/lib"

	"github.com/Masterminds/semver/v3"
	"github.com/spf13/cobra"
	"gomodules.xyz/semvers"
)

func NewCmdReleaseReadme() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "readme",
		Short:             "Generate release readme",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			table := GenerateTable()

			filename := filepath.Join(scriptRoot, "README.md")
			lib.WriteChangelogMarkdown(filename, "release-table.tpl", table)
		},
	}
	return cmd
}

func GenerateTable() api.ReleaseTable {
	var table api.ReleaseTable

	legacyfilename := filepath.Join(changelogRoot, "legacy_releases.json")
	data, err := os.ReadFile(legacyfilename)
	if !os.IsNotExist(err) {
		err = json.Unmarshal(data, &table)
		if err != nil {
			panic(err)
		}
	}

	entries, err := ioutil.ReadDir(changelogRoot)
	if err != nil {
		panic(err)
	}
	for _, fi := range entries {
		if !fi.IsDir() {
			continue
		}
		filename := filepath.Join(changelogRoot, fi.Name(), "CHANGELOG.json")
		if lib.Exists(filename) {
			var chlog api.Changelog
			data, err = os.ReadFile(filename)
			if err != nil {
				panic(err)
			}
			err = json.Unmarshal(data, &chlog)
			if err != nil {
				panic(err)
			}
			table.Releases = append(table.Releases, api.ReleaseSummary{
				Release:           chlog.Release,
				ReleaseDate:       chlog.ReleaseDate,
				KubernetesVersion: chlog.KubernetesVersion,
				ReleaseURL:        path.Join(chlog.ReleaseProjectURL, "releases", "tag", chlog.Release),
				ChangelogURL:      path.Join("/", api.ReleasesDir, chlog.Release, "README.md"),
				DocsURL:           chlog.DocsURL,
			})
		}
	}

	// Now keep the full releases and last rc
	var releases []api.ReleaseSummary
	for _, r := range table.Releases {
		v := semver.MustParse(r.Release)
		if v.Prerelease() == "" ||
			strings.HasPrefix(v.Prerelease(), "v") ||
			strings.HasPrefix(v.Prerelease(), "rc.") {
			releases = append(releases, r)
		}
	}
	sort.Slice(releases, func(i, j int) bool {
		return !semvers.CompareVersions(semver.MustParse(releases[i].Release), semver.MustParse(releases[j].Release))
	})
	for i := range releases {
		releases[i].ReleaseDate = releases[i].ReleaseDate.UTC()
	}
	table.Releases = releases

	return table
}
