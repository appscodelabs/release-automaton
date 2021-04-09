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
	"gomodules.xyz/semvers"
)

func NewCmdStashCreateRelease() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "create-release",
		Short:             "Create release file",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			rel := CreateStashReleaseFile()
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

func CreateStashReleaseFile() api.Release {
	prerelease := ""
	releaseNumber := "v2021.04.09" + prerelease
	updateVars := "release-automaton update-vars " +
		"--env-file=${WORKSPACE}/Makefile.env " +
		"--vars=STASH_VERSION=${STASHED_STASH_TAG} " +
		"--vars=STASH_CATALOG_VERSION=${STASH_CATALOG_VERSION} " +
		"--vars=CHART_REGISTRY=${CHART_REGISTRY} " +
		"--vars=CHART_REGISTRY_URL=${CHART_REGISTRY_URL}"
	updateEnvVars := []string{
		updateVars,
		"make add-license fmt",
	}

	return api.Release{
		ProductLine:       "Stash",
		Release:           releaseNumber,
		DocsURLTemplate:   "https://stash.run/docs/%s",
		KubernetesVersion: "1.14+",
		Projects: []api.IndependentProjects{
			{
				"github.com/stashed/apimachinery": api.Project{
					Tag: github.String("v0.12.2" + prerelease),
					ChartNames: []string{
						"stash-crds",
					},
				},
			},
			{
				"github.com/stashed/stash": api.Project{
					Key: "stash-community",
					Tag: github.String("v0.12.2" + prerelease),
					ChartNames: []string{
						"stash-community",
					},
				},
			},
			{
				"github.com/appscode/stash-enterprise": api.Project{
					Key: "stash-enterprise",
					Tag: github.String("v0.12.2" + prerelease),
					ChartNames: []string{
						"stash-enterprise",
						"stash-catalog",
					},
				},
			},
			{
				"github.com/stashed/cli": api.Project{
					// NOT a sub project anymore
					Key: "stash-cli",
					Tag: github.String("v0.12.2" + prerelease),
				},
			},
			{
				"github.com/stashed/postgres": api.Project{
					// NOT a sub project anymore
					Key: "stash-postgres",
					Tags: map[string]string{
						"9.6.19-v6": "release-9.6.19",
						"10.14-v6":  "release-10.14",
						"11.9-v6":   "release-11.9",
						"12.4-v6":   "release-12.4",
						"13.1-v3":   "release-13.1",
					},
				},
			},
			{
				"github.com/stashed/elasticsearch": api.Project{
					// NOT a sub project anymore
					Key: "stash-elasticsearch",
					Tags: map[string]string{
						"5.6.4-v8": "release-5.6.4",
						"6.2.4-v8": "release-6.2.4",
						"6.3.0-v8": "release-6.3.0",
						"6.4.0-v8": "release-6.4.0",
						"6.5.3-v8": "release-6.5.3",
						"6.8.0-v8": "release-6.8.0",
						"7.2.0-v8": "release-7.2.0",
						"7.3.2-v8": "release-7.3.2",
					},
				},
			},
			{
				"github.com/stashed/mongodb": api.Project{
					// NOT a sub project anymore
					Key: "stash-mongodb",
					Tags: map[string]string{
						"3.4.17-v7": "release-3.4.17",
						"3.4.22-v7": "release-3.4.22",
						"3.6.13-v7": "release-3.6.13",
						"3.6.8-v7":  "release-3.6.8",
						"4.0.11-v7": "release-4.0.11",
						"4.0.3-v7":  "release-4.0.3",
						"4.0.5-v7":  "release-4.0.5",
						"4.1.4-v7":  "release-4.1.4",
						"4.1.7-v7":  "release-4.1.7",
						"4.1.13-v7": "release-4.1.13",
						"4.2.3-v7":  "release-4.2.3",
					},
				},
			},
			{
				"github.com/stashed/mysql": api.Project{
					// NOT a sub project anymore
					Key: "stash-mysql",
					Tags: map[string]string{
						"5.7.25-v8": "release-5.7.25",
						"8.0.3-v8":  "release-8.0.3",
						"8.0.14-v8": "release-8.0.14",
						"8.0.21-v2": "release-8.0.21",
					},
				},
			},
			{
				"github.com/stashed/mariadb": api.Project{
					// NOT a sub project anymore
					Key: "stash-mariadb",
					Tags: map[string]string{
						"10.5.8-v2": "release-10.5.8",
					},
				},
			},
			{
				"github.com/stashed/percona-xtradb": api.Project{
					// NOT a sub project anymore
					Key: "stash-percona-xtradb",
					Tags: map[string]string{
						"5.7-v3": "release-5.7",
					},
				},
			},
			{
				"github.com/stashed/installer": api.Project{
					Key: "stash-installer",
					Tag: github.String(releaseNumber),
					Commands: []string{
						"./hack/scripts/import-crds.sh",
						"make update-charts CHART_VERSION=${RELEASE} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-stash-community CHART_VERSION=${STASHED_STASH_TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-stash-enterprise CHART_VERSION=${APPSCODE_STASH_ENTERPRISE_TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"release-automaton stash gen-catalog --release-file=${SCRIPT_ROOT}/releases/${RELEASE}/release.json --catalog-file=${WORKSPACE}/catalog/catalog.json",
						"make gen fmt",
						"helm dependency update charts/stash",
					},
				},
			},
			{
				"github.com/appscode/charts": api.Project{
					Charts: []string{
						"github.com/stashed/installer",
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
				"github.com/stashed/docs": api.Project{
					Key:           "stash",
					Tag:           github.String(releaseNumber),
					ReleaseBranch: "release-${TAG}",
					Commands: []string{
						"mv ${SCRIPT_ROOT}/releases/${RELEASE}/docs_changelog.md ${WORKSPACE}/docs/CHANGELOG-${RELEASE}.md",
					},
				},
			},
			{
				"github.com/stashed/website": api.Project{
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
			{
				"github.com/stashed/bundles": api.Project{
					Tag:           github.String(releaseNumber),
					ReleaseBranch: "release-${TAG}",
					Commands: []string{
						"release-automaton update-bundles --release-file=${SCRIPT_ROOT}/releases/${RELEASE}/release.json --workspace=${WORKSPACE} --charts-dir=charts",
					},
				},
			},
			{
				"github.com/bytebuilders/bundle-registry": api.Project{
					Charts: []string{
						"github.com/stashed/bundles",
					},
					Changelog: api.SkipChangelog,
				},
			},
		},
		ExternalProjects: map[string]api.ExternalProject{
			"github.com/kubedb/apimachinery": {},
			"github.com/kubedb/bundles":      {},
			"github.com/kubedb/cli":          {},
			"github.com/kubedb/docs":         {},
			"github.com/kubedb/elasticsearch": {
				Commands: updateEnvVars,
			},
			"github.com/kubedb/installer": {
				Commands: []string{
					"make fmt",
				},
			},
			"github.com/kubedb/memcached": {},
			"github.com/kubedb/mongodb": {
				Commands: updateEnvVars,
			},
			"github.com/kubedb/mysql": {
				Commands: updateEnvVars,
			},
			"github.com/kubedb/mariadb": {
				Commands: updateEnvVars,
			},
			"github.com/kubedb/operator": {},
			"github.com/kubedb/percona-xtradb": {
				Commands: updateEnvVars,
			},
			"github.com/kubedb/pg-leader-election": {},
			"github.com/kubedb/pgbouncer":          {},
			"github.com/kubedb/postgres": {
				Commands: updateEnvVars,
			},
			"github.com/kubedb/proxysql": {},
			"github.com/kubedb/redis":    {},
		},
	}
}
