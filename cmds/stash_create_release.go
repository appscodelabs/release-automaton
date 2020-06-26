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
	"sigs.k8s.io/yaml"
)

func NewCmdStashCreateRelease() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "create-release",
		Short:             "Create release file",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			rel := CreateReleaseFile()
			data, err := yaml.Marshal(rel)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(data))

			data, err = lib.MarshalJson(rel)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(data))
		},
	}
	return cmd
}

func CreateReleaseFile() api.Release {
	updateEnvVars := []string{
		"echo STASH_VERSION=${APPSCODE_CLOUD_STASH_TAG} > Makefile.env",
		"echo STASH_CATALOG_VERSION=${STASH_CATALOG_VERSION} >> Makefile.env",
	}
	return api.Release{
		ProductLine:       "Stash",
		Release:           "v2020.6.23",
		DocsURLTemplate:   "https://stash.run/docs/%s",
		KubernetesVersion: "1.12+",
		Projects: []api.IndependentProjects{
			{
				"github.com/appscode-cloud/apimachinery": api.Project{
					Tag: github.String("v0.10.0-alpha.3"),
				},
			},
			{
				"github.com/appscode-cloud/cli": api.Project{
					// NOT a sub project anymore
					// Key: "stash-cli",
					Tag: github.String("v0.10.0-alpha.3"),
				},
			},
			{
				"github.com/appscode-cloud/postgres": api.Project{
					Key: "stash-postgres",
					Tags: map[string]string{
						"9.6-v1":  "release-9.6",
						"10.2-v1": "release-10.2",
						"10.6-v1": "release-10.6",
						"11.1-v1": "release-11.1",
						"11.2-v1": "release-11.2",
					},
					Commands: []string{
						"make update-charts CHART_VERSION=${TAG}",
					},
				},
			},
			{
				"github.com/appscode-cloud/stash": api.Project{
					Tag: github.String("v0.10.0-alpha.3"),
				},
			},
			{
				"github.com/appscode-cloud/installer": api.Project{
					Tag: github.String("v0.10.0-alpha.3"),
					Commands: []string{
						"make chart-stash CHART_VERSION=${TAG}",
					},
				},
			},
			{
				"github.com/appscode-cloud/charts": api.Project{
					Charts: []string{
						"github.com/appscode-cloud/postgres",
						"github.com/appscode-cloud/installer",
					},
					Changelog: api.SkipChangelog,
				},
			},
			{
				"github.com/appscode-cloud/catalog": api.Project{
					Key:           "stash-catalog",
					Tag:           github.String("v2020.6.23"),
					ReleaseBranch: "release-${TAG}",
					Commands: []string{
						"release-automaton stash gen-catalog --release-file=${SCRIPT_ROOT}/CHANGELOG/${RELEASE}/release.json --catalog-file=${WORKSPACE}/catalog.json",
						"make gen fmt",
					},
				},
			},
			{
				// Must come before docs repo, so we can generate the docs_changelog.md
				"github.com/appscode-cloud/static-assets": api.Project{
					Commands: []string{
						"release-automaton update-assets --release-file=${SCRIPT_ROOT}/CHANGELOG/${RELEASE}/release.json --workspace=${WORKSPACE}",
					},
					Changelog: api.StandaloneWebsiteChangelog,
				},
			},
			{
				"github.com/appscode-cloud/docs": api.Project{
					Key:           "stash",
					Tag:           github.String("v2020.6.23"),
					ReleaseBranch: "release-${TAG}",
					Commands: []string{
						"mv ${SCRIPT_ROOT}/CHANGELOG/${RELEASE}/docs_changelog.md ${WORKSPACE}/docs/CHANGELOG-${RELEASE}.md",
					},
				},
			},
			{
				"github.com/appscode-cloud/website": api.Project{
					Tag:           github.String("v2020.6.23"),
					ReleaseBranch: "master",
					Commands: []string{
						"make set-assets-repo ASSETS_REPO_URL=https://github.com/appscode-cloud/static-assets",
						"make docs",
						"make set-version VERSION=${TAG}",
					},
					Changelog: api.SkipChangelog,
				},
			},
			// Bundle
			{
				"github.com/appscode-cloud/bundles": api.Project{
					Tag:           github.String("v2020.6.23"),
					ReleaseBranch: "release-${TAG}",
					Commands: []string{
						"release-automaton update-bundles --release-file=${SCRIPT_ROOT}/CHANGELOG/${RELEASE}/release.json --workspace=${WORKSPACE} --charts-dir=charts",
					},
				},
			},
			{
				"github.com/appscode-cloud/bb-bundles": api.Project{
					Charts: []string{
						"github.com/appscode-cloud/bundles",
					},
					Changelog: api.SkipChangelog,
				},
			},
		},
		ExternalProjects: map[string]api.ExternalProject{
			"github.com/kubedb/apimachinery": {},
			"github.com/kubedb/cli":          {},
			"github.com/kubedb/memcached":    {},
			"github.com/kubedb/operator":     {},
			"github.com/kubedb/pgbouncer":    {},
			"github.com/kubedb/proxysql":     {},
			"github.com/kubedb/redis":        {},
			"github.com/kubedb/elasticsearch": {
				Commands: updateEnvVars,
			},
			"github.com/kubedb/mongodb": {
				Commands: updateEnvVars,
			},
			"github.com/kubedb/mysql": {
				Commands: updateEnvVars,
			},
			"github.com/kubedb/percona-xtradb": {
				Commands: updateEnvVars,
			},
			"github.com/kubedb/postgres": {
				Commands: updateEnvVars,
			},
		},
	}
}
