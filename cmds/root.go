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
	"github.com/spf13/cobra"
	v "gomodules.xyz/x/version"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:               "release-automaton [command]",
		Short:             `release-automaton by AppsCode - Release often`,
		DisableAutoGenTag: true,
	}

	rootCmd.AddCommand(NewCmdRelease())
	rootCmd.AddCommand(NewCmdStash())
	rootCmd.AddCommand(NewCmdKubeDB())
	rootCmd.AddCommand(NewCmdKubeVault())
	rootCmd.AddCommand(NewCmdKubeform())
	rootCmd.AddCommand(NewCmdVoyager())
	rootCmd.AddCommand(NewCmdUpdateBundles())
	rootCmd.AddCommand(NewCmdUpdateAssets())
	rootCmd.AddCommand(NewCmdUpdateEnvVars())
	rootCmd.AddCommand(v.NewCmdVersion())
	return rootCmd
}
