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

func NewCmdAceCreateRelease() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "create-release",
		Short:             "Create release file",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			rel := CreateAceReleaseFile()
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

func CreateAceReleaseFile() api.Release {
	prerelease := ""
	releaseNumber := "v2025.12.15" + prerelease
	return api.Release{
		ProductLine: "ACE",
		Release:     releaseNumber,
		// DocsURLTemplate:   "https://appscode.com/docs/%s",
		KubernetesVersion: "1.28+",
		Projects: []api.IndependentProjects{
			{
				"github.com/appscode-cloud/ui-wizards": api.Project{
					Tag: TagP("v0.28.0", prerelease),
					ChartNames: []string{
						"kubedbcom-mongodb-editor-options",
					},
					Commands: []string{
						"make update-charts CHART_VERSION=${APPSCODE_CLOUD_UI_WIZARDS_TAG}",
						"make gen fmt",
					},
				},
			},
			{
				"github.com/kmodules/resource-metadata": api.Project{
					Tag: TagP("v0.40.0", prerelease),
					Commands: []string{
						"go run cmd/ui-updater/main.go --use-digest=false --chart.version=${APPSCODE_CLOUD_UI_WIZARDS_TAG}",
						"make fmt",
					},
				},
			},
			{
				"github.com/kmodules/image-packer": api.Project{
					Tag:           github.String(releaseNumber),
					ReleaseBranch: "release-${TAG}",
				},
				"github.com/appscode/website": api.Project{
					Tag:           github.String(releaseNumber),
					ReleaseBranch: "release-${TAG}",
					Commands: []string{
						"make assets",
					},
				},
				"github.com/kubeops/ui-server": api.Project{
					Tag: TagP("v0.0.68", prerelease),
				},
				"github.com/kubepack/lib-app": api.Project{
					Tag: TagP("v0.16.0", prerelease),
					Commands: []string{
						"make set-version VERSION=${APPSCODE_CLOUD_UI_WIZARDS_TAG}",
						"make fmt",
					},
				},
				/*
					"github.com/appscode-cloud/cluster-ui": api.Project{
						Tag: TagP("v0.3.0", prerelease),
						Commands: []string{
							"npm --no-git-tag-version --allow-same-version version ${TAG_WITHOUT_V_PREFIX}",
						},
					},
					"github.com/appscode-cloud/kubedb-ui": api.Project{
						Tag: TagP("v0.3.0", prerelease),
						Commands: []string{
							"npm --no-git-tag-version --allow-same-version version ${TAG_WITHOUT_V_PREFIX}",
						},
					},
					"github.com/appscode-cloud/accounts-ui": api.Project{
						Tag: TagP("v0.3.0", prerelease),
						Commands: []string{
							"npm --no-git-tag-version --allow-same-version version ${TAG_WITHOUT_V_PREFIX}",
						},
					},
				*/
			},
			{
				"github.com/appscode-cloud/b3": api.Project{
					Tag:           github.String(releaseNumber),
					ReleaseBranch: "release-${TAG}",
				},
			},
			{
				"github.com/kubeops/installer": api.Project{
					Key:           "kubeops-installer",
					Tag:           github.String(releaseNumber),
					ReleaseBranch: "release-${TAG}",
					ChartNames: []string{
						"kube-ui-server",
					},
					Commands: []string{
						"./hack/scripts/import-crds.sh",
						"make chart-kube-ui-server CHART_VERSION=${RELEASE} APP_VERSION=${KUBEOPS_UI_SERVER_TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"./hack/scripts/update-chart-dependencies.sh",
						"./hack/scripts/update-catalog.sh",
					},
				},
			},
			{
				"github.com/appscode-cloud/installer": api.Project{
					Key:           "kubedb-installer",
					Tag:           github.String(releaseNumber),
					ReleaseBranch: "release-${TAG}",
					ChartNames: []string{
						"opscenter-features",
					},
					Commands: []string{
						"./hack/scripts/import-crds.sh",
						"make update-charts CHART_VERSION=${RELEASE} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-acaas CHART_VERSION=${RELEASE} APP_VERSION=${RELEASE} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-accounts-ui CHART_VERSION=${RELEASE} APP_VERSION=${RELEASE} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-ace CHART_VERSION=${RELEASE} APP_VERSION=${RELEASE} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-ace-installer CHART_VERSION=${RELEASE} APP_VERSION=${RELEASE} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-billing CHART_VERSION=${RELEASE} APP_VERSION=${RELEASE} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-kubedb-ui-presets CHART_VERSION=${RELEASE} APP_VERSION=${RELEASE} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-marketplace-api CHART_VERSION=${RELEASE} APP_VERSION=${RELEASE} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-platform-api CHART_VERSION=${RELEASE} APP_VERSION=${RELEASE} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-service-gateway CHART_VERSION=${RELEASE} APP_VERSION=${RELEASE} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-service-gateway-presets CHART_VERSION=${RELEASE} APP_VERSION=${RELEASE} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-service-vault CHART_VERSION=${RELEASE} APP_VERSION=${RELEASE} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-stash-presets CHART_VERSION=${RELEASE} APP_VERSION=${RELEASE} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-website CHART_VERSION=${RELEASE} APP_VERSION=${RELEASE} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						// opscenter-features
						"make chart-opscenter-features CHART_VERSION=${RELEASE} APP_VERSION=${APPSCODE_CLOUD_UI_WIZARDS_TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"go run ./cmd/update-version/main.go",
						"./hack/scripts/update-chart-dependencies.sh",
						"./hack/scripts/update-catalog.sh",
						"make gen fmt",
					},
				},
			},

			{
				"github.com/appscode/charts": api.Project{
					ChartRepos: []string{
						"github.com/kubeops/installer",
						"github.com/appscode-cloud/installer",
					},
					Changelog: api.SkipChangelog,
				},
			},
		},
		ExternalProjects: map[string]api.ExternalProject{
			"github.com/kmodules/codespan-schema-checker":       {},
			"github.com/kmodules/metrics-configuration-checker": {},
			"github.com/kubepack/kubepack":                      {},
			"github.com/kubepack/lib-app":                       {},
			"github.com/appscode-cloud/outbox-syncer":           {},
		},
	}
}
