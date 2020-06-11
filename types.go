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

import "fmt"

type Project struct {
	Tag      *string           `json:"tag,omitempty"`
	Tags     map[string]string `json:"tags,omitempty"` // tag-> branch
	Commands []string          `json:"cmds,omitempty"`
}

type IndependentProjects map[string]Project

type Release struct {
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

	ChartMerged    ReplyType = "/chart-merged"
	ChartPublished ReplyType = "/chart-published"
)

type Reply struct {
	Type         ReplyType
	Tagged       *TaggedReplyData
	PR           *PullRequestReplyData
	ReadyToTag   *ReadyToTagReplyData
	CherryPicked *CherryPickedReplyData
	Go           *GoReplyData
	ChartMerged  *ChartMergedReplyData
}

func (r Reply) Repo() string {
	switch r.Type {
	case OkToRelease:
		return ""
	case Tagged:
		return r.Tagged.Repo
	case PR:
		return r.PR.Repo
	case Go:
		return r.Go.Repo
	case ReadyToTag:
		return r.ReadyToTag.Repo
	case CherryPicked:
		return r.CherryPicked.Repo
	case ChartMerged:
		return r.ChartMerged.Repo
	case ChartPublished:
		return ""
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
}

type ChartMergedReplyData struct {
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
	Release  string             `json:"release"`
	Projects []ProjectChangelog `json:"projects"`
}
