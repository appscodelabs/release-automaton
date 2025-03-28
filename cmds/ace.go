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
)

func NewCmdAce() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "ace",
		Short:             "ACE catalog commands",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	cmd.AddCommand(NewCmdAceCreateRelease())
	return cmd
}
