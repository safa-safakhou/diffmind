package ui

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/mohammad-safakhou/diffmind/internal/workspace/store"
)

func TestGitHubRepositoryEndpointStatesAndExactHosts(t *testing.T) {
	tests := []struct {
		repo        store.Repo
		state, base string
	}{
		{store.Repo{SourceType: "local", Path: t.TempDir()}, "local_only", ""},
		{store.Repo{SourceType: "git"}, "missing_remote", ""},
		{store.Repo{GitURL: "https://gitlab.com/acme/repo.git"}, "unsupported_provider", ""},
		{store.Repo{GitURL: "https://evil.example/github.com/acme/repo.git"}, "unsupported_provider", ""},
		{store.Repo{GitURL: "https://github.com.evil.example/acme/repo.git"}, "unsupported_provider", ""},
		{store.Repo{GitURL: "git@github.com:acme/repo.git"}, "ready", "https://api.github.com/repos/acme/repo"},
		{store.Repo{GitURL: "https://github.company.test/acme/repo.git", GitProvider: "github", GitAPIBase: "https://github.company.test/api/v3"}, "ready", "https://github.company.test/api/v3/repos/acme/repo"},
	}
	for _, tt := range tests {
		base, state, _ := githubRepositoryEndpoint(context.Background(), tt.repo)
		if state != tt.state || base != tt.base {
			t.Fatalf("repo=%+v state=%s base=%s", tt.repo, state, base)
		}
	}
	if _, _, ok := githubOwnerRepo("https://evil.example/github.com/acme/repo"); ok {
		t.Fatal("substring host accepted")
	}
}

func TestGitHubImportRetainsAPIEndpointForPRQueries(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "fixture-secret")
	var pulls atomic.Int32
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer fixture-secret" {
			t.Error("approved endpoint missing token")
		}
		switch r.URL.Path {
		case "/orgs/acme/repos":
			json.NewEncoder(w).Encode([]githubRepo{{Name: "repo", CloneURL: "https://github.company.test/acme/repo.git", DefaultBranch: "main"}})
		case "/repos/acme/repo/pulls":
			pulls.Add(1)
			w.Write([]byte("[]"))
		case "/repos/acme/repo/issues":
			w.Write([]byte("[]"))
		case "/repos/acme/repo/actions/runs":
			w.Write([]byte(`{"workflow_runs":[]}`))
		default:
			t.Errorf("unexpected path %s", r.URL.Path)
			http.NotFound(w, r)
		}
	}))
	defer api.Close()
	s := newAuthTestServer(t)
	project, _ := s.store.CreateProject(store.Project{Name: "provider"})
	req := importReposRequest{Provider: "github", Org: "acme", APIBase: api.URL, CloneTransport: "https", DryRun: true}
	preview, err := s.prepareImport(context.Background(), project.ID, req)
	if err != nil {
		t.Fatal(err)
	}
	req.DryRun = false
	req.PreviewDigest = preview.digest
	prepared, err := s.prepareImport(context.Background(), project.ID, req)
	if err != nil {
		t.Fatal(err)
	}
	results := prepared.apply(s, project.ID, req)
	if len(results) != 1 || results[0].Status != "imported" {
		t.Fatalf("import=%+v", results)
	}
	repo, _ := s.store.GetRepo(project.ID, results[0].RepoID)
	if repo.GitAPIBase != api.URL {
		t.Fatalf("API base=%s", repo.GitAPIBase)
	}
	result := githubOpenPullRequests(context.Background(), workspaceRepo{Repo: *repo})
	if result.Status != "ok" || result.OpenCount != 0 || pulls.Load() != 1 {
		t.Fatalf("PR result=%+v calls=%d", result, pulls.Load())
	}
	live := githubLiveStatus(context.Background(), *repo)
	if live.Status != "ok" {
		t.Fatalf("live status=%+v", live)
	}
}

func TestGitHubClientsRejectRedirectsWithoutForwardingCredentials(t *testing.T) {
	var leaked atomic.Int32
	other := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { leaked.Add(1) }))
	defer other.Close()
	redirect := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, other.URL, http.StatusFound) }))
	defer redirect.Close()
	var result []githubPull
	err := githubJSON(context.Background(), githubHTTPClient(time.Second), "secret", redirect.URL, &result)
	if err == nil || !strings.Contains(err.Error(), "redirect policy") || leaked.Load() != 0 {
		t.Fatalf("redirect err=%v leaked requests=%d", err, leaked.Load())
	}
	t.Setenv("GITHUB_TOKEN", "secret")
	_, err = githubOrgRepos(context.Background(), importReposRequest{Org: "acme", APIBase: redirect.URL})
	if err == nil || leaked.Load() != 0 {
		t.Fatalf("import redirect err=%v leaked requests=%d", err, leaked.Load())
	}
}

