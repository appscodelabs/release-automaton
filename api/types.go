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

package api

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
	"gomodules.xyz/semvers"
)

type ProjectMeta interface {
	GetCommands() []string
}

type ExternalProject struct {
	Commands []string `json:"commands,omitempty"`
}

func (p ExternalProject) GetCommands() []string {
	return p.Commands
}

type Project struct {
	Key           string            `json:"key,omitempty"`
	Tag           *string           `json:"tag,omitempty"`
	Tags          map[string]string `json:"tags,omitempty"` // tag-> branch
	ChartNames    []string          `json:"chartNames,omitempty"`
	ChartRepos    []string          `json:"charts,omitempty"`
	Commands      []string          `json:"commands,omitempty"`
	ReleaseBranch string            `json:"release_branch,omitempty"`
	ReadyToTag    bool              `json:"ready_to_tag,omitempty"`
	Changelog     ChangelogStatus   `json:"changelog,omitempty"`
	SubProjects   []string          `json:"sub_projects,omitempty"`
}

func (p Project) GetCommands() []string {
	return p.Commands
}

type ChangelogStatus string

const (
	AddToChangelog             ChangelogStatus = "" // by default show up in changelog
	SkipChangelog              ChangelogStatus = "Skip"
	StandaloneWebsiteChangelog ChangelogStatus = "StandaloneWebsite"
	SharedWebsiteChangelog     ChangelogStatus = "SharedWebsite"
)

type IndependentProjects map[string]Project

type Release struct {
	ProductLine       string `json:"product_line"`
	Release           string `json:"release"`
	DocsURLTemplate   string `json:"docs_url_template"` // "https://stash.run/docs/%s"
	KubernetesVersion string `json:"kubernetes_version"`
	// These projects can be released in sequence
	Projects         []IndependentProjects      `json:"projects"`
	ExternalProjects map[string]ExternalProject `json:"external_projects,omitempty"`
}

func (r Release) Validate() error {
	if r.Release == "" {
		return fmt.Errorf("missing release number")
	}
	v, err := StrictParseVersion(r.Release)
	if err != nil {
		return err
	}
	for _, projects := range r.Projects {
		for repoURL, project := range projects {
			// only check projects that uses semver tags (ie, does not match release number)
			if project.Tag != nil && r.Release != *project.Tag {
				projectVersion, err := StrictParseVersion(*project.Tag)
				if err != nil {
					return fmt.Errorf("invalid tag for repo %s: %s", repoURL, err)
				}

				if projectVersion.Patch() > 0 && projectVersion.Prerelease() != "" {
					return fmt.Errorf("%s tag %s is invalid because it is a patch release but includes a pre-release component", repoURL, *project.Tag)
				}

				if projectVersion.Major() != 0 || projectVersion.Minor() != 0 {
					if v.Prerelease() != projectVersion.Prerelease() {
						return fmt.Errorf("repo %s uses different prerelease version %s compared to product release number %s", repoURL, *project.Tag, r.Release)
					}
				}
			}
		}
	}
	return nil
}

// StrictParseVersion behaves as semver.StrictNewVersion, with as sole exception
// that it allows versions with a preceding "v" (i.e. v1.2.3).
// Ensure new releases are FluxCD compatible.
// xref: https://github.com/fluxcd/pkg/blob/main/version/version.go#L25-L33
func StrictParseVersion(v string) (*semver.Version, error) {
	vLessV := strings.TrimPrefix(v, "v")
	if _, err := semver.StrictNewVersion(vLessV); err != nil {
		return nil, errors.Wrapf(err, "invalid version %s", v)
	}
	return semver.NewVersion(v)
}

/*
- Only one pr per published_chart repo
- Different chart repos can have different prs
*/

type ReplyType string

const (
	OkToRelease ReplyType = "/ok-to-release"
	Done        ReplyType = "/done"

	Tagged       ReplyType = "/tagged"
	Go           ReplyType = "/go"
	ReadyToTag   ReplyType = "/ready-to-tag"
	CherryPicked ReplyType = "/cherry-picked"
	PR           ReplyType = "/pr"

	Chart          ReplyType = "/chart"
	ChartPublished ReplyType = "/chart-published"

	KrewManifest          ReplyType = "/krew-manifest"
	KrewManifestPublished ReplyType = "/krew-manifest-published"
)

type Replies map[ReplyType][]Reply

func MergeReplies(replies Replies, elems ...Reply) Replies {
	out := replies
	for idx := range elems {
		out = MergeReply(out, elems[idx])
	}
	return out
}

