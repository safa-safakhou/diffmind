package query

import (
	"context"
	"fmt"
	"github.com/mohammad-safakhou/diffmind/internal/workspace/archgraph"
	"github.com/mohammad-safakhou/diffmind/internal/workspace/store"
	"testing"
)

func TestCompareFlowsPaginationAndImmutableDirection(t *testing.T) {
	q, pid, rid := testQueryService(t)
	graph := &archgraph.ArchGraph{RunID: rid, Services: []*archgraph.ServiceNode{{Name: "orders", Known: true}}}
	for i := 0; i < 13; i++ {
		graph.Services[0].HTTPRoutes = append(graph.Services[0].HTTPRoutes, archgraph.EntitySummary{ID: fmt.Sprintf("route-%02d", i), Name: fmt.Sprintf("GET /orders/%02d", i)})
	}
	saveGraph(t, q, pid, rid, graph)
	a, err := q.CompareFlows(context.Background(), pid, rid, rid, "orders", 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if a.Total != 13 || len(a.Flows) != 10 || a.NextOffset == nil || *a.NextOffset != 10 || a.From.ID != rid || a.To.ID != rid {
		t.Fatalf("page: %+v", a)
	}
	b, err := q.CompareFlows(context.Background(), pid, rid, rid, "orders", 10, 10)
	if err != nil || len(b.Flows) != 3 || b.NextOffset != nil {
		t.Fatalf("last page %+v %v", b, err)
	}
	for _, flow := range a.Flows {
		if flow.Change != "unchanged" || !flow.Partial {
			t.Fatal("missing local evidence overclaimed")
		}
	}
	for _, args := range []struct {
		from, to, service string
		offset, limit     int
	}{{rid, rid, "orders", -1, 1}, {rid, rid, "orders", 0, 51}, {rid, "", "orders", 0, 1}, {rid, rid, "missing", 0, 1}} {
		if _, err := q.CompareFlows(context.Background(), pid, args.from, args.to, args.service, args.offset, args.limit); err == nil {
			t.Fatalf("invalid accepted: %+v", args)
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := q.CompareFlows(ctx, pid, rid, rid, "orders", 0, 1); err != context.Canceled {
		t.Fatalf("cancel: %v", err)
	}
	graph.Services[0].HTTPRoutes = append(graph.Services[0].HTTPRoutes, graph.Services[0].HTTPRoutes[0])
	saveGraph(t, q, pid, rid, graph)
	if _, err := q.CompareFlows(context.Background(), pid, rid, rid, "orders", 0, 1); err == nil {
		t.Fatal("ambiguous duplicate silently paired")
	}
}

func TestFlowDiffFiltersUnchangedBeforePagination(t *testing.T) {
	left := &archgraph.ArchGraph{Services: []*archgraph.ServiceNode{{Name: "orders"}}}
	for i := 0; i < 13; i++ {
		left.Services[0].HTTPRoutes = append(left.Services[0].HTTPRoutes, archgraph.EntitySummary{ID: fmt.Sprint(i), Name: fmt.Sprintf("GET /%02d", i)})
	}
	right := &archgraph.ArchGraph{Services: []*archgraph.ServiceNode{{Name: "orders", HTTPRoutes: append([]archgraph.EntitySummary{}, left.Services[0].HTTPRoutes...)}}}
	for i := 13; i < 15; i++ {
		right.Services[0].HTTPRoutes = append(right.Services[0].HTTPRoutes, archgraph.EntitySummary{ID: fmt.Sprint(i), Name: fmt.Sprintf("GET /%02d", i)})
	}
	before, after := &store.RunManifest{ID: "before"}, &store.RunManifest{ID: "after"}
	page, err := CompareFlowGraphs(context.Background(), "project", before, after, left, right, "orders", 0, 1)
	if err != nil {
		t.Fatal(err)
	}
	if page.Total != 2 || len(page.Flows) != 1 || page.Flows[0].Change != "added" || page.NextOffset == nil || *page.NextOffset != 1 {
		t.Fatalf("changed page: %+v", page)
	}
	page, err = CompareFlowGraphs(context.Background(), "project", before, after, left, right, "orders", 1, 1)
	if err != nil || len(page.Flows) != 1 || page.NextOffset != nil {
		t.Fatalf("last changed page: %+v %v", page, err)
	}
	page, err = CompareFlowGraphs(context.Background(), "project", before, after, left, left, "orders", 0, 10)
	if err != nil || page.Total != 0 || len(page.Flows) != 0 {
		t.Fatalf("unchanged comparison: %+v %v", page, err)
	}
	page, err = CompareFlowGraphs(context.Background(), "project", before, before, left, left, "orders", 0, 10)
	if err != nil || page.Total != 13 || len(page.Flows) != 10 {
		t.Fatalf("service exploration: %+v %v", page, err)
	}
}
