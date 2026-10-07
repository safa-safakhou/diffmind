package ui

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/mohammad-safakhou/diffmind/internal/workspace/query"
	"github.com/mohammad-safakhou/diffmind/internal/workspace/store"
)

type pullRequestFlows struct {
	ExploreRun string                  `json:"explore_run,omitempty"`
	Available  bool                    `json:"available"`
	BaseCommit string                  `json:"base_commit"`
	HeadCommit string                  `json:"head_commit"`
	BeforeRun  string                  `json:"before_run,omitempty"`
	AfterRun   string                  `json:"after_run,omitempty"`
	NextAction string                  `json:"next_action,omitempty"`
	Review     *query.FlowReviewResult `json:"review,omitempty"`
}

func (s *Server) handlePullRequestFlows(w http.ResponseWriter, r *http.Request) {
	pid := r.PathValue("pid")
	number, err := strconv.Atoi(r.PathValue("number"))
	if err != nil || number < 1 {
		writeErr(w, 400, errors.New("positive PR number required"))
		return
	}
	offset, err := queryInteger(r, "offset", 0)
	if err != nil {
		writeV1Result(w, nil, err)
		return
	}
	limit, err := queryInteger(r, "limit", 10)
	if err != nil {
		writeV1Result(w, nil, err)
		return
	}
	if offset < 0 || limit < 1 || limit > 50 {
		writeErr(w, 400, errors.New("offset must be nonnegative and limit between 1 and 50"))
		return
	}
	repo, err := s.store.GetRepo(pid, r.PathValue("repo_id"))
	if err != nil {
		s.writeStoreErr(w, err)
		return
	}
	endpoint, state, message := githubRepositoryEndpoint(r.Context(), *repo)
	if state != "ready" {
		writeErr(w, 400, errors.New(message))
		return
	}
	client := githubHTTPClient(30 * time.Second)
	token := githubAPIToken(r.Context(), endpoint)
	var pull githubPull
	if err := githubJSON(r.Context(), client, token, fmt.Sprintf("%s/pulls/%d", endpoint, number), &pull); err != nil {
		writeErr(w, 502, err)
		return
	}
	if pull.Base.SHA == "" || pull.Head.SHA == "" {
		writeErr(w, 502, errors.New("provider did not return PR base and head revisions"))
		return
	}
	var comparison struct {
		MergeBase struct {
			SHA string `json:"sha"`
		} `json:"merge_base_commit"`
	}
	// Compare immutable SHAs, not branch names that may move during the read.
	target := endpoint + "/compare/" + url.PathEscape(pull.Base.SHA) + "..." + url.PathEscape(pull.Head.SHA) + "?per_page=1"
	if err := githubJSON(r.Context(), client, token, target, &comparison); err != nil {
		writeErr(w, 502, err)
		return
	}
	if comparison.MergeBase.SHA == "" {
		writeErr(w, 502, errors.New("provider did not return a merge-base revision"))
		return
	}
	out := pullRequestFlows{BaseCommit: comparison.MergeBase.SHA, HeadCommit: pull.Head.SHA}
	if r.Method == http.MethodPost {
		if err := s.preparePRCapture(r.Context(), pid, *repo, out.BaseCommit, out.HeadCommit); err != nil {
			writeErr(w, 422, err)
			return
		}
	}
	if r.URL.Query().Get("from") == "" && r.URL.Query().Get("to") == "" {
		ready, err := s.capturedPRFlows(r.Context(), pid, *repo, &out, offset, limit)
		if err != nil {
			writeV1Result(w, nil, err)
			return
		}
		if ready {
			writeJSON(w, 200, out)
			return
		}
	}

	before, after, service, err := s.findPRFlowRuns(r.Context(), pid, *repo, out.BaseCommit, out.HeadCommit, r.URL.Query().Get("from"), r.URL.Query().Get("to"))
	if err != nil {
		writeV1Result(w, nil, err)
		return
	}
	out.BeforeRun, out.AfterRun, out.ExploreRun = before, after, after
	if before == "" || after == "" {
		out.NextAction = "Prepare PR flows to analyze the merge-base and PR head in isolated checkouts. This can take a few minutes and keeps company context fixed for both sides."
		writeJSON(w, 200, out)
		return
	}
	out.Review, err = s.query.CompareFlows(r.Context(), pid, before, after, service, offset, limit)
	if err != nil {
		writeV1Result(w, nil, err)
		return
	}
	out.Available = true
	out.Review.Notes = append(out.Review.Notes, "The selected repository matches the PR merge-base and head commits. Other repositories retain their recorded snapshot inputs; differences there may have independent causes.")
	writeJSON(w, 200, out)
}

func (s *Server) findPRFlowRuns(ctx context.Context, pid string, repo store.Repo, base, head, requestedFrom, requestedTo string) (string, string, string, error) {
	from, to, service := "", "", ""
	check := func(id, want string) (string, error) {
		if err := ctx.Err(); err != nil {
			return "", err
		}
		g, err := s.persistedArchGraphForRunFast(pid, id, &http.Request{URL: &url.URL{}})
		if err != nil {
			_, g, err = s.query.Load(pid, id)
		}
		if err != nil {
			return "", err
		}
		root := graphServiceForRepo(g, repo)
		node := graphService(g, root)
		if node == nil || node.AnalysisStatus == nil {
			return "", nil
		}
		status := node.AnalysisStatus
		if status.Dirty || status.State != "analyzed_clean" || !strings.EqualFold(status.AnalyzedRevision, want) {
			return "", nil
		}
		return root, nil
	}
	if requestedFrom != "" {
		root, err := check(requestedFrom, base)
		if err != nil {
			return "", "", "", err
		}
		if root == "" {
			return "", "", "", fmt.Errorf("before snapshot must contain a clean analysis of merge-base %s", base)
		}
		from, service = requestedFrom, root
	}
	if requestedTo != "" {
		root, err := check(requestedTo, head)
		if err != nil {
			return "", "", "", err
		}
		if root == "" {
			return "", "", "", fmt.Errorf("after snapshot must contain a clean analysis of PR head %s", head)
		}
		to = requestedTo
		if service != "" && service != root {
			return "", "", "", fmt.Errorf("service identity changed between snapshots; compare the renamed services separately")
		}
		service = root
	}
	// Bound automatic discovery; callers can pin older runs explicitly.
	for offset := 0; offset < 500 && (from == "" || to == ""); offset += 100 {
		if err := ctx.Err(); err != nil {
			return "", "", "", err
		}
		runs, err := s.query.GraphRuns(pid, offset, 100)
		if err != nil {
			return "", "", "", err
		}
		for _, run := range runs.Runs {
			if !run.GraphAvailable {
				continue
			}
			if from == "" {
				root, err := check(run.ID, base)
				if err == nil && root != "" && (service == "" || root == service) {
					from, service = run.ID, root
				}
			}
			if to == "" {
				root, err := check(run.ID, head)
				if err == nil && root != "" && (service == "" || root == service) {
					to, service = run.ID, root
				}
			}
			if from != "" && to != "" {
				break
			}
		}
		if runs.NextOffset == nil {
			break
		}
	}
	return from, to, service, ctx.Err()
}
