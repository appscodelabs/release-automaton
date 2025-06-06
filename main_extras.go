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

// nolint
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/appscodelabs/release-automaton/api"
	"github.com/appscodelabs/release-automaton/cmds"
	"github.com/appscodelabs/release-automaton/lib"

	"github.com/alessio/shellescape"
	"github.com/google/go-github/v45/github"
	"github.com/google/uuid"
	"github.com/hashicorp/go-getter"
	"github.com/kballard/go-shellquote"
	"gomodules.xyz/go-sh"
	shell "gomodules.xyz/go-sh"
	"sigs.k8s.io/yaml"
)

func main_Execute_ShellQuote() {
	// find . -type f -exec sed -i 's/from/to/' {} \;
	// sed -i 's|mysql-replication-mode-detector:.*|mysql-replication-mode-detector:v0.1.0-beta.4"|g' charts/kubedb-catalog/templates/mysql/*

	cmd := `find charts/kubedb-catalog/templates/mysql -type f -exec sed -i 's|mysql-replication-mode-detector:.*|mysql-replication-mode-detector:v0.1.0-beta.4"|g' {} \;`
	fields, err := shellquote.Split(cmd)
	if err != nil {
		panic(err)
	}
	fmt.Println(fields)

	args := make([]interface{}, len(fields)-1)
	for i := range fields[1:] {
		args[i] = fields[i+1]
	}

	session := sh.NewSession()
	session.SetDir("/home/tamal/go/src/kubedb.dev/installer")
	session.ShowCMD = true
	err = session.Command(fields[0], args...).Run()
	if err != nil {
		panic(err)
	}
	fmt.Println(shellescape.Quote(""))
}

func main_DetectVCSRoot() {
	vcs, err := lib.DetectVCSRoot("github.com/appscode/stash-enterprise")
	if err != nil {
		panic(err)
	}
	fmt.Println(vcs)
}

func main_RepoURL2EnvKey() {
	fmt.Println(lib.RepoURL2TagEnvKey("http://github.com/kubedb/mysql-replication-mode-detector"))
	fmt.Println(lib.RepoURL2TagEnvKey("github.com/kubedb/operator"))
	fmt.Println(lib.RepoURL2TagEnvKey("github.com/appscode/kubedb-enterprise"))
}

