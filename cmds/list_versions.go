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
	"fmt"
	"os"

	"github.com/appscodelabs/release-automaton/api"

	"github.com/spf13/cobra"
	"gomodules.xyz/semvers"
	"gomodules.xyz/sets"
)

/*
	release-automaton list-versions \
	  --release-file=/Users/tamal/go/src/kubedb.dev/CHANGELOG/releases/v2023.12.28/release.json
*/
func NewCmdListVersions() *cobra.Command {
	var relFile string
	cmd := &cobra.Command{
		Use:               "list-versions",
		Short:             "List versions",
		DisableAutoGenTag: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listReleaseVersions(relFile)
		},
	}

	cmd.Flags().StringVar(&relFile, "release-file", relFile, "Path to environment file")
	return cmd
}

func listReleaseVersions(relFile string) error {
	data, err := os.ReadFile(relFile)
	if err != nil {
		return err
	}

	var release api.Release
	err = json.Unmarshal(data, &release)
	if err != nil {
		return err
	}

	tags := sets.NewString()
	for _, projects := range release.Projects {
		for _, project := range projects {
			if project.Tag != nil {
				tags = tags.Insert(*project.Tag)
			}
		}
	}
	tagList := tags.List()
	semvers.SortVersions(tagList, semvers.CompareDesc)
	for _, tag := range tagList {
		fmt.Println(tag)
	}
	return nil
}
