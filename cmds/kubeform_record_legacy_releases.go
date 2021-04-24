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

func NewCmdKubeformRecordLegacyReleases() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "record-legacy-releases",
		Short:             "Writes legacy releases in releases/legacy_releases.json",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			table := CreateKubeformReleaseTable()

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

			filename := filepath.Join(scriptRoot, "README.md")
			lib.WriteChangelogMarkdown(filename, "release-table.tpl", table)
		},
	}
	return cmd
}

func CreateKubeformReleaseTable() api.ReleaseTable {
	return api.ReleaseTable{
		ProductLine: "Kubeform",
		Releases: []api.ReleaseSummary{
			{
				Release:           "v0.1.0",
				ReleaseDate:       MustTime(time.Parse(time.RFC3339, "2019-11-07T18:43:19Z")),
				KubernetesVersion: "1.12.x+",
				ReleaseURL:        "https://github.com/kubeform/kubeform/releases/tag/v0.1.0",
				ChangelogURL:      "https://github.com/kubeform/kubeform/releases/tag/v0.1.0",
				DocsURL:           "https://kubeform.com/docs/v0.1.0",
			},
		},
	}
}
