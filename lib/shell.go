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

package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kballard/go-shellquote"
	shell "gomodules.xyz/go-sh"
)

func Execute(sh *shell.Session, cmd string, env map[string]string) error {
	var cmdlets []string
	var appendOut bool
	var createOut bool
	var filename string
	if strings.Contains(cmd, ">>") {
		appendOut = true
		cmdlets = strings.SplitN(cmd, ">>", 2)
		filename = strings.TrimSpace(cmdlets[1])
		if !filepath.IsAbs(filename) {
			filename = filepath.Join(sh.Getwd(), filename)
		}
	} else if strings.Contains(cmd, ">") {
		createOut = true
		cmdlets = strings.SplitN(cmd, ">", 2)
		filename = strings.TrimSpace(cmdlets[1])
		if !filepath.IsAbs(filename) {
			filename = filepath.Join(sh.Getwd(), filename)
		}
	} else {
		cmdlets = []string{cmd}
	}

	fields, err := shellquote.Split(cmdlets[0])
	if err != nil {
		return err
	}
	if len(fields) == 0 {
		return fmt.Errorf("missing command: %s", cmd)
	}

	args := make([]interface{}, len(fields))
	for i := range fields[1:] {
		args[i] = fields[i+1]
	}
	if env == nil {
		args[len(fields)-1] = map[string]string{}
	} else {
		args[len(fields)-1] = env // pass env vars for this command only
	}

	s := sh.Command(fields[0], args...)
	if createOut {
		if !Exists(filename) {
			err := os.WriteFile(filename, []byte(""), 0o644)
			if err != nil {
				return err
			}
		} else {
			err := os.Truncate(filename, 0)
			if err != nil {
				return err
			}
		}
		return s.WriteStdout(filename)
	} else if appendOut {
		if !Exists(filename) {
			err := os.WriteFile(filename, []byte{}, 0o644)
			if err != nil {
				return err
			}
		}
		return s.AppendStdout(filename)
	}
	return s.Run()
}
