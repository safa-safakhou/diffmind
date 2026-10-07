package ui

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/mohammad-safakhou/diffmind/internal/workspace/store"
)

func TestPRPreparationPreservesDirtyCheckoutAndCompanyHistory(t *testing.T) {
	source := t.TempDir()
	git := func(args ...string) string {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = source
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git: %s %v", out, err)
		}
		return string(out)
	}
	git("init", "-q")
	git("config", "user.email", "fixture@example.com")
	git("config", "user.name", "Fixture")
	os.WriteFile(filepath.Join(source, "file.txt"), []byte("before"), 0600)
	git("add", ".")
	git("commit", "-qm", "before")
	base := gitOutput(context.Background(), source, "rev-parse", "HEAD")
	os.WriteFile(filepath.Join(source, "file.txt"), []byte("after"), 0600)
	git("commit", "-qam", "after")
	head := gitOutput(context.Background(), source, "rev-parse", "HEAD")
	os.WriteFile(filepath.Join(source, "file.txt"), []byte("unsaved"), 0600)
	dirty := git("status", "--porcelain")
	branch := git("symbolic-ref", "HEAD")
	s := newAuthTestServer(t)
	project, _ := s.store.CreateProject(store.Project{Name: "PR preparation"})
	repo, _ := s.store.CreateRepo(project.ID, store.Repo{Name: "orders", Path: source, SourceType: "local", Kind: "service_repo"})
	originalRun, _ := s.store.CreateRun(project.ID, store.RunManifest{Status: store.RunCompleted, StartedAt: time.Now()})
	analyzer := filepath.Join(t.TempDir(), "analyzer")
	script := `#!/usr/bin/python3
import sys,json,pathlib,subprocess
args=sys.argv
repo=args[args.index('--repo')+1];out=pathlib.Path(args[args.index('--out')+1])/'fixture';out.mkdir(parents=True)
sha=subprocess.check_output(['git','-C',repo,'rev-parse','HEAD'],text=True).strip()
(out/'run_manifest.json').write_text(json.dumps({'repo_path':repo,'repo_git_sha':sha,'team':'fixture'}))
(out/'exposures').mkdir()
(out/'exposures'/'http_route.json').write_text(json.dumps([{'id':'route','name':'GET /orders','kind':'http_route','details':{'path':'/orders','method':'GET'}}]))
`
	os.WriteFile(analyzer, []byte(script), 0700)
	s.SetAnalyzerBinary(analyzer)
	if err := s.preparePRCapture(context.Background(), project.ID, *repo, base, head); err != nil {
		t.Fatal(err)
	}
	c, b, a, err := s.loadPRCapture(project.ID, repo.ID, base, head)
	if err != nil {
		t.Fatal(err)
	}
	if c.ContextRun != originalRun.ID || b.Services[0].AnalysisStatus.AnalyzedRevision != base || a.Services[0].AnalysisStatus.AnalyzedRevision != head {
		t.Fatal("snapshot revisions/context do not match")
	}
	if git("status", "--porcelain") != dirty || git("symbolic-ref", "HEAD") != branch {
		t.Fatal("registered checkout changed")
	}
	content, _ := os.ReadFile(filepath.Join(source, "file.txt"))
	if string(content) != "unsaved" {
		t.Fatal("uncommitted file changed")
	}
	runs, _ := s.store.ListRuns(project.ID)
	if len(runs) != 1 || runs[0].ID != originalRun.ID {
		t.Fatal("PR snapshots contaminated company history")
	}
	updated, _ := s.store.GetRepo(project.ID, repo.ID)
	if updated.LastDiffMindRunID != "" {
		t.Fatal("repository current analysis was overwritten")
	}
	if _, _, _, err = s.loadPRCapture(project.ID, repo.ID, head, head); err == nil {
		t.Fatal("wrong merge-base accepted")
	}
	out := pullRequestFlows{BaseCommit: base, HeadCommit: head}
	ready, err := s.capturedPRFlows(context.Background(), project.ID, *repo, &out, 0, 10)
	if err != nil || !ready {
		t.Fatalf("comparison %v %v", ready, err)
	}
	if out.Review.Total != 0 || len(out.Review.Flows) != 0 || out.ExploreRun == "" {
		raw, _ := json.Marshal(out)
		t.Fatalf("unchanged PR flow should be excluded, with service exploration available: %s", raw)
	}
}
