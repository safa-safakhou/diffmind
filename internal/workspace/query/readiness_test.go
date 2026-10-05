package query

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mohammad-safakhou/diffmind/internal/workspace/store"
)

func TestReadinessSavedGraphSurvivesWorkMatrix(t *testing.T) {
	for _, status := range []string{"running", "completed", "partial", "failed", "cancelled", "interrupted"} {
		t.Run(status, func(t *testing.T) {
			q, pid, rid := testQueryService(t)
			ingestion, err := q.store.CreateIngestion(pid, store.Ingestion{Status: status, Phase: "analyzing"})
			if err != nil {
				t.Fatal(err)
			}
			out, err := q.Readiness(pid)
			if err != nil {
				t.Fatal(err)
			}
			if out.Graph != "queryable" || out.SavedRunID != rid || !out.Actions.Query || out.Work.ID != ingestion.ID || out.Work.Status != status {
				t.Fatalf("%+v", out)
			}
			if out.Actions.Refresh || out.Actions.Configure || out.ConnectionMode != "query_only" || out.Maintenance != "unknown" || out.PRHeadEligibility != "unknown" || out.Evidence.Coverage != "unverified" {
				t.Fatalf("query-only authority/coverage: %+v", out)
			}
			if out.SavedAt == nil || out.SavedAt.IsZero() {
				t.Fatal("missing snapshot age")
			}
		})
	}
}

func TestReadinessQueuedRefreshOverridesOldFailureAndHonorsPermissions(t *testing.T) {
	q, pid, _ := testQueryService(t)
	if _, err := q.store.CreateRepo(pid, store.Repo{Name: "one", Path: t.TempDir()}); err != nil {
		t.Fatal(err)
	}
	if _, err := q.store.CreateIngestion(pid, store.Ingestion{Status: "failed", Phase: "analyzing"}); err != nil {
		t.Fatal(err)
	}
	job, _, err := q.store.EnqueueJob(pid, "manual", "", "", 10)
	if err != nil {
		t.Fatal(err)
	}
	editor := q.WithReadinessPermissions(func(string) (ReadinessActions, error) { return ReadinessActions{Refresh: true, Configure: true}, nil })
	out, err := editor.Readiness(pid)
	if err != nil || out.Work.Status != "queued" || out.Work.ID != job.ID || out.Actions.Refresh || out.Actions.Configure || !out.Actions.Query || out.NextAction != "inspect_work" {
		t.Fatalf("queued: %+v %v", out, err)
	}
	if _, err := q.store.CancelJob(job.ID); err != nil {
		t.Fatal(err)
	}
	out, err = editor.Readiness(pid)
	if err != nil || !out.Actions.Refresh || !out.Actions.Configure {
		t.Fatalf("idle permissions: %+v %v", out, err)
	}
}

func TestReadinessEmptyAndDamagedGraphAreDistinct(t *testing.T) {
	st, err := store.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	p, err := st.CreateProject(store.Project{Name: "empty"})
	if err != nil {
		t.Fatal(err)
	}
	q := New(st)
	out, err := q.Readiness(p.ID)
	if err != nil || out.Graph != "missing" || out.Work.Status != "empty" || out.Actions.Query {
		t.Fatalf("empty: %+v %v", out, err)
	}
	run, err := st.CreateRun(p.ID, store.RunManifest{Status: "completed", StartedAt: time.Now()})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(st.RunDir(p.ID, run.ID), "graph.json"), []byte("null"), 0600); err != nil {
		t.Fatal(err)
	}
	out, err = q.Readiness(p.ID)
	if err != nil || out.Graph != "unavailable" || out.Actions.Query || out.Evidence.GraphState != "unavailable" {
		t.Fatalf("damaged: %+v %v", out, err)
	}
}

func TestReadinessAccessBeforeWorkAndArtifacts(t *testing.T) {
	q, pid, _ := testQueryService(t)
	denied := NewWithAccess(q.store, func(string) error { return store.ErrNotFound })
	out, err := denied.Readiness(pid)
	if out != nil || !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("access leak: %+v %v", out, err)
	}
}

func TestReadinessJobSuccessKeepsPartialIngestion(t *testing.T) {
	q, pid, _ := testQueryService(t)
	job, _, err := q.store.EnqueueJob(pid, "manual", "", "", 10)
	if err != nil {
		t.Fatal(err)
	}
	claimed, err := q.store.ClaimJob(time.Now())
	if err != nil || claimed.ID != job.ID {
		t.Fatal(err)
	}
	in, err := q.store.CreateIngestion(pid, store.Ingestion{JobID: job.ID, Status: "partial", Phase: "complete"})
	if err != nil {
		t.Fatal(err)
	}
	if err := q.store.FinishJob(job.ID, store.JobAttempt{Status: "succeeded", IngestionID: in.ID}, false, 0); err != nil {
		t.Fatal(err)
	}
	out, err := q.Readiness(pid)
	if err != nil || out.Work.Kind != "refresh_job" || out.Work.Status != "partial" || out.Work.Phase != "complete" || out.NextAction != "inspect_work" {
		t.Fatalf("partial job outcome: %+v %v", out, err)
	}
	in.Status = "completed"
	if err := q.store.SaveIngestion(pid, *in); err != nil {
		t.Fatal(err)
	}
	out, err = q.Readiness(pid)
	if err != nil || out.Work.Status != "completed" {
		t.Fatalf("successful outcome: %+v %v", out, err)
	}
}
