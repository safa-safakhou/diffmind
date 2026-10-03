package ui

import (
	"testing"

	"github.com/mohammad-safakhou/diffmind/internal/workspace/archgraph"
	"github.com/mohammad-safakhou/diffmind/internal/workspace/model"
	"github.com/mohammad-safakhou/diffmind/internal/workspace/store"
)

func TestAnalyzeCodebaseImpactFindsHighRiskSurfaces(t *testing.T) {
	pull := githubPull{Additions: 640, Deletions: 180, ChangedFiles: 5, Commits: 4}
	files := []githubPullFile{
		{Filename: "api/openapi.yaml", Additions: 80, Deletions: 30, Patch: "- /v1/orders\n+ /v2/orders"},
		{Filename: "db/migrations/004_drop_legacy.sql", Additions: 10, Deletions: 20, Patch: "+ALTER TABLE orders DROP COLUMN legacy_id"},
		{Filename: "internal/security/authorizer.go", Additions: 110, Deletions: 40},
		{Filename: "go.mod", Additions: 4, Deletions: 2},
		{Filename: "internal/orders/service.go", Additions: 436, Deletions: 88, Patch: "-func PublicOrder() {}"},
	}

	impact := analyzeCodebaseImpact(pull, files, false)
	if impact.RiskLevel != "high" && impact.RiskLevel != "critical" {
		t.Fatalf("risk level = %q (%d), want high or critical", impact.RiskLevel, impact.RiskScore)
	}
	for _, want := range []string{"api", "data", "security", "dependencies", "code"} {
		if !hasChangeCategory(impact.Categories, want) {
			t.Errorf("missing category %q in %+v", want, impact.Categories)
		}
	}
	for _, want := range []string{"contract_change", "destructive_schema", "security_boundary", "public_surface_removal"} {
		if !hasSemanticSignal(impact.Signals, want) {
			t.Errorf("missing signal %q in %+v", want, impact.Signals)
		}
	}
	if !containsExactString(impact.RiskReasons, "production code changed without test-file changes") {
		t.Fatalf("missing no-tests risk reason: %+v", impact.RiskReasons)
	}
}

func TestAnalyzeCodebaseImpactKeepsTestOnlyChangeLowRisk(t *testing.T) {
	pull := githubPull{Additions: 42, Deletions: 8, ChangedFiles: 2, Commits: 1}
	impact := analyzeCodebaseImpact(pull, []githubPullFile{
		{Filename: "internal/orders/service_test.go", Additions: 30, Deletions: 5},
		{Filename: "tests/orders.spec.js", Additions: 12, Deletions: 3},
	}, false)
	if impact.RiskLevel != "low" {
		t.Fatalf("risk level = %q (%d), want low", impact.RiskLevel, impact.RiskScore)
	}
	if len(impact.Categories) != 1 || impact.Categories[0].ID != "tests" {
		t.Fatalf("categories = %+v, want tests only", impact.Categories)
	}
}

