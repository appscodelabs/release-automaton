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
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/appscodelabs/release-automaton/api"
	"github.com/appscodelabs/release-automaton/lib"

	"github.com/Masterminds/semver"
	shell "github.com/codeskyblue/go-sh"
	"github.com/google/go-github/v32/github"
	"github.com/spf13/cobra"
	"github.com/tamalsaha/go-oneliners"
	"golang.org/x/mod/modfile"
	"gomodules.xyz/envsubst"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/yaml"
)

var (
	release        api.Release
	releaseFile    string
	releaseTracker string
	commentId      int64

	empty          = struct{}{}
	scriptRoot, _  = os.Getwd()
	changelogRoot  = filepath.Join(scriptRoot, "CHANGELOG")
	replies        api.Replies
	repoVersion    = map[string]string{}          // repo url -> version
	envVars        = map[string]string{}          // ENV var format(repo url) -> version
	modCache       = map[string]lib.GoImport{}    // module path -> repo
	tagged         = sets.NewString()             // already tagged repos
	merged         = map[api.MergeData]string{}   // (repo, branch) -> sha
	chartsMerged   = map[api.MergeData]struct{}{} // (repo, tag) -> empty
	chartPublished = sets.NewString()             // set(chart repo url)
	comments       []string
)

func NewCmdReleaseRun() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "run",
		Short:             "Run release process",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			runAutomaton()
		},
	}

	cmd.Flags().StringVar(&releaseFile, "release-file", "", "Path of release file (local file or url is accepted)")
	cmd.Flags().StringVar(&releaseTracker, "release-tracker", "", "URL of release tracker pull request")
	cmd.Flags().Int64Var(&commentId, "comment-id", 0, "Comment Id that triggered this run")
	return cmd
}

