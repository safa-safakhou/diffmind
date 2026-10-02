package query

import (
	"github.com/mohammad-safakhou/diffmind/internal/workspace/artifacts"
	"github.com/mohammad-safakhou/diffmind/internal/workspace/store"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRepositoryFreshnessReadsLocalCheckoutAndFailsUnknown(t *testing.T) {
	root := t.TempDir()
	git := func(args ...string) string {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git: %s %v", output, err)
		}
		return strings.TrimSpace(string(output))
	}
	git("init", "-q")
	file := filepath.Join(root, "source.go")
	if err := os.WriteFile(file, []byte("package source\n"), 0600); err != nil {
		t.Fatal(err)
	}
	git("add", ".")
	git("-c", "user.name=Test", "-c", "user.email=test@example.test", "commit", "-qm", "source")
	head := git("rev-parse", "HEAD")
	repo := store.Repo{Path: root, SourceType: "local", HeadSHA: head, DiffMindFreshness: "fresh"}
	latest := &artifacts.DiffMindRunInfo{RepoGitSHA: head}
	if got := RepositoryFreshness(repo, latest); got != "fresh" {
		t.Fatalf("fresh=%s", got)
	}
	os.WriteFile(file, []byte("package changed\n"), 0600)
	if got := RepositoryFreshness(repo, latest); got != "dirty" {
		t.Fatalf("dirty=%s", got)
	}
	git("add", ".")
	git("-c", "user.name=Test", "-c", "user.email=test@example.test", "commit", "-qm", "changed")
	if got := RepositoryFreshness(repo, latest); got != "stale" {
		t.Fatalf("stale=%s", got)
	}
	repo.Path = filepath.Join(root, "missing")
	if got := RepositoryFreshness(repo, latest); got != "unknown" {
		t.Fatalf("missing checkout reused stored SHA: %s", got)
	}
	if got := RepositoryFreshness(repo, nil); got != "unknown" {
		t.Fatalf("missing analysis=%s", got)
	}
}

func TestUnavailableAnalysisArtifactsKeepSavedGraphQueryable(t *testing.T) {
	q, pid, runID := testQueryService(t)
	_, err := q.store.CreateRepo(pid, store.Repo{Name: "catalog", Path: t.TempDir(), SourceType: "local", DiffMindFreshness: "fresh"})
	if err != nil {
		t.Fatal(err)
	}
	blocked := filepath.Join(t.TempDir(), "not-a-directory")
	if err = os.WriteFile(blocked, []byte("blocked"), 0600); err != nil {
		t.Fatal(err)
	}
	q.WithRunsDir(blocked)
	graphPath := filepath.Join(q.store.RunDir(pid, runID), "graph.json")
	before, _ := os.ReadFile(graphPath)
	summary, err := q.Summary(pid, runID)
	if err != nil {
		t.Fatal(err)
	}
	if summary.Evidence.GraphState != "saved" || summary.Evidence.AnalysisFreshness != "unknown" || summary.Evidence.Unknown != 1 {
		t.Fatalf("evidence=%+v", summary.Evidence)
	}
	_, graph, err := q.Load(pid, runID)
	if err != nil {
		t.Fatal(err)
	}
	if graph.Services[0].AnalysisStatus.CurrentState != "unknown" {
		t.Fatalf("saved service falsely fresh: %+v", graph.Services[0])
	}
	after, _ := os.ReadFile(graphPath)
	if string(before) != string(after) {
		t.Fatal("read-only freshness changed saved graph")
	}
}
