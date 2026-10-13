package ui

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mohammad-safakhou/diffmind/internal/workspace/artifacts"
	"github.com/mohammad-safakhou/diffmind/internal/workspace/orchestrator"
	"github.com/mohammad-safakhou/diffmind/internal/workspace/store"
)

func TestBackgroundGitReadsPreserveIndex(t *testing.T) {
	for _, operation := range []string{"inspection", "fingerprint"} {
		t.Run(operation, func(t *testing.T) {
			root := cleanTestRepo(t)
			indexPath := filepath.Join(root, ".git", "index")
			before, err := os.ReadFile(indexPath)
			if err != nil {
				t.Fatal(err)
			}
			oldTime := time.Unix(946684800, 0)
			if err := os.Chtimes(filepath.Join(root, "main.go"), oldTime, oldTime); err != nil {
				t.Fatal(err)
			}
			if operation == "inspection" {
				if info := inspectLocalGit(t.Context(), root, ""); info.Dirty || info.Head == "" {
					t.Fatalf("clean checkout inspection: %+v", info)
				}
			} else {
				s := newAuthTestServer(t)
				p, err := s.store.CreateProject(store.Project{Name: "Index preservation"})
				if err != nil {
					t.Fatal(err)
				}
				if fingerprint, err := s.analysisFingerprint(t.Context(), p.ID, store.Repo{Path: root}, orchestrator.DiffMindRunOptions{}, "engine"); err != nil || fingerprint == "" {
					t.Fatalf("clean checkout fingerprint: %q %v", fingerprint, err)
				}
			}
			after, err := os.ReadFile(indexPath)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(before, after) {
				t.Fatal("background Git read rewrote the index")
			}
		})
	}
}

func TestLocalDiffMindFreshnessTracksCheckoutAndDirtyState(t *testing.T) {
	dir := t.TempDir()
	runGitForTest(t, dir, "init")
	runGitForTest(t, dir, "config", "user.email", "test@example.com")
	runGitForTest(t, dir, "config", "user.name", "Test")
	path := filepath.Join(dir, "app.py")
	if err := os.WriteFile(path, []byte("print('one')\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGitForTest(t, dir, "add", "app.py")
	runGitForTest(t, dir, "commit", "-m", "one")
	head := gitOutput(t.Context(), dir, "rev-parse", "HEAD")
	repo := store.Repo{Path: dir, SourceType: "local"}
	latest := &artifacts.DiffMindRunInfo{RepoGitSHA: head}
	if got := diffmindFreshness(repo, latest); got != "fresh" {
		t.Fatalf("clean checkout freshness = %q", got)
	}
	if err := os.WriteFile(path, []byte("print('dirty')\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := diffmindFreshness(repo, latest); got != "dirty" {
		t.Fatalf("dirty checkout freshness = %q", got)
	}
	runGitForTest(t, dir, "add", "app.py")
	runGitForTest(t, dir, "commit", "-m", "two")
	if got := diffmindFreshness(repo, latest); got != "stale" {
		t.Fatalf("advanced checkout freshness = %q", got)
	}
}
