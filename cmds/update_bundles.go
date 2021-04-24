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
	"io/ioutil"
	"path/filepath"
	"sort"

	"github.com/appscodelabs/release-automaton/lib"

	"github.com/Masterminds/semver"
	"github.com/spf13/cobra"
	"gomodules.xyz/semvers"
	"gomodules.xyz/sets"
	yu "gomodules.xyz/x/encoding/yaml"
	ylib "gopkg.in/yaml.v2"
	kpapi "kubepack.dev/kubepack/apis/kubepack/v1alpha1"
	"sigs.k8s.io/yaml"
)

var (
	chartsDir = "charts"
)

/*
release-automaton update-bundles \
  --release-file=/home/tamal/go/src/github.com/tamalsaha/gh-release-automation-testing/v2020.6.16/release.json \
  --workspace=/home/tamal/go/src/stash.appscode.dev/bundles \
  --charts-dir=charts
*/
func NewCmdUpdateBundles() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "update-bundles",
		Short:             "Update Kubepack bundles",
		DisableAutoGenTag: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return updateBundles()
		},
	}

	cmd.Flags().StringVar(&releaseFile, "release-file", "", "Path of release file (local file or url is accepted)")
	cmd.Flags().StringVar(&repoWorkspace, "workspace", "", "Path to directory containing git repository")
	cmd.Flags().StringVar(&chartsDir, "charts-dir", chartsDir, "Directory containing bundles in the workspace")
	return cmd
}

func updateBundles() error {
	data, err := ioutil.ReadFile(releaseFile)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, &release)
	if err != nil {
		return err
	}

	dir := filepath.Join(repoWorkspace, chartsDir)

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	bundleSet := sets.NewString()

	for _, fi := range entries {
		if !fi.IsDir() {
			continue
		}
		chartFilename := filepath.Join(dir, fi.Name(), "Chart.yaml")
		if !lib.Exists(chartFilename) {
			continue
		}
		bundleFilename := filepath.Join(dir, fi.Name(), "templates", "bundle.yaml")
		if !lib.Exists(bundleFilename) {
			continue
		}

		// Update chart
		data, err := ioutil.ReadFile(chartFilename)
		if err != nil {
			return err
		}

		var ch ylib.MapSlice
		err = ylib.Unmarshal(data, &ch)
		if err != nil {
			return err
		}

		name, ok, err := yu.NestedString(ch, "name")
		if err != nil {
			return err
		}
		if !ok {
			continue
		}
		bundleSet.Insert(name)
		err = yu.SetNestedField(&ch, release.Release, "version")
		if err != nil {
			return err
		}
		err = yu.SetNestedField(&ch, release.Release, "appVersion")
		if err != nil {
			return err
		}

		data, err = ylib.Marshal(ch)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(chartFilename, data, 0644)
		if err != nil {
			return err
		}
	}

	for _, fi := range entries {
		if !fi.IsDir() {
			continue
		}
		chartFilename := filepath.Join(dir, fi.Name(), "Chart.yaml")
		if !lib.Exists(chartFilename) {
			continue
		}
		bundleFilename := filepath.Join(dir, fi.Name(), "templates", "bundle.yaml")
		if !lib.Exists(bundleFilename) {
			continue
		}

		// Update bundle.yaml
		data, err = ioutil.ReadFile(bundleFilename)
		if err != nil {
			return err
		}

		var bundle kpapi.Bundle
		err = yaml.Unmarshal(data, &bundle)
		if err != nil {
			return err
		}

		for _, pkg := range bundle.Spec.Packages {
			if pkg.Bundle != nil {
				if bundleSet.Has(pkg.Bundle.Name) {
					pkg.Bundle.Version = release.Release
				}
			} else if pkg.Chart != nil {
				if _, project, ok := findProjectByChart(pkg.Chart.Name, release); ok {
					if project.Tag != nil && !pkg.Chart.MultiSelect && len(pkg.Chart.Versions) == 1 {
						pkg.Chart.Versions[0].Version = *project.Tag
					}
					if len(project.Tags) > 0 && pkg.Chart.MultiSelect {
						sort.Slice(pkg.Chart.Versions, func(i, j int) bool {
							return semvers.CompareVersions(semver.MustParse(pkg.Chart.Versions[i].Version), semver.MustParse(pkg.Chart.Versions[j].Version))
						})
						latestDetail := pkg.Chart.Versions[len(pkg.Chart.Versions)-1]

						versions := make([]kpapi.VersionDetail, 0, len(project.Tags))
						for tag := range project.Tags {
							if detail, ok := findVersionDetail(tag, pkg.Chart.Versions); ok {
								detail.Version = tag
								versions = append(versions, detail)
							} else {
								nuDetail := latestDetail.DeepCopy() // use the latestDetail as "default" value
								nuDetail.Version = tag
								versions = append(versions, *nuDetail)
							}
						}
						sort.Slice(versions, func(i, j int) bool {
							return !semvers.CompareVersions(semver.MustParse(versions[i].Version), semver.MustParse(versions[j].Version))
						})
						pkg.Chart.Versions = versions
					}
				}
			}
		}

		data, err = yaml.Marshal(bundle)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(bundleFilename, data, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func findVersionDetail(tag string, versions []kpapi.VersionDetail) (kpapi.VersionDetail, bool) {
	for idx := range versions {
		vTag := semver.MustParse(tag)
		vDetail := semver.MustParse(versions[idx].Version)
		if vTag.Major() == vDetail.Major() && vTag.Minor() == vDetail.Minor() {
			return versions[idx], true
		}
	}
	return kpapi.VersionDetail{}, false
}
