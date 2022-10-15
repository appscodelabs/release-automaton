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
	prerelease := ""
	releaseNumber := "v2022.10.18" + prerelease
	return api.Release{
		ProductLine:       "KubeDB",
		Release:           releaseNumber,
		DocsURLTemplate:   "https://kubedb.com/docs/%s",
		KubernetesVersion: "1.18+",
		Projects: []api.IndependentProjects{
			{
				"github.com/kubedb/apimachinery": api.Project{
					Tag: github.String("v0.29.0" + prerelease),
				},
			},
			{
				"github.com/kubedb/cli": api.Project{
					Key: "kubedb-cli",
					Tag: github.String("v0.29.0" + prerelease),
				},
				"github.com/kubedb/elasticsearch": api.Project{
					Tag: github.String("v0.29.0" + prerelease),
				},
				"github.com/kubedb/mariadb": api.Project{
					Tag: github.String("v0.13.0" + prerelease),
				},
				"github.com/kubedb/memcached": api.Project{
					Tag: github.String("v0.22.0" + prerelease),
				},
				"github.com/kubedb/mongodb": api.Project{
					Tag: github.String("v0.22.0" + prerelease),
				},
				"github.com/kubedb/mysql": api.Project{
					Tag: github.String("v0.22.0" + prerelease),
				},
				"github.com/kubedb/postgres": api.Project{
					Tag: github.String("v0.29.0" + prerelease),
				},
				"github.com/kubedb/redis": api.Project{
					Tag: github.String("v0.22.0" + prerelease),
				},
				"github.com/kubedb/percona-xtradb": api.Project{
					Tag: github.String("v0.16.0" + prerelease),
				},
				"github.com/kubedb/pg-coordinator": api.Project{
					Tag: github.String("v0.13.0" + prerelease),
					// update catalog
				},
				"github.com/kubedb/mariadb-coordinator": api.Project{
					Tag: github.String("v0.9.0" + prerelease),
					// update catalog
				},
				"github.com/kubedb/mysql-coordinator": api.Project{
					Tag: github.String("v0.7.0" + prerelease),
					// update catalog
				},
				"github.com/kubedb/mysql-router-init": api.Project{
					Tag: github.String("v0.7.0" + prerelease),
					// update catalog
				},
				"github.com/kubedb/percona-xtradb-coordinator": api.Project{
					Tag: github.String("v0.2.0" + prerelease),
					// update catalog
				},
				"github.com/kubedb/redis-coordinator": api.Project{
					Tag: github.String("v0.8.0" + prerelease),
					// update catalog
				},
				"github.com/kubedb/replication-mode-detector": api.Project{
					Tag: github.String("v0.16.0" + prerelease),
					// update catalog
				},
				"github.com/kubedb/tests": api.Project{
					Tag: github.String("v0.14.0" + prerelease),
				},
			},
			{
				"github.com/kubedb/pgbouncer": api.Project{
					Tag: github.String("v0.16.0" + prerelease),
					Commands: []string{
						"release-automaton update-vars " +
							"--env-file=${WORKSPACE}/Makefile.env " +
							"--vars=POSTGRES_TAG=${KUBEDB_POSTGRES_TAG} ",
						"make add-license fmt",
					},
				},
				"github.com/kubedb/proxysql": api.Project{
					Tag: github.String("v0.16.0" + prerelease),
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
				"github.com/kubedb/provisioner": api.Project{
					Key: "kubedb-provisioner",
					Tag: github.String("v0.29.0" + prerelease),
					ChartNames: []string{
						"kubedb-provisioner",
					},
				},
				"github.com/kubedb/dashboard": api.Project{
					Key: "kubedb-dashboard",
					Tag: github.String("v0.5.0" + prerelease),
					ChartNames: []string{
						"kubedb-dashboard",
					},
				},
				"github.com/kubedb/schema-manager": api.Project{
					Key: "kubedb-schema-manager",
					Tag: github.String("v0.5.0" + prerelease),
					ChartNames: []string{
						"kubedb-schema-manager",
					},
				},
				"github.com/kubedb/ui-server": api.Project{
					// NOT a sub project anymore
					Key: "kubedb-ui-server",
					Tag: github.String("v0.5.0" + prerelease),
					ChartNames: []string{
						"kubedb-ui-server",
					},
				},
			},
			{
				"github.com/kubedb/ops-manager": api.Project{
					Key: "kubedb-ops-manager",
					Tag: github.String("v0.16.0" + prerelease),
					ChartNames: []string{
						"kubedb-ops-manager",
					},
				},
			},
			{
				"github.com/kubedb/autoscaler": api.Project{
					Key: "kubedb-autoscaler",
					Tag: github.String("v0.14.0" + prerelease),
					ChartNames: []string{
						"kubedb-autoscaler",
					},
				},
			},
			{
				"github.com/kubedb/webhook-server": api.Project{
					Key: "kubedb-webhook-server",
					Tag: github.String("v0.5.0" + prerelease),
					ChartNames: []string{
						"kubedb-webhook-server",
					},
				},
			},
			// Build Enterprise image
			{
				"github.com/kubedb/installer": api.Project{
					Key:           "kubedb-installer",
					Tag:           github.String(releaseNumber),
					ReleaseBranch: "release-${TAG}",
					ChartNames: []string{
						"kubedb-crds",
						"kubedb-catalog",
						"kubedb",
					},
					Commands: []string{
						"./hack/scripts/import-crds.sh",
						"go run ./hack/fmt/main.go --update-spec=spec.replicationModeDetector.image=kubedb/replication-mode-detector:${KUBEDB_REPLICATION_MODE_DETECTOR_TAG}",
						"go run ./hack/fmt/main.go --kind=MariaDBVersion --update-spec=spec.coordinator.image=kubedb/mariadb-coordinator:${KUBEDB_MARIADB_COORDINATOR_TAG}",
						"go run ./hack/fmt/main.go --kind=MySQLVersion --update-spec=spec.coordinator.image=kubedb/mysql-coordinator:${KUBEDB_MYSQL_COORDINATOR_TAG}",
						"go run ./hack/fmt/main.go --kind=MySQLVersion --update-spec=spec.routerInitContainer.image=kubedb/mysql-router-init:${KUBEDB_MYSQL_ROUTER_INIT_TAG}",
						"go run ./hack/fmt/main.go --kind=PerconaXtraDBVersion --update-spec=spec.coordinator.image=kubedb/percona-xtradb-coordinator:${KUBEDB_PERCONA_XTRADB_COORDINATOR_TAG}",
						"go run ./hack/fmt/main.go --kind=PostgresVersion --update-spec=spec.coordinator.image=kubedb/pg-coordinator:${KUBEDB_PG_COORDINATOR_TAG}",
						"go run ./hack/fmt/main.go --kind=RedisVersion --update-spec=spec.coordinator.image=kubedb/redis-coordinator:${KUBEDB_REDIS_COORDINATOR_TAG}",
						"make update-charts CHART_VERSION=${RELEASE} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-kubedb-provisioner CHART_VERSION=${KUBEDB_PROVISIONER_TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-kubedb-ops-manager CHART_VERSION=${KUBEDB_OPS_MANAGER_TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-kubedb-autoscaler CHART_VERSION=${KUBEDB_AUTOSCALER_TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-kubedb-dashboard CHART_VERSION=${KUBEDB_DASHBOARD_TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-kubedb-schema-manager CHART_VERSION=${KUBEDB_SCHEMA_MANAGER_TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-kubedb-ui-server CHART_VERSION=${KUBEDB_UI_SERVER_TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-kubedb-webhook-server CHART_VERSION=${KUBEDB_WEBHOOK_SERVER_TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"./hack/scripts/update-chart-dependencies.sh",
					},
				},
			},
			{
				"github.com/appscode/charts": api.Project{
					ChartRepos: []string{
						"github.com/kubedb/installer",
					},
					Changelog: api.SkipChangelog,
				},
			},
			{
				// Must come before docs repo, so we can generate the docs_changelog.md
				"github.com/appscode/static-assets": api.Project{
					Commands: []string{
						"release-automaton update-assets --hide --release-file=${SCRIPT_ROOT}/releases/${RELEASE}/release.json --workspace=${WORKSPACE}",
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
						semvers.IsPublicRelease(releaseNumber),
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
					ChartRepos: []string{
						"github.com/kubedb/bundles",
					},
					Changelog: api.SkipChangelog,
				},
			},
		},
		ExternalProjects: map[string]api.ExternalProject{},
	}
}
