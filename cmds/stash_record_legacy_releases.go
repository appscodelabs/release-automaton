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
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/appscodelabs/release-automaton/api"
	"github.com/appscodelabs/release-automaton/lib"

	"github.com/spf13/cobra"
)

func NewCmdStashRecordLegacyReleases() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "record-legacy-releases",
		Short:             "Writes legacy releases in releases/legacy_releases.json",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			table := CreateStashReleaseTable()

			data, err := lib.MarshalJson(table)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(data))

			err = os.MkdirAll(changelogRoot, 0755)
			if err != nil {
				panic(err)
			}

			legacyfilename := filepath.Join(changelogRoot, "legacy_releases.json")
			err = ioutil.WriteFile(legacyfilename, data, 0644)
			if err != nil {
				panic(err)
			}

			filename := filepath.Join(changelogRoot, "README.md")
			lib.WriteChangelogMarkdown(filename, "release-table.tpl", table)
		},
	}
	return cmd
}

func CreateStashReleaseTable() api.ReleaseTable {
	return api.ReleaseTable{
		ProductLine: "Stash",
		Releases: []api.ReleaseSummary{
			{
				Release:           "v0.9.0-rc.6",
				ReleaseDate:       MustTime(time.Parse(time.RFC3339, "2020-02-24T14:19:15Z")),
				KubernetesVersion: "1.11.x+",
				ReleaseURL:        "https://github.com/stashed/stash/releases/tag/v0.9.0-rc.6",
				ChangelogURL:      "https://github.com/stashed/stash/releases/tag/v0.9.0-rc.6",
				DocsURL:           "https://stash.run/docs/v0.9.0-rc.6",
			},
			{
				Release:           "0.8.3",
				ReleaseDate:       MustTime(time.Parse(time.RFC3339, "2019-02-19T02:57:32Z")),
				KubernetesVersion: "1.9.x+",
				ReleaseURL:        "https://github.com/stashed/stash/releases/tag/0.8.3",
				ChangelogURL:      "https://github.com/stashed/stash/releases/tag/0.8.3",
				DocsURL:           "https://stash.run/docs/0.8.3",
			},
			{
				Release:           "0.7.0",
				ReleaseDate:       MustTime(time.Parse(time.RFC3339, "2018-05-29T01:13:56Z")),
				KubernetesVersion: "1.8.x",
				ReleaseURL:        "https://github.com/stashed/stash/releases/tag/0.7.0",
				ChangelogURL:      "https://github.com/stashed/stash/releases/tag/0.7.0",
				DocsURL:           "https://stash.run/docs/0.7.0",
			},
			{
				Release:           "0.6.4",
				ReleaseDate:       MustTime(time.Parse(time.RFC3339, "2018-02-20T07:56:43Z")),
				KubernetesVersion: "1.7.x",
				ReleaseURL:        "https://github.com/stashed/stash/releases/tag/0.6.4",
				ChangelogURL:      "https://github.com/stashed/stash/releases/tag/0.6.4",
				DocsURL:           "https://stash.run/docs/0.6.4",
			},
			{
				Release:           "0.4.2",
				ReleaseDate:       MustTime(time.Parse(time.RFC3339, "2017-11-03T14:31:55Z")),
				KubernetesVersion: "1.5.x - 1.6.x",
				ReleaseURL:        "https://github.com/stashed/stash/releases/tag/0.4.2",
				ChangelogURL:      "https://github.com/stashed/stash/releases/tag/0.4.2",
				DocsURL:           "https://github.com/stashed/docs/tree/0.4.2/docs",
			},
		},
	}
}

func MustTime(t time.Time, e error) time.Time {
	if e != nil {
		panic(e)
	}
	return t
}
