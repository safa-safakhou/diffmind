package archgraph

import (
	"reflect"
	"strings"
)

// RelationshipEvidence explains provenance without asserting runtime traffic or
// PR revision eligibility. Historical facts without provenance remain unknown.
type RelationshipEvidence struct {
	Class         string `json:"class"`
	RunID         string `json:"run_id"`
	ObjectID      string `json:"object_id,omitempty"`
	Origin        string `json:"origin,omitempty"`
	PackID        string `json:"pack_id,omitempty"`
	PackVersion   string `json:"pack_version,omitempty"`
	Resolution    string `json:"resolution,omitempty"`
	Revision      any    `json:"revision,omitempty"`
	FileScope     any    `json:"file_scope,omitempty"`
	AnalysisScope any    `json:"analysis_scope,omitempty"`
	ScopeState    string `json:"scope_state"`
	Coverage      string `json:"coverage"`
	Limit         string `json:"limit"`
}

func DescribeRelationship(runID string, item EntitySummary) RelationshipEvidence {
	d := item.Details
	e := RelationshipEvidence{Class: "unknown", RunID: runID, ObjectID: item.ID,
		Origin: getString(d, "evidence_origin"), PackID: getString(d, "pack_id"), PackVersion: getString(d, "pack_version"),
		Resolution: getString(d, "resolution_reason"), Revision: d["repository_revision"], FileScope: d["source_locations"],
		ScopeState: "unknown", Coverage: "unverified", Limit: "Saved evidence is incomplete; it does not establish runtime traffic, absence of dependencies or PR-head eligibility."}
	if e.FileScope == nil {
		e.FileScope = d["locations"]
	}
	switch {
	case e.PackID != "" || getString(d, "match_basis") == "knowledge_pack" || strings.HasPrefix(getString(d, "plugin_source"), "knowledge_pack:"):
		e.Class = "pack_declared"
	case e.Origin == "manual" || e.Origin == "imported":
		e.Class = "declared"
	case e.Origin == "llm":
		e.Class = "source_inferred"
	case e.Origin == "runtime":
		// Imported runtime origins are labels, not validated traffic observations.
		e.Class = "runtime_origin_unverified"
	case e.Origin == "deterministic" || hasEvidenceItems(e.FileScope) || hasEvidenceItems(d["evidence"]) || getString(d, "match_basis") == "terraform_subscription":
		e.Class = "source_extracted"
	}
	return e
}

func hasEvidenceItems(value any) bool {
	if value == nil {
		return false
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map:
		return v.Len() > 0
	}
	return false
}

// WithRelationshipEvidence copies the projection, never rewriting a persisted
// snapshot or mutating a shared cached graph. Each linked object's origin is
// retained, so a mixed relationship does not inherit its strongest signal.
func WithRelationshipEvidence(g *ArchGraph) *ArchGraph {
	if g == nil {
		return nil
	}
	out := *g
	out.Edges = make([]*GraphEdge, len(g.Edges))
	scopes := map[string]any{}
	for _, service := range g.Services {
		if service != nil && service.AnalysisStatus != nil {
			scopes[service.Name] = service.AnalysisStatus.FileScope
		}
	}
	for i, edge := range g.Edges {
		if edge == nil {
			continue
		}
		copy := *edge
		copy.Evidence = make([]RelationshipEvidence, 0, max(1, len(edge.Details)))
		for _, item := range edge.Details {
			copy.Evidence = append(copy.Evidence, DescribeRelationship(g.RunID, item))
		}
		if len(copy.Evidence) == 0 {
			copy.Evidence = append(copy.Evidence, DescribeRelationship(g.RunID, EntitySummary{}))
		}
		for j := range copy.Evidence {
			if scope := scopes[edge.From]; scope != nil {
				copy.Evidence[j].AnalysisScope = scope
				copy.Evidence[j].ScopeState = "recorded"
			}
		}
		out.Edges[i] = &copy
	}
	return &out
}
