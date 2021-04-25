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
	"gomodules.xyz/semvers"
)

func NewCmdVoyagerCreateRelease() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "create-release",
		Short:             "Create release file",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			rel := CreateVoyagerReleaseFile()
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

func CreateVoyagerReleaseFile() api.Release {
	prerelease := "-rc.0"
	releaseNumber := "v2021.04.24" + prerelease
	return api.Release{
		ProductLine:       "Voyager",
		Release:           releaseNumber,
		DocsURLTemplate:   "https://voyagermesh.com/docs/%s",
		KubernetesVersion: "1.14+",
		Projects: []api.IndependentProjects{
			{
				"github.com/voyagermesh/voyager": api.Project{
					Tag: github.String("v13.0.0" + prerelease),
				},
			},
			{
				"github.com/voyagermesh/installer": api.Project{
					Key:           "voyager-installer",
					Tag:           github.String(releaseNumber),
					ReleaseBranch: "release-${TAG}",
					ChartNames: []string{
						"voyager",
						"voyager-crds",
					},
					Commands: []string{
						"./hack/scripts/import-crds.sh",
						"make update-charts CHART_VERSION=${RELEASE} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL} APP_VERSION=${VOYAGERMESH_VOYAGER_TAG}",
					},
				},
			},
			{
				"github.com/appscode/charts": api.Project{
					Charts: []string{
						"github.com/voyagermesh/installer",
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
				"github.com/voyagermesh/docs": api.Project{
					Key:           "voyager",
					Tag:           github.String(releaseNumber),
					ReleaseBranch: "release-${TAG}",
					Commands: []string{
						"mv ${SCRIPT_ROOT}/releases/${RELEASE}/docs_changelog.md ${WORKSPACE}/docs/CHANGELOG-${RELEASE}.md",
					},
				},
			},
			{
				"github.com/voyagermesh/website": api.Project{
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
		},
	}
}
