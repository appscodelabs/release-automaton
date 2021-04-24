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
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/appscodelabs/release-automaton/api"
	"github.com/appscodelabs/release-automaton/lib"

	"github.com/Masterminds/semver"
	"github.com/spf13/cobra"
	"gomodules.xyz/semvers"
)

func NewCmdVoyagerRecordLegacyReleases() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "record-legacy-releases",
		Short:             "Writes legacy releases in releases/legacy_releases.json",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			table, err := CreateVoyagerReleaseTable()
			if err != nil {
				panic(err)
			}

			data, err := lib.MarshalJson(table)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(data))

			err = os.MkdirAll(changelogRoot, 0755)
			if err != nil {
				panic(err)
			}

			legacyfilename := filepath.Join(changelogRoot, "legacy_releases.json")
			err = ioutil.WriteFile(legacyfilename, data, 0644)
			if err != nil {
				panic(err)
			}

			filename := filepath.Join(scriptRoot, "README.md")
			lib.WriteChangelogMarkdown(filename, "release-table.tpl", table)
		},
	}
	return cmd
}

func CreateVoyagerReleaseTable() (*api.ReleaseTable, error) {
	owner := "voyagermesh"
	repo := "voyager"

	gh := lib.NewGitHubClient()
	releases, err := lib.ListReleases(context.TODO(), gh, owner, repo)
	if err != nil {
		return nil, err
	}

	var summaries []api.ReleaseSummary
	for _, r := range releases {
		v := semver.MustParse(r.GetTagName())
		if v.Prerelease() == "" ||
			strings.HasPrefix(v.Prerelease(), "v") ||
			strings.HasPrefix(v.Prerelease(), "rc.") {
			summary := api.ReleaseSummary{
				Release:           r.GetTagName(),
				ReleaseDate:       r.GetCreatedAt().UTC(),
				KubernetesVersion: "",
				ReleaseURL:        r.GetHTMLURL(),
				ChangelogURL:      r.GetHTMLURL(),
				DocsURL:           fmt.Sprintf("https://voyagermesh.com/docs/%s", r.GetTagName()),
			}
			if v.LessThan(semver.MustParse("5.0.0")) {
				summary.DocsURL = fmt.Sprintf("https://github.com/%s/docs/tree/%s/docs", owner, r.GetTagName())
			}
			summaries = append(summaries, summary)
		}
	}
	sort.Slice(summaries, func(i, j int) bool {
		return !semvers.CompareVersions(semver.MustParse(summaries[i].Release), semver.MustParse(summaries[j].Release))
	})
	return &api.ReleaseTable{
		ProductLine: "Voyager",
		Releases:    summaries,
	}, nil
}
