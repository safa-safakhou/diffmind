package ui

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/mohammad-safakhou/diffmind/internal/workspace/query"
	"github.com/mohammad-safakhou/diffmind/internal/workspace/store"
)

func TestReadinessHTTPWorkspaceRoleParity(t *testing.T) {
	s, pid, private := accessFixture(t)
	jobs, err := s.store.ListJobs(pid)
	if err != nil {
		t.Fatal(err)
	}
	for _, job := range jobs {
		if _, err := s.store.CancelJob(job.ID); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := s.store.CreateRepo(pid, store.Repo{Name: "source", Path: t.TempDir()}); err != nil {
		t.Fatal(err)
	}
	for _, status := range []string{"running", "partial", "failed", "completed"} {
		ingestion, err := s.store.CreateIngestion(pid, store.Ingestion{Status: status, Phase: "analyzing"})
		if err != nil {
			t.Fatal(err)
		}
		for _, role := range []string{"viewer", "editor", "admin"} {
			h := s.Handler()
			response := accessRequest(h, "alice", role, "GET", "/api/v1/projects/"+pid+"/readiness", "")
			if response.Code != http.StatusOK {
				t.Fatal(response.Code, response.Body.String())
			}
			var got query.Readiness
			if err := json.Unmarshal(response.Body.Bytes(), &got); err != nil {
				t.Fatal(err)
			}
			if got.Graph != "queryable" || got.Runtime != "available" || got.Work.ID != ingestion.ID || got.Work.Status != status || !got.Actions.Query {
				t.Fatalf("%s %s: %+v", status, role, got)
			}
			canRefresh := role != "viewer" && status != "running"
			canConfigure := role == "admin" && status != "running"
			if got.Actions.Refresh != canRefresh || got.Actions.Configure != canConfigure {
				t.Fatalf("permissions %s %s: %+v", status, role, got.Actions)
			}
			response = accessRequest(h, "alice", role, "GET", "/api/projects/"+pid+"/workspace?graph=false", "")
			var workspace workspaceResponse
			if response.Code != 200 {
				t.Fatal(response.Code, response.Body.String())
			}
			if err := json.Unmarshal(response.Body.Bytes(), &workspace); err != nil {
				t.Fatal(err)
			}
			if workspace.Readiness.SavedRunID != got.SavedRunID || workspace.Readiness.Work != got.Work || workspace.Readiness.Actions != got.Actions || workspace.Readiness.Evidence != got.Evidence {
				t.Fatalf("HTTP/workspace mismatch: %+v %+v", got, workspace.Readiness)
			}
		}
		// Avoid attempting to replace an active ingestion in the following case.
		ingestion.Status = "completed"
		ingestion.FinishedAt = time.Now()
		if err := s.store.SaveIngestion(pid, *ingestion); err != nil {
			t.Fatal(err)
		}
	}
	denied := accessRequest(s.Handler(), "alice", "viewer", "GET", "/api/v1/projects/"+private+"/readiness", "")
	if denied.Code != 404 {
		t.Fatalf("denied response %d", denied.Code)
	}
	if _, err := s.store.PutProjectAccess(pid, 1, map[string]string{}); err != nil {
		t.Fatal(err)
	}
	denied = accessRequest(s.Handler(), "alice", "viewer", "GET", "/api/v1/projects/"+pid+"/readiness", "")
	if denied.Code != 404 {
		t.Fatalf("revoked response %d", denied.Code)
	}
}
