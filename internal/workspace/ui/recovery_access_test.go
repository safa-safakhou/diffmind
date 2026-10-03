package ui

import (
	"github.com/mohammad-safakhou/diffmind/internal/workspace/backup"
	"github.com/mohammad-safakhou/diffmind/internal/workspace/store"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestOfflineRestoreRequiresAccessReconciliation(t *testing.T) {
	for _, backend := range []string{"json", "sqlite"} {
		t.Run(backend, func(t *testing.T) {
			s := newAuthTestServer(t)
			s.SetAuthToken("recovery")
			s.SetTrustedProxySecret("proxy")
			if err := s.ConfigureProjectAccess("scoped"); err != nil {
				t.Fatal(err)
			}
			p, err := s.store.CreateProject(store.Project{Name: "restore-drill"})
			if err != nil {
				t.Fatal(err)
			}
			if _, err = s.store.PutProjectAccess(p.ID, 0, map[string]string{"departed": "viewer"}); err != nil {
				t.Fatal(err)
			}
			departed, credential, err := s.store.IssueProjectToken(p.ID, "assigned", "viewer", "admin", time.Hour)
			if err != nil {
				t.Fatal(err)
			}
			_, unrelated, err := s.store.IssueProjectToken(p.ID, "service", "viewer", "admin", time.Hour)
			if err != nil {
				t.Fatal(err)
			}
			job, _, err := s.store.EnqueueJob(p.ID, "manual", "", "", 256)
			if err != nil {
				t.Fatal(err)
			}
			if backend == "sqlite" {
				if _, err = store.MigrateQueue(s.store.HomeDir()); err != nil {
					t.Fatal(err)
				}
			}
			if _, err = s.store.ClaimJob(time.Now()); err != nil {
				t.Fatal(err)
			}
			archive := filepath.Join(t.TempDir(), "snapshot.tar.gz")
			report, err := backup.Create(s.store.HomeDir(), archive, "recovery-drill", backup.DefaultMaxBytes)
			if err != nil {
				t.Fatal(err)
			}
			if _, err = backup.Verify(archive, report.SHA256, backup.DefaultMaxBytes); err != nil {
				t.Fatal(err)
			}
			info, err := os.Stat(archive)
			if err != nil || info.Mode().Perm() != 0600 {
				t.Fatal("archive confidentiality")
			}
			if _, err = s.store.PutProjectAccess(p.ID, 1, map[string]string{}); err != nil {
				t.Fatal(err)
			}
			if _, err = s.store.RevokeProjectToken(p.ID, departed.ID, "admin"); err != nil {
				t.Fatal(err)
			}
			if w := bearerRequest(s.Handler(), credential, "GET", "/api/projects/"+p.ID, ""); w.Code != 401 {
				t.Fatalf("pre-restore revocation: %d", w.Code)
			}
			home := s.store.HomeDir()
			// Exact test-owned home; preserve rollback copy and restore original path.
			if err = os.Rename(home, home+"-rollback"); err != nil {
				t.Fatal(err)
			}
			if _, err = backup.Restore(archive, home, report.SHA256, backup.DefaultMaxBytes, false); err != nil {
				t.Fatal(err)
			}
			restored, err := store.New(home)
			if err != nil {
				t.Fatal(err)
			}
			s.store = restored
			// No listener is opened during reconciliation; these requests probe handlers.
			if w := bearerRequest(s.Handler(), credential, "GET", "/api/projects/"+p.ID, ""); w.Code != 200 {
				t.Fatalf("expected old credential resurrection: %d", w.Code)
			}
			if w := accessRequest(s.Handler(), "departed", "viewer", "GET", "/api/projects/"+p.ID, ""); w.Code != 200 {
				t.Fatalf("expected old membership resurrection: %d", w.Code)
			}
			policy, err := restored.GetProjectAccess(p.ID)
			if err != nil {
				t.Fatal(err)
			}
			if _, err = restored.PutProjectAccess(p.ID, policy.Revision, map[string]string{}); err != nil {
				t.Fatal(err)
			}
			if _, err = restored.RevokeProjectToken(p.ID, departed.ID, "admin"); err != nil {
				t.Fatal(err)
			}
			if w := bearerRequest(s.Handler(), credential, "GET", "/api/projects/"+p.ID, ""); w.Code != 401 {
				t.Fatalf("reconciled token: %d", w.Code)
			}
			if w := accessRequest(s.Handler(), "departed", "viewer", "GET", "/api/projects/"+p.ID, ""); w.Code != 404 {
				t.Fatalf("reconciled member: %d", w.Code)
			}
			for _, secret := range []string{unrelated, "recovery"} {
				if w := bearerRequest(s.Handler(), secret, "GET", "/api/projects/"+p.ID, ""); w.Code != 200 {
					t.Fatalf("preserved service/recovery grant: %d", w.Code)
				}
			}
			if err = restored.RecoverJobs(); err != nil {
				t.Fatal(err)
			}
			recovered, err := restored.GetJob(job.ID)
			if err != nil || recovered.Status == "running" {
				t.Fatalf("interrupted queue recovery: %+v %v", recovered, err)
			}
			if backend == "sqlite" {
				if _, err = restored.VerifyQueue(); err != nil {
					t.Fatal(err)
				}
			}
			t.Logf("verified backup, exact-path restore, resurrection, reconciliation and interrupted %s queue recovery", backend)
		})
	}
}
