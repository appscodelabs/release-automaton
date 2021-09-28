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

func NewCmdKubeformCreateRelease() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "create-release",
		Short:             "Create release file",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			rel := CreateKubeformReleaseFile()
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

func CreateKubeformReleaseFile() api.Release {
	prerelease := ""
	releaseNumber := "v2021.08.02" + prerelease
	tag := "v0.3.0"
	providers := []string{
		"alicloud",
		"aws",
		"azurerm",
		"civo",
		"datadog",
		"digitalocean",
		"dynatrace",
		"ec",
		"equinixmetal",
		"google",
		"grafana",
		"hcloud",
		"ibm",
		"linode",
		"mongodbatlas",
		"newrelic",
		"ovh",
		"pagerduty",
		"rediscloud",
		"upcloud",
		"vsphere",
		"vultr",
		"wavefront",
	}
	rel := api.Release{
		ProductLine:       "Kubeform",
		Release:           releaseNumber,
		DocsURLTemplate:   "https://kubeform.com/docs/%s",
		KubernetesVersion: "1.16+",
		Projects: []api.IndependentProjects{
			{
				// api projects
				//"github.com/kubeform/provider-aws-api": api.Project{
				//	Tag: github.String("v0.2.0" + prerelease),
				//},
			},
			{
				// controller projects
				//"github.com/kubeform/provider-aws-controller": api.Project{
				//	Key: "kubeform-aws",
				//	Tag: github.String("v0.2.0" + prerelease),
				//	ChartNames: []string{
				//		"kubeform-provider-aws",
				//	},
				//},
			},
			{
				"github.com/kubeform/installer": api.Project{
					Key:           "kubeform-installer",
					Tag:           github.String(releaseNumber),
					ReleaseBranch: "release-${TAG}",
					Commands: []string{
						"./hack/scripts/prepare-release.sh",
						"./hack/scripts/update-chart-dependencies.sh",
					},
				},
			},
			{
				"github.com/appscode/charts": api.Project{
					ChartRepos: []string{
						"github.com/kubeform/installer",
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
				"github.com/kubeform/kubeform": api.Project{
					Key:           "kubeform",
					Tag:           github.String(releaseNumber),
					ReleaseBranch: "release-${TAG}",
					Commands: []string{
						"mv ${SCRIPT_ROOT}/releases/${RELEASE}/docs_changelog.md ${WORKSPACE}/docs/CHANGELOG-${RELEASE}.md",
					},
				},
			},
			{
				"github.com/kubeform/website": api.Project{
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
				"github.com/kubeform/bundles": api.Project{
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
						"github.com/kubeform/bundles",
					},
					Changelog: api.SkipChangelog,
				},
			},
		},
	}

	apiProjects := rel.Projects[0]
	ctrlProjects := rel.Projects[1]
	for _, p := range providers {
		apiProjects[fmt.Sprintf("github.com/kubeform/provider-%s-api", p)] = api.Project{
			Tag: github.String(tag + prerelease),
		}
		ctrlProjects[fmt.Sprintf("github.com/kubeform/provider-%s-controller", p)] = api.Project{
			Key: fmt.Sprintf("kubeform-%s", p),
			Tag: github.String(tag + prerelease),
			ChartNames: []string{
				fmt.Sprintf("kubeform-provider-%s", p),
			},
		}
	}
	return rel
}
