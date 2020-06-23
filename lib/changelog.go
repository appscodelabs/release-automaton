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
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/appscodelabs/release-automaton/api"
	"github.com/appscodelabs/release-automaton/templates"

	"github.com/Masterminds/sprig"
)

func UpdateChangelog(dir string, release api.Release, repoURL, tag string, commits []api.Commit) {
	var status api.ChangelogStatus
	for _, projects := range release.Projects {
		for u, project := range projects {
			if u == repoURL {
				status = project.Changelog
				if status == api.SkipChangelog {
					return
				}
				break
			}
		}
	}

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
	chlog.ProductLine = release.ProductLine
	chlog.Release = release.Release
	chlog.ReleaseDate = time.Now().UTC()
	chlog.KubernetesVersion = release.KubernetesVersion

	var repoFound bool
	for repoIdx := range chlog.Projects {
		if chlog.Projects[repoIdx].URL == repoURL {
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
			URL: repoURL,
			Releases: []api.ReleaseChangelog{
				{
					Tag:     tag,
					Commits: commits,
				},
			},
		})
	}
	chlog.Sort()

	data, err = MarshalJson(chlog)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(filenameChlog, data, 0644)
	if err != nil {
		panic(err)
	}

	WriteChangelogMarkdown(filepath.Join(dir, "CHANGELOG.md"), "changelog.tpl", chlog)
	if status == api.StandaloneWebsiteChangelog {
		WriteChangelogMarkdown(filepath.Join(dir, "docs_changelog.md"), "standalone-changelog.tpl", chlog)
	} else if status == api.SharedWebsiteChangelog {
		WriteChangelogMarkdown(filepath.Join(dir, "docs_changelog.md"), "shared-changelog.tpl", chlog)
	}
}

func WriteChangelogMarkdown(filename string, tplname string, chlog api.Changelog) {
	tpl := template.Must(template.New("").Funcs(sprig.FuncMap()).Parse(string(templates.MustAsset(tplname))))
	var buf bytes.Buffer
	err := tpl.Execute(&buf, chlog)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(filename, buf.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
}
