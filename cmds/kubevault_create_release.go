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
	"gomodules.xyz/semvers"
)

func NewCmdKubeVaultCreateRelease() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "create-release",
		Short:             "Create release file",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			rel := CreateKubeVaultReleaseFile()
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

func CreateKubeVaultReleaseFile() api.Release {
	prerelease := ""
	releaseNumber := "v2025.5.26" + prerelease
	return api.Release{
		ProductLine:       "KubeVault",
		Release:           releaseNumber,
		DocsURLTemplate:   "https://kubevault.com/docs/%s",
		KubernetesVersion: "1.26+",
		Projects: []api.IndependentProjects{
			{
				"github.com/kubevault/apimachinery": api.Project{
					Tag: TagP("v0.21.0", prerelease),
				},
				"github.com/kubevault/unsealer": api.Project{
					Key: "kubevault-unsealer",
					Tag: TagP("v0.21.0", prerelease),
				},
			},
			{
				"github.com/kubevault/operator": api.Project{
					Key: "kubevault-operator",
					ChartNames: []string{
						"kubevault-operator",
					},
					Tag: TagP("v0.21.0", prerelease),
				},
				"github.com/kubevault/cli": api.Project{
					Key: "kubevault-cli",
					Tag: TagP("v0.21.0", prerelease),
				},
				// {
				// 	"github.com/kubevault/prometheus-exporter": api.Project{
				// 		Tag: TagP("v0.6.0" , prerelease),
				// 	},
				// },
			},
			{
				"github.com/kubevault/installer": api.Project{
					Key:           "kubevault-installer",
					Tag:           github.String(releaseNumber),
					ReleaseBranch: "release-${TAG}",
					ChartNames: []string{
						"kubevault-crds",
						"kubevault-catalog",
						"kubevault",
					},
					Commands: []string{
						"./hack/scripts/import-crds.sh",
						"go run ./hack/fmt/main.go --update-spec=spec.unsealer.image=kubevault/vault-unsealer:${KUBEVAULT_UNSEALER_TAG}",
						"make update-charts CHART_VERSION=${RELEASE} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-kubevault-operator CHART_VERSION=${KUBEVAULT_OPERATOR_TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-kubevault-webhook-server CHART_VERSION=${KUBEVAULT_OPERATOR_TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"./hack/scripts/update-chart-dependencies.sh",
						"./hack/scripts/update-catalog.sh",
					},
				},
			},
			{
				"github.com/appscode/charts": api.Project{
					ChartRepos: []string{
						"github.com/kubevault/installer",
					},
					Changelog: api.SkipChangelog,
				},
			},
			{
				// Must come before docs repo, so we can generate the docs_changelog.md
				"github.com/appscode/static-assets": api.Project{
					Commands: []string{
						"release-automaton update-assets --release-file=${SCRIPT_ROOT}/releases/${RELEASE}/release.json --workspace=${WORKSPACE}",
					},
					Changelog: api.StandaloneWebsiteChangelog,
				},
			},
			{
				"github.com/kubevault/kubevault": api.Project{
					Key:           "kubevault",
					Tag:           github.String(releaseNumber),
					ReleaseBranch: "release-${TAG}",
					Commands: []string{
						"mv ${SCRIPT_ROOT}/releases/${RELEASE}/docs_changelog.md ${WORKSPACE}/docs/CHANGELOG-${RELEASE}.md",
					},
				},
			},
			{
				"github.com/kubevault/website": api.Project{
					Tag:           github.String(releaseNumber),
					ReleaseBranch: "master",
					Commands: lib.AppendIf(
						[]string{
							"make set-assets-repo ASSETS_REPO_URL=https://github.com/appscode/static-assets",
							"make docs",
						},
						semvers.IsPublicRelease(releaseNumber),
						"make set-version VERSION=${TAG}",
					),
					Changelog: api.SkipChangelog,
				},
			},
			// Bundle
			// {
			// 	"github.com/kubevault/bundles": api.Project{
			// 		Tag:           github.String(releaseNumber),
			// 		ReleaseBranch: "release-${TAG}",
			// 		Commands: []string{
			// 			"release-automaton update-bundles --release-file=${SCRIPT_ROOT}/releases/${RELEASE}/release.json --workspace=${WORKSPACE} --charts-dir=charts",
			// 		},
			// 	},
			// },
			// {
			// 	"github.com/bytebuilders/bundle-registry": api.Project{
			// 		ChartRepos: []string{
			// 			"github.com/kubevault/bundles",
			// 		},
			// 		Changelog: api.SkipChangelog,
			// 	},
			// },
		},
		ExternalProjects: map[string]api.ExternalProject{},
	}
}
