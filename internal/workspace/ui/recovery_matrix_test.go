package ui

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/mohammad-safakhou/diffmind/internal/workspace/agentapi"
	"github.com/mohammad-safakhou/diffmind/internal/workspace/store"
)

func TestRecoveryMatrixActualHTTPAndManagementAdapter(t *testing.T) {
	s := newAuthTestServer(t)
	s.SetTrustedProxySecret("proxy")
	s.SetAuthToken("recovery")
	if err := s.ConfigureProjectAccess("scoped"); err != nil {
		t.Fatal(err)
	}
	p, err := s.store.CreateProject(store.Project{Name: "Recovery matrix"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = s.store.PutProjectAccess(p.ID, 0, map[string]string{"member": "viewer"}); err != nil {
		t.Fatal(err)
	}
	if _, err = s.putProjectLimits(p.ID, 0, 1, 1); err != nil {
		t.Fatal(err)
	}
	if _, _, err = s.enqueueRefresh(p.ID, "scheduled", "", ""); err != nil {
		t.Fatal(err)
	}
	provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("private-provider-response-secret"))
	}))
	defer provider.Close()
	for _, tc := range []struct {
		name, user, role string
		in               agentapi.Input
		code             int
		category         string
	}{
		{"invalid root", "admin", "admin", agentapi.Input{Operation: "import_repositories", Body: map[string]any{"provider": "local", "root": filepath.Join(t.TempDir(), "missing"), "dry_run": true}}, 400, "validation"},
		{"invalid configuration", "admin", "admin", agentapi.Input{Operation: "start_ingestion", Body: map[string]any{"concurrency": -1}}, 400, "validation"},
		{"denied authority", "member", "viewer", agentapi.Input{Operation: "delete_project", Confirm: "delete_project"}, 403, "access_or_availability"},
		{"ungranted existing project", "outsider", "viewer", agentapi.Input{Operation: "get_project"}, 404, "access_or_availability"},
		{"authentication unavailable", "", "", agentapi.Input{Operation: "get_project"}, 401, "access_or_availability"},
		{"access revision conflict", "admin", "admin", agentapi.Input{Operation: "set_access", Body: map[string]any{"revision": 0, "members": map[string]string{}}}, 409, "conflict"},
		{"queue capacity", "admin", "admin", agentapi.Input{Operation: "enqueue_refresh"}, 429, "capacity"},
		{"provider authentication", "admin", "admin", agentapi.Input{Operation: "import_repositories", Body: map[string]any{"provider": "github", "org": "public-fixture", "api_base": provider.URL, "dry_run": true}}, 502, "environment_or_provider"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tc.in.Selectors = map[string]string{"pid": p.ID}
			op, ok := agentapi.Find(tc.in.Operation)
			if !ok {
				t.Fatal("unknown fixture operation", tc.in.Operation)
			}
			request, err := agentapi.Request(context.Background(), tc.in, op.Method == "GET")
			if err != nil {
				t.Fatal(err)
			}
			headers := http.Header{}
			if tc.user != "" {
				headers.Set(proxySecretHeader, "proxy")
				headers.Set(proxyUserHeader, tc.user)
				headers.Set(proxyRoleHeader, tc.role)
			}
			request.Header = headers.Clone()
			direct := httptest.NewRecorder()
			s.Handler().ServeHTTP(direct, request)
			var body struct {
				Error    string            `json:"error"`
				Recovery agentapi.Recovery `json:"recovery"`
			}
			if err = json.Unmarshal(direct.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			request, err = agentapi.Request(context.Background(), tc.in, op.Method == "GET")
			if err != nil {
				t.Fatal(err)
			}
			result, err := s.invokeAgentOperation(context.Background(), &mcp.CallToolRequest{Extra: &mcp.RequestExtra{Header: headers}}, request)
			if err != nil {
				t.Fatal(err)
			}
			if direct.Code != tc.code || result.Status != tc.code || result.Recovery == nil || !reflect.DeepEqual(*result.Recovery, body.Recovery) || body.Recovery.Category != tc.category || body.Recovery.Retryable || body.Recovery.NextAction == "" {
				t.Fatalf("HTTP %d %s; adapter %+v; expected %d %s", direct.Code, direct.Body, result, tc.code, tc.category)
			}
			if strings.Contains(direct.Body.String(), "private-provider-response-secret") {
				t.Fatal("provider diagnostics leaked")
			}
		})
	}
	projects, err := s.store.ListProjects()
	if err != nil || len(projects) != 1 {
		t.Fatal("failed mutation changed registrations")
	}
	jobs, err := s.store.ListJobs(p.ID)
	if err != nil || len(jobs) != 1 {
		t.Fatal("failed mutation changed queued work")
	}
}

func TestMissingExecutableIsAcceptedWorkThenObservableFailure(t *testing.T) {
	s := newAuthTestServer(t)
	t.Setenv("DIFFMIND_BINARY", filepath.Join(t.TempDir(), "missing-analyzer"))
	p, err := s.store.CreateProject(store.Project{Name: "Missing executable"})
	if err != nil {
		t.Fatal(err)
	}
	source := t.TempDir()
	if err = os.WriteFile(filepath.Join(source, "main.go"), []byte("package main\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err = s.store.CreateRepo(p.ID, store.Repo{Name: "source", Path: source, SourceType: "local", Kind: "service_repo"}); err != nil {
		t.Fatal(err)
	}
	w := ingestionPost(s, p.ID, "", "{}")
	if w.Code != http.StatusAccepted {
		t.Fatal(w.Code, w.Body)
	}
	state := awaitIngestionIdle(t, s, p.ID)
	if state.Status != store.IngestionFailed || len(state.Errors) == 0 {
		t.Fatalf("failure not observable: %+v", state)
	}
	repos, err := s.store.ListRepos(p.ID)
	if err != nil || len(repos) != 1 || len(state.Request) == 0 {
		t.Fatal("failure discarded source context")
	}
}