func TestGraphServiceForRepoUsesStableIdentityThenPathThenName(t *testing.T) {
	graph := &archgraph.ArchGraph{Services: []*archgraph.ServiceNode{
		{Name: "orders-service", RepoID: "repo-orders", RepoPath: "/repos/orders"},
		{Name: "billing_api", RepoPath: "/repos/billing"},
	}}
	cases := []struct {
		name string
		repo store.Repo
		want string
	}{
		{name: "repo id", repo: store.Repo{ID: "repo-orders", Name: "different", Path: "/elsewhere"}, want: "orders-service"},
		{name: "path", repo: store.Repo{ID: "unknown", Name: "different", Path: "/repos/billing"}, want: "billing_api"},
		{name: "normalized name", repo: store.Repo{Name: "billing-api", Path: "/elsewhere"}, want: "billing_api"},
		{name: "not represented", repo: store.Repo{Name: "search-index"}, want: ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := graphServiceForRepo(graph, tc.repo); got != tc.want {
				t.Fatalf("graphServiceForRepo() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestClassifyChangedFile(t *testing.T) {
	cases := map[string]string{
		"proto/orders.proto":                      "api",
		"db/migrations/001_orders.sql":            "data",
		"deploy/terraform/main.tf":                "infrastructure",
		"internal/auth/permissions.go":            "security",
		".example/config/production/secrets.yaml": "security",
		"package-lock.json":                       "dependencies",
		"config/application-production.yaml":      "configuration",
		"internal/orders/service_test.go":         "tests",
		"docs/operations.md":                      "documentation",
		"internal/orders/service.go":              "code",
	}
	for path, want := range cases {
		if got := classifyChangedFile(path); got != want {
			t.Errorf("classifyChangedFile(%q) = %q, want %q", path, got, want)
		}
	}
}

func TestChangedEntrypointsUsesChangedLinesNotWholeFile(t *testing.T) {
	service := &archgraph.ServiceNode{HTTPRoutes: []archgraph.EntitySummary{
		{ID: "first", Kind: "http_endpoint", Name: "GET /items", Details: map[string]any{"source_locations": []any{map[string]any{"file": "controller.go", "start_line": float64(10), "end_line": float64(15)}}}},
		{ID: "second", Kind: "http_endpoint", Name: "PATCH /items/{id}", Details: map[string]any{"source_locations": []any{map[string]any{"file": "controller.go", "start_line": float64(30), "end_line": float64(35)}}}},
	}}
	files := []githubPullFile{{Filename: "controller.go", Patch: "@@ -30,2 +30,2 @@\n context\n-old\n+new"}}
	for i := range service.HTTPRoutes {
		service.HTTPRoutes[i].Details["repository_revision"] = map[string]any{"commit": "head"}
	}

	got := changedEntrypoints(service, files, "head")
	if len(got) != 1 || got[0].ID != "second" || got[0].Match != "changed_line" {
		t.Fatalf("changed entrypoints = %+v, want only second with changed-line evidence", got)
	}
}

func TestExactChangedSurfaceCallersRejectsUnrelatedServiceEdge(t *testing.T) {
	graph := &archgraph.ArchGraph{Services: []*archgraph.ServiceNode{
		{Name: "api"}, {Name: "exact-client", AnalysisStatus: &archgraph.RepositoryAnalysisStatus{State: "analyzed_clean", AnalyzedRevision: "caller-head"}}, {Name: "other-client", AnalysisStatus: &archgraph.RepositoryAnalysisStatus{State: "analyzed_clean", AnalyzedRevision: "caller-head"}},
	}, Edges: []*archgraph.GraphEdge{
		{From: "exact-client", To: "api", Type: "http", Details: []archgraph.EntitySummary{{Name: "PATCH /items/${itemId}", Details: map[string]any{"evidence_origin": "deterministic", "repository_revision": map[string]any{"commit": "caller-head"}}}}},
		{From: "other-client", To: "api", Type: "http", Details: []archgraph.EntitySummary{{Name: "GET /other", Details: map[string]any{"evidence_origin": "deterministic", "repository_revision": map[string]any{"commit": "caller-head"}}}}},
	}}
	changed := []changedEntrypoint{{ID: "patch-item", Name: "PATCH /items/{id}", Match: "changed_line"}}

	got := exactChangedSurfaceCallers(graph, "api", changed)
	if len(got) != 1 || len(got["exact-client"]) != 1 {
		t.Fatalf("exact callers = %+v, want exact-client only", got)
	}
	if _, exists := got["other-client"]; exists {
		t.Fatalf("unrelated service was promoted to an exact caller: %+v", got)
	}
}

func TestPullRequestGraphFreshnessUsesSnapshotRevision(t *testing.T) {
	cases := []struct {
		name     string
		revision graphRevision
		head     string
		want     string
	}{
		{name: "matching", revision: graphRevision{Commit: "abc"}, head: "abc", want: "fresh"},
		{name: "different", revision: graphRevision{Commit: "abc"}, head: "def", want: "stale"},
		{name: "dirty snapshot", revision: graphRevision{Commit: "abc", Dirty: true}, head: "abc", want: "dirty"},
		{name: "missing revision", head: "abc", want: "unknown"},
		{name: "missing PR head", revision: graphRevision{Commit: "abc"}, want: "unknown"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := pullRequestGraphFreshness(tc.revision, tc.head); got != tc.want {
				t.Fatalf("freshness = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestChangedEntrypointsRequiresHeadCoordinates(t *testing.T) {
	// The diff's old line 10 is now an unchanged endpoint; the actual change
	// moved to line 30 after an earlier insertion.
	files := []githubPullFile{{Filename: "controller.go", Patch: "@@ -10 +30 @@\n-old\n+new"}}
	for _, tc := range []struct {
		name     string
		revision map[string]any
		want     string
		count    int
	}{
		{name: "head", revision: map[string]any{"commit": "head"}, want: "changed_line", count: 1},
		{name: "base snapshot", revision: map[string]any{"commit": "base"}, want: "file_scope", count: 2},
		{name: "unknown snapshot", want: "file_scope", count: 2},
		{name: "dirty head", revision: map[string]any{"commit": "head", "dirty": true}, want: "file_scope", count: 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			service := &archgraph.ServiceNode{HTTPRoutes: []archgraph.EntitySummary{
				{ID: "unrelated", Name: "GET /unrelated", Details: map[string]any{"source_locations": []model.Location{{File: "controller.go", StartLine: 10, EndLine: 10}}, "repository_revision": tc.revision}},
				{ID: "changed", Name: "GET /changed", Details: map[string]any{"source_locations": []model.Location{{File: "controller.go", StartLine: 30, EndLine: 30}}, "repository_revision": tc.revision}},
			}}
			got := changedEntrypoints(service, files, "head")
			if len(got) != tc.count {
				t.Fatalf("entrypoints = %+v, want %d", got, tc.count)
			}
			for _, entrypoint := range got {
				if entrypoint.Match != tc.want || (tc.want == "changed_line" && entrypoint.ID != "changed") {
					t.Fatalf("incorrect coordinate/revision match: %+v", entrypoint)
				}
			}
			if tc.want == "file_scope" {
				graph := &archgraph.ArchGraph{Services: []*archgraph.ServiceNode{{Name: "api"}, {Name: "client"}}, Edges: []*archgraph.GraphEdge{{From: "client", To: "api", Details: []archgraph.EntitySummary{{Name: "GET /changed"}}}}}
				if callers := exactChangedSurfaceCallers(graph, "api", got); len(callers) != 0 {
					t.Fatalf("unverified snapshot promoted callers: %+v", callers)
				}
			}
		})
	}
}

func TestChangedEntrypointsMissingPatchStaysFileScope(t *testing.T) {
	service := &archgraph.ServiceNode{HTTPRoutes: []archgraph.EntitySummary{{ID: "endpoint", Details: map[string]any{"source_locations": []model.Location{{File: "controller.go", StartLine: 10}}, "repository_revision": map[string]any{"commit": "head"}}}}}
	got := changedEntrypoints(service, []githubPullFile{{Filename: "controller.go"}}, "head")
	if len(got) != 1 || got[0].Match != "file_scope" {
		t.Fatalf("missing patch = %+v", got)
	}
}

func TestChangedEntrypointsMissingLineRangeStaysFileScope(t *testing.T) {
	service := &archgraph.ServiceNode{HTTPRoutes: []archgraph.EntitySummary{{ID: "endpoint", Details: map[string]any{"source_locations": []model.Location{{File: "controller.go"}}, "repository_revision": map[string]any{"commit": "head"}}}}}
	got := changedEntrypoints(service, []githubPullFile{{Filename: "controller.go", Patch: "@@ -10 +10 @@\n-old\n+new"}}, "head")
	if len(got) != 1 || got[0].Match != "file_scope" {
		t.Fatalf("missing line range = %+v", got)
	}
}

func TestChangedEntrypointsHeadDoesNotUseDeletedBaseLines(t *testing.T) {
	service := &archgraph.ServiceNode{HTTPRoutes: []archgraph.EntitySummary{{ID: "untouched", Details: map[string]any{"source_locations": []model.Location{{File: "controller.go", StartLine: 30}}, "repository_revision": map[string]any{"commit": "head"}}}}}
	got := changedEntrypoints(service, []githubPullFile{{Filename: "controller.go", Patch: "@@ -30 +29,0 @@\n-deleted"}}, "head")
	if len(got) != 0 {
		t.Fatalf("deleted base line matched head endpoint: %+v", got)
	}
}

func TestExactChangedSurfaceCallersPreservesPathIdentity(t *testing.T) {
	graph := &archgraph.ArchGraph{Services: []*archgraph.ServiceNode{{Name: "api"}, {Name: "exact"}, {Name: "case"}, {Name: "slash"}, {Name: "empty"}}, Edges: []*archgraph.GraphEdge{
		{From: "exact", To: "api", Details: []archgraph.EntitySummary{{Name: "get /Items/:id", Details: map[string]any{"evidence_origin": "deterministic", "repository_revision": map[string]any{"commit": "caller-head"}}}}},
		{From: "case", To: "api", Details: []archgraph.EntitySummary{{Name: "GET /items/{id}", Details: map[string]any{"evidence_origin": "deterministic", "repository_revision": map[string]any{"commit": "caller-head"}}}}},
		{From: "slash", To: "api", Details: []archgraph.EntitySummary{{Name: "GET /Items/{id}/", Details: map[string]any{"evidence_origin": "deterministic", "repository_revision": map[string]any{"commit": "caller-head"}}}}},
		{From: "empty", To: "api", Details: []archgraph.EntitySummary{{Name: "GET /elsewhere", Details: map[string]any{"evidence_origin": "deterministic", "repository_revision": map[string]any{"commit": "caller-head"}}}}},
	}}
	for _, service := range graph.Services[1:] {
		service.AnalysisStatus = &archgraph.RepositoryAnalysisStatus{State: "analyzed_clean", AnalyzedRevision: "caller-head"}
	}
	// Older artifacts may lack IDs; empty IDs must never establish a match.
	got := exactChangedSurfaceCallers(graph, "api", []changedEntrypoint{{Name: "GET /Items/{id}", Match: "changed_line"}})
	if len(got) != 1 || len(got["exact"]) != 1 {
		t.Fatalf("callers = %+v, want only exact", got)
	}
}

func TestDeclaredOrUnknownRelationshipsCannotBecomeExactPRCallers(t *testing.T) {
	for _, details := range []map[string]any{nil, {"pack_id": "approved-pack", "source_locations": []model.Location{{File: "app.go", StartLine: 1}}}, {"evidence_origin": "manual"}, {"evidence_origin": "llm"}, {"evidence_origin": "runtime"}} {
		if details == nil {
			details = map[string]any{}
		}
		details["repository_revision"] = map[string]any{"commit": "caller-head"}
		graph := &archgraph.ArchGraph{Services: []*archgraph.ServiceNode{{Name: "api"}, {Name: "caller", AnalysisStatus: &archgraph.RepositoryAnalysisStatus{State: "analyzed_clean", AnalyzedRevision: "caller-head"}}}, Edges: []*archgraph.GraphEdge{{From: "caller", To: "api", Details: []archgraph.EntitySummary{{Name: "GET /items", Details: details}}}}}
		got := exactChangedSurfaceCallers(graph, "api", []changedEntrypoint{{Name: "GET /items", Match: "changed_line"}})
		if len(got) != 0 {
			t.Fatalf("promoted origin: %+v", details)
		}
	}
}

func TestExactPRCallerRequiresItsOwnCleanSavedRevision(t *testing.T) {
	for _, tc := range []struct {
		name     string
		status   *archgraph.RepositoryAnalysisStatus
		revision map[string]any
		want     int
	}{
		{"clean matching caller", &archgraph.RepositoryAnalysisStatus{State: "analyzed_clean", AnalyzedRevision: "caller-head"}, map[string]any{"commit": "caller-head"}, 1},
		{"missing status", nil, map[string]any{"commit": "caller-head"}, 0},
		{"dirty analysis", &archgraph.RepositoryAnalysisStatus{State: "analyzed_dirty", AnalyzedRevision: "caller-head", Dirty: true}, map[string]any{"commit": "caller-head"}, 0},
		{"unknown fact revision", &archgraph.RepositoryAnalysisStatus{State: "analyzed_clean", AnalyzedRevision: "caller-head"}, nil, 0},
		{"old fact revision", &archgraph.RepositoryAnalysisStatus{State: "analyzed_clean", AnalyzedRevision: "caller-head"}, map[string]any{"commit": "older-caller"}, 0},
		{"dirty fact", &archgraph.RepositoryAnalysisStatus{State: "analyzed_clean", AnalyzedRevision: "caller-head"}, map[string]any{"commit": "caller-head", "dirty": true}, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			graph := &archgraph.ArchGraph{Services: []*archgraph.ServiceNode{{Name: "api"}, {Name: "caller", AnalysisStatus: tc.status}}, Edges: []*archgraph.GraphEdge{{From: "caller", To: "api", Details: []archgraph.EntitySummary{{Name: "GET /items", Details: map[string]any{"evidence_origin": "deterministic", "repository_revision": tc.revision}}}}}}
			got := exactChangedSurfaceCallers(graph, "api", []changedEntrypoint{{Name: "GET /items", Match: "changed_line"}})
			if len(got) != tc.want {
				t.Fatalf("callers = %+v, want %d", got, tc.want)
			}
		})
	}
}

func TestServiceGraphRevisionRejectsMixedSnapshots(t *testing.T) {
	service := &archgraph.ServiceNode{HTTPRoutes: []archgraph.EntitySummary{
		{Details: map[string]any{"repository_revision": map[string]any{"commit": "head"}}},
		{Details: map[string]any{"repository_revision": map[string]any{"commit": "old"}}},
	}}
	if got := pullRequestGraphFreshness(serviceGraphRevision(service), "head"); got != "unknown" {
		t.Fatalf("mixed revisions considered %s", got)
	}
	service.HTTPRoutes[1].Details["repository_revision"] = map[string]any{"commit": "head", "dirty": true}
	if got := pullRequestGraphFreshness(serviceGraphRevision(service), "head"); got != "dirty" {
		t.Fatalf("dirty revision considered %s", got)
	}
}

func TestCompanyImpactScoreExcludesUnverifiedGraphCandidates(t *testing.T) {
	for _, impact := range []companyImpact{
		{Available: true, PotentialServices: []impactedService{{Name: "candidate"}}},
		{Available: true, DirectServices: 1},
		{Available: true, ScoreEligible: true},
	} {
		if got := companyImpactScore(impact); got != 0 {
			t.Fatalf("score = %d, want 0 for unverified/empty evidence: %+v", got, impact)
		}
	}
	verified := companyImpact{Available: true, ScoreEligible: true, DirectServices: 1, Teams: []string{"team"}}
	if got := companyImpactScore(verified); got != 28 {
		t.Fatalf("score = %d, want 28 for verified evidence", got)
	}
}

func hasChangeCategory(categories []changeCategory, id string) bool {
	for _, category := range categories {
		if category.ID == id {
			return true
		}
	}
	return false
}

func hasSemanticSignal(signals []semanticSignal, kind string) bool {
	for _, signal := range signals {
		if signal.Kind == kind {
			return true
		}
	}
	return false
}

func containsExactString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
