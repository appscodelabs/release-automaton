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

	"github.com/google/go-github/v32/github"
	"github.com/spf13/cobra"
)

func NewCmdKubeVaultCreateRelease() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "create-release",
		Short:             "Create release file",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			rel := CreateKubeVaultReleaseFile()
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
	return api.Release{
		ProductLine:       "KubeVault",
		Release:           "v2020.07.09-beta.0",
		DocsURLTemplate:   "https://kubevault.com/docs/%s",
		KubernetesVersion: "1.12+",
		Projects: []api.IndependentProjects{
			{
				"github.com/kubevault/operator": api.Project{
					Key: "vault-operator",
					ChartNames: []string{
						"vault-catalog",
					},
					Tag: github.String("v0.4.0-beta.0"),
				},
			},
			{
				"github.com/kubevault/unsealer": api.Project{
					Tag: github.String("v0.4.0-beta.0"),
				},
			},
			{
				"github.com/kubevault/cli": api.Project{
					Tag: github.String("v0.4.0-beta.0"),
				},
			},
			{
				"github.com/kubevault/csi-driver": api.Project{
					Key: "csi-vault",
					Tag: github.String("v0.4.0-beta.0"),
				},
			},
			{
				"github.com/kubevault/installer": api.Project{
					Tag: github.String("v0.4.0-beta.0"),
					Commands: []string{
						"make update-charts CHART_VERSION=${TAG}",
					},
				},
			},
			{
				"github.com/appscode/charts": api.Project{
					Charts: []string{
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
				"github.com/kubevault/docs": api.Project{
					Key:           "kubevault",
					Tag:           github.String("v2020.07.09-beta.0"),
					ReleaseBranch: "release-${TAG}",
					Commands: []string{
						"mv ${SCRIPT_ROOT}/releases/${RELEASE}/docs_changelog.md ${WORKSPACE}/docs/CHANGELOG-${RELEASE}.md",
					},
				},
			},
			{
				"github.com/kubevault/website": api.Project{
					Tag:           github.String("v2020.07.09-beta.0"),
					ReleaseBranch: "master",
					Commands: []string{
						"make set-assets-repo ASSETS_REPO_URL=https://github.com/appscode/static-assets",
						"make docs",
						"make set-version VERSION=${TAG}",
					},
					Changelog: api.SkipChangelog,
				},
			},
			// Bundle
			{
				"github.com/kubevault/bundles": api.Project{
					Tag:           github.String("v2020.07.09-beta.0"),
					ReleaseBranch: "release-${TAG}",
					Commands: []string{
						"release-automaton update-bundles --release-file=${SCRIPT_ROOT}/releases/${RELEASE}/release.json --workspace=${WORKSPACE} --charts-dir=charts",
					},
				},
			},
			{
				"github.com/bytebuilders/bundle-registry": api.Project{
					Charts: []string{
						"github.com/kubevault/bundles",
					},
					Changelog: api.SkipChangelog,
				},
			},
			//{
			//	"github.com/kubevault/vault_exporter": api.Project{
			//		Tag: github.String("v0.4.0-beta.0"),
			//	},
			//},
		},
		ExternalProjects: map[string]api.ExternalProject{},
	}
}
