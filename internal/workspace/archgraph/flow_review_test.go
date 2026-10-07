package archgraph

import (
	"encoding/json"
	"github.com/mohammad-safakhou/diffmind/protocol"
	"strings"
	"testing"
)

func reviewFixture() *ArchGraph {
	route := EntitySummary{ID: "route", Name: "PUT /orders/{id}", Kind: "http_endpoint", Details: map[string]any{"method": "PUT", "path": "/orders/{id}", "source_locations": []any{map[string]any{"file": "OrdersController.java", "start_line": 20}}}}
	write := EntitySummary{ID: "write", Name: "Update orders", Kind: "db_operation", Details: map[string]any{"table": "orders", "operation": "UPDATE"}}
	return &ArchGraph{RunID: "before", Services: []*ServiceNode{{Name: "orders", Known: true, HTTPRoutes: []EntitySummary{route}, Dependencies: []EntitySummary{write}, Connections: []ConnectionSummary{{FromID: route.ID, FromName: route.Name, FromType: route.Kind, ToID: write.ID, ToName: write.Name, ToType: write.Kind, EntrypointID: route.ID, Summary: "Update orders", Nodes: []protocol.FlowNode{{ID: "n1", Ref: "route", Role: "entrypoint"}, {ID: "n2", Symbol: "OrdersService.update", Role: "handler"}, {ID: "n3", Ref: "write", Role: "dependency"}}, Edges: []protocol.FlowEdge{{From: "n1", To: "n2", Kind: "call"}, {From: "n2", To: "n3", Kind: "write"}}}}}, {Name: "storefront", Known: true}}, Edges: []*GraphEdge{{From: "storefront", To: "orders", Type: "http", Details: []EntitySummary{{Name: "update", Details: map[string]any{"method": "PUT", "path": "/orders/{id}"}}}}}}
}
func reviewed(t *testing.T, a, b *ArchGraph) EntryFlowReview {
	t.Helper()
	left, right := EntryFlows(a, "orders"), EntryFlows(b, "orders")
	for key, aa := range left {
		bb := right[key]
		var br *EntrypointRef
		if len(bb) > 0 {
			br = &bb[0]
		}
		return ReviewEntryFlow(a, b, key, &aa[0], br, FlowOptions{Depth: 4, MaxNodes: 200})
	}
	t.Fatal("missing entry")
	return EntryFlowReview{}
}
func TestFlowReviewRetainsStepsOperationsCallersAndOriginalEvidence(t *testing.T) {
	a := reviewFixture()
	before, _ := json.Marshal(a)
	r := reviewed(t, a, cloneComparison(a))
	if r.Change != "unchanged" || r.Partial {
		t.Fatalf("self: %+v", r)
	}
	if len(r.Before.Callers) != 1 {
		t.Fatal("exact caller missing")
	}
	data, _ := json.Marshal(r.Before)
	for _, want := range []string{"OrdersService.update", "UPDATE", "OrdersController.java"} {
		if !strings.Contains(string(data), want) {
			t.Fatalf("missing %s", want)
		}
	}
	if len(r.Before.Graph.Nodes) != 4 {
		t.Fatalf("expanded nodes: %+v", r.Before.Graph.Nodes)
	}
	after, _ := json.Marshal(a)
	if string(before) != string(after) {
		t.Fatal("mutated original snapshot")
	}
}
func TestFlowReviewDetectsNewDatabaseWriteAndDeletedEntrypoint(t *testing.T) {
	a := reviewFixture()
	b := cloneComparison(a)
	svc := b.Services[0]
	svc.Dependencies = append(svc.Dependencies, EntitySummary{ID: "audit", Name: "Write audit_log", Kind: "db_operation", Details: map[string]any{"table": "audit_log", "operation": "INSERT"}})
	c := svc.Connections[0]
	c.ToID = "audit"
	c.ToName = "Write audit_log"
	c.Nodes = nil
	c.Edges = nil
	svc.Connections = append(svc.Connections, c)
	r := reviewed(t, a, b)
	if r.Change != "modified" || r.Counts["added"] == 0 {
		t.Fatalf("new write: %+v", r)
	}
	found := false
	for _, c := range r.Changes {
		if c.Kind == "node" && c.Label == "Write audit_log" && c.Change == "added" {
			found = true
		}
	}
	if !found {
		t.Fatal("missing added write")
	}
	b.Services[0].HTTPRoutes = nil
	r = reviewed(t, a, b)
	if r.Change != "removed" || r.After != nil || len(r.Before.Callers) != 1 {
		t.Fatal("removed surface callers lost")
	}
}
func TestFlowReviewIgnoresRegeneratedObjectAndDAGIDs(t *testing.T) {
	a := reviewFixture()
	b := cloneComparison(a)
	svc := b.Services[0]
	svc.HTTPRoutes[0].ID = "new-route"
	svc.Dependencies[0].ID = "new-write"
	c := &svc.Connections[0]
	c.FromID = "new-route"
	c.ToID = "new-write"
	c.EntrypointID = "new-route"
	c.Nodes = []protocol.FlowNode{{ID: "new-1", Ref: "new-route", Role: "entrypoint"}, {ID: "new-2", Symbol: "OrdersService.update", Role: "handler"}, {ID: "new-3", Ref: "new-write", Role: "dependency"}}
	c.Edges = []protocol.FlowEdge{{From: "new-1", To: "new-2", Kind: "call"}, {From: "new-2", To: "new-3", Kind: "write"}}
	r := reviewed(t, a, b)
	if r.Change != "unchanged" {
		t.Fatalf("generated identity noise: %+v", r.Changes)
	}
}
func TestFlowReviewDoesNotInventOrderAndMarksMissingEvidence(t *testing.T) {
	a := reviewFixture()
	a.Services[0].Connections[0].Edges = nil
	r := reviewed(t, a, cloneComparison(a))
	for _, n := range r.Before.Graph.Nodes {
		if n.Label == "OrdersService.update" {
			t.Fatal("unordered steps presented as path")
		}
	}
	a.Services[0].Connections = nil
	r = reviewed(t, a, cloneComparison(a))
	if !r.Partial {
		t.Fatal("missing local flow not marked partial")
	}
	if len(r.Before.Graph.Nodes) < 2 {
		t.Fatal("declaration missing from diagram")
	}
}
func TestFlowReviewConditionsAndEvidenceOnlyChanges(t *testing.T) {
	a := reviewFixture()
	b := cloneComparison(a)
	b.Services[0].HTTPRoutes[0].Details["source_locations"] = []any{map[string]any{"file": "OrdersController.java", "start_line": 99}}
	r := reviewed(t, a, b)
	if len(r.Changes) != 1 || !r.Changes[0].EvidenceOnly {
		t.Fatalf("evidence classified incorrectly: %+v", r.Changes)
	}
	b.Services[0].Connections[0].Condition = map[string]any{"summary": "only if confirmed"}
	r = reviewed(t, a, b)
	found := false
	for _, c := range r.Changes {
		if c.Kind == "connection" && !c.EvidenceOnly {
			found = true
		}
	}
	if !found {
		t.Fatal("condition change lost")
	}
}

