package ui

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mohammad-safakhou/diffmind/internal/workspace/archgraph"
	"github.com/mohammad-safakhou/diffmind/internal/workspace/store"
)

func TestPRFlowsMatchMergeBaseRatherThanMovingBaseAndRetainRemovedSurfaces(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "")
	t.Setenv("GH_TOKEN", "")
	providerFailure := false
	provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if providerFailure {
			http.Error(w, "provider unavailable", 503)
			return
		}
		switch r.URL.Path {
		case "/repos/acme/orders/pulls/7":
			w.Write([]byte(`{"number":7,"head":{"sha":"head-commit"},"base":{"sha":"moving-base"}}`))
		case "/repos/acme/orders/compare/moving-base...head-commit":
			w.Write([]byte(`{"merge_base_commit":{"sha":"merge-base"}}`))
		default:
			t.Errorf("unexpected provider path %s", r.URL.Path)
			http.NotFound(w, r)
		}
	}))
	defer provider.Close()
	s := newAuthTestServer(t)
	project, _ := s.store.CreateProject(store.Project{Name: "Flow review"})
	repo, err := s.store.CreateRepo(project.ID, store.Repo{Name: "orders", GitProvider: "github", GitURL: "https://github.fixture.test/acme/orders.git", GitAPIBase: provider.URL})
	if err != nil {
		t.Fatal(err)
	}
	save := func(commit string, dirty, removed bool) string {
		run, err := s.store.CreateRun(project.ID, store.RunManifest{Status: store.RunCompleted, StartedAt: time.Now()})
		if err != nil {
			t.Fatal(err)
		}
		svc := &archgraph.ServiceNode{Name: "orders", RepoID: repo.ID, Known: true, AnalysisStatus: &archgraph.RepositoryAnalysisStatus{State: "analyzed_clean", AnalyzedRevision: commit, Dirty: dirty}, HTTPRoutes: []archgraph.EntitySummary{{ID: "old", Name: "DELETE /orders/{id}", Kind: "http_endpoint"}}}
		if removed {
			svc.HTTPRoutes = nil
		}
		body, _ := json.Marshal(&archgraph.ArchGraph{RunID: run.ID, Services: []*archgraph.ServiceNode{svc}})
		if err := os.WriteFile(filepath.Join(s.store.RunDir(project.ID, run.ID), "graph.json"), body, 0600); err != nil {
			t.Fatal(err)
		}
		return run.ID
	}
	wrongBase := save("moving-base", false, false)
	cleanBase := save("merge-base", false, false)
	cleanHead := save("head-commit", false, true)
	dirtyHead := save("head-commit", true, false)
	url := "/api/projects/" + project.ID + "/pull-requests/" + repo.ID + "/7/flows"
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, httptest.NewRequest("GET", url, nil))
	if w.Code != 200 {
		t.Fatalf("%d %s", w.Code, w.Body)
	}
	var result pullRequestFlows
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if !result.Available || result.BeforeRun != cleanBase || result.AfterRun != cleanHead || len(result.Review.Flows) != 1 || result.Review.Flows[0].Change != "removed" {
		t.Fatalf("review %+v", result)
	}
	for _, bad := range []string{url + "?from=" + wrongBase + "&to=" + cleanHead, url + "?from=" + cleanBase + "&to=" + dirtyHead, url + "?offset=-1", url + "?limit=51"} {
		w = httptest.NewRecorder()
		s.Handler().ServeHTTP(w, httptest.NewRequest("GET", bad, nil))
		if w.Code != 400 {
			t.Fatalf("%s status=%d body=%s", bad, w.Code, w.Body)
		}
	}
	// A missing clean base cannot be replaced with the moving branch's snapshot.
	if err := s.store.DeleteRun(project.ID, cleanBase); err != nil {
		t.Fatal(err)
	}
	w = httptest.NewRecorder()
	s.Handler().ServeHTTP(w, httptest.NewRequest("GET", url, nil))
	if w.Code != 200 {
		t.Fatal(w.Body)
	}
	result = pullRequestFlows{}
	json.Unmarshal(w.Body.Bytes(), &result)
	if result.Available || result.BeforeRun != "" || result.NextAction == "" {
		t.Fatal("missing base overclaimed")
	}
	providerFailure = true
	w = httptest.NewRecorder()
	s.Handler().ServeHTTP(w, httptest.NewRequest("GET", url, nil))
	if w.Code != 502 {
		t.Fatalf("provider failure %d", w.Code)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, _, _, err := s.findPRFlowRuns(ctx, project.ID, *repo, "merge-base", "head-commit", "", ""); err != context.Canceled {
		t.Fatalf("cancellation %v", err)
	}
}
