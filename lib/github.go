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
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/appscodelabs/release-automaton/api"

	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
	"k8s.io/apimachinery/pkg/util/sets"
)

const (
	defaultSecondaryRetryDelay = time.Minute
	defaultServerErrorDelay    = 5 * time.Second
	maxSecondaryRetryDelay     = 15 * time.Minute
	maxRateLimitRetryAttempts  = 8
)

func NewGitHubClient() (*github.Client, error) {
	token, found := os.LookupEnv(api.GitHubTokenKey)
	if !found {
		return nil, fmt.Errorf("%s env var is not set", api.GitHubTokenKey)
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.TODO(), ts)

	baseTransport := tc.Transport
	if baseTransport == nil {
		baseTransport = http.DefaultTransport
	}
	tc.Transport = &rateLimitTransport{base: baseTransport}

	return github.NewClient(tc), nil
}

func ListLabelsByIssue(ctx context.Context, gh *github.Client, owner, repo string, number int) (sets.Set[string], error) {
	opt := &github.ListOptions{
		PerPage: 100,
	}

	result := sets.New[string]()
	for {
		labels, resp, err := gh.Issues.ListLabelsByIssue(ctx, owner, repo, number, opt)
		if err != nil {
			if ge, ok := err.(*github.ErrorResponse); ok && ge.Response.StatusCode == http.StatusNotFound {
				log.Println(err)
			} else {
				return nil, err
			}
		}

		for _, entry := range labels {
			result.Insert(entry.GetName())
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return result, nil
}

func ListReleases(ctx context.Context, gh *github.Client, owner, repo string) ([]*github.RepositoryRelease, error) {
	opt := &github.ListOptions{
		PerPage: 100,
	}

	var result []*github.RepositoryRelease
	for {
		releases, resp, err := gh.Repositories.ListReleases(ctx, owner, repo, opt)
		if err != nil {
			if ge, ok := err.(*github.ErrorResponse); ok && ge.Response.StatusCode == http.StatusNotFound {
				log.Println(err)
			} else {
				return nil, err
			}
		}

		result = append(result, releases...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return result, nil
}

func ListReviews(ctx context.Context, gh *github.Client, owner, repo string, number int) ([]*github.PullRequestReview, error) {
	opt := &github.ListOptions{
		PerPage: 100,
	}

	var result []*github.PullRequestReview
	for {
		reviews, resp, err := gh.PullRequests.ListReviews(ctx, owner, repo, number, opt)
		if err != nil {
			if ge, ok := err.(*github.ErrorResponse); ok && ge.Response.StatusCode == http.StatusNotFound {
				log.Println(err)
			} else {
				return nil, err
			}
		}

		result = append(result, reviews...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return result, nil
}

func ListPullRequestComment(ctx context.Context, gh *github.Client, owner, repo string, number int) ([]*github.PullRequestComment, error) {
	opt := &github.PullRequestListCommentsOptions{
		Sort:      "created",
		Direction: "asc",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	var result []*github.PullRequestComment
	for {
		comments, resp, err := gh.PullRequests.ListComments(ctx, owner, repo, number, opt)
		if err != nil {
			if ge, ok := err.(*github.ErrorResponse); ok && ge.Response.StatusCode == http.StatusNotFound {
				log.Println(err)
			} else {
				return nil, err
			}
		}

		result = append(result, comments...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return result, nil
}

func ListComments(ctx context.Context, gh *github.Client, owner, repo string, number int) ([]*github.IssueComment, error) {
	opt := &github.IssueListCommentsOptions{
		Sort:      github.String("created"),
		Direction: github.String("asc"),
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	var result []*github.IssueComment
	for {
		comments, resp, err := gh.Issues.ListComments(ctx, owner, repo, number, opt)
		if err != nil {
			if ge, ok := err.(*github.ErrorResponse); ok && ge.Response.StatusCode == http.StatusNotFound {
				log.Println(err)
			} else {
				return nil, err
			}
		}

		result = append(result, comments...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return result, nil
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

	head := req.GetHead()
	if !strings.ContainsRune(head, ':') {
		head = owner + ":" + req.GetHead()
	}
	prs, _, err := gh.PullRequests.List(context.TODO(), owner, repo, &github.PullRequestListOptions{
		State: "open",
		Head:  head,
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
		// "A pull request already Exists" error should NEVER happen since we already checked for existence
		if err != nil {
			return nil, err
		}
		//if e2, ok := err.(*github.ErrorResponse); ok {
		//	var matched bool
		//	for _, entry := range e2.Errors {
		//		if strings.HasPrefix(entry.Message, "A pull request already Exists") {
		//			matched = true
		//			break
		//		}
		//	}
		//	if !matched {
		//		return nil, err
		//	}
		//	// else ignore error because pr already Exists
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

func ClosePR(gh *github.Client, owner string, repo string, head, base string) (*github.PullRequest, error) {
	if !strings.ContainsRune(head, ':') {
		head = owner + ":" + head
	}
	prs, _, err := gh.PullRequests.List(context.TODO(), owner, repo, &github.PullRequestListOptions{
		State: "open",
		Head:  head,
		Base:  base,
		ListOptions: github.ListOptions{
			PerPage: 1,
		},
	})
	if err != nil {
		return nil, err
	}

	if len(prs) == 1 {
		pr, _, err := gh.PullRequests.Edit(context.TODO(), owner, repo, prs[0].GetNumber(), &github.PullRequest{
			State: github.String("closed"),
		})
		return pr, err
	}

	return nil, fmt.Errorf("pr not found")
}

func LabelPR(gh *github.Client, owner string, repo, head, base string, labels ...string) error {
	labelSet := sets.NewString(labels...)
	var result *github.PullRequest

	if !strings.ContainsRune(head, ':') {
		head = owner + ":" + head
	}
	prs, _, err := gh.PullRequests.List(context.TODO(), owner, repo, &github.PullRequestListOptions{
		State: "open",
		Head:  head,
		Base:  base,
		ListOptions: github.ListOptions{
			PerPage: 1,
		},
	})
	if err != nil {
		return err
	}
	if len(prs) == 0 {
		return fmt.Errorf("no open pr found")
	}

	result = prs[0]
	for _, label := range result.Labels {
		labelSet.Delete(label.GetName())
	}
	if labelSet.Len() > 0 {
		_, _, err := gh.Issues.AddLabelsToIssue(context.TODO(), owner, repo, result.GetNumber(), labelSet.UnsortedList())
		if err != nil {
			return err
		}
	}
	return nil
}

func RemoveLabel(gh *github.Client, owner string, repo string, number int, label string) error {
	_, err := gh.Issues.RemoveLabelForIssue(context.TODO(), owner, repo, number, label)
	if ge, ok := err.(*github.ErrorResponse); ok {
		if ge.Response.StatusCode == http.StatusNotFound {
			return nil
		}
	}
	return err
}

func ListTags2(ctx context.Context, gh *github.Client, owner, repo string) ([]*github.RepositoryTag, error) {
	opt := &github.ListOptions{
		PerPage: 100,
	}

	var result []*github.RepositoryTag
	for {
		reviews, resp, err := gh.Repositories.ListTags(ctx, owner, repo, opt)
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

type rateLimitTransport struct {
	base http.RoundTripper
}

func (t *rateLimitTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.GetBody == nil {
		if f, ok := req.Body.(*os.File); ok {
			name := f.Name()
			req.GetBody = func() (io.ReadCloser, error) {
				return os.Open(name)
			}
		}
	}

	canRetryBody := req.Body == nil || req.Body == http.NoBody || req.GetBody != nil

	for attempt := 0; ; attempt++ {
		currReq := req
		if attempt > 0 {
			currReq = req.Clone(req.Context())
			if req.GetBody != nil {
				body, err := req.GetBody()
				if err != nil {
					return nil, err
				}
				currReq.Body = body
			}
		}

		resp, err := t.base.RoundTrip(currReq)
		if err != nil {
			return nil, err
		}
		if resp == nil {
			return resp, nil
		}

		kind := classifyResponse(resp)
		if kind == kindNone {
			return resp, nil
		}
		if attempt >= maxRateLimitRetryAttempts {
			log.Printf("GitHub API still failing after %d retries (status=%d); giving up", attempt, resp.StatusCode)
			return resp, nil
		}
		if !canRetryBody {
			return resp, nil
		}

		delay, source := retryDelay(resp, kind, attempt)
		logRetry(resp, kind, delay, source)

		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()

		timer := time.NewTimer(delay)
		select {
		case <-req.Context().Done():
			if !timer.Stop() {
				<-timer.C
			}
			return nil, req.Context().Err()
		case <-timer.C:
		}
	}
}

type responseKind int

const (
	kindNone responseKind = iota
	kindPrimaryRateLimit
	kindSecondaryRateLimit
	kindServerError
)

// classifyResponse decides whether a response is retryable per
// https://docs.github.com/en/rest/using-the-rest-api/rate-limits-for-the-rest-api.
//
// A 403/429 is only treated as a rate limit when GitHub provides a Retry-After
// header or X-RateLimit-Remaining=0. Bare 403s are permission errors and
// retrying them is wasteful.
func classifyResponse(resp *http.Response) responseKind {
	switch resp.StatusCode {
	case http.StatusForbidden, http.StatusTooManyRequests:
		hasRetryAfter := resp.Header.Get("Retry-After") != ""
		exhausted := resp.Header.Get("X-RateLimit-Remaining") == "0"
		switch {
		case exhausted:
			return kindPrimaryRateLimit
		case hasRetryAfter:
			return kindSecondaryRateLimit
		default:
			return kindNone
		}
	case http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return kindServerError
	default:
		return kindNone
	}
}

func retryDelay(resp *http.Response, kind responseKind, attempt int) (time.Duration, string) {
	if d := parseRetryAfter(resp.Header.Get("Retry-After")); d > 0 {
		return d, "Retry-After header"
	}

	switch kind {
	case kindPrimaryRateLimit:
		if reset := parseUnixTime(resp.Header.Get("X-RateLimit-Reset")); !reset.IsZero() {
			return max(time.Until(reset)+time.Second, time.Second), "X-RateLimit-Reset header"
		}
		return secondaryBackoff(attempt), "rate-limit backoff (no reset header)"
	case kindSecondaryRateLimit:
		return secondaryBackoff(attempt), "secondary rate-limit backoff"
	case kindServerError:
		return serverErrorBackoff(attempt), "exponential backoff after server error"
	}
	return time.Second, ""
}

// secondaryBackoff implements the GitHub-recommended "at least 1 minute,
// increase exponentially" backoff for secondary rate limits.
func secondaryBackoff(attempt int) time.Duration {
	d := time.Duration(float64(defaultSecondaryRetryDelay) * math.Pow(2, float64(attempt)))
	return max(min(d, maxSecondaryRetryDelay), defaultSecondaryRetryDelay)
}

func serverErrorBackoff(attempt int) time.Duration {
	d := time.Duration(float64(defaultServerErrorDelay) * math.Pow(2, float64(attempt)))
	return max(min(d, maxSecondaryRetryDelay), time.Second)
}

func logRetry(resp *http.Response, kind responseKind, delay time.Duration, source string) {
	rounded := delay.Round(time.Second)
	resource := resp.Header.Get("X-RateLimit-Resource")
	suffix := ""
	if resource != "" {
		suffix = " resource=" + resource
	}

	switch kind {
	case kindPrimaryRateLimit:
		log.Printf("GitHub API primary rate limit hit (%d%s); waiting %s before retry (%s)", resp.StatusCode, suffix, rounded, source)
	case kindSecondaryRateLimit:
		log.Printf("GitHub API secondary rate limit hit (%d%s); waiting %s before retry (%s)", resp.StatusCode, suffix, rounded, source)
	case kindServerError:
		log.Printf("GitHub API server error (%d); waiting %s before retry (%s)", resp.StatusCode, rounded, source)
	}
}

func parseRetryAfter(value string) time.Duration {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0
	}
	if seconds, err := strconv.Atoi(value); err == nil {
		if seconds < 1 {
			seconds = 1
		}
		return time.Duration(seconds) * time.Second
	}
	if ts, err := http.ParseTime(value); err == nil {
		return max(time.Until(ts), time.Second)
	}
	return 0
}

func parseUnixTime(value string) time.Time {
	if value == "" {
		return time.Time{}
	}
	sec, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return time.Time{}
	}
	return time.Unix(sec, 0).UTC()
}
