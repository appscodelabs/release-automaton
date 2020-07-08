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

func NewCmdKubeVaultRecordLegacyReleases() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "record-legacy-releases",
		Short:             "Writes legacy releases in releases/legacy_releases.json",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			table := CreateKubeVaultReleaseTable()

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

func CreateKubeVaultReleaseTable() api.ReleaseTable {
	return api.ReleaseTable{
		ProductLine: "KubeVault",
		Releases: []api.ReleaseSummary{
			{
				Release:           "v0.3.0",
				ReleaseDate:       MustTime(time.Parse(time.RFC3339, "2020-01-13T16:50:26Z")),
				KubernetesVersion: "1.12.x+",
				ReleaseURL:        "https://github.com/kubevault/operator/releases/tag/v0.3.0",
				ChangelogURL:      "https://github.com/kubevault/operator/releases/tag/v0.3.0",
				DocsURL:           "https://kubevault.com/docs/v0.3.0",
			},
			{
				Release:           "0.2.0",
				ReleaseDate:       MustTime(time.Parse(time.RFC3339, "2019-03-01T10:50:26Z")),
				KubernetesVersion: "1.11.x+",
				ReleaseURL:        "https://github.com/kubevault/operator/releases/tag/0.2.0",
				ChangelogURL:      "https://github.com/kubevault/operator/releases/tag/0.2.0",
				DocsURL:           "https://kubevault.com/docs/0.2.0",
			},
			{
				Release:           "0.1.0",
				ReleaseDate:       MustTime(time.Parse(time.RFC3339, "2019-03-01T10:49:06Z")),
				KubernetesVersion: "1.11.x",
				ReleaseURL:        "https://github.com/kubevault/operator/releases/tag/0.1.0",
				ChangelogURL:      "https://github.com/kubevault/operator/releases/tag/0.1.0",
				DocsURL:           "https://kubevault.com/docs/0.1.0",
			},
		},
	}
}
