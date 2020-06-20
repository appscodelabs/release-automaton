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

	"github.com/appscodelabs/release-automaton/api"

	"github.com/spf13/cobra"
)

func NewCmdStashCreateCatalog() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "create-catalog",
		Short:             "Create Stash catalog",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			catalog := CreateCatalogData()
			catalog.Sort()

			data, err := json.MarshalIndent(catalog, "", "  ")
			if err != nil {
				panic(err)
			}
			fmt.Println(string(data))
		},
	}
	return cmd
}

func CreateCatalogData() api.StashCatalog {
	return api.StashCatalog{
		Addons: []api.Addon{
			{
				Name: "postgres",
				Versions: []string{
					"9.6",
					"10.2",
					"10.6",
					"11.1",
					"11.2",
				},
			},
			{
				Name: "mongodb",
				Versions: []string{
					"3.4.17",
					"3.4.22",
					"3.6.8",
					"3.6.13",
					"4.0.3",
					"4.0.5",
					"4.0.11",
					"4.1.4",
					"4.1.7",
					"4.1.13",
					"4.2.3",
				},
			},
			{
				Name: "elasticsearch",
				Versions: []string{
					"5.6.4",
					"6.2.4",
					"6.3.0",
					"6.4.0",
					"6.5.3",
					"6.8.0",
					"7.2.0",
					"7.3.2",
				},
			},
			{
				Name: "mysql",
				Versions: []string{
					"5.7.25",
					"8.0.3",
					"8.0.14",
				},
			},
			{
				Name: "percona-xtradb",
				Versions: []string{
					"5.7",
				},
			},
		},
	}
}
