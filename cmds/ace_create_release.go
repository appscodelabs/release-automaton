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

func NewCmdAceCreateRelease() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "create-release",
		Short:             "Create release file",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			rel := CreateAceReleaseFile()
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

func CreateAceReleaseFile() api.Release {
	prerelease := ""
	releaseNumber := "v2025.3.14" + prerelease
	return api.Release{
		ProductLine: "ACE",
		Release:     releaseNumber,
		// DocsURLTemplate:   "https://appscode.com/docs/%s",
		KubernetesVersion: "1.26+",
		Projects: []api.IndependentProjects{
			{
				"github.com/appscode-cloud/ui-wizards": api.Project{
					Tag: TagP("v0.14.0", prerelease),
					ChartNames: []string{
						"kubedbcom-mongodb-editor-options",
					},
					Commands: []string{
						"make update-charts CHART_VERSION=${APPSCODE_CLOUD_UI_WIZARDS_TAG}",
						"make gen fmt",
					},
				},
			},
			{
				"github.com/kmodules/resource-metadata": api.Project{
					Tag: TagP("v0.10.0", prerelease),
					Commands: []string{
						"go run cmd/ui-updater/main.go --chart.version=${APPSCODE_CLOUD_UI_WIZARDS_TAG}",
						"make fmt",
					},
				},
			},
			{
				"github.com/kmodules/image-packer": api.Project{
					Tag:           github.String(releaseNumber),
					ReleaseBranch: "release-${TAG}",
				},
				"github.com/kubeops/ui-server": api.Project{
					Tag: TagP("v0.0.54", prerelease),
				},
				"github.com/appscode-cloud/b3": api.Project{
					Tag:           github.String(releaseNumber),
					ReleaseBranch: "release-${TAG}",
				},
				"github.com/kubepack/lib-app": api.Project{
					Tag: TagP("v0.0.54", prerelease),
					Commands: []string{
						`sed -i 's/CHART_VERSION=\${CHART_VERSION:-v[0-9]*\.[0-9]*\.[0-9]*}/CHART_VERSION=\${CHART_VERSION:-${APPSCODE_CLOUD_UI_WIZARDS_TAG}}/g' *.sh`,
						"make fmt",
					},
				},

				"github.com/appscode-cloud/cluster-ui": api.Project{
					Tag: TagP("v0.3.0", prerelease),
					Commands: []string{
						"npm --no-git-tag-version --allow-same-version version ${TAG_WITHOUT_V_PREFIX}",
					},
				},
				"github.com/appscode-cloud/kubedb-ui": api.Project{
					Tag: TagP("v0.3.0", prerelease),
					Commands: []string{
						"npm --no-git-tag-version --allow-same-version version ${TAG_WITHOUT_V_PREFIX}",
					},
				},
				"github.com/appscode-cloud/accounts-ui": api.Project{
					Tag: TagP("v0.3.0", prerelease),
					Commands: []string{
						"npm --no-git-tag-version --allow-same-version version ${TAG_WITHOUT_V_PREFIX}",
					},
				},
			},
			{
				"github.com/kubeops/installer": api.Project{
					Key:           "kubeops-installer",
					Tag:           github.String(releaseNumber),
					ReleaseBranch: "release-${TAG}",
					ChartNames: []string{
						"kubedb-ui-server",
					},
					Commands: []string{
						"./hack/scripts/import-crds.sh",
						"make chart-kube-ui-server CHART_VERSION=${RELEASE} APP_VERSION=${KUBEOPS_UI_SERVER_TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"./hack/scripts/update-chart-dependencies.sh",
						"./hack/scripts/update-catalog.sh",
					},
				},
			},
			{
				"github.com/appscode-cloud/installer": api.Project{
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

						"go run ./catalog/kubedb/fmt/main.go --update-spec=spec.replicationModeDetector.image=ghcr.io/kubedb/replication-mode-detector:${KUBEDB_REPLICATION_MODE_DETECTOR_TAG}",
						"go run ./catalog/kubedb/fmt/main.go --kind=MariaDBVersion --update-spec=spec.archiver.walg.image=${KUBEDB_MARIADB_ARCHIVER_TAG}",
						"go run ./catalog/kubedb/fmt/main.go --kind=MariaDBVersion --update-spec=spec.coordinator.image=ghcr.io/kubedb/mariadb-coordinator:${KUBEDB_MARIADB_COORDINATOR_TAG}",
						"go run ./catalog/kubedb/fmt/main.go --kind=MySQLVersion --update-spec=spec.archiver.walg.image=${KUBEDB_MYSQL_ARCHIVER_TAG}",
						"go run ./catalog/kubedb/fmt/main.go --kind=MySQLVersion --update-spec=spec.coordinator.image=ghcr.io/kubedb/mysql-coordinator:${KUBEDB_MYSQL_COORDINATOR_TAG}",
						"go run ./catalog/kubedb/fmt/main.go --kind=MySQLVersion --update-spec=spec.routerInitContainer.image=ghcr.io/kubedb/mysql-router-init:${KUBEDB_MYSQL_ROUTER_INIT_TAG}",
						"go run ./catalog/kubedb/fmt/main.go --kind=PerconaXtraDBVersion --update-spec=spec.coordinator.image=ghcr.io/kubedb/percona-xtradb-coordinator:${KUBEDB_PERCONA_XTRADB_COORDINATOR_TAG}",
						"go run ./catalog/kubedb/fmt/main.go --kind=PostgresVersion --update-spec=spec.archiver.walg.image=${KUBEDB_POSTGRES_ARCHIVER_TAG}",
						"go run ./catalog/kubedb/fmt/main.go --kind=PostgresVersion --update-spec=spec.coordinator.image=ghcr.io/kubedb/pg-coordinator:${KUBEDB_PG_COORDINATOR_TAG}",
						"go run ./catalog/kubedb/fmt/main.go --kind=RedisVersion --update-spec=spec.coordinator.image=ghcr.io/kubedb/redis-coordinator:${KUBEDB_REDIS_COORDINATOR_TAG}",
						"go run ./catalog/kubedb/fmt/main.go --kind=SinglestoreVersion --update-spec=spec.coordinator.image=ghcr.io/kubedb/singlestore-coordinator:${KUBEDB_SINGLESTORE_COORDINATOR_TAG}",
						"go run ./catalog/kubedb/fmt/main.go --kind=MSSQLServerVersion --update-spec=spec.archiver.walg.image=ghcr.io/kubedb/mssqlserver-archiver:${KUBEDB_MSSQLSERVER_ARCHIVER_TAG}",
						"go run ./catalog/kubedb/fmt/main.go --kind=MSSQLServerVersion --update-spec=spec.coordinator.image=ghcr.io/kubedb/mssql-coordinator:${KUBEDB_MSSQL_COORDINATOR_TAG}",

						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=elasticsearch-dashboard-backup --update-spec=spec.image=ghcr.io/kubedb/dashboard-restic-plugin:${KUBEDB_DASHBOARD_RESTIC_PLUGIN_TAG}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=elasticsearch-dashboard-restore --update-spec=spec.image=ghcr.io/kubedb/dashboard-restic-plugin:${KUBEDB_DASHBOARD_RESTIC_PLUGIN_TAG}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=elasticsearch-backup --update-spec=spec.image=ghcr.io/kubedb/elasticsearch-restic-plugin:${KUBEDB_ELASTICSEARCH_RESTIC_PLUGIN_TAG}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=elasticsearch-restore --update-spec=spec.image=ghcr.io/kubedb/elasticsearch-restic-plugin:${KUBEDB_ELASTICSEARCH_RESTIC_PLUGIN_TAG}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=kubedbmanifest-backup --update-spec=spec.image=ghcr.io/kubedb/kubedb-manifest-plugin:${KUBEDB_KUBEDB_MANIFEST_PLUGIN_TAG}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=kubedbmanifest-restore --update-spec=spec.image=ghcr.io/kubedb/kubedb-manifest-plugin:${KUBEDB_KUBEDB_MANIFEST_PLUGIN_TAG}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=mariadb-backup --update-spec=spec.image=ghcr.io/kubedb/mariadb-restic-plugin:${KUBEDB_MARIADB_RESTIC_PLUGIN_TAG}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=mariadb-csi-snapshotter --update-spec=spec.image=ghcr.io/kubedb/mariadb-csi-snapshotter-plugin:${KUBEDB_MARIADB_CSI_SNAPSHOTTER_PLUGIN_TAG}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=mariadb-restore --update-spec=spec.image=ghcr.io/kubedb/mariadb-restic-plugin:${KUBEDB_MARIADB_RESTIC_PLUGIN_TAG}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=mongodb-backup --update-spec=spec.image=ghcr.io/kubedb/mongodb-restic-plugin:${KUBEDB_MONGODB_RESTIC_PLUGIN_TAG}_$${DB_VERSION}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=mongodb-csi-snapshotter --update-spec=spec.image=ghcr.io/kubedb/mongodb-csi-snapshotter-plugin:${KUBEDB_MONGODB_CSI_SNAPSHOTTER_PLUGIN_TAG}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=mongodb-restore --update-spec=spec.image=ghcr.io/kubedb/mongodb-restic-plugin:${KUBEDB_MONGODB_RESTIC_PLUGIN_TAG}_$${DB_VERSION}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=mysql-backup --update-spec=spec.image=ghcr.io/kubedb/mysql-restic-plugin:${KUBEDB_MYSQL_RESTIC_PLUGIN_TAG}_$${DB_VERSION}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=mysql-csi-snapshotter --update-spec=spec.image=ghcr.io/kubedb/mysql-csi-snapshotter-plugin:${KUBEDB_MYSQL_CSI_SNAPSHOTTER_PLUGIN_TAG}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=mysql-restore --update-spec=spec.image=ghcr.io/kubedb/mysql-restic-plugin:${KUBEDB_MYSQL_RESTIC_PLUGIN_TAG}_$${DB_VERSION}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=opensearch-backup --update-spec=spec.image=ghcr.io/kubedb/elasticsearch-restic-plugin:${KUBEDB_ELASTICSEARCH_RESTIC_PLUGIN_TAG}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=opensearch-restore --update-spec=spec.image=ghcr.io/kubedb/elasticsearch-restic-plugin:${KUBEDB_ELASTICSEARCH_RESTIC_PLUGIN_TAG}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=postgres-backup --update-spec=spec.image=ghcr.io/kubedb/postgres-restic-plugin:${KUBEDB_POSTGRES_RESTIC_PLUGIN_TAG}_$${DB_VERSION}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=mysql-physical-backup --update-spec=spec.image=ghcr.io/kubedb/xtrabackup-restic-plugin:${KUBEDB_XTRABACKUP_RESTIC_PLUGIN_TAG}_$${DB_VERSION}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=mysql-physical-restore --update-spec=spec.image=ghcr.io/kubedb/xtrabackup-restic-plugin:${KUBEDB_XTRABACKUP_RESTIC_PLUGIN_TAG}_$${DB_VERSION}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=postgres-csi-snapshotter --update-spec=spec.image=ghcr.io/kubedb/postgres-csi-snapshotter-plugin:${KUBEDB_POSTGRES_CSI_SNAPSHOTTER_PLUGIN_TAG}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=postgres-physical-backup --update-spec=spec.image=ghcr.io/kubedb/postgres-restic-plugin:${KUBEDB_POSTGRES_RESTIC_PLUGIN_TAG}_$${DB_VERSION}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=postgres-physical-backup-restore --update-spec=spec.image=ghcr.io/kubedb/postgres-restic-plugin:${KUBEDB_POSTGRES_RESTIC_PLUGIN_TAG}_16.1",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=postgres-restore --update-spec=spec.image=ghcr.io/kubedb/postgres-restic-plugin:${KUBEDB_POSTGRES_RESTIC_PLUGIN_TAG}_$${DB_VERSION}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=redis-backup --update-spec=spec.image=ghcr.io/kubedb/redis-restic-plugin:${KUBEDB_REDIS_RESTIC_PLUGIN_TAG}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=redis-restore --update-spec=spec.image=ghcr.io/kubedb/redis-restic-plugin:${KUBEDB_REDIS_RESTIC_PLUGIN_TAG}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=singlestore-backup --update-spec=spec.image=ghcr.io/kubedb/singlestore-restic-plugin:${KUBEDB_SINGLESTORE_RESTIC_PLUGIN_TAG}_$${DB_VERSION}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=singlestore-restore --update-spec=spec.image=ghcr.io/kubedb/singlestore-restic-plugin:${KUBEDB_SINGLESTORE_RESTIC_PLUGIN_TAG}_$${DB_VERSION}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=zookeeper-backup --update-spec=spec.image=ghcr.io/kubedb/zookeeper-restic-plugin:${KUBEDB_ZOOKEEPER_RESTIC_PLUGIN_TAG}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=zookeeper-restore --update-spec=spec.image=ghcr.io/kubedb/zookeeper-restic-plugin:${KUBEDB_ZOOKEEPER_RESTIC_PLUGIN_TAG}",

						"make update-charts CHART_VERSION=${RELEASE} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-kubedb-crd-manager CHART_VERSION=${KUBEDB_CRD_MANAGER_TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-kubedb-provisioner CHART_VERSION=${KUBEDB_PROVISIONER_TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-kubedb-ops-manager CHART_VERSION=${KUBEDB_OPS_MANAGER_TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-kubedb-autoscaler CHART_VERSION=${KUBEDB_AUTOSCALER_TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-kubedb-dashboard CHART_VERSION=${KUBEDB_KIBANA_TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-kubedb-schema-manager CHART_VERSION=${KUBEDB_SCHEMA_MANAGER_TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-kubedb-ui-server CHART_VERSION=${KUBEDB_UI_SERVER_TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-kubedb-webhook-server CHART_VERSION=${KUBEDB_WEBHOOK_SERVER_TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						// crossplane
						"make chart-kubedb-provider-aws CHART_VERSION=${RELEASE} APP_VERSION=${KUBEDB_PROVIDER_AWS_TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-kubedb-provider-azure CHART_VERSION=${RELEASE} APP_VERSION=${KUBEDB_PROVIDER_AZURE_TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-kubedb-provider-gcp CHART_VERSION=${RELEASE} APP_VERSION=${KUBEDB_PROVIDER_GCP_TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",

						"./hack/scripts/update-chart-dependencies.sh",
						"./hack/scripts/update-catalog.sh",
					},
				},
			},
			{
				"github.com/appscode/charts": api.Project{
					ChartRepos: []string{
						"github.com/kubeops/installer",
						"github.com/appscode-cloud/installer",
					},
					Changelog: api.SkipChangelog,
				},
			},
		},
		ExternalProjects: map[string]api.ExternalProject{
			"github.com/kmodules/codespan-schema-checker":       {},
			"github.com/kmodules/metrics-configuration-checker": {},
			"github.com/kubepack/kubepack":                      {},
			"github.com/kubepack/lib-app":                       {},
		},
	}
}
