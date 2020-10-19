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
	"strings"

	"github.com/appscodelabs/release-automaton/api"
	"github.com/appscodelabs/release-automaton/lib"

	"github.com/Masterminds/semver"
	stringz "github.com/appscode/go/strings"
	saapi "github.com/appscode/static-assets/api"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/yaml"
)

var (
	repoWorkspace string
)

/*
release-automaton update-assets \
  --release-file=${SCRIPT_ROOT}/v2020.6.16/release.json \
  --workspace=${WORKSPACE}
*/
func NewCmdUpdateAssets() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "update-assets",
		Short:             "Update static assets",
		DisableAutoGenTag: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return updateAssets()
		},
	}

	cmd.Flags().StringVar(&releaseFile, "release-file", "", "Path of release file (local file or url is accepted)")
	cmd.Flags().StringVar(&repoWorkspace, "workspace", "", "Path to directory containing git repository")
	return cmd
}

func updateAssets() error {
	data, err := ioutil.ReadFile(releaseFile)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, &release)
	if err != nil {
		return err
	}

	for _, projects := range release.Projects {
		for _, project := range projects {
			if project.Key == "" {
				continue
			}
			err = updateAsset(release, project)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func updateAsset(release api.Release, project api.Project) error {
	filename := filepath.Join(repoWorkspace, "data", "products", project.Key+".json")
	if !lib.Exists(filename) {
		// Avoid missing product files like stash-catalog key
		return nil
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	var prod saapi.Product
	err = yaml.Unmarshal(data, &prod)
	if err != nil {
		return err
	}

	if project.Tag != nil {
		tag := *project.Tag
		if !findProductVersion(tag, prod.Versions) {
			nuV := saapi.ProductVersion{
				Version:  tag,
				HostDocs: true,
				Show:     showDocs(tag),
			}
			nuV.Info = generateInfo(prod, release)
			prod.Versions = append(prod.Versions, nuV)
		}

		for subKey, ref := range prod.SubProjects {
			if subKey == "stash-cli" {
				// NOT a sub project anymore
				continue
			}
			if _, subProject, ok := findProjectByKey(subKey, release); ok {
				subTags := sets.NewString()
				if subProject.Tag != nil {
					subTags.Insert(*subProject.Tag)
				} else if len(subProject.Tags) > 0 {
					subTags.Insert(lib.Keys(subProject.Tags)...)
				}
				for _, subTag := range subTags.UnsortedList() {
					for idx, mapping := range ref.Mappings {
						if stringz.Contains(mapping.SubProjectVersions, subTag) {
							vs, err := api.SortVersions(sets.NewString(mapping.Versions...).Insert(*project.Tag).UnsortedList())
							if err != nil {
								return err
							}
							mapping.Versions = vs
							ref.Mappings[idx] = mapping
							subTags.Delete(subTag)
							break
						}
					}
				}
				if subTags.Len() > 0 {
					subVersions, err := api.SortVersions(subTags.UnsortedList())
					if err != nil {
						return err
					}
					ref.Mappings = append(ref.Mappings, saapi.Mapping{
						Versions: []string{
							*project.Tag,
						},
						SubProjectVersions: subVersions,
					})
				}
				sort.Slice(ref.Mappings, func(i, j int) bool {
					return !api.CompareVersions(
						semver.MustParse(ref.Mappings[i].Versions[0]),
						semver.MustParse(ref.Mappings[j].Versions[0]),
					)
				})
			}
			prod.SubProjects[subKey] = ref
		}
	} else if project.Tags != nil {
		for tag := range project.Tags {
			if !findProductVersion(tag, prod.Versions) {
				prod.Versions = append(prod.Versions, saapi.ProductVersion{
					Version:  tag,
					HostDocs: true,
					Show:     showDocs(tag),
				})
			}
		}
	}
	prod.Versions, prod.LatestVersion = sortProductVersions(prod.Versions)

	data, err = lib.MarshalJson(prod)
	if err != nil {
		panic(err)
	}
	return ioutil.WriteFile(filename, data, 0644)
}

//nolint:unparam
func findProjectByKey(key string, release api.Release) (string, api.Project, bool) {
	for _, projects := range release.Projects {
		for repoURL, project := range projects {
			if project.Key == key {
				return repoURL, project, true
			}
		}
	}
	return "", api.Project{}, false
}

//nolint:unparam
func findProjectByChart(chartName string, release api.Release) (string, api.Project, bool) {
	for _, projects := range release.Projects {
		for repoURL, project := range projects {
			if stringz.Contains(project.ChartNames, chartName) {
				return repoURL, project, true
			}
		}
	}
	return "", api.Project{}, false
}

func findProductVersion(x string, versions []saapi.ProductVersion) bool {
	for _, v := range versions {
		if v.Version == x {
			return true
		}
	}
	return false
}

func showDocs(version string) bool {
	if version == api.BranchMaster {
		return false
	}
	v := semver.MustParse(version)
	return v.Prerelease() == "" ||
		strings.HasPrefix(v.Prerelease(), "rc.") ||
		strings.HasPrefix(v.Prerelease(), "v")
}

func sortProductVersions(versions []saapi.ProductVersion) ([]saapi.ProductVersion, string) {
	var m saapi.ProductVersion

	data := versions
	for i := range versions {
		if versions[i].Version == api.BranchMaster {
			m = versions[i]
			data = append(versions[:i], versions[i+1:]...)
			break
		}
	}

	// sort
	sort.Slice(data, func(i, j int) bool {
		return !api.CompareVersions(semver.MustParse(data[i].Version), semver.MustParse(data[j].Version))
	})
	latestVersion := data[0].Version
	for i := range data {
		v := semver.MustParse(data[i].Version)
		if strings.HasPrefix(v.Prerelease(), "alpha.") || strings.HasPrefix(v.Prerelease(), "beta.") {
			continue
		}
		// Use the latest non alpha/beta release
		latestVersion = data[i].Version
		break
	}

	// inject to the top
	if m.Version == api.BranchMaster {
		data = append([]saapi.ProductVersion{m}, data...)
	}
	return data, latestVersion
}

func generateInfo(p saapi.Product, release api.Release) map[string]interface{} {
	info := make(map[string]interface{})

	for _, projects := range release.Projects {
		for _, project := range projects {
			if project.Key == "" || project.Key == p.Key {
				continue
			}

			key := project.Key
			if strings.HasPrefix(key, p.Key+"-") {
				key = strings.TrimPrefix(key, p.Key+"-")
			}

			if project.Tag != nil {
				info[key] = *project.Tag
			} else if len(project.Tags) > 0 {
				info[key] = lib.Keys(project.Tags)
			}
		}
	}

	if len(info) == 0 {
		return nil
	}
	return info
}
