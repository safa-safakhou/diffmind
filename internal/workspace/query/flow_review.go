package query

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/mohammad-safakhou/diffmind/internal/workspace/archgraph"
	"github.com/mohammad-safakhou/diffmind/internal/workspace/model"
	"github.com/mohammad-safakhou/diffmind/internal/workspace/store"
)

type FlowReviewResult struct {
	EnvironmentChanges []archgraph.EnvironmentChange `json:"environment_changes"`
	ProjectID          string                        `json:"project_id"`
	From               GraphRun                      `json:"from"`
	To                 GraphRun                      `json:"to"`
	Service            string                        `json:"service"`
	Flows              []archgraph.EntryFlowReview   `json:"flows"`
	Total              int                           `json:"total"`
	NextOffset         *int                          `json:"next_offset,omitempty"`
	InputsBefore       []SnapshotInput               `json:"inputs_before"`
	InputsAfter        []SnapshotInput               `json:"inputs_after"`
	Notes              []string                      `json:"notes"`
}

func (s *Service) CompareFlows(ctx context.Context, projectID, from, to, service string, offset, limit int) (*FlowReviewResult, error) {
	if strings.TrimSpace(from) == "" || strings.TrimSpace(to) == "" || strings.TrimSpace(service) == "" {
		return nil, fmt.Errorf("from, to and service are required")
	}
	if limit == 0 {
		limit = 10
	}
	if limit > 50 {
		return nil, fmt.Errorf("flow review limit must be between 1 and 50")
	}
	if _, _, _, err := pageBounds(offset, limit, 0); err != nil {
		return nil, err
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	p, err := s.ResolveProject(projectID)
	if err != nil {
		return nil, err
	}
	a, left, err := s.loadGraph(p.ID, from)
	if err != nil {
		return nil, err
	}
	b, right, err := s.loadGraph(p.ID, to)
	if err != nil {
		return nil, err
	}
	return CompareFlowGraphs(ctx, p.ID, a, b, left, right, service, offset, limit)
}

// CompareFlowGraphs shares comparison semantics with isolated PR snapshots.
// The supplied graphs and manifests are immutable captured inputs.
func CompareFlowGraphs(ctx context.Context, projectID string, a, b *store.RunManifest, left, right *archgraph.ArchGraph, service string, offset, limit int) (*FlowReviewResult, error) {
	if limit == 0 {
		limit = 10
	}
	if limit < 1 || limit > 50 || offset < 0 {
		return nil, fmt.Errorf("offset must be nonnegative and limit between 1 and 50")
	}
	found := false
	for _, g := range []*archgraph.ArchGraph{left, right} {
		for _, svc := range g.Services {
			if svc != nil && svc.Name == service {
				found = true
			}
		}
	}
	if !found {
		return nil, fmt.Errorf("%w: %s", ErrServiceNotFound, service)
	}
	entriesA, entriesB := archgraph.EntryFlows(left, service), archgraph.EntryFlows(right, service)
	var beforeMetrics, afterMetrics *model.RepoMetrics
	for _, svc := range left.Services {
		if svc != nil && svc.Name == service {
			beforeMetrics = svc.RepoMetrics
		}
	}
	for _, svc := range right.Services {
		if svc != nil && svc.Name == service {
			afterMetrics = svc.RepoMetrics
		}
	}
	keys := map[string]bool{}
	for k := range entriesA {
		keys[k] = true
	}
	for k := range entriesB {
		keys[k] = true
	}
	ordered := make([]string, 0, len(keys))
	for k := range keys {
		ordered = append(ordered, k)
	}
	sort.Slice(ordered, func(i, j int) bool {
		// New and removed entrypoints are useful starting points in a PR review.
		changed := func(key string) bool { return len(entriesA[key]) == 0 || len(entriesB[key]) == 0 }
		a, b := changed(ordered[i]), changed(ordered[j])
		if a != b {
			return a
		}
		return ordered[i] < ordered[j]
	})
	reviews := make([]archgraph.EntryFlowReview, 0, len(ordered))
	out := &FlowReviewResult{ProjectID: projectID, From: graphRun(*a), To: graphRun(*b), Service: service, Flows: []archgraph.EntryFlowReview{}, InputsBefore: snapshotInputs(a, left), InputsAfter: snapshotInputs(b, right), Notes: []string{
		"Saved execution evidence, not observed runtime traffic. Missing or truncated local flows do not prove absence of impact.",
		"Entrypoints pair by exact service, kind and name. Renames remain removal plus addition; ambiguous duplicate entrypoints cannot be compared.",
		"Each flow follows at most 4 service hops and 200 nodes; one-hop callers are included only when the saved edge matches the selected entrypoint.",
	}}
	out.From.GraphAvailable = true
	out.To.GraphAvailable = true
	for _, key := range ordered {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		aa, bb := entriesA[key], entriesB[key]
		if len(aa) > 1 || len(bb) > 1 {
			return nil, fmt.Errorf("ambiguous entrypoint identity %s; refine extracted entrypoint names", key)
		}
		var ar, br *archgraph.EntrypointRef
		if len(aa) > 0 {
			ar = &aa[0]
		}
		if len(bb) > 0 {
			br = &bb[0]
		}
		review := archgraph.ReviewEntryFlow(left, right, key, ar, br, archgraph.FlowOptions{Depth: 4, MaxNodes: 200, Expand: "steps"})
		if a.ID != b.ID && review.Change == "unchanged" {
			continue
		}
		changes := review.Changes[:0]
		for _, change := range review.Changes {
			if !change.EvidenceOnly {
				changes = append(changes, change)
			}
		}
		review.Changes = changes
		reviews = append(reviews, review)
	}
	start, end, next, _ := pageBounds(offset, limit, len(reviews))
	out.Total, out.NextOffset = len(reviews), next
	if beforeMetrics != nil && afterMetrics != nil && beforeMetrics.DependencyInventory != nil && afterMetrics.DependencyInventory != nil {
		out.EnvironmentChanges = archgraph.CompareEnvironments(beforeMetrics, afterMetrics)
	} else {
		out.Notes = append(out.Notes, "Dependency-version evidence is missing in at least one snapshot; compatibility changes cannot be assessed.")
	}
	if beforeMetrics != nil && afterMetrics != nil && beforeMetrics.DetectorRevision != afterMetrics.DetectorRevision {
		out.Notes = append(out.Notes, "Detector revisions differ; extraction changes cannot be attributed solely to repository changes.")
	}
	if len(out.EnvironmentChanges) > 0 {
		out.Notes = append(out.Notes, "Dependency declarations or selected versions changed. This is separate from structural flow changes and does not establish runtime compatibility.")
	}
	out.Flows = append(out.Flows, reviews[start:end]...)
	if out.From.PackSetDigest != out.To.PackSetDigest {
		out.Notes = append(out.Notes, "Knowledge-pack inputs differ; inspect evidence before attributing these changes to a PR.")
	}
	return out, nil
}
