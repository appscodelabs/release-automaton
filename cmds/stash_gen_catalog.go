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
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/appscodelabs/release-automaton/api"
	"github.com/appscodelabs/release-automaton/lib"

	"github.com/Masterminds/semver"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
	"stash.appscode.dev/catalog"
)

var (
	catalogFile string
)

/*
release-automaton stash gen-catalog \
  --release-file=/home/tamal/go/src/github.com/tamalsaha/gh-release-automation-testing/v2020.6.16/release.json \
  --catalog-file=/home/tamal/go/src/stash.appscode.dev/catalog/catalog.json
*/
func NewCmdStashGenCatalog() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "gen-catalog",
		Short:             "Generate Stash catalog",
		DisableAutoGenTag: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateCatalog()
		},
	}

	cmd.Flags().StringVar(&releaseFile, "release-file", "", "Path of release file (local file or url is accepted)")
	cmd.Flags().StringVar(&catalogFile, "catalog-file", "", "Path to Stash catalog file")
	return cmd
}

func generateCatalog() error {
	data, err := ioutil.ReadFile(releaseFile)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, &release)
	if err != nil {
		return err
	}

	var catalog catalog.StashCatalog
	data, err = ioutil.ReadFile(catalogFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &catalog)
	if err != nil {
		return err
	}

	v := semver.MustParse(release.Release)
	if strings.HasPrefix(v.Prerelease(), "alpha.") || strings.HasPrefix(v.Prerelease(), "beta.") {
		catalog.ChartRegistry = api.TestChartRegistry
		catalog.ChartRegistryURL = api.TestChartRegistryURL
	} else {
		catalog.ChartRegistry = api.StableChartRegistry
		catalog.ChartRegistryURL = api.StableChartRegistryURL
	}
	for i, addon := range catalog.Addons {
		if versions, ok := findAddonVersions(addon.Name); ok {
			catalog.Addons[i].Versions = versions
		}
	}
	catalog.Sort()

	b, err := json.MarshalIndent(catalog, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(catalogFile, b, 0644)
	return err
}

func findAddonVersions(addon string) ([]string, bool) {
	for _, projects := range release.Projects {
		for repoURL, project := range projects {
			if !strings.HasSuffix(repoURL, "/"+addon) {
				continue
			}
			if project.Tags != nil {
				return lib.Keys(project.Tags), true
			}
		}
	}
	return nil, false

}
