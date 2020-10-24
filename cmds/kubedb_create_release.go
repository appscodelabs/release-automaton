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
	return api.Release{
		ProductLine:       "KubeDB",
		Release:           "v2020.10.24-beta.0",
		DocsURLTemplate:   "https://kubedb.com/docs/%s",
		KubernetesVersion: "1.12+",
		Projects: []api.IndependentProjects{
			{
				"github.com/kubedb/apimachinery": api.Project{
					Tag: github.String("v0.14.0-beta.4"),
				},
			},
			{
				"github.com/kubedb/pg-leader-election": api.Project{
					Tag: github.String("v0.2.0-beta.4"),
					// update catalog
				},
			},
			{
				"github.com/kubedb/cli": api.Project{
					Key: "kubedb-cli",
					Tag: github.String("v0.14.0-beta.4"),
				},
				"github.com/kubedb/elasticsearch": api.Project{
					Tag: github.String("v0.14.0-beta.4"),
				},
				"github.com/kubedb/memcached": api.Project{
					Tag: github.String("v0.7.0-beta.4"),
				},
				"github.com/kubedb/mongodb": api.Project{
					Tag: github.String("v0.7.0-beta.4"),
				},
				"github.com/kubedb/mysql": api.Project{
					Tag: github.String("v0.7.0-beta.4"),
				},
				"github.com/kubedb/postgres": api.Project{
					Tag: github.String("v0.14.0-beta.4"),
				},
				"github.com/kubedb/redis": api.Project{
					Tag: github.String("v0.7.0-beta.4"),
				},
				"github.com/kubedb/percona-xtradb": api.Project{
					Tag: github.String("v0.1.0-beta.4"),
				},
				"github.com/kubedb/mysql-replication-mode-detector": api.Project{
					Tag: github.String("v0.1.0-beta.4"),
					// update catalog
				},
			},
			{
				"github.com/kubedb/pgbouncer": api.Project{
					Tag: github.String("v0.1.0-beta.4"),
					Commands: []string{
						"release-automaton update-vars " +
							"--env-file=${WORKSPACE}/Makefile.env " +
							"--vars=POSTGRES_TAG=${KUBEDB_POSTGRES_TAG} ",
						"make add-license fmt",
					},
				},
				"github.com/kubedb/proxysql": api.Project{
					Tag: github.String("v0.1.0-beta.4"),
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
					Tag: github.String("v0.14.0-beta.4"),
					ChartNames: []string{
						"kubedb",
						"kubedb-catalog",
					},
				},
			},
			{
				"github.com/appscode/kubedb-enterprise": api.Project{
					Key: "kubedb-enterprise",
					Tag: github.String("v0.1.0-beta.4"),
					ChartNames: []string{
						"kubedb-enterprise",
					},
				},
			},
			// Build Enterprise image
			{
				"github.com/kubedb/installer": api.Project{
					Key: "kubedb-installer",
					Tag: github.String("v0.14.0-beta.4"),
					Commands: []string{
						"make chart-kubedb CHART_VERSION=${TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-kubedb-catalog CHART_VERSION=${TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						`sed -i 's|mysql-replication-mode-detector:.*|mysql-replication-mode-detector:${KUBEDB_MYSQL_REPLICATION_MODE_DETECTOR_TAG}"|g' charts/kubedb-catalog/templates/mysql/*`,
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
					Tag:           github.String("v2020.10.24-beta.0"),
					ReleaseBranch: "release-${TAG}",
					Commands: []string{
						"mv ${SCRIPT_ROOT}/releases/${RELEASE}/docs_changelog.md ${WORKSPACE}/docs/CHANGELOG-${RELEASE}.md",
					},
				},
			},
			{
				"github.com/kubedb/website": api.Project{
					Tag:           github.String("v2020.10.24-beta.0"),
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
				"github.com/kubedb/bundles": api.Project{
					Tag:           github.String("v2020.10.24-beta.0"),
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
