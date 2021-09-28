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

	"github.com/appscodelabs/release-automaton/api"
	"github.com/appscodelabs/release-automaton/lib"

	"github.com/google/go-github/v35/github"
	"github.com/spf13/cobra"
)

func NewCmdConsoleCreateRelease() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "create-release",
		Short:             "Create release file",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			rel := CreateConsoleReleaseFile()
			err := rel.Validate()
			if err != nil {
				panic(err)
			}
			data, err := lib.MarshalJson(rel)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(data))
		},
	}
	return cmd
}

func CreateConsoleReleaseFile() api.Release {
	prerelease := "-rc.0"
	releaseNumber := "v2021.09.27" + prerelease
	return api.Release{
		ProductLine: "Console",
		Release:     releaseNumber,
		// DocsURLTemplate:   "https://appscode.com/docs/%s",
		KubernetesVersion: "1.16+",
		Projects: []api.IndependentProjects{
			{
				"github.com/bytebuilders/ui-wizards": api.Project{
					Tag: github.String("v0.2.0" + prerelease),
					ChartNames: []string{
						"kubedbcom-mongodb-editor-options",
					},
					Commands: []string{
						"make update-charts CHART_VERSION=${BYTEBUILDERS_UI_WIZARDS_TAG} CHART_REGISTRY=${UI_REGISTRY} CHART_REGISTRY_URL=${UI_REGISTRY_URL}",
						"make gen fmt",
					},
				},
			},
			{
				"github.com/bytebuilders/bundle-registry": api.Project{
					ChartRepos: []string{
						"github.com/bytebuilders/ui-wizards",
					},
					Changelog: api.SkipChangelog,
				},
			},
			{
				"github.com/kmodules/resource-metadata": api.Project{
					Tag: github.String("v0.6.0" + prerelease),
					Commands: []string{
						"go run cmd/ui-updater/*.go --chart.registry-url=${UI_REGISTRY_URL} --chart.version=${BYTEBUILDERS_UI_WIZARDS_TAG}",
						"go run cmd/resourcedescriptor-fmt/*.go",
					},
				},
			},
			{
				"github.com/appscode/cluster-ui": api.Project{
					Tag: github.String(releaseNumber),
				},
				"github.com/appscode/kubedb-ui": api.Project{
					Tag: github.String(releaseNumber),
				},
				"github.com/appscode/accounts-ui": api.Project{
					Tag: github.String(releaseNumber),
				},
			},
			{
				"github.com/appscode/gitea": api.Project{
					Tag: github.String(releaseNumber),
				},
			},
			// installer
			// deploy to QA
			// deploy to BB
		},
		ExternalProjects: map[string]api.ExternalProject{
			"github.com/kmodules/codespan-schema-checker":       {},
			"github.com/kmodules/metrics-configuration-checker": {},
			"github.com/kubepack/kubepack":                      {},
			"github.com/kubepack/lib-app":                       {},
			"github.com/appscode/products":                      {},
		},
	}
}
