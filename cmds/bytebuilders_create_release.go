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

	"github.com/google/go-github/v45/github"
	"github.com/spf13/cobra"
)

func NewCmdByteBuildersCreateRelease() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "create-release",
		Short:             "Create release file",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			rel := CreateByteBuildersReleaseFile()
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

func CreateByteBuildersReleaseFile() api.Release {
	prerelease := ""
	releaseNumber := "v2022.03.09" + prerelease
	return api.Release{
		ProductLine: "ByteBuilders",
		Release:     releaseNumber,
		// DocsURLTemplate:   "https://appscode.com/docs/%s",
		KubernetesVersion: "1.16+",
		Projects: []api.IndependentProjects{
			{
				"github.com/bytebuilders/ui-wizards": api.Project{
					Tag: github.String("v0.4.0" + prerelease),
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
					Tag: github.String("v0.10.0" + prerelease),
					Commands: []string{
						"go run cmd/ui-updater/main.go --chart.registry-url=${UI_REGISTRY_URL} --chart.version=${BYTEBUILDERS_UI_WIZARDS_TAG}",
						"make fmt",
					},
				},
			},
			{
				"github.com/bytebuilders/cluster-ui": api.Project{
					Tag: github.String("v0.3.0" + prerelease),
					Commands: []string{
						"npm --no-git-tag-version --allow-same-version version ${TAG_WITHOUT_V_PREFIX}",
					},
				},
				"github.com/bytebuilders/kubedb-ui": api.Project{
					Tag: github.String("v0.3.0" + prerelease),
					Commands: []string{
						"npm --no-git-tag-version --allow-same-version version ${TAG_WITHOUT_V_PREFIX}",
					},
				},
				"github.com/bytebuilders/accounts-ui": api.Project{
					Tag: github.String("v0.3.0" + prerelease),
					Commands: []string{
						"npm --no-git-tag-version --allow-same-version version ${TAG_WITHOUT_V_PREFIX}",
					},
				},
			},
			{
				"github.com/bytebuilders/b3": api.Project{
					Tag:           github.String(releaseNumber),
					ReleaseBranch: "release-${TAG}",
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
			"github.com/bytebuilders/products":                  {},
		},
	}
}