func runAutomaton() {
	sh := shell.NewSession()
	sh.ShowCMD = true
	sh.PipeFail = true
	sh.PipeStdErrors = true

	err := os.RemoveAll(api.Workspace)
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
				if project.Key != "" {
					envVars[key2EnvKey(project.Key)] = *project.Tag
				}
			}
			if project.ReadyToTag {
				replies = api.MergeReplies(replies, api.Reply{
					Type: api.ReadyToTag,
					ReadyToTag: &api.ReadyToTagReplyData{
						Repo:           repoURL,
						MergeCommitSHA: "",
					},
				})
			}
		}
	}

	releaseOwner, releaseRepo, releasePR := lib.ParsePullRequestURL(releaseTracker)

	gh := lib.NewGitHubClient()
	pr, _, err := gh.PullRequests.Get(context.TODO(), releaseOwner, releaseRepo, releasePR)
	if err != nil {
		panic(err)
	}
	if pr.GetDraft() {
		fmt.Println("Release tracker pr is currently in draft mode")
		return
	}
	if pr.GetState() != "open" {
		fmt.Println("Release tracker pr is not open")
		return
	}
	if !lib.PRApproved(gh, releaseOwner, releaseRepo, releasePR) {
		fmt.Println("PR must be approved to continue")
		return
	}

	// Build state
	prComments, err := lib.ListComments(context.TODO(), gh, releaseOwner, releaseRepo, releasePR)
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
		replies = api.MergeReplies(replies, lib.ParseComment(comment.GetBody())...)
	}
	for _, reply := range replies[api.Go] {
		modCache[reply.Go.ModulePath] = lib.GoImport{
			RepoRoot: reply.Go.Repo,
			VCSRoot:  reply.Go.VCSRoot,
		}
	}

	if _, ok := replies[api.OkToRelease]; !ok {
		fmt.Println("Not /ok-to-release yet")
		return
	}
	if _, ok := replies[api.Done]; ok {
		fmt.Println("Already done!")
		return
	}

	existingLabels, err := lib.ListLabelsByIssue(context.TODO(), gh, releaseOwner, releaseRepo, releasePR)
	if err != nil {
		panic(err)
	}
	if existingLabels.Has(api.LabelLocked) {
		fmt.Println("Already locked, exiting ...")
		return
	}

	_, _, err = gh.Issues.AddLabelsToIssue(context.TODO(), releaseOwner, releaseRepo, releasePR, []string{
		api.LabelLocked,
	})
	if err != nil {
		panic(err)
	}
	defer func() {
		err = lib.RemoveLabel(gh, releaseOwner, releaseRepo, releasePR, api.LabelLocked)
		if err != nil {
			panic(err)
		}
	}()

	for groupIdx, projects := range release.Projects {
		firstGroup := groupIdx == 0

		// regenerate caches as previous group might have changed stuff
		for _, reply := range replies[api.Tagged] {
			tagged.Insert(reply.Tagged.Repo)
		}
		for _, reply := range replies[api.ReadyToTag] {
			merged[api.MergeData{
				Repo: reply.ReadyToTag.Repo,
				Ref:  api.BranchMaster,
			}] = reply.ReadyToTag.MergeCommitSHA
		}
		for _, reply := range replies[api.CherryPicked] {
			merged[api.MergeData{
				Repo: reply.CherryPicked.Repo,
				Ref:  reply.CherryPicked.Branch,
			}] = reply.CherryPicked.MergeCommitSHA
		}
		for _, reply := range replies[api.Chart] {
			chartsMerged[api.MergeData{
				Repo: reply.Chart.Repo,
				Ref:  reply.Chart.Tag,
			}] = empty
		}
		for _, reply := range replies[api.ChartPublished] {
			chartPublished.Insert(reply.ChartPublished.Repo)
		}

		if ProjectsDone(projects) {
			continue
		}

		chartsReadyToPublish := sets.NewString()
		chartsYetToMerge := map[api.MergeData]struct{}{}

		notTagged := sets.NewString()
		openPRs := sets.NewString()
		for repoURL, project := range projects {
			if len(project.Charts) == 0 {
				// Skip if invoked by /chart comment
				for _, reply := range lib.ParseComment(prComments[len(prComments)-1].GetBody()) {
					if reply.Type == api.Chart {
						return
					}
				}

				if !tagged.Has(repoURL) {
					notTagged.Insert(repoURL)
				}
			} else {
				for _, chartRepo := range project.Charts {
					if tags, ok := findRepoTags(chartRepo); ok {
						for _, tag := range tags {
							mergeKey := api.MergeData{
								Repo: chartRepo,
								Ref:  tag,
							}
							if _, ok := chartsMerged[mergeKey]; !ok {
								chartsYetToMerge[mergeKey] = empty
							}
						}
					}
				}

				if len(chartsYetToMerge) == 0 && !chartPublished.Has(repoURL) {
					chartsReadyToPublish.Insert(repoURL)
				}
			}
		}

		var readyToTag sets.String
		if firstGroup {
			readyToTag = notTagged
			notTagged = sets.NewString() // make it empty
		} else {
			readyToTag = sets.NewString()

			// check repos that are /ready-to-tag
			for _, data := range replies[api.ReadyToTag] {
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
			for _, data := range replies[api.PR] {
				repoURL := data.PR.Repo
				if notTagged.Has(repoURL) {
					notTagged.Delete(repoURL)
					openPRs.Insert(repoURL)
				}
			}
		}

		// Now, open pr for notTagged
		for _, repoURL := range notTagged.UnsortedList() {
			oneliners.FILE()
			project := projects[repoURL]
			if project.Tag == nil && len(project.Tags) == 0 {
				err = PrepareExternalProject(gh, sh, releaseTracker, repoURL, project)
				chlog := lib.LoadChangelog(filepath.Join(changelogRoot, release.Release), release)
				if project.Changelog == api.StandaloneWebsiteChangelog {
					lib.WriteChangelogMarkdown(filepath.Join(changelogRoot, release.Release, "docs_changelog.md"), "standalone-changelog.tpl", chlog)
				} else if project.Changelog == api.SharedWebsiteChangelog {
					lib.WriteChangelogMarkdown(filepath.Join(changelogRoot, release.Release, "docs_changelog.md"), "shared-changelog.tpl", chlog)
				}
				if lib.AnyRepoModified(scriptRoot, sh) {
					err = lib.CommitAnyRepo(scriptRoot, sh, "", "Update changelog")
					if err != nil {
						panic(err)
					}
					err = lib.PushAnyRepo(scriptRoot, sh, false)
					if err != nil {
						panic(err)
					}
				}
			} else {
				err = PrepareProject(gh, sh, releaseTracker, repoURL, project)
			}
			if err != nil {
				panic(err)
			}
		}

		// Tag the repos in readyToTag
		for _, repoURL := range readyToTag.UnsortedList() {
			oneliners.FILE()
			err = ReleaseProject(sh, releaseTracker, repoURL, projects[repoURL])
			if err != nil {
				panic(err)
			}
		}

		// Publish chart registry
		for _, repoURL := range chartsReadyToPublish.UnsortedList() {
			oneliners.FILE()
			owner, repo := lib.ParseRepoURL(repoURL)
			err = lib.LabelPR(gh, owner, repo, fmt.Sprintf("%s@%s", release.ProductLine, release.Release), api.BranchMaster, "automerge")
			if err != nil {
				panic(err)
			}
		}

		oneliners.FILE("COMMENTS>>>>", strings.Join(comments, "\n"))
		if len(comments) > 0 {
			comments = lib.UniqComments(comments)
			_, _, err := gh.Issues.CreateComment(context.TODO(), releaseOwner, releaseRepo, releasePR, &github.IssueComment{
				Body: github.String(strings.Join(comments, "\n")),
			})
			if err != nil {
				panic(err)
			}
			return // Let next execution to pick up
		}

		if openPRs.Len() > 0 {
			fmt.Println("Waiting for prs to close:")
			for _, pr := range openPRs.List() {
				fmt.Println(">>> " + pr)
			}
			return
		}

		if len(chartsYetToMerge) > 0 {
			fmt.Println("Waiting for charts to be merged:")
			for data := range chartsYetToMerge {
				fmt.Println(">>> ", data)
			}
			return
		}

		if chartsReadyToPublish.Len() > 0 {
			fmt.Println("Waiting for charts to be published:")
			for _, repoURL := range chartsReadyToPublish.List() {
				fmt.Println(">>> ", repoURL)
			}
			return
		}
	}

	openPRs := sets.NewString()
	// skip repos where prs have been opened
	for _, data := range replies[api.PR] {
		openPRs.Insert(data.PR.Repo)
	}
	for repoURL, project := range release.ExternalProjects {
		if !openPRs.Has(repoURL) {
			// Do not want external projects to report back, so releaseTracker is not set.
			err = PrepareExternalProject(gh, sh, "", repoURL, project)
			if err != nil {
				break
			}
		}
	}

	oneliners.FILE("COMMENTS>>>>", strings.Join(comments, "\n"))
	{
		comments = append(comments, string(api.Done))
		comments = lib.UniqComments(comments)
		_, _, err := gh.Issues.CreateComment(context.TODO(), releaseOwner, releaseRepo, releasePR, &github.IssueComment{
			Body: github.String(strings.Join(comments, "\n")),
		})
		if err != nil {
			panic(err)
		}
		return // Let next execution to pick up
	}
}

