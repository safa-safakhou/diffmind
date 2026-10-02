package ui

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/mohammad-safakhou/diffmind/internal/workspace/store"
)

func TestImportReviewRejectsChangedScopeWithoutMutation(t *testing.T) {
	s := newAuthTestServer(t)
	project, _ := s.store.CreateProject(store.Project{Name: "review"})
	root := t.TempDir()
	mkdirGit(t, filepath.Join(root, "one"))
	req := importReposRequest{Provider: "local", Root: root, DryRun: true}
	preview, err := s.prepareImport(context.Background(), project.ID, req)
	if err != nil {
		t.Fatal(err)
	}
	req.DryRun = false
	req.PreviewDigest = preview.digest
	mkdirGit(t, filepath.Join(root, "two"))
	for _, pipeline := range []bool{false, true} {
		var body []byte
		if pipeline {
			body, _ = json.Marshal(ingestionRequest{Import: &req})
		} else {
			body, _ = json.Marshal(req)
		}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
		r.SetPathValue("pid", project.ID)
		if pipeline {
			s.handleStartIngestion(w, r)
		} else {
			s.handleImportRepos(w, r)
		}
		if w.Code != http.StatusConflict {
			t.Fatalf("pipeline %v: %d %s", pipeline, w.Code, w.Body.String())
		}
		repos, _ := s.store.ListRepos(project.ID)
		if len(repos) != 0 {
			t.Fatalf("changed preview registered %d repos", len(repos))
		}
	}
}

func TestImportReviewTracksFiltersAndAnalysisPaths(t *testing.T) {
	s := newAuthTestServer(t)
	project, _ := s.store.CreateProject(store.Project{Name: "review"})
	root := t.TempDir()
	repo := filepath.Join(root, "one")
	mkdirGit(t, repo)
	req := importReposRequest{Provider: "local", Root: root, DryRun: true}
	preview, err := s.prepareImport(context.Background(), project.ID, req)
	if err != nil {
		t.Fatal(err)
	}
	req.DryRun = false
	req.PreviewDigest = preview.digest
	changed := req
	changed.Include = "one"
	if _, err = s.prepareImport(context.Background(), project.ID, changed); !errors.Is(err, errImportReviewChanged) {
		t.Fatalf("filter change: %v", err)
	}
	if err = os.WriteFile(filepath.Join(repo, "diffmind-configuration.yaml"), []byte("schema: diffmind.config.v1\npaths:\n  exclude: [examples/**]\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err = s.prepareImport(context.Background(), project.ID, req); !errors.Is(err, errImportReviewChanged) {
		t.Fatalf("analysis scope change: %v", err)
	}
}

func TestImportReviewFreezesAcceptedCandidates(t *testing.T) {
	s := newAuthTestServer(t)
	project, _ := s.store.CreateProject(store.Project{Name: "review"})
	root := t.TempDir()
	mkdirGit(t, filepath.Join(root, "one"))
	req := importReposRequest{Provider: "local", Root: root, DryRun: true}
	preview, err := s.prepareImport(context.Background(), project.ID, req)
	if err != nil {
		t.Fatal(err)
	}
	req.DryRun = false
	req.PreviewDigest = preview.digest
	accepted, err := s.prepareImport(context.Background(), project.ID, req)
	if err != nil {
		t.Fatal(err)
	}
	mkdirGit(t, filepath.Join(root, "two"))
	results := accepted.apply(s, project.ID, req)
	repos, _ := s.store.ListRepos(project.ID)
	if len(results) != 1 || len(repos) != 1 || repos[0].Name != "one" {
		t.Fatalf("accepted scope expanded: %+v %+v", results, repos)
	}
}
