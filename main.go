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

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/appscodelabs/release-automaton/lib"
	"github.com/appscodelabs/release-automaton/templates"

	"github.com/Masterminds/semver"
	"github.com/Masterminds/sprig"
	shell "github.com/codeskyblue/go-sh"
	"github.com/google/go-github/v32/github"
	"github.com/google/uuid"
	"github.com/hashicorp/go-getter"
	"github.com/keighl/metabolize"
	flag "github.com/spf13/pflag"
	"github.com/tamalsaha/go-oneliners"
	"golang.org/x/mod/modfile"
	"golang.org/x/oauth2"
	"gomodules.xyz/envsubst"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/yaml"
)

func NewGitHubClient() *github.Client {
	token, found := os.LookupEnv(GitHubTokenKey)
	if !found {
		log.Fatalln(GitHubTokenKey + " env var is not set")
	}

	// Create the http client.
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.TODO(), ts)

	return github.NewClient(tc)
}

func ListTags2(ctx context.Context, client *github.Client, owner, repo string) ([]*github.RepositoryTag, error) {
	opt := &github.ListOptions{
		PerPage: 100,
	}

	var result []*github.RepositoryTag
	for {
		reviews, resp, err := client.Repositories.ListTags(ctx, owner, repo, opt)
		if err != nil {
			break
		}
		result = append(result, reviews...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return result, nil
}

func ListTags(sh *shell.Session) ([]string, error) {
	data, err := sh.Command("git", "tag").Output()
	if err != nil {
		return nil, err
	}
	return strings.Fields(string(data)), nil
}

func ListReviews(ctx context.Context, client *github.Client, owner, repo string, number int) ([]*github.PullRequestReview, error) {
	opt := &github.ListOptions{
		PerPage: 100,
	}

	var result []*github.PullRequestReview
	for {
		reviews, resp, err := client.PullRequests.ListReviews(ctx, owner, repo, number, opt)
		if err != nil {
			break
		}
		result = append(result, reviews...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return result, nil
}

func ListPullRequestComment(ctx context.Context, client *github.Client, owner, repo string, number int) ([]*github.PullRequestComment, error) {
	opt := &github.PullRequestListCommentsOptions{
		Sort:      "created",
		Direction: "asc",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	var result []*github.PullRequestComment
	for {
		comments, resp, err := client.PullRequests.ListComments(ctx, owner, repo, number, opt)
		if err != nil {
			break
		}
		result = append(result, comments...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return result, nil
}

func ListComments(ctx context.Context, client *github.Client, owner, repo string, number int) ([]*github.IssueComment, error) {
	opt := &github.IssueListCommentsOptions{
		Sort:      github.String("created"),
		Direction: github.String("asc"),
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	var result []*github.IssueComment
	for {
		comments, resp, err := client.Issues.ListComments(ctx, owner, repo, number, opt)
		if err != nil {
			break
		}
		result = append(result, comments...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return result, nil
}

func RemoteBranchExists(sh *shell.Session, branch string) bool {
	data, err := sh.Command("git", "ls-remote", "--heads", "origin", branch).Output()
	if err != nil {
		panic(err)
	}
	return len(bytes.TrimSpace(data)) > 0
}

func RemoteTagExists(sh *shell.Session, tag string) bool {
	// git ls-remote --exit-code --tags origin <tag>
	err := sh.Command("git", "ls-remote", "--exit-code", "--tags", "origin", tag).Run()
	return err == nil
}

func GetRemoteTag(sh *shell.Session, tag string) string {
	// git ls-remote --exit-code --tags origin <tag>
	data, err := sh.Command("git", "ls-remote", "--exit-code", "--tags", "origin", tag).Output()
	if err != nil {
		return ""
	}
	return strings.Fields(string(data))[0]
}

type ConditionFunc func(*shell.Session, string) bool

func MeetsCondition(fn ConditionFunc, sh *shell.Session, items ...string) bool {
	for _, item := range items {
		if !fn(sh, item) {
			return false
		}
	}
	return true
}

func FirstCommit(sh *shell.Session) string {
	// git rev-list --max-parents=0 HEAD
	// ref: https://stackoverflow.com/a/5189296
	data, err := sh.Command("git", "rev-list", "--max-parents=0", "HEAD").Output()
	if err != nil {
		panic(err)
	}
	commits := strings.Fields(string(data))
	return commits[len(commits)-1]
}

func LatestCommit(sh *shell.Session) string {
	// // git show -s --format=%H
	data, err := sh.Command("git", "show", "-s", "--format=%H").Output()
	if err != nil {
		panic(err)
	}
	commits := strings.Fields(string(data))
	return commits[0]
}

func PrepareProject(gh *github.Client, sh *shell.Session, releaseTracker, repoURL string, project Project) error {
	if project.Tags != nil && project.Tag != nil {
		return fmt.Errorf("repo %s is provided an invalid project configuration which uses both tag and tags", repoURL)
	}

	// pushd, popd
	wdOrig := sh.Getwd()
	defer sh.SetDir(wdOrig)

	owner, repo := ParseRepoURL(repoURL)

	// TODO: cache git repo
	wdCur := filepath.Join(gitRoot, owner)
	err := os.MkdirAll(wdCur, 0755)
	if err != nil {
		return err
	}

	if !exists(filepath.Join(wdCur, repo)) {
		sh.SetDir(wdCur)

		err = sh.Command("git",
			"clone",
			// "--no-tags", //TODO: ok?
			"--no-recurse-submodules",
			//"--depth=1",
			//"--no-single-branch",
			fmt.Sprintf("https://%s:%s@%s.git", os.Getenv(GitHubUserKey), os.Getenv(GitHubTokenKey), repoURL),
		).Run()
		if err != nil {
			return err
		}
	}
	wdCur = filepath.Join(wdCur, repo)
	sh.SetDir(wdCur)

	modPath := DetectGoMod(wdCur)
	if modPath != "" {
		gm := GoImport{
			RepoRoot: repoURL,
		}
		vcs, err := DetectVCSRoot(modPath)
		if err != nil {
			panic(err)
		}
		if vcs != repoURL {
			gm.VCSRoot = vcs
		}
		modCache[modPath] = gm
	}

	tags := project.Tags
	if project.Tag != nil {
		tags = map[string]string{
			*project.Tag: "master", // pr always opened against master branch
		}
	}

	// All remote tags exist, so only add Go module path if needed.
	if MeetsCondition(RemoteTagExists, sh, lib.Keys(tags)...) {
		var ok bool
		// Make sure /tagged, /cherry-picked comments exist
		if project.Tag != nil {
			sha := GetRemoteTag(sh, *project.Tag)
			if replies, ok = AppendReplyIfMissing(replies, Reply{
				Type: ReadyToTag,
				ReadyToTag: &ReadyToTagReplyData{
					Repo:           repoURL,
					MergeCommitSHA: sha,
				},
			}); ok {
				comments = append(comments, fmt.Sprintf("%s %s %s", ReadyToTag, repoURL, sha))
			}
		}
		if project.Tags != nil {
			for tag, branch := range project.Tags {
				sha := GetRemoteTag(sh, tag)
				if replies, ok = AppendReplyIfMissing(replies, Reply{
					Type: CherryPicked,
					CherryPicked: &CherryPickedReplyData{
						Repo:           repoURL,
						Branch:         branch,
						MergeCommitSHA: sha,
					},
				}); ok {
					comments = append(comments, fmt.Sprintf("%s %s %s %s", CherryPicked, repoURL, branch, sha))
				}
			}
		}

		if replies, ok = AppendReplyIfMissing(replies, Reply{
			Type: Tagged,
			Tagged: &TaggedReplyData{
				Repo: repoURL,
			},
		}); ok {
			comments = append(comments, fmt.Sprintf("%s %s", Tagged, repoURL))
		}

		if modPath != "" {
			AppendGo(modPath)
			modPath = ""
		}
		return nil
	}

	usesCherryPick := project.Tags != nil && project.Tag == nil
	if usesCherryPick {
		tags[releaseNumber] = "master" // if cherry pick is used, there must be an extra pr against the master branch
	}

	for tag, branch := range tags {
		if usesCherryPick {
			// remote branch must already exist
			if !RemoteBranchExists(sh, branch) {
				return fmt.Errorf("repo %s is missing branch for tag %s", repoURL, tag)
			}
		}

		// -----------------------

		vars := lib.MergeMaps(map[string]string{
			repoURL2EnvKey(repoURL): tag,
			"TAG":                   tag,
			"RELEASE":               releaseNumber,
			"RELEASE_TRACKER":       releaseTracker,
		}, envVars)

		headBranch := fmt.Sprintf("%s-%s", releaseNumber, branch)

		err = sh.Command("git", "checkout", branch).Run()
		if err != nil {
			return err
		}

		err = sh.Command("git", "checkout", "-b", headBranch).Run()
		if err != nil {
			return err
		}

		if exists(filepath.Join(wdCur, "go.mod")) {
			// Update Go mod
			UpdateGoMod(wdCur)
			if RepoModified(sh) {
				err = sh.Command("go", "mod", "tidy").Run()
				if err != nil {
					return err
				}
				err = sh.Command("go", "mod", "vendor").Run()
				if err != nil {
					return err
				}
			}
		}

		for _, cmd := range project.Commands {
			cmd, err = envsubst.EvalMap(cmd, vars)
			if err != nil {
				return err
			}
			fields := strings.Fields(cmd)
			if len(fields) > 0 {
				args := make([]interface{}, len(fields)-1)
				for i := range fields[1:] {
					args[i] = fields[i+1]
				}

				err = sh.Command(fields[0], args...).Run()
				if err != nil {
					return err
				}
			}
		}

		if RepoModified(sh) {
			messages := []string{
				"ProductLine: " + release.ProductLine,
				"Release: " + releaseNumber,
			}
			if !usesCherryPick || branch != "master" {
				// repos that use cherry pick, a pr is opened against the master branch
				// That pr MUST NOT report back to release tracker.
				messages = append(messages, "Release-tracker: "+releaseTracker)
			}
			err = CommitRepo(sh, tag, messages...)
			if err != nil {
				return err
			}
			err = PushRepo(sh, true)
			if err != nil {
				return err
			}

			// open pr against project repo
			prBody := fmt.Sprintf(`ProductLine: %s
Release: %s
Release-tracker: %s`, release.ProductLine, releaseNumber, releaseTracker)
			pr, err := CreatePR(gh, owner, repo, &github.NewPullRequest{
				Title:               github.String(fmt.Sprintf("Prepare for release %s", tag)),
				Head:                github.String(headBranch),
				Base:                github.String(branch),
				Body:                github.String(prBody),
				MaintainerCanModify: github.Bool(true),
				Draft:               github.Bool(false),
			}, "automerge")
			if err != nil {
				panic(err)
			}

			// add comments to release repo
			comments = append(comments, fmt.Sprintf("%s %s", PR, pr.GetHTMLURL()))
		} else {
			comments = append(comments, fmt.Sprintf("%s %s %s", ReadyToTag, repoURL, LatestCommit(sh)))
			// TODO: add to replies map?
		}

		if modPath != "" {
			AppendGo(modPath)
			modPath = ""
		}
	}

	return nil
}

func AppendGo(modPath string) {
	gm := modCache[modPath]
	comments = append(comments, fmt.Sprintf(`%s %s %s %s`, Go, gm.RepoRoot, modPath, gm.VCSRoot))
}

func ReleaseProject(gh *github.Client, sh *shell.Session, releaseTracker, repoURL string, project Project) error {
	if project.Tags != nil && project.Tag != nil {
		return fmt.Errorf("repo %s is provided an invalid project configuration which uses both tag and tags", repoURL)
	}

	// pushd, popd
	wdOrig := sh.Getwd()
	defer sh.SetDir(wdOrig)

	owner, repo := ParseRepoURL(repoURL)

	// TODO: cache git repo
	wdCur := filepath.Join(gitRoot, owner)
	err := os.MkdirAll(wdCur, 0755)
	if err != nil {
		return err
	}

	if !exists(filepath.Join(wdCur, repo)) {
		sh.SetDir(wdCur)

		err = sh.Command("git",
			"clone",
			// "--no-tags", //TODO: ok?
			"--no-recurse-submodules",
			//"--depth=1",
			//"--no-single-branch",
			fmt.Sprintf("https://%s:%s@%s.git", os.Getenv(GitHubUserKey), os.Getenv(GitHubTokenKey), repoURL),
		).Run()
		if err != nil {
			return err
		}
	}
	wdCur = filepath.Join(wdCur, repo)
	sh.SetDir(wdCur)

	modPath := DetectGoMod(wdCur)
	if modPath != "" {
		gm := GoImport{
			RepoRoot: repoURL,
		}
		vcs, err := DetectVCSRoot(modPath)
		if err != nil {
			panic(err)
		}
		if vcs != repoURL {
			gm.VCSRoot = vcs
		}
		modCache[modPath] = gm
	}

	tags := project.Tags
	if project.Tag != nil {
		tags = map[string]string{
			*project.Tag: "", // branch unknown
		}
	}

	// All remote tags exist, so only add Go module path if needed.
	if MeetsCondition(RemoteTagExists, sh, lib.Keys(tags)...) {
		// make sure /tagged is appended for next group of projects in this run
		// and added to comments for next run
		var ok bool
		if replies, ok = AppendReplyIfMissing(replies, Reply{
			Type: Tagged,
			Tagged: &TaggedReplyData{
				Repo: repoURL,
			},
		}); ok {
			comments = append(comments, fmt.Sprintf("%s %s", Tagged, repoURL))
		}

		if modPath != "" {
			AppendGo(modPath)
		}
		return nil
	}

	usesCherryPick := project.Tags != nil && project.Tag == nil

	for tag, branch := range tags {
		vTag, err := semver.NewVersion(tag)
		if err != nil {
			return err
		}

		// detect branch
		if usesCherryPick {
			// remote branch must already exist
			if !RemoteBranchExists(sh, branch) {
				return fmt.Errorf("repo %s is missing branch for tag %s", repoURL, tag)
			}
		} else {
			if vTag.Patch() > 0 { // PATCH release
				if vTag.Prerelease() != "" {
					panic(fmt.Errorf("version %s is invalid because it is a patch release but includes a pre-release component", tag))
				}

				patchBranch := fmt.Sprintf("release-%d.%d.%d", vTag.Major(), vTag.Minor(), vTag.Patch())
				if RemoteBranchExists(sh, patchBranch) {
					branch = patchBranch
				} else {
					minorBranch := fmt.Sprintf("release-%d.%d", vTag.Major(), vTag.Minor())
					if RemoteBranchExists(sh, minorBranch) {
						branch = minorBranch
					}
				}
				if branch == "" {
					return fmt.Errorf("repo %s is missing branch for tag %s", repoURL, tag)
				}
				tags[tag] = branch
			} else {
				branch = fmt.Sprintf("release-%d.%d", vTag.Major(), vTag.Minor())
				tags[tag] = branch
			}
		}

		// -----------------------

		if usesCherryPick || vTag.Patch() > 0 {
			err = sh.Command("git", "checkout", branch).Run()
			if err != nil {
				return err
			}

			if sha, found := MergedCommitSHA(repoURL, branch, usesCherryPick); found {
				// git reset --hard cedc856
				err = sh.Command("git", "reset", "--hard", sha).Run()
				if err != nil {
					return err
				}
			}

			err = TagRepo(sh, tag, "ProductLine: "+release.ProductLine, "Release: "+releaseNumber, "Release-tracker: "+releaseTracker)
			if err != nil {
				return err
			}
			err = PushRepo(sh, true)
			if err != nil {
				return err
			}
		} else if vTag.Patch() == 0 {
			if RemoteBranchExists(sh, branch) {
				err = sh.Command("git", "checkout", branch).Run()
				if err != nil {
					return err
				}
				ref := "master"
				if sha, found := MergedCommitSHA(repoURL, branch, usesCherryPick); found {
					ref = sha
				}
				err = sh.Command("git", "merge", ref).Run()
				if err != nil {
					return err
				}
				err = TagRepo(sh, tag, "ProductLine: "+release.ProductLine, "Release: "+releaseNumber, "Release-tracker: "+releaseTracker)
				if err != nil {
					return err
				}
				err = PushRepo(sh, true)
				if err != nil {
					return err
				}
			} else {
				err = sh.Command("git", "checkout", "master").Run()
				if err != nil {
					return err
				}
				if sha, found := MergedCommitSHA(repoURL, branch, usesCherryPick); found {
					// git reset --hard $sha
					err = sh.Command("git", "reset", "--hard", sha).Run()
					if err != nil {
						return err
					}
				}
				err = sh.Command("git", "checkout", "-b", branch).Run()
				if err != nil {
					return err
				}
				err = TagRepo(sh, tag, "ProductLine: "+release.ProductLine, "Release: "+releaseNumber, "Release-tracker: "+releaseTracker)
				if err != nil {
					return err
				}
				err = PushRepo(sh, true)
				if err != nil {
					return err
				}
			}
		}

		// add comments to release repo
		{
			comments = append(comments, fmt.Sprintf("%s %s", Tagged, repoURL))
			if modPath != "" {
				AppendGo(modPath)
			}
		}

		tags, err := ListTags(sh)
		if err != nil {
			return err
		}
		tagSet := sets.NewString(tags...)
		tagSet.Insert(tag)

		vs := make([]*semver.Version, tagSet.Len())
		for i, r := range tags {
			v, err := semver.NewVersion(r)
			if err != nil {
				return fmt.Errorf("error parsing version: %s", err)
			}
			vs[i] = v
		}
		sort.Sort(lib.SemverCollection(vs))

		var tagIdx = -1
		for idx, vs := range vs {
			if vs.Equal(vTag) {
				tagIdx = idx
				break
			}
		}

		var commits []Commit
		if tagIdx == 0 {
			commits = ListCommits(sh, FirstCommit(sh), vs[tagIdx].Original())
		} else {
			commits = ListCommits(sh, vs[tagIdx-1].Original(), vs[tagIdx].Original())
		}
		UpdateChangelog(filepath.Join(scriptRoot, releaseNumber), repoURL, tag, commits)
		if AnyRepoModified(scriptRoot, sh) {
			err = CommitAnyRepo(scriptRoot, sh, "", "Update changelog")
			if err != nil {
				return err
			}
			err = PushAnyRepo(scriptRoot, sh, false)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func MergedCommitSHA(repoURL, branch string, useCherryPick bool) (string, bool) {
	key := MergeData{
		Repo:   repoURL,
		Branch: branch,
	}
	if !useCherryPick {
		key.Branch = "master"
	}
	sha, ok := merged[key]
	return sha, ok
}

func ResetRepo(sh *shell.Session) error {
	// git add --all; git stash; git stash drop
	err := sh.Command("git", "add", "--all").Run()
	if err != nil {
		return err
	}
	return sh.Command("git", "stash").Run()
}

func AnyRepoModified(wd string, sh *shell.Session) bool {
	wdorig := sh.Getwd()
	defer sh.SetDir(wdorig)

	sh.SetDir(wd)
	return RepoModified(sh)
}

func RepoModified(sh *shell.Session) bool {
	err := sh.Command("git", "add", "--all").Run()
	if err != nil {
		panic(err)
	}
	// https://stackoverflow.com/questions/10385551/get-exit-code-go
	err = sh.Command("git", "diff", "--exit-code", "-s", "HEAD").Run()
	return err != nil
}

func CommitAnyRepo(wd string, sh *shell.Session, tag string, messages ...string) error {
	wdorig := sh.Getwd()
	defer sh.SetDir(wdorig)

	sh.SetDir(wd)
	return CommitRepo(sh, tag, messages...)
}

func CommitRepo(sh *shell.Session, tag string, messages ...string) error {
	err := sh.Command("git", "add", "--all").Run()
	if err != nil {
		return err
	}
	//  git commit -a -s -m "Prepare for release %tag"
	args := []interface{}{
		"commit", "-a", "-s",
	}
	if tag != "" {
		args = append(args, "-m", "Prepare for release "+tag)
	}
	for _, msg := range messages {
		args = append(args, "-m", msg)
	}
	return sh.Command("git", args...).Run()
}

func PushAnyRepo(wd string, sh *shell.Session, pushTag bool) error {
	wdorig := sh.Getwd()
	defer sh.SetDir(wdorig)

	sh.SetDir(wd)
	return PushRepo(sh, pushTag)
}

func PushRepo(sh *shell.Session, pushTag bool) error {
	args := []interface{}{"push", "-u", "origin", "HEAD"}
	if pushTag {
		args = append(args, "--tags")
	}
	return sh.Command("git", args...).Run()
}

func TagRepo(sh *shell.Session, tag string, messages ...string) error {
	args := []interface{}{
		"tag", "-a", tag, "-m", tag,
	}
	for _, msg := range messages {
		args = append(args, "-m", msg)
	}
	return sh.Command("git", args...).Run()
}

func ListCommits(sh *shell.Session, start, end string) []Commit {
	// git log --oneline --ancestry-path start..end | cat
	// ref: https://stackoverflow.com/a/44344164/244009
	data, err := sh.Command("git", "log", "--oneline", "--ancestry-path", fmt.Sprintf("%s..%s", start, end)).Output()
	if err != nil {
		panic(err)
	}
	var commits []Commit
	for _, line := range strings.Split(string(data), "\n") {
		if line == "" {
			continue
		}
		idx := strings.IndexRune(line, ' ')
		if idx != -1 {
			commits = append(commits, Commit{
				SHA:     line[:idx],
				Subject: line[idx+1:],
			})
		}
	}
	return commits
}

func ParseComment(s string) []Reply {
	var out []Reply
	for _, line := range strings.Split(s, "\n") {
		if reply := ParseReply(line); reply != nil {
			out = append(out, *reply)
		}
	}
	return out
}

func ParsePullRequestURL(prURL string) (string, string, int) {
	if !strings.Contains(prURL, "://") {
		prURL = "https://" + prURL
	}

	u, err := url.Parse(prURL)
	if err != nil {
		panic(err)
	}
	parts := strings.Split(u.Path, "/")
	if u.Hostname() != "github.com" || len(parts) != 5 || parts[3] != "pull" {
		panic(fmt.Errorf("invalid or unsupported release tracker url: %s", prURL))
	}

	owner := parts[1]
	repo := parts[2]
	prNumber, err := strconv.Atoi(parts[4])
	if err != nil {
		panic(err)
	}
	return owner, repo, prNumber
}

func ParseRepoURL(repoURL string) (string, string) {
	if !strings.Contains(repoURL, "://") {
		repoURL = "https://" + repoURL
	}

	u, err := url.Parse(repoURL)
	if err != nil {
		panic(err)
	}
	parts := strings.Split(u.Path, "/")
	if u.Hostname() != "github.com" || len(parts) != 3 {
		panic(fmt.Errorf("invalid or unsupported repo url: %s", repoURL))
	}

	owner := parts[1]
	repo := parts[2]
	return owner, repo
}

//
//func WriteChangelogMarkdown(filename string, chlog Changelog) {
//	funcMap := template.FuncMap{
//		"trimPrefix": strings.TrimPrefix,
//	}
//	tpl := template.Must(template.New("").Funcs(funcMap).Parse(string(templates.MustAsset("changelog.tpl"))))
//	var buf bytes.Buffer
//	err := tpl.Execute(&buf, chlog)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(buf.String())
//	err = os.MkdirAll(filepath.Dir(filename), 0755)
//	if err != nil {
//		panic(err)
//	}
//	err = ioutil.WriteFile(filename, buf.Bytes(), 0644)
//	if err != nil {
//		panic(err)
//	}
//}

func UpdateChangelog(dir, url, tag string, commits []Commit) {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		panic(err)
	}

	filenameChlog := filepath.Join(dir, "CHANGELOG.json")

	var chlog Changelog
	data, err := ioutil.ReadFile(filenameChlog)
	if err == nil {
		err = json.Unmarshal(data, &chlog)
		if err != nil {
			panic(err)
		}
	}

	chlog.Release = releaseNumber

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
				chlog.Projects[repoIdx].Releases = append(chlog.Projects[repoIdx].Releases, ReleaseChangelog{
					Tag:     tag,
					Commits: commits,
				})
			}
			repoFound = true
			break
		}
	}
	if !repoFound {
		chlog.Projects = append(chlog.Projects, ProjectChangelog{
			URL: url,
			Releases: []ReleaseChangelog{
				{
					Tag:     tag,
					Commits: commits,
				},
			},
		})
	}
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

func WriteChangelogMarkdown(dir string, chlog Changelog) {
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

func ProjectsTagged(projects IndependentProjects) bool {
	for repoURL := range projects {
		if !tagged.Has(repoURL) {
			return false
		}
	}
	return true
}

func ProjectCherryPicked(repoURL string, project Project) bool {
	if project.Tags == nil {
		return false
	}

	data := MergeData{Repo: repoURL}
	for _, branch := range project.Tags {
		data.Branch = branch
		if _, ok := merged[data]; !ok {
			return false
		}
	}
	return true
}

// https://developer.github.com/v3/pulls/reviews/#create-a-review-for-a-pull-request
func PRApproved(gh *github.Client, owner string, repo string, prNumber int) bool {
	reviews, err := ListReviews(context.TODO(), gh, owner, repo, prNumber)
	if err != nil {
		panic(err)
	}
	for _, review := range reviews {
		if review.GetState() == "REQUEST_CHANGES" {
			return false
		}
	}
	for _, review := range reviews {
		if review.GetState() == "APPROVED" {
			return true
		}
	}
	return false
}

func CreatePR(gh *github.Client, owner string, repo string, req *github.NewPullRequest, labels ...string) (*github.PullRequest, error) {
	labelSet := sets.NewString(labels...)
	var result *github.PullRequest

	prs, _, err := gh.PullRequests.List(context.TODO(), owner, repo, &github.PullRequestListOptions{
		State: "open",
		Head:  req.GetHead(),
		Base:  req.GetBase(),
		ListOptions: github.ListOptions{
			PerPage: 1,
		},
	})
	if err != nil {
		return nil, err
	}

	if len(prs) == 0 {
		result, _, err = gh.PullRequests.Create(context.TODO(), owner, repo, req)
		// "A pull request already exists" error should NEVER happen since we already checked for existence
		if err != nil {
			return nil, err
		}
		//if e2, ok := err.(*github.ErrorResponse); ok {
		//	var matched bool
		//	for _, entry := range e2.Errors {
		//		if strings.HasPrefix(entry.Message, "A pull request already exists") {
		//			matched = true
		//			break
		//		}
		//	}
		//	if !matched {
		//		return nil, err
		//	}
		//	// else ignore error because pr already exists
		//	// else should NEVER happen since we already checked for existence
		//} else if err != nil {
		//	return nil, err
		//}
	} else {
		result = prs[0]
		for _, label := range result.Labels {
			labelSet.Delete(label.GetName())
		}
	}

	if labelSet.Len() > 0 {
		_, _, err := gh.Issues.AddLabelsToIssue(context.TODO(), owner, repo, result.GetNumber(), labelSet.UnsortedList())
		if err != nil {
			return nil, err
		}
	}

	return result, err
}

func DetectGoMod(dir string) string {
	filename := filepath.Join(dir, "go.mod")
	if !exists(filename) {
		return ""
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	gomod, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		panic(err)
	}
	path := gomod.Module.Mod.Path
	if _, ok := modCache[path]; !ok {
		return path
	}
	return ""
}

func UpdateGoMod(dir string) {
	filename := filepath.Join(dir, "go.mod")
	if !exists(filename) {
		return
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	f, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		panic(err)
	}

	// Add replaces first because it may be coming from forked repo during testing automaton
	for _, x := range f.Replace {
		if gm, ok := modCache[x.Old.Path]; ok && gm.VCSRoot != "" { // meaning using forked repo
			err = f.DropReplace(x.Old.Path, x.Old.Version)
			if err != nil {
				panic(err)
			}
			if v, ok := repoVersion[gm.RepoRoot]; ok {
				err = f.AddReplace(x.Old.Path, x.Old.Version, gm.RepoRoot, v)
				if err != nil {
					panic(err)
				}
			}
		}
	}

	for _, x := range f.Require {
		if gm, ok := modCache[x.Mod.Path]; ok {
			if v, ok := repoVersion[gm.RepoRoot]; ok {
				if gm.VCSRoot != "" {
					// using forked repo, so we need to use replace statement to get the newly tagged code
					// This path should only be taken during testing
					err = f.AddReplace(x.Mod.Path, "", gm.RepoRoot, v)
					if err != nil {
						panic(err)
					}
				} else {
					// we tagged the vcs repo, so we can just use require statement
					// this path should be taken in actual releases, since we tag the vcs repo
					err = f.AddRequire(x.Mod.Path, v)
					if err != nil {
						panic(err)
					}
				}
			}
		}
	}

	data, err = f.Format()
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		panic(err)
	}
}

func repoURL2EnvKey(repoURL string) string {
	if !strings.Contains(repoURL, "://") {
		repoURL = "https://" + repoURL
	}

	u, err := url.Parse(repoURL)
	if err != nil {
		panic(err)
	}
	return strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(path.Join(u.Path, "tag"), "/", "_"), "-", "_"))
}

const (
	gitRoot        = "/tmp/workspace"
	GitHubUserKey  = "GITHUB_USER"
	GitHubTokenKey = "LGTM_GITHUB_TOKEN"
)

var (
	release        Release
	releaseFile    string
	releaseNumber  string
	releaseTracker string
	commentId      int64

	scriptRoot, _ = os.Getwd()
	sessionID     = uuid.New().String()
	replies       Replies
	repoVersion   = map[string]string{}    // repo url -> version
	envVars       = map[string]string{}    // ENV var format(repo url) -> version
	modCache      = map[string]GoImport{}  // module path -> repo
	tagged        = sets.NewString()       // already tagged repos
	merged        = map[MergeData]string{} // (repo, branch) -> sha
	comments      []string
)

func init() {
	flag.StringVar(&releaseFile, "release-file", "", "Path of release file (local file or url is accepted)")
	flag.StringVar(&releaseNumber, "release", "", "Release number")
	flag.StringVar(&releaseTracker, "release-tracker", "", "URL of release tracker pull request")
	flag.Int64Var(&commentId, "comment-id", 0, "Comment Id that triggered this run")
}

func main() {
	flag.Parse()

	sh := shell.NewSession()
	sh.ShowCMD = true
	sh.PipeFail = true
	sh.PipeStdErrors = true

	err := os.RemoveAll(gitRoot)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadFile(releaseFile)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(data, &release)
	if err != nil {
		panic(err)
	}

	for _, projects := range release.Projects {
		for repoURL, project := range projects {
			if project.Tag != nil {
				repoVersion[repoURL] = *project.Tag
				envVars[repoURL2EnvKey(repoURL)] = *project.Tag
			}
		}
	}

	releaseOwner, releaseRepo, releasePR := ParsePullRequestURL(releaseTracker)

	gh := NewGitHubClient()
	pr, _, err := gh.PullRequests.Get(context.TODO(), releaseOwner, releaseRepo, releasePR)
	if err != nil {
		panic(err)
	}
	if pr.GetDraft() {
		fmt.Println("Release tracker pr is currently in draft mode")
		os.Exit(0)
	}
	if pr.GetState() != "open" {
		fmt.Println("Release tracker pr is not open")
		os.Exit(0)
	}
	if !PRApproved(gh, releaseOwner, releaseRepo, releasePR) {
		fmt.Println("PR must be approved to continue")
		os.Exit(0)
	}

	// Build state
	prComments, err := ListComments(context.TODO(), gh, releaseOwner, releaseRepo, releasePR)
	if err != nil {
		panic(err)
	}
	if commentId > 0 {
		// This is done to avoid using any comments that was added after this action was triggered
		idx := -1
		for i, comment := range prComments {
			if comment.GetID() == commentId {
				idx = i
				break
			}
		}
		if idx > -1 {
			prComments = prComments[:idx+1]
		}
	}

	for _, comment := range prComments {
		replies = MergeReplies(replies, ParseComment(comment.GetBody())...)
	}
	for _, reply := range replies[Go] {
		modCache[reply.Go.ModulePath] = GoImport{
			RepoRoot: reply.Go.Repo,
			VCSRoot:  reply.Go.VCSRoot,
		}
	}

	if _, ok := replies[OkToRelease]; !ok {
		fmt.Println("Not /ok-to-release yet")
		os.Exit(0)
	}

	for groupIdx, projects := range release.Projects {
		firstGroup := groupIdx == 0

		// regenerate caches as previous group might have changed stuff
		for _, reply := range replies[Tagged] {
			tagged.Insert(reply.Tagged.Repo)
		}
		for _, reply := range replies[ReadyToTag] {
			merged[MergeData{
				Repo:   reply.ReadyToTag.Repo,
				Branch: "master",
			}] = reply.ReadyToTag.MergeCommitSHA
		}
		for _, reply := range replies[CherryPicked] {
			merged[MergeData{
				Repo:   reply.CherryPicked.Repo,
				Branch: reply.CherryPicked.Branch,
			}] = reply.CherryPicked.MergeCommitSHA
		}

		if ProjectsTagged(projects) {
			continue
		}

		notTagged := sets.NewString()
		for repoURL := range projects {
			if !tagged.Has(repoURL) {
				notTagged.Insert(repoURL)
			}
		}

		var readyToTag sets.String
		if firstGroup {
			readyToTag = notTagged
			notTagged = sets.NewString() // make it empty
		} else {
			readyToTag = sets.NewString()

			// check repos that are /ready-to-tag
			for _, data := range replies[ReadyToTag] {
				repoURL := data.ReadyToTag.Repo
				if notTagged.Has(repoURL) {
					readyToTag.Insert(repoURL)
					notTagged.Delete(repoURL)
				}
			}

			// check repos where all branches have been cherry picked
			for _, repoURL := range notTagged.UnsortedList() {
				if ProjectCherryPicked(repoURL, projects[repoURL]) {
					readyToTag.Insert(repoURL)
					notTagged.Delete(repoURL)
				}
			}

			// skip repos where prs have been opened
			for _, data := range replies[PR] {
				repoURL := data.PR.Repo
				if notTagged.Has(repoURL) {
					notTagged.Delete(repoURL)
				}
			}
		}

		// Now, open pr for notTagged
		for _, repoURL := range notTagged.UnsortedList() {
			oneliners.FILE()
			err = PrepareProject(gh, sh, releaseTracker, repoURL, projects[repoURL])
			if err != nil {
				panic(err)
			}
		}

		// Tag the repos in readyToTag
		for _, repoURL := range readyToTag.UnsortedList() {
			oneliners.FILE()
			err = ReleaseProject(gh, sh, releaseTracker, repoURL, projects[repoURL])
			if err != nil {
				panic(err)
			}
		}

		oneliners.FILE("COMMENTS>>>>", strings.Join(comments, "\n"))
		if len(comments) > 0 {
			_, _, err := gh.Issues.CreateComment(context.TODO(), releaseOwner, releaseRepo, releasePR, &github.IssueComment{
				Body: github.String(strings.Join(comments, "\n")),
			})
			if err != nil {
				panic(err)
			}
			os.Exit(0) // Let next execution to pick up
		}
	}
}

func CreateReleaseFile() Release {
	return Release{
		ProductLine: "Stash",
		Projects: []IndependentProjects{
			{
				"github.com/appscode-cloud/apimachinery": Project{
					Tag: github.String("v0.10.0-alpha.2"),
				},
			},
			{
				"github.com/appscode-cloud/stash": Project{
					Tag: github.String("v0.10.0-alpha.2"),
				},
				"github.com/appscode-cloud/postgres": Project{
					Tags: map[string]string{
						"9.6-v1":  "release-9.6",
						"10.2-v1": "release-10.2",
						"10.6-v1": "release-10.6",
						"11.1-v1": "release-11.1",
						"11.2-v1": "release-11.2",
					},
				},
				"github.com/appscode-cloud/cli": Project{
					Tag: github.String("v0.10.0-alpha.2"),
				},
			},
		},
	}
}

func main_PrintReleaseFile() {
	rel := CreateReleaseFile()
	data, err := yaml.Marshal(rel)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))

	data, err = json.MarshalIndent(rel, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))

	// https://github.com/tamalsaha/gh-release-automation-testing/pull/8
}

func main_ParsePullRequestURL() {
	ParsePullRequestURL("https://github.com/appscodelabs/gh-release-automation-testing/pull/21")
	ParseRepoURL("https://github.com/appscodelabs/gh-release-automation-testing")

	ParsePullRequestURL("github.com/appscodelabs/gh-release-automation-testing/pull/21")
	ParseRepoURL("github.com/appscode-cloud/apimachinery")
}

func mm() {
	localfile := filepath.Join(os.TempDir(), sessionID, "release.txt")
	opts := func(c *getter.Client) error {
		pwd, err := os.Getwd()
		if err != nil {
			return err
		}
		c.Pwd = pwd
		return nil
	}
	err := getter.GetFile(localfile, releaseFile, opts)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadFile(localfile)
	if err != nil {
		panic(err)
	}

	var r Release
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

	gh := NewGitHubClient()

	pr, err := CreatePR(gh, owner, repo, &github.NewPullRequest{
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

	commits := ListCommits(sh, "0ab9faa68308cd646e1e63271950cf75e3cf62c0", "v0.9.0-rc.6")
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

	tags, err := ListTags(sh)

	fmt.Println(tags, err)
}

func main_RepoModified() {
	sh := shell.NewSession()
	sh.ShowCMD = true
	sh.PipeFail = true
	sh.PipeStdErrors = true

	if RepoModified(sh) {
		fmt.Println("Something to commit")
	}
}

func main_WriteChangelogMarkdown() {
	chlog := Changelog{
		Release: "v2020.6.12",
		Projects: []ProjectChangelog{
			{
				URL: "github.com/stashed/apimachinery",
				Releases: []ReleaseChangelog{
					{
						Tag: "v0.1.0",
						Commits: []Commit{
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
	WriteChangelogMarkdown("/tmp", chlog)
}

func mainMergeReply() {
	var replies Replies

	r := Reply{
		Type: Tagged,
		Tagged: &TaggedReplyData{
			Repo: "a/b",
		},
	}

	replies = MergeReply(replies, r)
	fmt.Println(replies)
}

type GoImport struct {
	RepoRoot string
	VCSRoot  string
}

func (g GoImport) String() string {
	if g.VCSRoot == "" {
		return g.RepoRoot
	}
	return g.RepoRoot + " " + g.VCSRoot
}

type MetaData struct {
	GoImport string `meta:"go-import"`
}

// ref: https://gist.github.com/inotnako/c4a82f6723f6ccea5d83c5d3689373dd
// ref: https://github.com/keighl/metabolize
// ref: https://github.com/rsc/go-import-redirector/blob/master/main.go#L134
//ref: https://github.com/appscodelabs/gh-release-automation-testing/issues/22
func main_DetectVCSRoot() {
	// res, _ := http.Get("https://stash.appscode.dev/cli?go-get=1")
	// res, _ := http.Get("https://k8s.io/api?go-get=1")
	// "https://github.com/cloudevents/sdk-go/blob/master/samples/kafka?go-get=1"

	r, err := DetectVCSRoot("stash.appscode.dev/cli")
	if err != nil {
		panic(err)
	}
	fmt.Println(r)
}

func DetectVCSRoot(repoURL string) (string, error) {
	if !strings.Contains(repoURL, "://") {
		repoURL = "https://" + repoURL
	}

	uRepo, err := url.Parse(repoURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse repo url: %v", err)
	}
	qRepo := uRepo.Query()
	qRepo.Set("go-get", "1")
	uRepo.RawQuery = qRepo.Encode()

	res, err := http.Get(uRepo.String())
	if err != nil {
		return "", err
	}
	data := new(MetaData)

	err = metabolize.Metabolize(res.Body, data)
	if err != nil {
		return "", err
	}

	// GoImport: stash.appscode.dev/cli git https://github.com/stashed/cli
	if data.GoImport == "" {
		return "", fmt.Errorf("%s is missing go-import meta tag", uRepo.String())
	}
	fmt.Printf("GoImport: %s\n", data.GoImport)

	parts := strings.Fields(data.GoImport)
	if len(parts) != 3 {
		return "", fmt.Errorf("%s contains badly formatted go-import meta tag %s", uRepo.String(), data.GoImport)
	}

	uVCS, err := url.Parse(parts[2])
	if err != nil {
		return "", fmt.Errorf("failed to parse VCS root %s: %v", parts[2], err)
	}
	//uVCS.Scheme = ""
	vcsURL := path.Join(uVCS.Hostname(), uVCS.Path)
	return strings.TrimSuffix(vcsURL, path.Ext(vcsURL)), nil
}
