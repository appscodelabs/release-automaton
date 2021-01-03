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

func NewCmdKubeDBCreateRelease() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "create-release",
		Short:             "Create release file",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			rel := CreateKubeDBReleaseFile()
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

func CreateKubeDBReleaseFile() api.Release {
	prerelease := "-rc.0"
	releaseNumber := "v2021.01.02" + prerelease
	return api.Release{
		ProductLine:       "KubeDB",
		Release:           releaseNumber,
		DocsURLTemplate:   "https://kubedb.com/docs/%s",
		KubernetesVersion: "1.14+",
		Projects: []api.IndependentProjects{
			{
				"github.com/kubedb/apimachinery": api.Project{
					Tag: github.String("v0.16.0" + prerelease),
				},
			},
			{
				"github.com/kubedb/pg-leader-election": api.Project{
					Tag: github.String("v0.4.0" + prerelease),
					// update catalog
				},
			},
			{
				"github.com/kubedb/cli": api.Project{
					Key: "kubedb-cli",
					Tag: github.String("v0.16.0" + prerelease),
				},
				"github.com/kubedb/elasticsearch": api.Project{
					Tag: github.String("v0.16.0" + prerelease),
				},
				"github.com/kubedb/memcached": api.Project{
					Tag: github.String("v0.9.0" + prerelease),
				},
				"github.com/kubedb/mongodb": api.Project{
					Tag: github.String("v0.9.0" + prerelease),
				},
				"github.com/kubedb/mysql": api.Project{
					Tag: github.String("v0.9.0" + prerelease),
				},
				"github.com/kubedb/postgres": api.Project{
					Tag: github.String("v0.16.0" + prerelease),
				},
				"github.com/kubedb/redis": api.Project{
					Tag: github.String("v0.9.0" + prerelease),
				},
				"github.com/kubedb/percona-xtradb": api.Project{
					Tag: github.String("v0.3.0" + prerelease),
				},
				"github.com/kubedb/replication-mode-detector": api.Project{
					Tag: github.String("v0.3.0" + prerelease),
					// update catalog
				},
			},
			{
				"github.com/kubedb/pgbouncer": api.Project{
					Tag: github.String("v0.3.0" + prerelease),
					Commands: []string{
						"release-automaton update-vars " +
							"--env-file=${WORKSPACE}/Makefile.env " +
							"--vars=POSTGRES_TAG=${KUBEDB_POSTGRES_TAG} ",
						"make add-license fmt",
					},
				},
				"github.com/kubedb/proxysql": api.Project{
					Tag: github.String("v0.3.0" + prerelease),
					Commands: []string{
						"release-automaton update-vars " +
							"--env-file=${WORKSPACE}/Makefile.env " +
							"--vars=MYSQL_TAG=${KUBEDB_MYSQL_TAG} " +
							"--vars=PERCONA_XTRADB_TAG=${KUBEDB_PERCONA_XTRADB_TAG} ",
						"make add-license fmt",
					},
				},
			},
			{
				"github.com/kubedb/operator": api.Project{
					Key: "kubedb-community",
					Tag: github.String("v0.16.0" + prerelease),
					ChartNames: []string{
						"kubedb",
						"kubedb-catalog",
					},
				},
			},
			{
				"github.com/appscode/kubedb-enterprise": api.Project{
					Key: "kubedb-enterprise",
					Tag: github.String("v0.3.0" + prerelease),
					ChartNames: []string{
						"kubedb-enterprise",
					},
				},
			},
			{
				"github.com/appscode/kubedb-autoscaler": api.Project{
					Key: "kubedb-autocaler",
					Tag: github.String("v0.0.1" + prerelease),
					ChartNames: []string{
						"kubedb-autoscaler",
					},
				},
			},
			// Build Enterprise image
			{
				"github.com/kubedb/installer": api.Project{
					Key: "kubedb-installer",
					Tag: github.String("v0.16.0" + prerelease),
					Commands: []string{
						"make chart-kubedb CHART_VERSION=${KUBEDB_OPERATOR_TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-kubedb-catalog CHART_VERSION=${KUBEDB_OPERATOR_TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-kubedb-enterprise CHART_VERSION=${APPSCODE_KUBEDB_ENTERPRISE_TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-kubedb-autoscaler CHART_VERSION=${APPSCODE_KUBEDB_AUTOSCALER_TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						// https://stackoverflow.com/a/48290678
						`find charts/kubedb-catalog/templates/mongodb -type f -exec sed -i 's|replication-mode-detector:.*|replication-mode-detector:${KUBEDB_REPLICATION_MODE_DETECTOR_TAG}"|g' {} \;`,
						`find charts/kubedb-catalog/templates/mysql -type f -exec sed -i 's|replication-mode-detector:.*|replication-mode-detector:${KUBEDB_REPLICATION_MODE_DETECTOR_TAG}"|g' {} \;`,
					},
				},
			},
			{
				"github.com/appscode/charts": api.Project{
					Charts: []string{
						"github.com/kubedb/installer",
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
				"github.com/kubedb/docs": api.Project{
					Key:           "kubedb",
					Tag:           github.String(releaseNumber),
					ReleaseBranch: "release-${TAG}",
					Commands: []string{
						"mv ${SCRIPT_ROOT}/releases/${RELEASE}/docs_changelog.md ${WORKSPACE}/docs/CHANGELOG-${RELEASE}.md",
					},
				},
			},
			{
				"github.com/kubedb/website": api.Project{
					Tag:           github.String(releaseNumber),
					ReleaseBranch: "master",
					Commands: lib.AppendIf(
						[]string{
							"make set-assets-repo ASSETS_REPO_URL=https://github.com/appscode/static-assets",
							"make docs",
						},
						api.IsPublicRelease(releaseNumber),
						"make set-version VERSION=${TAG}",
					),
					Changelog: api.SkipChangelog,
				},
			},
			// Bundle
			{
				"github.com/kubedb/bundles": api.Project{
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
						"github.com/kubedb/bundles",
					},
					Changelog: api.SkipChangelog,
				},
			},
		},
		ExternalProjects: map[string]api.ExternalProject{},
	}
}