func main_CreateKubeDBReleaseTable() {
	table := cmds.CreateKubeDBReleaseTable()
	data, err := json.MarshalIndent(table, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}

func main_RemoveLabel() {
	gh := lib.NewGitHubClient()
	u := "https://github.com/appscodelabs/release-automaton-demo/pull/1"
	fmt.Println(u)

	var err error

	_, _, err = gh.Issues.AddLabelsToIssue(context.TODO(), "appscodelabs", "release-automaton-demo", 1, []string{
		"xyz",
	})
	if ge, ok := err.(*github.ErrorResponse); ok {
		panic(ge)
	}
	_, _, err = gh.Issues.AddLabelsToIssue(context.TODO(), "appscodelabs", "release-automaton-demo", 1, []string{
		"xyz",
	})
	if ge, ok := err.(*github.ErrorResponse); ok {
		panic(ge)
	}
	err = lib.RemoveLabel(gh, "appscodelabs", "release-automaton-demo", 1, "xyz")
	if err != nil {
		panic(err)
	}
	err = lib.RemoveLabel(gh, "appscodelabs", "release-automaton-demo", 1, "xyz")
	if err != nil {
		panic(err)
	}
	err = lib.RemoveLabel(gh, "appscodelabs", "release-automaton-demo", 1, "xyz")
	if err != nil {
		panic(err)
	}
}

func main_Execute() {
	sh := shell.NewSession()
	sh.ShowCMD = true
	sh.PipeFail = true
	sh.PipeStdErrors = true

	err := lib.Execute(sh, "echo A=B > /tmp/abc.txt", nil)
	if err != nil {
		panic(err)
	}
	err = lib.Execute(sh, "echo C=D >> /tmp/abc.txt", nil)
	if err != nil {
		panic(err)
	}
}

func main_UpdateChangelog() {
	dir := "/home/tamal/go/src/github.com/tamalsaha/release-automaton-demo/CHANGELOG/v2020.6.23"

	releaseFile := filepath.Join(dir, "release.json")
	data, err := os.ReadFile(releaseFile)
	if err != nil {
		panic(err)
	}
	var release api.Release
	err = yaml.Unmarshal(data, &release)
	if err != nil {
		panic(err)
	}

	repoURL := "github.com/appscode-cloud/static-assets"
	tag := "v1.0.0"
	commits := []api.Commit{
		{
			SHA:     "DFGHJK45",
			Subject: "This is a test",
		},
	}
	lib.UpdateChangelog(dir, release, repoURL, tag, commits)
}

func main_ParsePullRequestURL() {
	lib.ParsePullRequestURL("https://github.com/appscodelabs/gh-release-automation-testing/pull/21")
	lib.ParseRepoURL("https://github.com/appscodelabs/gh-release-automation-testing")

	lib.ParsePullRequestURL("github.com/appscodelabs/gh-release-automation-testing/pull/21")
	lib.ParseRepoURL("github.com/appscode-cloud/apimachinery")
}

func mm() {
	sessionID := uuid.New().String()
	localfile := filepath.Join(os.TempDir(), sessionID, "release.txt")
	opts := func(c *getter.Client) error {
		pwd, err := os.Getwd()
		if err != nil {
			return err
		}
		c.Pwd = pwd
		return nil
	}
	releaseFile := ""
	err := getter.GetFile(localfile, releaseFile, opts)
	if err != nil {
		panic(err)
	}

	data, err := os.ReadFile(localfile)
	if err != nil {
		panic(err)
	}

	var r api.Release
	err = yaml.UnmarshalStrict(data, &r)
	if err != nil {
		panic(err)
	}
	fmt.Println(r)

	sh := shell.NewSession()
	sh.ShowCMD = true
	sh.PipeFail = true
	sh.PipeStdErrors = true

	err = sh.Command("env").Run()
	if err != nil {
		panic(err)
	}
}

func main_CreatePR() {
	// https://github.com/tamalsaha/gh-release-automation-testing/pull/new/pr-4

	owner := "tamalsaha"
	repo := "gh-release-automation-testing"

	gh := lib.NewGitHubClient()

	pr, err := lib.CreatePR(gh, owner, repo, &github.NewPullRequest{
		Title:               github.String("Test pr api"),
		Head:                github.String("pr-4"),
		Base:                github.String("master"),
		Body:                github.String("XYZ"),
		Issue:               nil,
		MaintainerCanModify: github.Bool(true),
		Draft:               github.Bool(false),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(pr)
}

func main_ListCommits() {
	sh := shell.NewSession()
	sh.ShowCMD = true
	sh.PipeFail = true
	sh.PipeStdErrors = true

	sh.SetDir("/home/tamal/go/src/stash.appscode.dev/stash")

	commits := lib.ListCommits(sh, "0ab9faa68308cd646e1e63271950cf75e3cf62c0", "v0.9.0-rc.6")
	for _, commit := range commits {
		fmt.Println(commit.SHA, commit.Subject)
	}
}

func main_CherryPick() {
	sh := shell.NewSession()
	sh.ShowCMD = true
	sh.PipeFail = true
	sh.PipeStdErrors = true

	sh.SetDir("/home/tamal/go/src/github.com/appscode-cloud/pg")

	data, err := sh.Command("git", "show", "-s", "--format=%b").Output()
	if err != nil {
		panic(err)
	}
	a := []byte("ProductLine: Stash")
	fmt.Println(string(data), a)
}

func main_ShellGetwd() {
	sh := shell.NewSession()
	sh.ShowCMD = true
	sh.PipeFail = true
	sh.PipeStdErrors = true

	fmt.Println(sh.Getwd())

	sh.SetDir("/home/tamal/go/src/stash.appscode.dev/stash")

	fmt.Println(sh.Getwd())
}

func main_ListTags() {
	sh := shell.NewSession()
	sh.ShowCMD = true
	sh.PipeFail = true
	sh.PipeStdErrors = true

	sh.SetDir("/home/tamal/go/src/github.com/appscodelabs/release-automaton")

	tags, err := lib.ListTags(sh)

	fmt.Println(tags, err)
}

func main_RepoModified() {
	sh := shell.NewSession()
	sh.ShowCMD = true
	sh.PipeFail = true
	sh.PipeStdErrors = true

	if lib.RepoModified(sh) {
		fmt.Println("Something to commit")
	}
}

func main_WriteChangelogMarkdown() {
	chlog := api.Changelog{
		Release: "v2020.6.12",
		Projects: []api.ProjectChangelog{
			{
				URL: "github.com/stashed/apimachinery",
				Releases: []api.ReleaseChangelog{
					{
						Tag: "v0.1.0",
						Commits: []api.Commit{
							{SHA: "d12b3d4b42a91166081d514a4a03e226b60a1b1f", Subject: "Update to Kubernetes v1.18.3 (#21)"},
							{SHA: "1956a31259db988fdd047aa9c45f08cc23d866d8", Subject: "Update to Kubernetes v1.18.3"},
							{SHA: "c39660025dd0610992691f468174d8edc089f678", Subject: "Unwrap top level api folder (#20)"},
							{SHA: "5ba03fb5ea9064e6de7a172bd8d0a0d76df5f0d5", Subject: "Update to Kubernetes v1.18.3 (#19)"},
							{SHA: "abeb620e309283ab8c5a1eace065912242022aef", Subject: "Update to Kubernetes v1.18.3"},
							{SHA: "6fdf8a609b831d361e768ab08cdac1949947f3d9", Subject: "Enable https://kodiakhq.com (#13)"},
							{SHA: "479258eda8cd2c0fe2c5c024532d6f64dc45c092", Subject: "Update dev scripts (#12)"},
						},
					},
				},
			},
		},
	}
	lib.WriteChangelogMarkdown(filepath.Join("/tmp", "CHANGELOG.md"), "changelog.tpl", chlog)
}

func mainMergeReply() {
	var replies api.Replies

	r := api.Reply{
		Type: api.Tagged,
		Tagged: &api.TaggedReplyData{
			Repo: "a/b",
		},
	}

	replies = api.MergeReply(replies, r)
	fmt.Println(replies)
}
