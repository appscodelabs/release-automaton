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
	prerelease := "-rc.0"
	releaseNumber := "v2025.12.9" + prerelease
	return api.Release{
		ProductLine:       "KubeDB",
		Release:           releaseNumber,
		DocsURLTemplate:   "https://kubedb.com/docs/%s",
		KubernetesVersion: "1.26+",
		Projects: []api.IndependentProjects{
			{
				"github.com/kubedb/apimachinery": api.Project{Tag: TagP("v0.60.0", prerelease)},
			},
			{
				"github.com/kubedb/db-client-go": api.Project{Tag: TagP("v0.15.0", prerelease)},
			},
			{
				"github.com/kubedb/cli": api.Project{
					Key: "kubedb-cli",
					Tag: TagP("v0.60.0", prerelease),
				},
				"github.com/kubedb/crd-manager":                api.Project{Tag: TagP("v0.15.0", prerelease)},
				"github.com/kubedb/clickhouse":                 api.Project{Tag: TagP("v0.15.0", prerelease)},
				"github.com/kubedb/druid":                      api.Project{Tag: TagP("v0.15.0", prerelease)},
				"github.com/kubedb/ferretdb":                   api.Project{Tag: TagP("v0.15.0", prerelease)},
				"github.com/kubedb/pgpool":                     api.Project{Tag: TagP("v0.15.0", prerelease)},
				"github.com/kubedb/rabbitmq":                   api.Project{Tag: TagP("v0.15.0", prerelease)},
				"github.com/kubedb/singlestore":                api.Project{Tag: TagP("v0.15.0", prerelease)},
				"github.com/kubedb/singlestore-coordinator":    api.Project{Tag: TagP("v0.15.0", prerelease)},
				"github.com/kubedb/solr":                       api.Project{Tag: TagP("v0.15.0", prerelease)},
				"github.com/kubedb/zookeeper":                  api.Project{Tag: TagP("v0.15.0", prerelease)},
				"github.com/kubedb/elasticsearch":              api.Project{Tag: TagP("v0.60.0", prerelease)},
				"github.com/kubedb/kafka":                      api.Project{Tag: TagP("v0.31.0", prerelease)},
				"github.com/kubedb/mariadb":                    api.Project{Tag: TagP("v0.44.0", prerelease)},
				"github.com/kubedb/mariadb-coordinator":        api.Project{Tag: TagP("v0.40.0", prerelease)},
				"github.com/kubedb/memcached":                  api.Project{Tag: TagP("v0.53.0", prerelease)},
				"github.com/kubedb/mongodb":                    api.Project{Tag: TagP("v0.53.0", prerelease)},
				"github.com/kubedb/mssqlserver":                api.Project{Tag: TagP("v0.15.0", prerelease)},
				"github.com/kubedb/mssql-coordinator":          api.Project{Tag: TagP("v0.15.0", prerelease)},
				"github.com/kubedb/mysql":                      api.Project{Tag: TagP("v0.53.0", prerelease)},
				"github.com/kubedb/mysql-coordinator":          api.Project{Tag: TagP("v0.38.0", prerelease)},
				"github.com/kubedb/mysql-router-init":          api.Project{Tag: TagP("v0.38.0", prerelease)},
				"github.com/kubedb/percona-xtradb":             api.Project{Tag: TagP("v0.47.0", prerelease)},
				"github.com/kubedb/percona-xtradb-coordinator": api.Project{Tag: TagP("v0.33.0", prerelease)},
				"github.com/kubedb/pg-coordinator":             api.Project{Tag: TagP("v0.44.0", prerelease)},
				"github.com/kubedb/postgres":                   api.Project{Tag: TagP("v0.60.0", prerelease)},
				"github.com/kubedb/redis":                      api.Project{Tag: TagP("v0.53.0", prerelease)},
				"github.com/kubedb/redis-coordinator":          api.Project{Tag: TagP("v0.39.0", prerelease)},
				"github.com/kubedb/replication-mode-detector":  api.Project{Tag: TagP("v0.47.0", prerelease)},
				"github.com/kubedb/tests":                      api.Project{Tag: TagP("v0.45.0", prerelease)},
				"github.com/kubedb/cassandra":                  api.Project{Tag: TagP("v0.13.0", prerelease)},
				"github.com/kubedb/hazelcast":                  api.Project{Tag: TagP("v0.6.0", prerelease)},
				"github.com/kubedb/oracle":                     api.Project{Tag: TagP("v0.6.0", prerelease)},
				"github.com/kubedb/oracle-coordinator":         api.Project{Tag: TagP("v0.6.0", prerelease)},
				"github.com/kubedb/db2":                        api.Project{Tag: TagP("v0.1.0", prerelease)},
				"github.com/kubedb/db2-coordinator":            api.Project{Tag: TagP("v0.1.0", prerelease)},
				"github.com/kubedb/hanadb":                     api.Project{Tag: TagP("v0.1.0", prerelease)},
				"github.com/kubedb/milvus":                     api.Project{Tag: TagP("v0.1.0", prerelease)},
				"github.com/kubedb/neo4j":                      api.Project{Tag: TagP("v0.1.0", prerelease)},
				"github.com/kubedb/qdrant":                     api.Project{Tag: TagP("v0.1.0", prerelease)},
				"github.com/kubedb/weaviate":                   api.Project{Tag: TagP("v0.1.0", prerelease)},
				// kubestash plugins
				"github.com/kubedb/dashboard-restic-plugin":         api.Project{Tag: TagP("v0.18.0", prerelease)},
				"github.com/kubedb/elasticsearch-restic-plugin":     api.Project{Tag: TagP("v0.23.0", prerelease)},
				"github.com/kubedb/kubedb-manifest-plugin":          api.Project{Tag: TagP("v0.23.0", prerelease)},
				"github.com/kubedb/mariadb-archiver":                api.Project{Tag: TagP("v0.20.0", prerelease)},
				"github.com/kubedb/mongodb-csi-snapshotter-plugin":  api.Project{Tag: TagP("v0.21.0", prerelease)},
				"github.com/kubedb/mariadb-restic-plugin":           api.Project{Tag: TagP("v0.18.0", prerelease)},
				"github.com/kubedb/mongodb-restic-plugin":           api.Project{Tag: TagP("v0.23.0", prerelease)},
				"github.com/kubedb/mysql-archiver":                  api.Project{Tag: TagP("v0.21.0", prerelease)},
				"github.com/kubedb/mysql-csi-snapshotter-plugin":    api.Project{Tag: TagP("v0.21.0", prerelease)},
				"github.com/kubedb/mysql-restic-plugin":             api.Project{Tag: TagP("v0.23.0", prerelease)},
				"github.com/kubedb/xtrabackup-restic-plugin":        api.Project{Tag: TagP("v0.8.0", prerelease)},
				"github.com/kubedb/postgres-archiver":               api.Project{Tag: TagP("v0.21.0", prerelease)},
				"github.com/kubedb/postgres-csi-snapshotter-plugin": api.Project{Tag: TagP("v0.21.0", prerelease)},
				"github.com/kubedb/postgres-restic-plugin":          api.Project{Tag: TagP("v0.23.0", prerelease)},
				"github.com/kubedb/redis-restic-plugin":             api.Project{Tag: TagP("v0.23.0", prerelease)},
				"github.com/kubedb/singlestore-restic-plugin":       api.Project{Tag: TagP("v0.18.0", prerelease)},
				"github.com/kubedb/zookeeper-restic-plugin":         api.Project{Tag: TagP("v0.16.0", prerelease)},
				"github.com/kubedb/mssqlserver-walg-plugin":         api.Project{Tag: TagP("v0.14.0", prerelease)},
				"github.com/kubedb/mssqlserver-archiver":            api.Project{Tag: TagP("v0.14.0", prerelease)},
				"github.com/kubedb/gitops":                          api.Project{Tag: TagP("v0.8.0", prerelease)},
				"github.com/kubedb/kubedb-verifier":                 api.Project{Tag: TagP("v0.11.0", prerelease)},
				"github.com/kubedb/ignite":                          api.Project{Tag: TagP("v0.7.0", prerelease)},
				"github.com/kubedb/cassandra-medusa-plugin":         api.Project{Tag: TagP("v0.7.0", prerelease)},
				// crossplane
				"github.com/kubedb/provider-aws":   api.Project{Tag: TagP("v0.21.0", prerelease)},
				"github.com/kubedb/provider-azure": api.Project{Tag: TagP("v0.21.0", prerelease)},
				"github.com/kubedb/provider-gcp":   api.Project{Tag: TagP("v0.21.0", prerelease)},
			},
			{
				"github.com/kubedb/mariadb-csi-snapshotter-plugin": api.Project{Tag: TagP("v0.20.0", prerelease)},
				"github.com/kubedb/pgbouncer": api.Project{
					Tag: TagP("v0.47.0", prerelease),
					Commands: []string{
						"release-automaton update-vars " +
							"--env-file=${WORKSPACE}/Makefile.env " +
							"--vars=POSTGRES_TAG=${KUBEDB_POSTGRES_TAG} ",
						"make add-license fmt",
					},
				},
				"github.com/kubedb/proxysql": api.Project{
					Tag: TagP("v0.47.0", prerelease),
					Commands: []string{
						"release-automaton update-vars " +
							"--env-file=${WORKSPACE}/Makefile.env " +
							"--vars=MYSQL_TAG=${KUBEDB_MYSQL_TAG} " +
							"--vars=PERCONA_XTRADB_TAG=${KUBEDB_PERCONA_XTRADB_TAG} ",
						"make add-license fmt",
					},
				},
				"github.com/kubedb/kibana": api.Project{
					Key: "kubedb-dashboard",
					Tag: TagP("v0.36.0", prerelease),
					ChartNames: []string{
						"kubedb-dashboard",
					},
				},
			},
			{
				"github.com/kubedb/provisioner": api.Project{
					Key: "kubedb-provisioner",
					Tag: TagP("v0.60.0", prerelease),
					ChartNames: []string{
						"kubedb-provisioner",
					},
				},
				"github.com/kubedb/schema-manager": api.Project{
					Key: "kubedb-schema-manager",
					Tag: TagP("v0.36.0", prerelease),
					ChartNames: []string{
						"kubedb-schema-manager",
					},
				},
				"github.com/kubedb/ui-server": api.Project{
					// NOT a sub project anymore
					Key: "kubedb-ui-server",
					Tag: TagP("v0.36.0", prerelease),
					ChartNames: []string{
						"kubedb-ui-server",
					},
				},
			},
			{
				"github.com/kubedb/ops-manager": api.Project{
					Key: "kubedb-ops-manager",
					Tag: TagP("v0.47.0", prerelease),
					ChartNames: []string{
						"kubedb-ops-manager",
					},
				},
			},
			{
				"github.com/kubedb/autoscaler": api.Project{
					Key: "kubedb-autoscaler",
					Tag: TagP("v0.45.0", prerelease),
					ChartNames: []string{
						"kubedb-autoscaler",
					},
				},
			},
			{
				"github.com/kubedb/webhook-server": api.Project{
					Key: "kubedb-webhook-server",
					Tag: TagP("v0.36.0", prerelease),
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

						"go run ./catalog/kubedb/fmt/main.go --update-spec=spec.replicationModeDetector.image=ghcr.io/kubedb/replication-mode-detector:${KUBEDB_REPLICATION_MODE_DETECTOR_TAG}",
						"go run ./catalog/kubedb/fmt/main.go --kind=MariaDBVersion --update-spec=spec.archiver.walg.image=${KUBEDB_MARIADB_ARCHIVER_TAG}",
						"go run ./catalog/kubedb/fmt/main.go --kind=MariaDBVersion --update-spec=spec.coordinator.image=ghcr.io/kubedb/mariadb-coordinator:${KUBEDB_MARIADB_COORDINATOR_TAG}",
						"go run ./catalog/kubedb/fmt/main.go --kind=MySQLVersion --update-spec=spec.archiver.walg.image=${KUBEDB_MYSQL_ARCHIVER_TAG}",
						"go run ./catalog/kubedb/fmt/main.go --kind=MySQLVersion --update-spec=spec.coordinator.image=ghcr.io/kubedb/mysql-coordinator:${KUBEDB_MYSQL_COORDINATOR_TAG}",
						"go run ./catalog/kubedb/fmt/main.go --kind=MySQLVersion --update-spec=spec.routerInitContainer.image=ghcr.io/kubedb/mysql-router-init:${KUBEDB_MYSQL_ROUTER_INIT_TAG}",
						"go run ./catalog/kubedb/fmt/main.go --kind=OracleVersion --update-spec=spec.coordinator.image=ghcr.io/kubedb/oracle-coordinator:${KUBEDB_ORACLE_COORDINATOR_TAG}",
						"go run ./catalog/kubedb/fmt/main.go --kind=DB2Version --update-spec=spec.coordinator.image=ghcr.io/kubedb/db2-coordinator:${KUBEDB_DB2_COORDINATOR_TAG}",
						"go run ./catalog/kubedb/fmt/main.go --kind=PerconaXtraDBVersion --update-spec=spec.coordinator.image=ghcr.io/kubedb/percona-xtradb-coordinator:${KUBEDB_PERCONA_XTRADB_COORDINATOR_TAG}",
						"go run ./catalog/kubedb/fmt/main.go --kind=PostgresVersion --update-spec=spec.archiver.walg.image=${KUBEDB_POSTGRES_ARCHIVER_TAG}",
						"go run ./catalog/kubedb/fmt/main.go --kind=PostgresVersion --update-spec=spec.coordinator.image=ghcr.io/kubedb/pg-coordinator:${KUBEDB_PG_COORDINATOR_TAG}",
						"go run ./catalog/kubedb/fmt/main.go --kind=RedisVersion --update-spec=spec.coordinator.image=ghcr.io/kubedb/redis-coordinator:${KUBEDB_REDIS_COORDINATOR_TAG}",
						"go run ./catalog/kubedb/fmt/main.go --kind=SinglestoreVersion --update-spec=spec.coordinator.image=ghcr.io/kubedb/singlestore-coordinator:${KUBEDB_SINGLESTORE_COORDINATOR_TAG}",
						"go run ./catalog/kubedb/fmt/main.go --kind=MSSQLServerVersion --update-spec=spec.archiver.walg.image=ghcr.io/kubedb/mssqlserver-archiver:${KUBEDB_MSSQLSERVER_ARCHIVER_TAG}",
						"go run ./catalog/kubedb/fmt/main.go --kind=MSSQLServerVersion --update-spec=spec.coordinator.image=ghcr.io/kubedb/mssql-coordinator:${KUBEDB_MSSQL_COORDINATOR_TAG}",

						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=cassandra-backup --update-spec=spec.image=ghcr.io/kubedb/cassandra-medusa-plugin:${KUBEDB_CASSANDRA_MEDUSA_PLUGIN_TAG}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=cassandra-restore --update-spec=spec.image=ghcr.io/kubedb/cassandra-medusa-plugin:${KUBEDB_CASSANDRA_MEDUSA_PLUGIN_TAG}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=elasticsearch-dashboard-backup --update-spec=spec.image=ghcr.io/kubedb/dashboard-restic-plugin:${KUBEDB_DASHBOARD_RESTIC_PLUGIN_TAG}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=elasticsearch-dashboard-restore --update-spec=spec.image=ghcr.io/kubedb/dashboard-restic-plugin:${KUBEDB_DASHBOARD_RESTIC_PLUGIN_TAG}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=elasticsearch-backup --update-spec=spec.image=ghcr.io/kubedb/elasticsearch-restic-plugin:${KUBEDB_ELASTICSEARCH_RESTIC_PLUGIN_TAG}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=elasticsearch-restore --update-spec=spec.image=ghcr.io/kubedb/elasticsearch-restic-plugin:${KUBEDB_ELASTICSEARCH_RESTIC_PLUGIN_TAG}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=kubedbmanifest-backup --update-spec=spec.image=ghcr.io/kubedb/kubedb-manifest-plugin:${KUBEDB_KUBEDB_MANIFEST_PLUGIN_TAG}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=kubedbmanifest-restore --update-spec=spec.image=ghcr.io/kubedb/kubedb-manifest-plugin:${KUBEDB_KUBEDB_MANIFEST_PLUGIN_TAG}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=mariadb-backup --update-spec=spec.image=ghcr.io/kubedb/mariadb-restic-plugin:${KUBEDB_MARIADB_RESTIC_PLUGIN_TAG}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=mariadb-csi-snapshotter --update-spec=spec.image=ghcr.io/kubedb/mariadb-csi-snapshotter-plugin:${KUBEDB_MARIADB_CSI_SNAPSHOTTER_PLUGIN_TAG}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=mariadb-restore --update-spec=spec.image=ghcr.io/kubedb/mariadb-restic-plugin:${KUBEDB_MARIADB_RESTIC_PLUGIN_TAG}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=mariadb-physical-backup --update-spec=spec.image=ghcr.io/kubedb/mariadb-restic-plugin:${KUBEDB_MARIADB_RESTIC_PLUGIN_TAG}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=mariadb-physical-restore --update-spec=spec.image=ghcr.io/kubedb/mariadb-restic-plugin:${KUBEDB_MARIADB_RESTIC_PLUGIN_TAG}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=mongodb-backup --update-spec=spec.image=ghcr.io/kubedb/mongodb-restic-plugin:${KUBEDB_MONGODB_RESTIC_PLUGIN_TAG}_$${DB_VERSION}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=mongodb-csi-snapshotter --update-spec=spec.image=ghcr.io/kubedb/mongodb-csi-snapshotter-plugin:${KUBEDB_MONGODB_CSI_SNAPSHOTTER_PLUGIN_TAG}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=mongodb-restore --update-spec=spec.image=ghcr.io/kubedb/mongodb-restic-plugin:${KUBEDB_MONGODB_RESTIC_PLUGIN_TAG}_$${DB_VERSION}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=mssqlserver-backup --update-spec=spec.image=ghcr.io/kubedb/mssqlserver-walg-plugin:${KUBEDB_MSSQLSERVER_WALG_PLUGIN_TAG}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=mssqlserver-restore --update-spec=spec.image=ghcr.io/kubedb/mssqlserver-walg-plugin:${KUBEDB_MSSQLSERVER_WALG_PLUGIN_TAG}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=mysql-backup --update-spec=spec.image=ghcr.io/kubedb/mysql-restic-plugin:${KUBEDB_MYSQL_RESTIC_PLUGIN_TAG}_$${DB_VERSION}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=mysql-csi-snapshotter --update-spec=spec.image=ghcr.io/kubedb/mysql-csi-snapshotter-plugin:${KUBEDB_MYSQL_CSI_SNAPSHOTTER_PLUGIN_TAG}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=mysql-restore --update-spec=spec.image=ghcr.io/kubedb/mysql-restic-plugin:${KUBEDB_MYSQL_RESTIC_PLUGIN_TAG}_$${DB_VERSION}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=opensearch-dashboard-backup --update-spec=spec.image=ghcr.io/kubedb/dashboard-restic-plugin:${KUBEDB_DASHBOARD_RESTIC_PLUGIN_TAG}",
						"go run ./catalog/kubestash/fmt/main.go --kind=Function --name=opensearch-dashboard-restore --update-spec=spec.image=ghcr.io/kubedb/dashboard-restic-plugin:${KUBEDB_DASHBOARD_RESTIC_PLUGIN_TAG}",
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
						"make chart-kubedb-autoscaler CHART_VERSION=${KUBEDB_AUTOSCALER_TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-kubedb-crd-manager CHART_VERSION=${KUBEDB_CRD_MANAGER_TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-kubedb-dashboard CHART_VERSION=${KUBEDB_KIBANA_TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-kubedb-gitops CHART_VERSION=${KUBEDB_GITOPS_TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-kubedb-ops-manager CHART_VERSION=${KUBEDB_OPS_MANAGER_TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
						"make chart-kubedb-provisioner CHART_VERSION=${KUBEDB_PROVISIONER_TAG} CHART_REGISTRY=${CHART_REGISTRY} CHART_REGISTRY_URL=${CHART_REGISTRY_URL}",
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
			// {
			// 	"github.com/kubedb/bundles": api.Project{
			// 		Tag:           github.String(releaseNumber),
			// 		ReleaseBranch: "release-${TAG}",
			// 		Commands: []string{
			// 			"release-automaton update-bundles --release-file=${SCRIPT_ROOT}/releases/${RELEASE}/release.json --workspace=${WORKSPACE} --charts-dir=charts",
			// 		},
			// 	},
			// },
			// {
			// 	"github.com/bytebuilders/bundle-registry": api.Project{
			// 		ChartRepos: []string{
			// 			"github.com/kubedb/bundles",
			// 		},
			// 		Changelog: api.SkipChangelog,
			// 	},
			// },
		},
		ExternalProjects: map[string]api.ExternalProject{},
	}
}
