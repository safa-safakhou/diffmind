package archgraph

import (
	"encoding/json"
	"sort"
	"strings"
)

// FlowSnapshot retains diagram topology and the source evidence behind it.
// IDs in the diagram are semantic identities so regenerated extractor IDs do
// not make identical snapshots look different.
type FlowSnapshot struct {
	Entry       EntrypointRef       `json:"entry"`
	Graph       *FlowView           `json:"graph"`
	Connections []ConnectionSummary `json:"connections"`
	Callers     []*GraphEdge        `json:"callers"`
}
type FlowChange struct {
	Kind         string `json:"kind"`
	Key          string `json:"key"`
	Label        string `json:"label"`
	Change       string `json:"change"`
	EvidenceOnly bool   `json:"evidence_only,omitempty"`
	Before       any    `json:"before,omitempty"`
	After        any    `json:"after,omitempty"`
}
type EntryFlowReview struct {
	Key     string         `json:"key"`
	Service string         `json:"service"`
	Name    string         `json:"name"`
	Kind    string         `json:"kind"`
	Change  string         `json:"change"`
	Before  *FlowSnapshot  `json:"before"`
	After   *FlowSnapshot  `json:"after"`
	Changes []FlowChange   `json:"changes"`
	Counts  map[string]int `json:"counts"`
	Partial bool           `json:"partial"`
}

// EntryFlows groups entrypoint identities; ambiguous duplicates are rejected
// by the caller rather than silently choosing an arbitrary occurrence.
func EntryFlows(g *ArchGraph, service string) map[string][]EntrypointRef {
	out := map[string][]EntrypointRef{}
	for _, svc := range g.Services {
		if svc == nil || svc.Name != service {
			continue
		}
		groups := []struct {
			kind  string
			items []EntitySummary
		}{
			{"http_endpoint", svc.HTTPRoutes}, {"rpc_endpoint", svc.RPCEndpoints},
			{"queue_consumer", svc.QueueConsumers}, {"scheduled_job", svc.ScheduledJobs},
			{"webhook", svc.Webhooks}, {"cli_command", svc.CLICommands},
		}
		for _, group := range groups {
			for _, item := range group.items {
				name := firstNonEmpty(item.Name, item.ID)
				key := jsonKey([]string{svc.Name, group.kind, name})
				out[key] = append(out[key], EntrypointRef{Service: svc.Name, Team: svc.Team, ID: firstNonEmpty(item.ID, item.Name), Kind: group.kind, Name: name})
			}
		}
	}
	return out
}

func ReviewEntryFlow(before, after *ArchGraph, key string, a, b *EntrypointRef, opts FlowOptions) EntryFlowReview {
	entry := a
	if b != nil {
		entry = b
	}
	out := EntryFlowReview{Key: key, Service: entry.Service, Name: entry.Name, Kind: entry.Kind, Change: "unchanged", Changes: []FlowChange{}, Counts: map[string]int{"added": 0, "removed": 0, "modified": 0}}
	if a != nil {
		out.Before = buildFlowSnapshot(before, *a, opts)
	}
	if b != nil {
		out.After = buildFlowSnapshot(after, *b, opts)
	}
	left, right := flowSnapshotFacts(out.Before), flowSnapshotFacts(out.After)
	keys := map[string]bool{}
	for k := range left {
		keys[k] = true
	}
	for k := range right {
		keys[k] = true
	}
	for k := range keys {
		av, aok := left[k]
		bv, bok := right[k]
		c := FlowChange{Key: k, Kind: firstNonEmpty(bv.Kind, av.Kind), Label: firstNonEmpty(bv.Label, av.Label), Before: av.After, After: bv.After}
		switch {
		case !aok:
			c.Change = "added"
		case !bok:
			c.Change = "removed"
		case reviewJSON(av.After) != reviewJSON(bv.After):
			c.Change = "modified"
		default:
			continue
		}
		if c.Change == "modified" {
			c.EvidenceOnly = reviewJSON(withoutProvenance(c.Before)) == reviewJSON(withoutProvenance(c.After))
		}
		out.Changes = append(out.Changes, c)
		if !c.EvidenceOnly {
			out.Counts[c.Change]++
		}
	}
	sort.Slice(out.Changes, func(i, j int) bool { return out.Changes[i].Key < out.Changes[j].Key })
	switch {
	case a == nil:
		out.Change = "added"
	case b == nil:
		out.Change = "removed"
	case out.Counts["added"]+out.Counts["removed"]+out.Counts["modified"] > 0:
		out.Change = "modified"
	}
	for _, snap := range []*FlowSnapshot{out.Before, out.After} {
		if snap != nil && snap.Graph.Status != "complete" {
			out.Partial = true
		}
	}
	return out
}

