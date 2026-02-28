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

func NewCmdVirtualSecretsCreateRelease() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "create-release",
		Short:             "Create release file",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			rel := CreateVirtualSecretsReleaseFile()
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

func CreateVirtualSecretsReleaseFile() api.Release {
	prerelease := ""
	releaseNumber := "v2026.2.27" + prerelease
	return api.Release{
		ProductLine:       "VirtualSecrets",
		Release:           releaseNumber,
		DocsURLTemplate:   "https://virtual-secrets.dev/docs/%s",
		KubernetesVersion: "1.28+",
		Projects: []api.IndependentProjects{
			{
				"github.com/virtual-secrets/apimachinery": api.Project{
					Tag: TagP("v0.1.0", prerelease),
				},
			},
			{
				"github.com/virtual-secrets/server": api.Project{
					Tag: TagP("v0.3.0", prerelease),
				},
				"github.com/virtual-secrets/csi-provider": api.Project{
					Tag: TagP("v0.1.0", prerelease),
				},
			},
			{
				"github.com/virtual-secrets/installer": api.Project{
					Key:           "virtual-secrets-installer",
					Tag:           github.String(releaseNumber),
					ReleaseBranch: "release-${TAG}",
					ChartNames: []string{
						"virtual-secrets-server",
					},
					Commands: []string{
						"./hack/scripts/import-crds.sh",
						"make chart-virtual-secrets-server CHART_VERSION=${RELEASE} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL} APP_VERSION=${VIRTUAL_SECRETS_SERVER_TAG}",
						"make chart-secrets-store-csi-driver-provider-virtual-secrets CHART_VERSION=${RELEASE} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL} APP_VERSION=${VIRTUAL_SECRETS_CSI_PROVIDER_TAG}",
						"./hack/scripts/update-catalog.sh",
					},
				},
			},
			{
				"github.com/appscode/charts": api.Project{
					ChartRepos: []string{
						"github.com/virtual-secrets/installer",
					},
					Changelog: api.SkipChangelog,
				},
			},
			/*
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
					"github.com/virtual-secrets/docs": api.Project{
						Key:           "virtualsecrets",
						Tag:           github.String(releaseNumber),
						ReleaseBranch: "release-${TAG}",
						Commands: []string{
							"mv ${SCRIPT_ROOT}/releases/${RELEASE}/docs_changelog.md ${WORKSPACE}/docs/CHANGELOG-${RELEASE}.md",
						},
					},
				},
				{
					"github.com/virtual-secrets/website": api.Project{
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
			*/
		},
	}
}