func MergeReply(replies Replies, r Reply) Replies {
	if replies == nil {
		replies = map[ReplyType][]Reply{}
	}
	rts := replies[r.Type]

	idx := -1
	for i, existing := range rts {
		if existing.Key() == r.Key() {
			idx = i
			break
		}
	}
	if idx > -1 {
		rts = append(rts[:idx], rts[idx+1:]...)
	}
	rts = append(rts, r)
	replies[r.Type] = rts

	return replies
}

func AppendReplyIfMissing(replies Replies, r Reply) (Replies, bool) {
	if replies == nil {
		replies = map[ReplyType][]Reply{}
	}
	rts := replies[r.Type]

	for _, existing := range rts {
		if existing.Key() == r.Key() {
			return replies, false
		}
	}
	replies[r.Type] = append(rts, r)

	return replies, true
}

type Reply struct {
	Type                  ReplyType
	Tagged                *TaggedReplyData
	PR                    *PullRequestReplyData
	ReadyToTag            *ReadyToTagReplyData
	CherryPicked          *CherryPickedReplyData
	Go                    *GoReplyData
	Chart                 *ChartReplyData
	ChartPublished        *ChartPublishedReplyData
	KrewManifest          *KrewManifestReplyData
	KrewManifestPublished *KrewManifestPublishedReplyData
}

type ReplyKey struct {
	Repo string
	B    string
}

func (r Reply) Key() ReplyKey {
	switch r.Type {
	case OkToRelease:
		fallthrough
	case Done:
		return ReplyKey{}
	case Tagged:
		return ReplyKey{Repo: r.Tagged.Repo}
	case PR:
		return ReplyKey{Repo: r.PR.Repo, B: strconv.Itoa(r.PR.Number)}
	case Go:
		return ReplyKey{Repo: r.Go.Repo}
	case ReadyToTag:
		return ReplyKey{Repo: r.ReadyToTag.Repo}
	case CherryPicked:
		return ReplyKey{Repo: r.CherryPicked.Repo, B: r.CherryPicked.Branch}
	case Chart:
		return ReplyKey{Repo: r.Chart.Repo, B: r.Chart.Tag}
	case ChartPublished:
		return ReplyKey{Repo: r.ChartPublished.Repo}
	case KrewManifest:
		return ReplyKey{Repo: r.KrewManifest.Repo}
	case KrewManifestPublished:
		return ReplyKey{Repo: r.KrewManifestPublished.Repo}
	default:
		panic(fmt.Errorf("unknown reply type %s", r.Type))
	}
}

type TaggedReplyData struct {
	Repo string
}

type PullRequestReplyData struct {
	Repo   string
	Number int
}

type ReadyToTagReplyData struct {
	Repo           string
	MergeCommitSHA string
}

type CherryPickedReplyData struct {
	Repo           string
	Branch         string
	MergeCommitSHA string
}

type MergeData struct {
	Repo string
	Ref  string
}

func (d MergeData) String() string {
	return fmt.Sprintf("%s@%s", d.Repo, d.Ref)
}

type GoReplyData struct {
	Repo       string
	ModulePath string
	VCSRoot    string
}

type ChartReplyData struct {
	Repo string
	Tag  string
}

type ChartPublishedReplyData struct {
	Repo string
}

type KrewManifestReplyData struct {
	Repo string
	Tag  string
}

type KrewManifestPublishedReplyData struct {
	Repo string
}

type Commit struct {
	SHA     string
	Subject string
}

type ReleaseChangelog struct {
	Tag     string   `json:"tag"`
	Commits []Commit `json:"commits"`
}

type ProjectChangelog struct {
	URL      string             `json:"url"`
	Releases []ReleaseChangelog `json:"releases"`
}

type Changelog struct {
	ProductLine       string             `json:"product_line"`
	Release           string             `json:"release"`
	ReleaseDate       time.Time          `json:"release_date"`
	ReleaseProjectURL string             `json:"release_project_url"`
	DocsURL           string             `json:"docs_url"`
	KubernetesVersion string             `json:"kubernetes_version,omitempty"`
	Projects          []ProjectChangelog `json:"projects"`
}

func (chlog *Changelog) Sort() {
	sort.Slice(chlog.Projects, func(i, j int) bool { return chlog.Projects[i].URL < chlog.Projects[j].URL })
	for idx, projects := range chlog.Projects {
		sort.Slice(projects.Releases, func(i, j int) bool {
			vi, _ := semver.NewVersion(projects.Releases[i].Tag)
			vj, _ := semver.NewVersion(projects.Releases[j].Tag)
			return semvers.CompareVersions(vi, vj)
		})
		chlog.Projects[idx] = projects
	}
}

type ReleaseSummary struct {
	Release           string
	ReleaseDate       time.Time
	KubernetesVersion string
	ReleaseURL        string
	ChangelogURL      string
	DocsURL           string
}

type ReleaseTable struct {
	ProductLine string
	Releases    []ReleaseSummary
}