func buildFlowSnapshot(g *ArchGraph, entry EntrypointRef, opts FlowOptions) *FlowSnapshot {
	flow, _ := BuildFlowView(g, entry.Service, entry.ID, opts)
	out := &FlowSnapshot{Entry: entry, Graph: flow, Connections: []ConnectionSummary{}, Callers: []*GraphEdge{}}
	svc := serviceMap(g)[entry.Service]
	// Only local evidence matched to this entrypoint belongs in this report.
	out.Connections = append(out.Connections, matchingConnections(svc, entry.ID)...)
	entities := map[string]EntitySummary{}
	for _, s := range g.Services {
		if s == nil {
			continue
		}
		for _, items := range [][]EntitySummary{s.HTTPRoutes, s.RPCEndpoints, s.QueueConsumers, s.ScheduledJobs, s.Webhooks, s.CLICommands, s.Dependencies} {
			for _, e := range items {
				entities[flowObjectNodeID(s.Name, firstNonEmpty(e.ID, normalizeLookup(e.Name)))] = e
			}
		}
	}
	// A declaration without an extracted local flow still deserves a visible
	// entrypoint, while the underlying flow remains explicitly partial.
	entryNodeID := flowObjectNodeID(entry.Service, entry.ID)
	present := false
	for _, n := range flow.Nodes {
		if n.ID == entryNodeID {
			present = true
		}
	}
	if !present {
		flow.Nodes = append(flow.Nodes, FlowNode{ID: entryNodeID, Service: entry.Service, Kind: entry.Kind, Label: entry.Name})
	}
	expandReviewSteps(flow, entry.Service, out.Connections, entities, opts.MaxNodes)
	ids := map[string]string{}
	for i := range flow.Nodes {
		n := &flow.Nodes[i]
		if e, ok := entities[n.ID]; ok {
			n.Details = map[string]any{"summary": e.Summary}
			for key, value := range e.Details {
				n.Details[key] = value
			}
			if e.Kind != "" {
				n.Kind = e.Kind
			}
		}
		if n.ID == entryNodeID {
			n.Kind = entry.Kind
		}
		old := n.ID
		n.ID = jsonKey([]string{n.Service, n.Kind, n.Label})
		ids[old] = n.ID
	}
	valid := flow.Edges[:0]
	for _, e := range flow.Edges {
		if ids[e.From] == "" || ids[e.To] == "" {
			continue
		}
		e.From = ids[e.From]
		e.To = ids[e.To]
		valid = append(valid, e)
	}
	flow.Edges = valid
	sort.Slice(flow.Nodes, func(i, j int) bool { return flow.Nodes[i].ID < flow.Nodes[j].ID })
	sort.Slice(flow.Edges, func(i, j int) bool { return reviewJSON(flow.Edges[i]) < reviewJSON(flow.Edges[j]) })
	flow.Stats.NodeCount = len(flow.Nodes)
	flow.Stats.EdgeCount = len(flow.Edges)
	for _, edge := range g.Edges {
		if edge == nil || edge.To != entry.Service {
			continue
		}
		match, status := matchServiceEntrypoint(edge, svc)
		if edge.Type == "queue_consume" {
			match, status = matchQueueConsumer(edge.From, svc)
		}
		if match == entry.ID && strings.HasPrefix(status, "exact_") {
			out.Callers = append(out.Callers, edge)
		}
	}
	sort.Slice(out.Callers, func(i, j int) bool { return reviewJSON(out.Callers[i]) < reviewJSON(out.Callers[j]) })
	return out
}

func flowSnapshotFacts(s *FlowSnapshot) map[string]FlowChange {
	out := map[string]FlowChange{}
	if s == nil {
		return out
	}
	// Group repeated identities as multisets so no evidence is discarded.
	values := map[string][]any{}
	labels := map[string]string{}
	kinds := map[string]string{}
	add := func(kind, key, label string, value any) {
		id := kind + ":" + key
		values[id] = append(values[id], value)
		labels[id] = label
		kinds[id] = kind
	}
	for _, n := range s.Graph.Nodes {
		if n.Kind == "service" && n.Service == s.Entry.Service {
			continue // Context, not an entrypoint step being added or removed.
		}
		add("node", n.ID, n.Label, map[string]any{"kind": n.Kind, "label": n.Label, "service": n.Service, "details": n.Details})
	}
	for _, e := range s.Graph.Edges {
		add("edge", jsonKey([]string{e.From, e.To, e.Kind}), e.Kind, e)
	}
	for _, d := range s.Graph.DataDependencies {
		add("data", jsonKey([]string{d.Service, d.From, d.To}), d.From+" → "+d.To, d.Dependencies)
	}
	for _, c := range s.Connections {
		add("connection", jsonKey([]string{c.FromName, c.ToName, c.FromType, c.ToType}), firstNonEmpty(c.ToName, c.ToID), map[string]any{"summary": c.Summary, "kind": c.Kind, "reachability": c.Reachability, "condition": c.Condition, "data_dependencies": c.DataDependencies, "side_effects": c.SideEffects, "topology": reviewConnectionTopology(c)})
	}
	for _, e := range s.Callers {
		details := []any{}
		for _, d := range e.Details {
			details = append(details, objectFact(d))
		}
		add("caller", jsonKey([]string{e.From, e.To, e.Type}), e.From, map[string]any{"from": e.From, "to": e.To, "type": e.Type, "details": sortedValues(details), "evidence": e.Evidence})
	}
	for key, v := range values {
		out[key] = FlowChange{Kind: kinds[key], Label: labels[key], After: sortedValues(v)}
	}
	return out
}
func reviewJSON(v any) string { b, _ := json.Marshal(v); return string(b) }
