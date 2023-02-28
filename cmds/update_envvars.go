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
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var envFile string

/*
	release-automaton update-vars \
	  --env-file=/home/tamal/go/src/kubedb.dev/postgres/Makefile.env \
	  --vars=charts
*/
func NewCmdUpdateEnvVars() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "update-vars",
		Short:             "Update Env variables",
		DisableAutoGenTag: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return updateEnvVars()
		},
	}

	cmd.Flags().StringVar(&envFile, "env-file", envFile, "Path to environment file")
	cmd.Flags().StringToStringVar(&envVars, "vars", envVars, "Key-Value pairs of variables")
	return cmd
}

func updateEnvVars() error {
	data, err := os.ReadFile(envFile)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	for i := range lines {
		for k, v := range envVars {
			if strings.HasPrefix(lines[i], k+"=") {
				lines[i] = fmt.Sprintf("%s=%s", k, v)
			}
		}
	}

	return os.WriteFile(envFile, []byte(strings.Join(lines, "\n")), 0o644)
}
