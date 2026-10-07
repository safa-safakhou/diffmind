package ui

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/mohammad-safakhou/diffmind/internal/workspace/archgraph"
	"github.com/mohammad-safakhou/diffmind/internal/workspace/artifacts"
	"github.com/mohammad-safakhou/diffmind/internal/workspace/orchestrator"
	"github.com/mohammad-safakhou/diffmind/internal/workspace/query"
	"github.com/mohammad-safakhou/diffmind/internal/workspace/store"
)

// PR snapshots live outside company graph history: reviewing a branch must not
// replace the project graph or change the registered checkout/analysis state.
type prFlowCapture struct {
	Base       string            `json:"base"`
	Head       string            `json:"head"`
	ContextRun string            `json:"context_run"`
	Before     store.RunManifest `json:"before"`
	After      store.RunManifest `json:"after"`
}

var prRevisionPattern = regexp.MustCompile(`^[a-fA-F0-9]{40,64}$`)

func (s *Server) prCaptureDir(pid, rid, head string) string {
	identity, _ := orchestrator.AnalyzerIdentity(s.analyzerExecutable())
	sum := sha256.Sum256([]byte(pid + "\x00" + rid + "\x00" + head + "\x00" + identity))
	return filepath.Join(s.store.HomeDir(), "pr-reviews", hex.EncodeToString(sum[:]))
}

func (s *Server) loadPRCapture(pid, rid, base, head string) (*prFlowCapture, *ArchGraph, *ArchGraph, error) {
	dir := s.prCaptureDir(pid, rid, head)
	raw, err := os.ReadFile(filepath.Join(dir, "capture.json"))
	if err != nil {
		return nil, nil, nil, err
	}
	var c prFlowCapture
	if err = json.Unmarshal(raw, &c); err != nil {
		return nil, nil, nil, err
	}
	if c.Head != head || (base != "" && c.Base != base) {
		return nil, nil, nil, errors.New("PR revisions changed; prepare flows again")
	}
	read := func(name string) (*ArchGraph, error) {
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return nil, err
		}
		var g ArchGraph
		err = json.Unmarshal(data, &g)
		return &g, err
	}
	before, err := read("before.json")
	if err != nil {
		return nil, nil, nil, err
	}
	after, err := read("after.json")
	return &c, before, after, err
}

func (s *Server) capturedPRFlows(ctx context.Context, pid string, repo store.Repo, out *pullRequestFlows, offset, limit int) (bool, error) {
	c, before, after, err := s.loadPRCapture(pid, repo.ID, out.BaseCommit, out.HeadCommit)
	if err != nil {
		return false, nil
	}
	service := graphServiceForRepo(after, repo)
	if service == "" {
		return false, errors.New("prepared PR service is missing")
	}
	out.ExploreRun = c.ContextRun
	out.Review, err = query.CompareFlowGraphs(ctx, pid, &c.Before, &c.After, before, after, service, offset, limit)
	if err != nil {
		return false, err
	}
	out.Available = true
	out.BeforeRun = c.Before.ID
	out.AfterRun = c.After.ID
	out.Review.Notes = append(out.Review.Notes, "Only the selected repository changes between these isolated snapshots. Company context is held at "+c.ContextRun+"; nearby callers may have older or dirty evidence.", "PR analysis uses tracked source files and deterministic extraction. Untracked local hints and project graph supplements are not applied.")
	return true, nil
}