func TestFlowReviewProvenanceChangesDoNotMarkBehaviorChanged(t *testing.T) {
	a := reviewFixture()
	b := cloneComparison(a)
	b.Services[0].HTTPRoutes[0].Details["repository_revision"] = map[string]any{"commit": "new-head"}
	a.Services[0].HTTPRoutes[0].Details["metadata"] = map[string]any{"legacy_id": "baseline-id"}
	b.Services[0].HTTPRoutes[0].Details["metadata"] = map[string]any{"legacy_id": "capture-specific-id"}
	b.Services[0].HTTPRoutes[0].Details["source_locations"] = []any{map[string]any{"file": "OrdersController.java", "start_line": 99}}
	r := reviewed(t, a, b)
	if r.Change != "unchanged" || r.Counts["modified"] != 0 || len(r.Changes) == 0 || !r.Changes[0].EvidenceOnly {
		t.Fatalf("provenance presented as flow change: %+v", r)
	}
}

func TestFlowReviewMissingLocalEvidenceDoesNotBorrowRepositoryDependencies(t *testing.T) {
	a := reviewFixture()
	a.Services[0].Connections = nil
	a.Edges = append(a.Edges, &GraphEdge{From: "orders", To: "storefront", Type: "http", Details: []EntitySummary{{Name: "GET /unrelated"}}})
	r := reviewed(t, a, cloneComparison(a))
	for _, node := range r.Before.Graph.Nodes {
		if node.Service == "storefront" {
			t.Fatal("unrelated service presented as entrypoint flow")
		}
	}
	if !r.Partial {
		t.Fatal("missing entrypoint evidence must be partial")
	}
}