func PrepareProject(gh *github.Client, sh *shell.Session, releaseTracker, repoURL string, project api.Project) error {
	if project.Tags != nil && project.Tag != nil {
		return fmt.Errorf("repo %s is provided an invalid project configuration which uses both tag and tags", repoURL)
	}

	// pushd, popd
	wdOrig := sh.Getwd()
	defer sh.SetDir(wdOrig)

	owner, repo := lib.ParseRepoURL(repoURL)

	// TODO: cache git repo
	wdCur := filepath.Join(api.Workspace, owner)
	err := os.MkdirAll(wdCur, 0755)
	if err != nil {
		return err
	}

	if !lib.Exists(filepath.Join(wdCur, repo)) {
		sh.SetDir(wdCur)

		err = sh.Command("git",
			"clone",
			// "--no-tags", //TODO: ok?
			"--no-recurse-submodules",
			//"--depth=1",
			//"--no-single-branch",
			fmt.Sprintf("https://%s:%s@%s.git", os.Getenv(api.GitHubUserKey), os.Getenv(api.GitHubTokenKey), repoURL),
		).Run()
		if err != nil {
			return err
		}
	}
	wdCur = filepath.Join(wdCur, repo)
	sh.SetDir(wdCur)

	modPath := DetectGoMod(wdCur)
	if modPath != "" {
		gm := lib.GoImport{
			RepoRoot: repoURL,
		}
		vcs, err := lib.DetectVCSRoot(modPath)
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
			*project.Tag: api.BranchMaster, // pr always opened against master branch
		}
	}

	// All remote tags exist, so only add Go module path if needed.
	if lib.MeetsCondition(lib.RemoteTagExists, sh, lib.Keys(tags)...) {
		var ok bool
		// Make sure /tagged, /cherry-picked comments exist
		if project.Tag != nil {
			sha := lib.GetRemoteTag(sh, *project.Tag)
			if replies, ok = api.AppendReplyIfMissing(replies, api.Reply{
				Type: api.ReadyToTag,
				ReadyToTag: &api.ReadyToTagReplyData{
					Repo:           repoURL,
					MergeCommitSHA: sha,
				},
			}); ok {
				comments = append(comments, fmt.Sprintf("%s %s %s", api.ReadyToTag, repoURL, sha))
			}
		}
		if project.Tags != nil {
			for tag, branch := range project.Tags {
				sha := lib.GetRemoteTag(sh, tag)
				if replies, ok = api.AppendReplyIfMissing(replies, api.Reply{
					Type: api.CherryPicked,
					CherryPicked: &api.CherryPickedReplyData{
						Repo:           repoURL,
						Branch:         branch,
						MergeCommitSHA: sha,
					},
				}); ok {
					comments = append(comments, fmt.Sprintf("%s %s %s %s", api.CherryPicked, repoURL, branch, sha))
				}
			}
		}

		if replies, ok = api.AppendReplyIfMissing(replies, api.Reply{
			Type: api.Tagged,
			Tagged: &api.TaggedReplyData{
				Repo: repoURL,
			},
		}); ok {
			comments = append(comments, fmt.Sprintf("%s %s", api.Tagged, repoURL))
		}

		if modPath != "" {
			AppendGo(modPath)
		}
		return nil
	}

	usesCherryPick := project.Tags != nil && project.Tag == nil
	if usesCherryPick {
		tags[release.Release] = api.BranchMaster // if cherry pick is used, there must be an extra pr against the master branch
	}

	for _, pair := range lib.ToOrderedPair(tags) {
		tag, branch := pair.Key, pair.Value

		if usesCherryPick {
			// remote branch must already exist
			if !lib.RemoteBranchExists(sh, branch) {
				return fmt.Errorf("repo %s is missing branch for tag %s", repoURL, tag)
			}
		}

		// -----------------------

		vars := lib.MergeMaps(map[string]string{
			repoURL2EnvKey(repoURL): tag,
			"SCRIPT_ROOT":           scriptRoot,
			"WORKSPACE":             sh.Getwd(),
			"TAG":                   tag,
			"PRODUCT_LINE":          release.ProductLine,
			"RELEASE":               release.Release,
			"RELEASE_TRACKER":       releaseTracker,
		}, envVars)

		headBranch := fmt.Sprintf("%s-%s", release.Release, branch)

		err = sh.Command("git", "checkout", branch).Run()
		if err != nil {
			return err
		}

		err = sh.Command("git", "checkout", "-b", headBranch).Run()
		if err != nil {
			return err
		}

		if lib.Exists(filepath.Join(wdCur, "go.mod")) {
			// Update Go mod
			UpdateGoMod(wdCur)
			if lib.RepoModified(sh) {
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
			err = lib.Execute(sh, cmd)
			if err != nil {
				return err
			}
		}

		if lib.RepoModified(sh) {
			messages := []string{
				"ProductLine: " + release.ProductLine,
				"Release: " + release.Release,
			}
			if !usesCherryPick || branch != api.BranchMaster {
				// repos that use cherry pick, a pr is opened against the master branch
				// That pr MUST NOT report back to release tracker.
				messages = append(messages, "Release-tracker: "+releaseTracker)
			}
			err = lib.CommitRepo(sh, tag, messages...)
			if err != nil {
				return err
			}
			err = lib.PushRepo(sh, true)
			if err != nil {
				return err
			}

			// open pr against project repo
			pr, err := lib.CreatePR(gh, owner, repo, &github.NewPullRequest{
				Title:               github.String(fmt.Sprintf("Prepare for release %s", tag)),
				Head:                github.String(headBranch),
				Base:                github.String(branch),
				Body:                github.String(lib.LastCommitBody(sh, true)),
				MaintainerCanModify: github.Bool(true),
				Draft:               github.Bool(false),
			}, "automerge")
			if err != nil {
				panic(err)
			}

			// add comments to release repo
			comments = append(comments, fmt.Sprintf("%s %s", api.PR, pr.GetHTMLURL()))
		} else {
			comments = append(comments, fmt.Sprintf("%s %s %s", api.ReadyToTag, repoURL, lib.LastCommitSHA(sh)))
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
	comments = append(comments, fmt.Sprintf(`%s %s %s %s`, api.Go, gm.RepoRoot, modPath, gm.VCSRoot))
}

func ReleaseProject(sh *shell.Session, releaseTracker, repoURL string, project api.Project) error {
	if project.Tags != nil && project.Tag != nil {
		return fmt.Errorf("repo %s is provided an invalid project configuration which uses both tag and tags", repoURL)
	}

	// pushd, popd
	wdOrig := sh.Getwd()
	defer sh.SetDir(wdOrig)

	owner, repo := lib.ParseRepoURL(repoURL)

	// TODO: cache git repo
	wdCur := filepath.Join(api.Workspace, owner)
	err := os.MkdirAll(wdCur, 0755)
	if err != nil {
		return err
	}

	if !lib.Exists(filepath.Join(wdCur, repo)) {
		sh.SetDir(wdCur)

		err = sh.Command("git",
			"clone",
			// "--no-tags", //TODO: ok?
			"--no-recurse-submodules",
			//"--depth=1",
			//"--no-single-branch",
			fmt.Sprintf("https://%s:%s@%s.git", os.Getenv(api.GitHubUserKey), os.Getenv(api.GitHubTokenKey), repoURL),
		).Run()
		if err != nil {
			return err
		}
	}
	wdCur = filepath.Join(wdCur, repo)
	sh.SetDir(wdCur)

	modPath := DetectGoMod(wdCur)
	if modPath != "" {
		gm := lib.GoImport{
			RepoRoot: repoURL,
		}
		vcs, err := lib.DetectVCSRoot(modPath)
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
	if lib.MeetsCondition(lib.RemoteTagExists, sh, lib.Keys(tags)...) {
		// make sure /tagged is appended for next group of projects in this run
		// and added to comments for next run
		var ok bool
		if replies, ok = api.AppendReplyIfMissing(replies, api.Reply{
			Type: api.Tagged,
			Tagged: &api.TaggedReplyData{
				Repo: repoURL,
			},
		}); ok {
			comments = append(comments, fmt.Sprintf("%s %s", api.Tagged, repoURL))
		}

		if modPath != "" {
			AppendGo(modPath)
		}
		return nil
	}

	usesCherryPick := project.Tags != nil && project.Tag == nil

	for _, pair := range lib.ToOrderedPair(tags) {
		tag, branch := pair.Key, pair.Value

		vTag, err := semver.NewVersion(tag)
		if err != nil {
			return err
		}

		// detect branch
		if usesCherryPick {
			// remote branch must already exist
			if !lib.RemoteBranchExists(sh, branch) {
				return fmt.Errorf("repo %s is missing branch for tag %s", repoURL, tag)
			}
		} else {
			if project.ReleaseBranch != "" {
				vars := lib.MergeMaps(map[string]string{
					repoURL2EnvKey(repoURL): tag,
					"SCRIPT_ROOT":           scriptRoot,
					"WORKSPACE":             sh.Getwd(),
					"TAG":                   tag,
					"PRODUCT_LINE":          release.ProductLine,
					"RELEASE":               release.Release,
					"RELEASE_TRACKER":       releaseTracker,
				}, envVars)
				branch, err = envsubst.EvalMap(project.ReleaseBranch, vars)
				if err != nil {
					return err
				}
				tags[tag] = branch
			} else if vTag.Patch() > 0 { // PATCH release
				if vTag.Prerelease() != "" {
					panic(fmt.Errorf("version %s is invalid because it is a patch release but includes a pre-release component", tag))
				}

				patchBranch := fmt.Sprintf("release-%d.%d.%d", vTag.Major(), vTag.Minor(), vTag.Patch())
				if lib.RemoteBranchExists(sh, patchBranch) {
					branch = patchBranch
				} else {
					minorBranch := fmt.Sprintf("release-%d.%d", vTag.Major(), vTag.Minor())
					if lib.RemoteBranchExists(sh, minorBranch) {
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

		if usesCherryPick || (vTag.Patch() > 0 && project.ReleaseBranch == "") {
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

			err = lib.TagRepo(sh, tag, "ProductLine: "+release.ProductLine, "Release: "+release.Release, "Release-tracker: "+releaseTracker)
			if err != nil {
				return err
			}
			err = lib.PushRepo(sh, true)
			if err != nil {
				return err
			}
		} else if vTag.Patch() == 0 || project.ReleaseBranch != "" {
			if lib.RemoteBranchExists(sh, branch) {
				err = sh.Command("git", "checkout", branch).Run()
				if err != nil {
					return err
				}
				if branch != api.BranchMaster {
					ref := api.BranchMaster
					if sha, found := MergedCommitSHA(repoURL, branch, usesCherryPick); found {
						ref = sha
					}
					err = sh.Command("git", "merge", ref).Run()
					if err != nil {
						return err
					}
				}
				err = lib.TagRepo(sh, tag, "ProductLine: "+release.ProductLine, "Release: "+release.Release, "Release-tracker: "+releaseTracker)
				if err != nil {
					return err
				}
				err = lib.PushRepo(sh, true)
				if err != nil {
					return err
				}
			} else {
				err = sh.Command("git", "checkout", api.BranchMaster).Run()
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
				err = lib.TagRepo(sh, tag, "ProductLine: "+release.ProductLine, "Release: "+release.Release, "Release-tracker: "+releaseTracker)
				if err != nil {
					return err
				}
				err = lib.PushRepo(sh, true)
				if err != nil {
					return err
				}
			}
		}

		if project.Changelog == api.AddToChangelog {
			tags, err := lib.ListTags(sh)
			if err != nil {
				return err
			}
			tagSet := sets.NewString(tags...)
			tagSet.Insert(tag)

			vs := make([]*semver.Version, 0, tagSet.Len())
			for _, x := range tagSet.UnsortedList() {
				v := semver.MustParse(x)
				// filter out lower importance tags
				if api.AtLeastAsImp(vTag, v) {
					vs = append(vs, v)
				}
			}
			sort.Sort(api.SemverCollection(vs))

			var tagIdx = -1
			for idx, vs := range vs {
				if vs.Equal(vTag) {
					tagIdx = idx
					break
				}
			}

			var commits []api.Commit
			if tagIdx == 0 {
				commits = lib.ListCommits(sh, lib.FirstCommit(sh), vs[tagIdx].Original())
			} else {
				commits = lib.ListCommits(sh, vs[tagIdx-1].Original(), vs[tagIdx].Original())
			}
			lib.UpdateChangelog(filepath.Join(changelogRoot, release.Release), release, repoURL, tag, commits)
			if lib.AnyRepoModified(scriptRoot, sh) {
				err = lib.CommitAnyRepo(scriptRoot, sh, "", "Update changelog")
				if err != nil {
					return err
				}
				err = lib.PushAnyRepo(scriptRoot, sh, false)
				if err != nil {
					return err
				}
			}
		}
	}

	// add comments to release repo
	{
		comments = append(comments, fmt.Sprintf("%s %s", api.Tagged, repoURL))
		if modPath != "" {
			AppendGo(modPath)
		}
	}

	return nil
}

func PrepareExternalProject(gh *github.Client, sh *shell.Session, releaseTracker, repoURL string, project api.ProjectMeta) error {
	// pushd, popd
	wdOrig := sh.Getwd()
	defer sh.SetDir(wdOrig)

	owner, repo := lib.ParseRepoURL(repoURL)

	// TODO: cache git repo
	wdCur := filepath.Join(api.Workspace, owner)
	err := os.MkdirAll(wdCur, 0755)
	if err != nil {
		return err
	}

	if !lib.Exists(filepath.Join(wdCur, repo)) {
		sh.SetDir(wdCur)

		err = sh.Command("git",
			"clone",
			"--no-tags",
			"--no-recurse-submodules",
			"--depth=1",
			"--no-single-branch",
			fmt.Sprintf("https://%s:%s@%s.git", os.Getenv(api.GitHubUserKey), os.Getenv(api.GitHubTokenKey), repoURL),
		).Run()
		if err != nil {
			return err
		}
	}
	wdCur = filepath.Join(wdCur, repo)
	sh.SetDir(wdCur)

	// -----------------------

	vars := lib.MergeMaps(map[string]string{
		"SCRIPT_ROOT":     scriptRoot,
		"WORKSPACE":       sh.Getwd(),
		"PRODUCT_LINE":    release.ProductLine,
		"RELEASE":         release.Release,
		"RELEASE_TRACKER": releaseTracker,
	}, envVars)

	headBranch := fmt.Sprintf("%s-%s", release.ProductLine, release.Release)

	err = sh.Command("git", "checkout", api.BranchMaster).Run()
	if err != nil {
		return err
	}

	err = sh.Command("git", "checkout", "-b", headBranch).Run()
	if err != nil {
		return err
	}

	if lib.Exists(filepath.Join(wdCur, "go.mod")) {
		// Update Go mod
		UpdateGoMod(wdCur)
		if lib.RepoModified(sh) {
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

	for _, cmd := range project.GetCommands() {
		cmd, err = envsubst.EvalMap(cmd, vars)
		if err != nil {
			return err
		}

		err = lib.Execute(sh, cmd)
		if err != nil {
			return err
		}
	}

	if lib.RepoModified(sh) {
		messages := []string{
			fmt.Sprintf("Update for release %s@%s", release.ProductLine, release.Release),
			"ProductLine: " + release.ProductLine,
			"Release: " + release.Release,
		}
		if releaseTracker != "" {
			messages = append(messages, "Release-tracker: "+releaseTracker)
		}
		err = lib.CommitRepo(sh, "", messages...)
		if err != nil {
			return err
		}
		err = lib.PushRepo(sh, true)
		if err != nil {
			return err
		}

		// open pr against project repo
		pr, err := lib.CreatePR(gh, owner, repo, &github.NewPullRequest{
			Title:               github.String(messages[0]),
			Head:                github.String(headBranch),
			Base:                github.String(api.BranchMaster),
			Body:                github.String(strings.Join(messages[1:], "\n")),
			MaintainerCanModify: github.Bool(true),
			Draft:               github.Bool(false),
		}, "automerge")
		if err != nil {
			panic(err)
		}

		// add comments to release repo
		comments = append(comments, fmt.Sprintf("%s %s", api.PR, pr.GetHTMLURL()))
	}
	return nil
}

func MergedCommitSHA(repoURL, branch string, useCherryPick bool) (string, bool) {
	key := api.MergeData{
		Repo: repoURL,
		Ref:  branch,
	}
	if !useCherryPick {
		key.Ref = api.BranchMaster
	}
	sha, ok := merged[key]
	return sha, ok
}

func ProjectsDone(projects api.IndependentProjects) bool {
	for repoURL, project := range projects {
		done := (len(project.Charts) == 0 && tagged.Has(repoURL)) ||
			(len(project.Charts) > 0 && chartPublished.Has(repoURL))
		if !done {
			return false
		}
	}
	return true
}

func ProjectCherryPicked(repoURL string, project api.Project) bool {
	if project.Tags == nil {
		return false
	}

	data := api.MergeData{Repo: repoURL}
	for _, branch := range project.Tags {
		data.Ref = branch
		if _, ok := merged[data]; !ok {
			return false
		}
	}
	return true
}

func DetectGoMod(dir string) string {
	filename := filepath.Join(dir, "go.mod")
	if !lib.Exists(filename) {
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
	if !lib.Exists(filename) {
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
	return toEnvKey(path.Join(u.Path, "tag"))
}

func key2EnvKey(key string) string {
	return toEnvKey(path.Join(key, "version"))
}

func toEnvKey(key string) string {
	key = strings.Trim(key, "/")
	key = strings.ReplaceAll(key, "/", "_")
	key = strings.ReplaceAll(key, "-", "_")
	return strings.ToUpper(key)
}

func findRepoTags(reg string) ([]string, bool) {
	for _, projects := range release.Projects {
		for repoURL, project := range projects {
			if repoURL != reg {
				continue
			}
			if project.Tag != nil {
				return []string{*project.Tag}, true
			}
			if project.Tags != nil {
				return lib.Keys(project.Tags), true
			}
		}
	}
	return nil, false
}