func (s *Server) preparePRCapture(ctx context.Context, pid string, repo store.Repo, base, head string) error {
	if !prRevisionPattern.MatchString(base) || !prRevisionPattern.MatchString(head) {
		return errors.New("provider must supply complete hexadecimal revisions")
	}
	key := s.prCaptureDir(pid, repo.ID, head)
	s.prPreparationMu.Lock()
	if s.prPreparing == nil {
		s.prPreparing = map[string]bool{}
	}
	if s.prPreparing[key] {
		s.prPreparationMu.Unlock()
		return errors.New("this PR is already being prepared; refresh flows after it finishes")
	}
	s.prPreparing[key] = true
	s.prPreparationMu.Unlock()
	defer func() { s.prPreparationMu.Lock(); delete(s.prPreparing, key); s.prPreparationMu.Unlock() }()
	if _, _, _, err := s.loadPRCapture(pid, repo.ID, base, head); err == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	release, err := s.acquireRepository(ctx, pid)
	if err != nil {
		return err
	}
	defer release()
	contextRun := s.latestCompletedWorkspaceRun(pid)
	if contextRun == nil {
		return errors.New("a completed company graph is needed before preparing PR flows")
	}
	dirs := map[string]string{}
	for _, ref := range contextRun.Repos {
		registered, e := s.store.GetRepo(pid, ref.RepoID)
		if e != nil || registered.Kind == "infra_repo" || ref.DiffMindRunID == "" {
			continue
		}
		if info, ok := artifacts.DiffMindRunByID(s.diffmindRunsDir, ref.DiffMindRunID); ok && artifacts.RunMatchesRepo(info, registered.Name, registered.ID, registered.Path) {
			dirs[registered.Name] = filepath.Join(s.diffmindRunsDir, ref.DiffMindRunID)
		}
	}
	source := firstNonEmpty(repo.Path, repo.ClonePath)
	if source == "" {
		return errors.New("repository local checkout is unavailable")
	}
	temp, err := os.MkdirTemp("", "diffmind-pr-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(temp)
	checkout := filepath.Join(temp, "source")
	// A private clone avoids worktree metadata changes and user branch switches.
	if err = gitCommand(ctx, "", "", "-c", "core.hooksPath=/dev/null", "clone", "--no-hardlinks", "--no-checkout", "--", source, checkout); err != nil {
		return fmt.Errorf("isolated clone: %w", err)
	}
	remote := firstNonEmpty(repo.GitURL, gitOutput(ctx, source, "remote", "get-url", "origin"))
	for _, sha := range []string{base, head} {
		if gitOutput(ctx, checkout, "rev-parse", "--verify", sha+"^{commit}") != "" {
			continue
		}
		if remote == "" {
			return errors.New("PR revision is not local and no remote is configured")
		}
		if err = gitCommand(ctx, checkout, remote, "-c", "core.hooksPath=/dev/null", "fetch", "--no-tags", "--", remote, sha); err != nil {
			return fmt.Errorf("fetch PR revision: %w", err)
		}
	}
	captured := prFlowCapture{Base: base, Head: head, ContextRun: contextRun.ID}
	analyze := func(sha, side string) (*ArchGraph, store.RunManifest, error) {
		if e := gitCommand(ctx, checkout, "", "-c", "core.hooksPath=/dev/null", "checkout", "--detach", sha); e != nil {
			return nil, store.RunManifest{}, e
		}
		if got := gitOutput(ctx, checkout, "rev-parse", "HEAD"); got != sha {
			return nil, store.RunManifest{}, errors.New("isolated checkout revision does not match PR")
		}
		output := filepath.Join(temp, side)
		if e := orchestrator.RunDiffMindContext(ctx, s.analyzerExecutable(), checkout, orchestrator.DiffMindRunOptions{OutDir: output, Workers: 2}, s.log); e != nil {
			return nil, store.RunManifest{}, e
		}
		runs, e := artifacts.DiscoverDiffMindRuns(output)
		if e != nil || len(runs) == 0 {
			return nil, store.RunManifest{}, errors.New("analyzer produced no PR snapshot")
		}
		// Preserve registered service identity; only its artifact directory changes.
		dirs[repo.Name] = filepath.Join(output, runs[0].RunID)
		id := "pr-" + side + "-" + sha[:12]
		graph := archgraph.Build(id, dirs)
		root := graphServiceForRepo(graph, repo)
		node := graphService(graph, root)
		if node == nil || node.AnalysisStatus == nil || node.AnalysisStatus.Dirty || node.AnalysisStatus.AnalyzedRevision != sha {
			return nil, store.RunManifest{}, errors.New("PR snapshot is not a clean analysis of the requested revision")
		}
		repos, _ := s.store.ListRepos(pid)
		for _, svc := range graph.Services {
			for _, r := range repos {
				if svc.Name == r.Name {
					svc.RepoID = r.ID
					svc.Team = r.Team
				}
			}
		}
		now := time.Now().UTC()
		return graph, store.RunManifest{ID: id, ProjectID: pid, Status: store.RunCompleted, StartedAt: now, FinishedAt: now, ServiceCount: len(graph.Services), EdgeCount: len(graph.Edges)}, nil
	}
	before, m, err := analyze(base, "before")
	if err != nil {
		return err
	}
	captured.Before = m
	after, m, err := analyze(head, "after")
	if err != nil {
		return err
	}
	captured.After = m
	if err = os.MkdirAll(filepath.Dir(key), 0700); err != nil {
		return err
	}
	staging, err := os.MkdirTemp(filepath.Dir(key), "capture-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(staging)
	for name, value := range map[string]any{"before.json": before, "after.json": after, "capture.json": captured} {
		data, e := json.Marshal(value)
		if e != nil {
			return e
		}
		if e = os.WriteFile(filepath.Join(staging, name), data, 0600); e != nil {
			return e
		}
	}
	// Existing captures of an older merge-base are replaced only after both sides succeed.
	if err = os.RemoveAll(key); err != nil {
		return err
	}
	return os.Rename(staging, key)
}
