package ui

import (
	"net/http/httptest"
	"testing"
)

func TestReviewedLegacyToScopedSetup(t *testing.T) {
	s, visible, private := accessFixture(t)
	if err := s.ConfigureProjectAccess("legacy"); err != nil {
		t.Fatal(err)
	}
	if w := accessRequest(s.Handler(), "ungranted", "viewer", "GET", "/api/projects/"+private, ""); w.Code != 200 {
		t.Fatalf("legacy compatibility: %d", w.Code)
	}
	if err := s.ConfigureProjectAccess("scoped"); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		user, role, path string
		status           int
	}{
		{"ungranted", "viewer", "/api/projects/" + visible, 404},
		{"alice", "viewer", "/api/projects/" + visible, 200},
		{"alice", "editor", "/api/projects/" + visible, 200},
		{"alice", "editor", "/api/projects/" + private, 404},
		{"admin", "admin", "/api/projects/" + private, 200},
	} {
		if w := accessRequest(s.Handler(), tc.user, tc.role, "GET", tc.path, ""); w.Code != tc.status {
			t.Fatalf("%s %s: %d", tc.user, tc.path, w.Code)
		}
	}
	request := httptest.NewRequest("GET", "/api/projects/"+private, nil)
	request.Header.Set("Authorization", "Bearer recovery")
	response := httptest.NewRecorder()
	s.Handler().ServeHTTP(response, request)
	if response.Code != 200 {
		t.Fatalf("recovery credential: %d", response.Code)
	}
	if w := accessRequest(s.Handler(), "alice", "viewer", "POST", "/api/v1/projects/"+visible+"/refresh-jobs", "{}"); w.Code != 403 {
		t.Fatalf("viewer mutation: %d", w.Code)
	}
	if w := accessRequest(s.Handler(), "alice", "editor", "POST", "/api/v1/projects/"+visible+"/refresh-jobs", "{}"); w.Code != 202 {
		t.Fatalf("editor refresh: %d", w.Code)
	}
}
