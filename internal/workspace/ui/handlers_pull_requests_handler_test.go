package ui

import (
	"testing"

	"github.com/mohammad-safakhou/diffmind/internal/workspace/archgraph"
	"github.com/mohammad-safakhou/diffmind/internal/workspace/model"
)

func TestChangedEntrypointsHandlerBodyAndEvidencePriority(t *testing.T) {
	service := &archgraph.ServiceNode{HTTPRoutes: []archgraph.EntitySummary{
		{ID: "login", Name: "POST /login", Details: map[string]any{
			"repository_revision": map[string]any{"commit": "head"},
			"source_locations": []model.Location{
				{File: "routes.py", StartLine: 4, EndLine: 4},
				{File: "handlers.py", StartLine: 5, EndLine: 9},
			},
		}},
		{ID: "other", Name: "GET /other", Details: map[string]any{
			"repository_revision": map[string]any{"commit": "head"},
			"source_locations":    []model.Location{{File: "handlers.py", StartLine: 12, EndLine: 16}},
		}},
	}}
	files := []githubPullFile{
		{Filename: "routes.py"}, // missing patch cannot mask precise handler evidence
		{Filename: "handlers.py", Patch: "@@ -6 +6 @@\n-old\n+new"},
	}
	got := changedEntrypoints(service, files, "head")
	if len(got) != 1 || got[0].ID != "login" || got[0].File != "handlers.py" || got[0].Match != "changed_line" {
		t.Fatalf("changed entrypoints = %+v, want own handler only", got)
	}
	for _, revision := range []map[string]any{{"commit": "old"}, {"commit": "head", "dirty": true}} {
		service.HTTPRoutes[0].Details["repository_revision"] = revision
		got = changedEntrypoints(service, files, "head")
		for _, entry := range got {
			if entry.ID == "login" && entry.Match != "file_scope" {
				t.Fatalf("unverified handler promoted to exact: %+v", entry)
			}
		}
	}
}
