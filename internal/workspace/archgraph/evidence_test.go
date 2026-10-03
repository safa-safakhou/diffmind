package archgraph

import (
	"encoding/json"
	"testing"
)

func TestRelationshipOriginsRemainSeparateAndReadOnly(t *testing.T) {
	g := &ArchGraph{RunID: "saved", Services: []*ServiceNode{{Name: "a", AnalysisStatus: &RepositoryAnalysisStatus{FileScope: map[string]any{"exclude": []string{"examples/**"}}}}}, Edges: []*GraphEdge{{From: "a", To: "unresolved", Details: []EntitySummary{
		{ID: "declared", Details: map[string]any{"pack_id": "company", "pack_version": "1", "source_locations": []any{map[string]any{"file": "app.go"}}, "repository_revision": map[string]any{"commit": "old", "dirty": true}}},
		{ID: "resolved", Details: map[string]any{"evidence_origin": "deterministic", "resolution_reason": "approved_alias"}},
		{ID: "unknown"}, {ID: "llm", Details: map[string]any{"evidence_origin": "llm"}},
	}}, nil, {From: "a", To: "b"}}}
	before, _ := json.Marshal(g)
	out := WithRelationshipEvidence(g)
	want := []string{"pack_declared", "source_extracted", "unknown", "source_inferred"}
	for i, class := range want {
		e := out.Edges[0].Evidence[i]
		if e.Class != class || e.RunID != "saved" || e.Coverage != "unverified" || e.ScopeState != "recorded" {
			t.Fatalf("origin %d: %+v", i, e)
		}
	}
	if out.Edges[0].Evidence[1].Resolution != "approved_alias" || out.Edges[2].Evidence[0].Class != "unknown" {
		t.Fatal("lost identity resolution or unknown origin")
	}
	after, _ := json.Marshal(g)
	if string(before) != string(after) {
		t.Fatal("mutated saved graph")
	}
	overview := Overview(out)
	if len(overview.Edges[0].Evidence) != 4 || len(overview.Edges[0].Details) != 0 {
		t.Fatal("overview lost origin or retained full object detail")
	}
}

func TestSummaryRetainsLegacyOriginWithoutMutatingInput(t *testing.T) {
	input := map[string]any{"name": "declared", "plugin_source": "knowledge_pack:org@1/rule", "locations": []any{map[string]any{"file": "source.go"}}}
	summary := toSummary(input)
	if DescribeRelationship("r", summary).Class != "pack_declared" {
		t.Fatal("lost legacy declaration provenance")
	}
	if input["details"] != nil {
		t.Fatal("mutated input")
	}
}
