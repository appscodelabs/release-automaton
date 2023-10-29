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

func NewCmdKubeStashCreateRelease() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "create-release",
		Short:             "Create release file",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			rel := CreateKubeStashReleaseFile()
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

func CreateKubeStashReleaseFile() api.Release {
	prerelease := "-rc.0"
	releaseNumber := "v2023.10.30" + prerelease
	return api.Release{
		ProductLine:       "KubeStash",
		Release:           releaseNumber,
		DocsURLTemplate:   "https://kubestash.com/docs/%s",
		KubernetesVersion: "1.20+",
		Projects: []api.IndependentProjects{
			{
				"github.com/kubestash/apimachinery": api.Project{Tag: github.String("v0.2.0" + prerelease)},
			},
			{
				"github.com/kubestash/kubestash":          api.Project{Tag: github.String("v0.2.0" + prerelease)},
				"github.com/kubestash/kubedb-manifest":    api.Project{Tag: github.String("v0.2.0" + prerelease)},
				"github.com/kubestash/mongodb":            api.Project{Tag: github.String("v0.1.0" + prerelease)},
				"github.com/kubestash/pvc":                api.Project{Tag: github.String("v0.1.0" + prerelease)},
				"github.com/kubestash/workload":           api.Project{Tag: github.String("v0.1.0" + prerelease)},
				"github.com/kubestash/cli":                api.Project{Tag: github.String("v0.1.0" + prerelease)},
				"github.com/kubestash/mysql":              api.Project{Tag: github.String("v0.1.0" + prerelease)},
				"github.com/kubestash/postgres":           api.Project{Tag: github.String("v0.1.0" + prerelease)},
				"github.com/kubestash/kubedump":           api.Project{Tag: github.String("v0.1.0" + prerelease)},
				"github.com/kubestash/volume-snapshotter": api.Project{Tag: github.String("v0.1.0" + prerelease)},
				"github.com/kubestash/elasticsearch":      api.Project{Tag: github.String("v0.1.0" + prerelease)},
				"github.com/kubestash/redis":              api.Project{Tag: github.String("v0.1.0" + prerelease)},
			},
			{
				"github.com/kubestash/installer": api.Project{
					Key:           "kubestash-installer",
					Tag:           github.String(releaseNumber),
					ReleaseBranch: "release-${TAG}",
					ChartNames: []string{
						"kubestash",
						"kubestash-operator",
					},
					Commands: []string{
						"./hack/scripts/import-crds.sh",
						"make update-charts CHART_VERSION=${RELEASE} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-kubestash-operator CHART_VERSION=${KUBESTASH_KUBESTASH_TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"./hack/scripts/update-chart-dependencies.sh",
					},
				},
			},
			{
				"github.com/appscode/charts": api.Project{
					ChartRepos: []string{
						"github.com/kubestash/installer",
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
					"github.com/kubestash/docs": api.Project{
						Key:           "voyager",
						Tag:           github.String(releaseNumber),
						ReleaseBranch: "release-${TAG}",
						Commands: []string{
							"mv ${SCRIPT_ROOT}/releases/${RELEASE}/docs_changelog.md ${WORKSPACE}/docs/CHANGELOG-${RELEASE}.md",
						},
					},
				},
				{
					"github.com/kubestash/website": api.Project{
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
