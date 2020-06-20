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
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/appscodelabs/release-automaton/api"
	"github.com/appscodelabs/release-automaton/templates"

	"github.com/Masterminds/sprig"
)

func UpdateChangelog(dir, url, tag string, commits []api.Commit) {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		panic(err)
	}

	filenameChlog := filepath.Join(dir, "CHANGELOG.json")

	var chlog api.Changelog
	data, err := ioutil.ReadFile(filenameChlog)
	if err == nil {
		err = json.Unmarshal(data, &chlog)
		if err != nil {
			panic(err)
		}
	}

	var repoFound bool
	for repoIdx := range chlog.Projects {
		if chlog.Projects[repoIdx].URL == url {
			repoFound = true

			var tagFound bool
			for tagIdx := range chlog.Projects[repoIdx].Releases {
				if chlog.Projects[repoIdx].Releases[tagIdx].Tag == tag {
					chlog.Projects[repoIdx].Releases[tagIdx].Commits = commits
					tagFound = true
					break
				}
			}
			if !tagFound {
				chlog.Projects[repoIdx].Releases = append(chlog.Projects[repoIdx].Releases, api.ReleaseChangelog{
					Tag:     tag,
					Commits: commits,
				})
			}
			repoFound = true
			break
		}
	}
	if !repoFound {
		chlog.Projects = append(chlog.Projects, api.ProjectChangelog{
			URL: url,
			Releases: []api.ReleaseChangelog{
				{
					Tag:     tag,
					Commits: commits,
				},
			},
		})
	}
	chlog.Sort()

	b, err := json.MarshalIndent(chlog, "", "  ")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(filenameChlog, b, 0644)
	if err != nil {
		panic(err)
	}

	WriteChangelogMarkdown(dir, chlog)
}

func WriteChangelogMarkdown(dir string, chlog api.Changelog) {
	tpl := template.Must(template.New("").Funcs(sprig.FuncMap()).Parse(string(templates.MustAsset("changelog.tpl"))))
	var buf bytes.Buffer
	err := tpl.Execute(&buf, chlog)
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.String())
	err = ioutil.WriteFile(filepath.Join(dir, "CHANGELOG.md"), buf.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
}
