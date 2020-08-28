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

func NewCmdStashCreateRelease() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "create-release",
		Short:             "Create release file",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			rel := CreateStashReleaseFile()
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
	releaseNumber := "v2020.08.27"
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
		KubernetesVersion: "1.12+",
		Projects: []api.IndependentProjects{
			{
				"github.com/stashed/apimachinery": api.Project{
					Tag: github.String("v0.10.0"),
				},
			},
			{
				"github.com/stashed/stash": api.Project{
					Key: "stash-community",
					Tag: github.String("v0.10.0"),
					ChartNames: []string{
						"stash",
					},
				},
			},
			{
				"github.com/appscode/stash-enterprise": api.Project{
					Key: "stash-enterprise",
					Tag: github.String("v0.10.0"),
					ChartNames: []string{
						"stash-enterprise",
					},
				},
			},
			{
				"github.com/stashed/cli": api.Project{
					// NOT a sub project anymore
					Key: "stash-cli",
					Tag: github.String("v0.10.0"),
				},
			},
			{
				"github.com/stashed/postgres": api.Project{
					Key: "stash-postgres",
					ChartNames: []string{
						"stash-postgres",
					},
					Tags: map[string]string{
						"9.6-v1":  "release-9.6",
						"10.2-v1": "release-10.2",
						"10.6-v1": "release-10.6",
						"11.1-v1": "release-11.1",
						"11.2-v1": "release-11.2",
					},
					Commands: []string{
						"make update-charts CHART_VERSION=${TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
					},
				},
			},
			{
				"github.com/stashed/elasticsearch": api.Project{
					Key: "stash-elasticsearch",
					ChartNames: []string{
						"stash-elasticsearch",
					},
					Tags: map[string]string{
						"5.6.4-v1": "release-5.6.4",
						"6.2.4-v1": "release-6.2.4",
						"6.3.0-v1": "release-6.3.0",
						"6.4.0-v1": "release-6.4.0",
						"6.5.3-v1": "release-6.5.3",
						"6.8.0-v1": "release-6.8.0",
						"7.2.0-v1": "release-7.2.0",
						"7.3.2-v1": "release-7.3.2",
					},
					Commands: []string{
						"make update-charts CHART_VERSION=${TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
					},
				},
			},
			{
				"github.com/stashed/mongodb": api.Project{
					Key: "stash-mongodb",
					ChartNames: []string{
						"stash-mongodb",
					},
					Tags: map[string]string{
						"3.4.1-v1":  "release-3.4.17",
						"3.4.2-v1":  "release-3.4.22",
						"3.6.1-v1":  "release-3.6.13",
						"3.6.8-v1":  "release-3.6.8",
						"4.0.11-v1": "release-4.0.11",
						"4.0.3-v1":  "release-4.0.3",
						"4.0.5-v1":  "release-4.0.5",
						"4.1.1-v1":  "release-4.1.13",
						"4.1.4-v1":  "release-4.1.4",
						"4.1.7-v1":  "release-4.1.7",
						"4.2.3-v1":  "release-4.2.3",
					},
					Commands: []string{
						"make update-charts CHART_VERSION=${TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
					},
				},
			},
			{
				"github.com/stashed/mysql": api.Project{
					Key: "stash-mysql",
					ChartNames: []string{
						"stash-mysql",
					},
					Tags: map[string]string{
						"5.7.25-v1": "release-5.7.25",
						"8.0.14-v1": "release-8.0.14",
						"8.0.3-v1":  "release-8.0.3",
					},
					Commands: []string{
						"make update-charts CHART_VERSION=${TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
					},
				},
			},
			{
				"github.com/stashed/percona-xtradb": api.Project{
					Key: "stash-percona-xtradb",
					ChartNames: []string{
						"stash-percona-xtradb",
					},
					Tags: map[string]string{
						"5.7-v1": "release-5.7",
					},
					Commands: []string{
						"make update-charts CHART_VERSION=${TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
					},
				},
			},
			{
				"github.com/stashed/installer": api.Project{
					Tag: github.String("v0.10.0"),
					Commands: []string{
						"make update-charts CHART_VERSION=${TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
					},
				},
			},
			{
				"github.com/appscode/charts": api.Project{
					Charts: []string{
						"github.com/stashed/postgres",
						"github.com/stashed/elasticsearch",
						"github.com/stashed/mongodb",
						"github.com/stashed/mysql",
						"github.com/stashed/percona-xtradb",
						"github.com/stashed/installer",
					},
					Changelog: api.SkipChangelog,
				},
			},
			{
				"github.com/stashed/catalog": api.Project{
					Key:           "stash-catalog",
					Tag:           github.String(releaseNumber),
					ReleaseBranch: "release-${TAG}",
					Commands: []string{
						"release-automaton stash gen-catalog --release-file=${SCRIPT_ROOT}/releases/${RELEASE}/release.json --catalog-file=${WORKSPACE}/catalog.json",
						"make gen fmt",
					},
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
			"github.com/kubedb/installer": {},
			"github.com/kubedb/memcached": {},
			"github.com/kubedb/mongodb": {
				Commands: updateEnvVars,
			},
			"github.com/kubedb/mysql": {
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
