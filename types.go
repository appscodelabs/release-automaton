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
	"fmt"
	"strconv"
	"strings"
)

type Project struct {
	Tag      *string           `json:"tag,omitempty"`
	Tags     map[string]string `json:"tags,omitempty"` // tag-> branch
	Commands []string          `json:"cmds,omitempty"`
}

type IndependentProjects map[string]Project

type Release struct {
	ProductLine string `json:"productLine"`
	// These projects can be released in sequence
	Projects []IndependentProjects `json:"projects"`
}

type ReplyType string

const (
	OkToRelease  ReplyType = "/ok-to-release"
	Tagged       ReplyType = "/tagged"
	Go           ReplyType = "/go"
	ReadyToTag   ReplyType = "/ready-to-tag"
	CherryPicked ReplyType = "/cherry-picked"
	PR           ReplyType = "/pr"

	Chart          ReplyType = "/chart-merged"
	ChartPublished ReplyType = "/chart-published"
)

type Replies map[ReplyType][]Reply

func MergeReplies(replies Replies, elems ...Reply) Replies {
	var out = replies
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

// https://www.calhoun.io/when-nil-isnt-equal-to-nil/

//func (replies Replies) Find(repoURL string) (*Reply, bool) {
//	for i := range replies {
//		if replies[i].Repo() == repoURL {
//			return &replies[i], true
//		}
//	}
//	return nil, false
//}
//
//func (replies *Replies) Append(r Reply) {
//	*replies = append(*replies, r)
//}
//
//func (replies *Replies) AppendIfMissing(r Reply) bool {
//	for _, entry := range *replies {
//		if entry.Repo() == r.Repo() {
//			return false
//		}
//	}
//	*replies = append(*replies, r)
//	return true
//}

type Reply struct {
	Type         ReplyType
	Tagged       *TaggedReplyData
	PR           *PullRequestReplyData
	ReadyToTag   *ReadyToTagReplyData
	CherryPicked *CherryPickedReplyData
	Go           *GoReplyData
	Chart        *ChartReplyData
}

type ReplyKey struct {
	Repo string
	B    string
}

func (r Reply) Key() ReplyKey {
	switch r.Type {
	case OkToRelease:
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
		return ReplyKey{}
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
	Repo   string
	Branch string
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
	Release  string             `json:"release"`
	Projects []ProjectChangelog `json:"projects"`
}

func ParseReply(s string) *Reply {
	fields := strings.Fields(s)
	if len(fields) == 0 {
		return nil
	}

	rt := ReplyType(fields[0])
	params := fields[1:]

	switch rt {
	case OkToRelease:
		fallthrough
	case ChartPublished:
		if len(params) > 0 {
			panic(fmt.Errorf("unsupported parameters with reply %s", s))
		}
		return &Reply{Type: rt}
	case Tagged:
		if len(params) != 1 {
			panic(fmt.Errorf("unsupported parameters with reply %s", s))
		}
		return &Reply{Type: rt, Tagged: &TaggedReplyData{
			Repo: params[0],
		}}
	case Go:
		if len(params) != 2 && len(params) != 3 {
			panic(fmt.Errorf("unsupported parameters with reply %s", s))
		}
		data := &GoReplyData{
			Repo:       params[0],
			ModulePath: params[1],
		}
		if len(params) == 3 {
			data.VCSRoot = params[2]
		}
		return &Reply{Type: rt, Go: data}
	case PR:
		if len(params) != 1 {
			panic(fmt.Errorf("unsupported parameters with reply %s", s))
		}
		owner, repo, prNumber := ParsePullRequestURL(params[0])
		return &Reply{Type: rt, PR: &PullRequestReplyData{
			Repo:   fmt.Sprintf("github.com/%s/%s", owner, repo),
			Number: prNumber,
		}}
	case ReadyToTag:
		if len(params) != 2 {
			panic(fmt.Errorf("unsupported parameters with reply %s", s))
		}
		return &Reply{Type: rt, ReadyToTag: &ReadyToTagReplyData{
			Repo:           params[0],
			MergeCommitSHA: params[1],
		}}
	case CherryPicked:
		if len(params) != 3 {
			panic(fmt.Errorf("unsupported parameters with reply %s", s))
		}
		return &Reply{Type: rt, CherryPicked: &CherryPickedReplyData{
			Repo:           params[0],
			Branch:         params[1],
			MergeCommitSHA: params[2],
		}}
	case Chart:
		if len(params) != 2 {
			panic(fmt.Errorf("unsupported parameters with reply %s", s))
		}
		return &Reply{Type: rt, Chart: &ChartReplyData{
			Repo: params[0],
			Tag:  params[1],
		}}
	default:
		fmt.Printf("unknown reply type found in %s\n", s)
		return nil
	}
}
