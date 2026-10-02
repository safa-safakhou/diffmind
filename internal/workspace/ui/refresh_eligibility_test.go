package ui

import (
	"context"
	"github.com/mohammad-safakhou/diffmind/internal/workspace/store"
	"testing"
	"time"
)

func TestAutomaticRefreshEligibilityPersistsCadenceAndFailureBackoff(t *testing.T) {
	now := time.Now().UTC()
	tests := []struct {
		name     string
		history  []store.RefreshJob
		interval time.Duration
		want     string
		next     time.Duration
	}{
		{"empty", nil, 15 * time.Minute, "", 0},
		{"recent success", []store.RefreshJob{{Status: "succeeded", UpdatedAt: now.Add(-time.Minute)}}, 15 * time.Minute, "not_due", 14 * time.Minute},
		{"overdue success", []store.RefreshJob{{Status: "succeeded", UpdatedAt: now.Add(-16 * time.Minute)}}, 15 * time.Minute, "", 0},
		{"explicit startup with manual cadence", []store.RefreshJob{{Status: "succeeded", UpdatedAt: now}}, 0, "", 0},
		{"failed terminal job", []store.RefreshJob{{Status: "failed", UpdatedAt: now}}, 15 * time.Minute, "failure_backoff", 30 * time.Minute},
		{"repeated failures", []store.RefreshJob{{Status: "failed", UpdatedAt: now}, {Status: "failed"}, {Status: "succeeded"}}, 15 * time.Minute, "failure_backoff", time.Hour},
		{"short interval remains bounded", []store.RefreshJob{{Status: "failed", UpdatedAt: now}}, time.Second, "failure_backoff", 2 * time.Minute},
		{"overdue failure", []store.RefreshJob{{Status: "failed", UpdatedAt: now.Add(-time.Hour)}}, 15 * time.Minute, "", 0},
		{"pending uses queue coalescing", []store.RefreshJob{{Status: "running", UpdatedAt: now}}, 15 * time.Minute, "", 0},
		{"attempt completion authoritative", []store.RefreshJob{{Status: "succeeded", UpdatedAt: now.Add(-time.Hour), Attempts: []store.JobAttempt{{FinishedAt: now}}}}, 15 * time.Minute, "not_due", 15 * time.Minute},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reason, next := automaticRefreshEligibility(tt.history, tt.interval, now)
			if reason != tt.want {
				t.Fatalf("reason=%s want=%s", reason, tt.want)
			}
			if reason != "" && next.Sub(now) != tt.next {
				t.Fatalf("next=%s want=%s", next.Sub(now), tt.next)
			}
		})
	}
	history := make([]store.RefreshJob, 20)
	for i := range history {
		history[i] = store.RefreshJob{Status: "failed", UpdatedAt: now}
	}
	_, next := automaticRefreshEligibility(history, 24*time.Hour, now)
	if next.Sub(now) != 6*time.Hour {
		t.Fatalf("backoff cap=%s", next.Sub(now))
	}
}

func TestReconnectSkipsRecentSuccessButManualRefreshStillRuns(t *testing.T) {
	s := newAuthTestServer(t)
	project, err := s.store.CreateProject(store.Project{Name: "reconnect"})
	if err != nil {
		t.Fatal(err)
	}
	job, _, err := s.store.EnqueueJob(project.ID, "fleet_refresh", "", "", 32)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = s.store.ClaimJob(time.Now().UTC()); err != nil {
		t.Fatal(err)
	}
	if err = s.store.FinishJob(job.ID, store.JobAttempt{Status: "succeeded"}, false, 0); err != nil {
		t.Fatal(err)
	}
	restarted := New(s.store, s.runs, s.diffmindRunsDir, "127.0.0.1", 8090, s.log)
	t.Cleanup(restarted.StopOperations)
	restarted.refreshProject = func(_ context.Context, pid string) ProjectRefreshResult { return ProjectRefreshResult{ProjectID: pid} }
	if err = restarted.ConfigureRefresh(RefreshConfig{Interval: 15 * time.Minute, OnStart: true}); err != nil {
		t.Fatal(err)
	}
	results, err := restarted.refreshProjects(context.Background(), "startup")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Skipped != "not_due" || results[0].NextEligible == nil {
		t.Fatalf("reconnect=%+v", results)
	}
	history, _ := s.store.ListJobs(project.ID)
	if len(history) != 1 {
		t.Fatalf("startup created %d jobs", len(history))
	}
	results, err = restarted.refreshAllProjects(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Skipped != "" {
		t.Fatalf("manual=%+v", results)
	}
	history, _ = s.store.ListJobs(project.ID)
	if len(history) != 2 {
		t.Fatalf("manual created %d jobs", len(history))
	}
}

func TestCompletedIngestionSatisfiesReconnectCadence(t *testing.T) {
	s := newAuthTestServer(t)
	project, _ := s.store.CreateProject(store.Project{Name: "manual ingestion"})
	now := time.Now().UTC()
	if err := s.store.SaveIngestion(project.ID, store.Ingestion{ID: "recent-ingestion", ProjectID: project.ID, Status: store.IngestionCompleted, FinishedAt: now}); err != nil {
		t.Fatal(err)
	}
	history, err := s.automaticRefreshHistory(project.ID)
	if err != nil {
		t.Fatal(err)
	}
	reason, _ := automaticRefreshEligibility(history, 15*time.Minute, now)
	if reason != "not_due" {
		t.Fatalf("recent ingestion=%s", reason)
	}
}