func TestGitHubFeedbackDoesNotReflectProviderSecrets(t *testing.T) {
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"secret-token reflected"}`))
	}))
	defer api.Close()
	var result []githubPull
	err := githubJSON(context.Background(), githubHTTPClient(time.Second), "secret-token", api.URL, &result)
	if err == nil || strings.Contains(err.Error(), "secret-token") || !strings.Contains(err.Error(), "authentication failed") {
		t.Fatalf("unsafe auth error=%v", err)
	}
}

func TestGitHubAPIConfigurationValidation(t *testing.T) {
	for _, raw := range []string{"https://api.github.com", "https://github.company.test/api/v3", "http://127.0.0.1:8080"} {
		if err := validateGitHubAPIBase(raw); err != nil {
			t.Fatalf("valid %q: %v", raw, err)
		}
	}
	for _, raw := range []string{"http://other.example", "https://token@api.github.com", "https://api.github.com?token=secret", "/api/v3", "ftp://api.github.com"} {
		if err := validateGitHubAPIBase(raw); err == nil {
			t.Fatalf("invalid %q accepted", raw)
		}
	}
	for _, raw := range []string{"https://github.com.evil.test/org/repo", "https://evil.test/github.com/org/repo"} {
		if inferGitProvider(raw) == "github" {
			t.Fatalf("incorrect provider inferred for %s", raw)
		}
	}
	_, state, _ := githubRepositoryEndpoint(context.Background(), store.Repo{GitURL: "https://github.com.evil.test/org/repo", GitProvider: "github"})
	if state != "provider_unavailable" {
		t.Fatalf("legacy guessed endpoint state=%s", state)
	}
}

func TestChangingRepositoryURLRequiresAPIEndpointReview(t *testing.T) {
	s := newAuthTestServer(t)
	project, _ := s.store.CreateProject(store.Project{Name: "provider edit"})
	repo, _ := s.store.CreateRepo(project.ID, store.Repo{Name: "service", GitURL: "https://github.old.test/acme/service.git", GitProvider: "github", GitAPIBase: "https://github.old.test/api/v3"})
	update := func(body string) *store.Repo {
		t.Helper()
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPatch, "/", bytes.NewBufferString(body))
		r.SetPathValue("pid", project.ID)
		r.SetPathValue("rid", repo.ID)
		s.handlePatchRepo(w, r)
		if w.Code != 200 {
			t.Fatalf("patch=%d %s", w.Code, w.Body.String())
		}
		updated, _ := s.store.GetRepo(project.ID, repo.ID)
		return updated
	}
	same := update(`{"git_url":"https://github.old.test/acme/service.git"}`)
	if same.GitAPIBase == "" {
		t.Fatal("unchanged URL lost approved endpoint")
	}
	changed := update(`{"git_url":"https://github.new.test/acme/service.git"}`)
	if changed.GitAPIBase != "" {
		t.Fatal("new source inherited old API approval")
	}
	reviewed := update(`{"git_url":"https://github.new.test/acme/service.git","git_api_base":"https://github.new.test/api/v3"}`)
	if reviewed.GitAPIBase != "https://github.new.test/api/v3" {
		t.Fatalf("reviewed endpoint=%s", reviewed.GitAPIBase)
	}
}

func TestCustomAPIHostnameSelectsMatchingCredential(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "")
	t.Setenv("GH_TOKEN", "")
	bin := t.TempDir()
	script := filepath.Join(bin, "gh")
	if err := os.WriteFile(script, []byte("#!/bin/sh\nprintf '%s' \"$*\"\n"), 0700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin)
	if got := githubAPIToken(context.Background(), "https://code.company.test/api/v3/repos/org/service"); got != "auth token --hostname code.company.test" {
		t.Fatalf("custom credentials selected %q", got)
	}
	if got := githubAPIToken(context.Background(), "https://api.github.com/repos/org/service"); got != "auth token" {
		t.Fatalf("public credentials selected %q", got)
	}
}
