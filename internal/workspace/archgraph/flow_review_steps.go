package archgraph

import (
	"encoding/json"
	"github.com/mohammad-safakhou/diffmind/protocol"
)

// Expanded steps are drawn only if the extracted DAG actually connects the
// selected exposure to its dependency. Never manufacture a sequential path
// from an unordered list of symbols.
func expandReviewSteps(flow *FlowView, service string, connections []ConnectionSummary, entities map[string]EntitySummary, budget int) {
	known := map[string]bool{}
	for _, n := range flow.Nodes {
		known[n.ID] = true
	}
	for _, c := range connections {
		var nodes []protocol.FlowNode
		var edges []protocol.FlowEdge
		decodeReview(c.Nodes, &nodes)
		decodeReview(c.Edges, &edges)
		ids := map[string]string{}
		extra := []FlowNode{}
		for _, n := range nodes {
			id := flowObjectNodeID(service, n.Ref)
			e, hasRef := entities[id]
			kind, label := n.Role, n.Symbol
			if hasRef {
				kind, label = e.Kind, e.Name
			} else {
				id = flowObjectNodeID(service, jsonKey([]string{n.Role, n.Symbol}))
				if label == "" {
					label = n.Ref
				}
				if label == "" {
					continue
				}
			}
			ids[n.ID] = id
			if !known[id] {
				extra = append(extra, FlowNode{ID: id, Service: service, Kind: firstNonEmpty(kind, "step"), Label: label})
			}
		}
		from, to := flowObjectNodeID(service, c.FromID), flowObjectNodeID(service, c.ToID)
		adjacency := map[string][]string{}
		for _, e := range edges {
			if ids[e.From] != "" && ids[e.To] != "" {
				adjacency[ids[e.From]] = append(adjacency[ids[e.From]], ids[e.To])
			}
		}
		seen := map[string]bool{from: true}
		queue := []string{from}
		for i := 0; i < len(queue); i++ {
			for _, next := range adjacency[queue[i]] {
				if !seen[next] {
					seen[next] = true
					queue = append(queue, next)
				}
			}
		}
		if !seen[to] {
			continue
		}
		if len(known)+len(extra) > budget {
			flow.Status = "truncated"
			flow.Quality = append(flow.Quality, "Expanded local steps exceed the diagram node budget; condensed connection retained.")
			continue
		}
		for _, n := range extra {
			if !known[n.ID] {
				flow.Nodes = append(flow.Nodes, n)
				known[n.ID] = true
			}
		}
		kept := flow.Edges[:0]
		for _, e := range flow.Edges {
			if e.From == from && e.To == to && e.MatchStatus == "local_flow" {
				continue
			}
			kept = append(kept, e)
		}
		flow.Edges = kept
		for _, e := range edges {
			if ids[e.From] == "" || ids[e.To] == "" {
				continue
			}
			var condition map[string]any
			decodeReview(e.Condition, &condition)
			flow.Edges = append(flow.Edges, FlowEdge{From: ids[e.From], To: ids[e.To], Kind: firstNonEmpty(e.Kind, "local"), Reachability: string(e.Reachability), Condition: condition, MatchStatus: "local_flow"})
		}
	}
}

func decodeReview(v, out any) { body, _ := json.Marshal(v); _ = json.Unmarshal(body, out) }

// Canonicalize generated DAG node IDs while retaining all recorded edge
// conditions, loop/fanout metadata and unknown fields in the evidence.
func reviewConnectionTopology(c ConnectionSummary) any {
	var nodes, edges []map[string]any
	decodeReview(c.Nodes, &nodes)
	decodeReview(c.Edges, &edges)
	ids := map[string]string{}
	for _, n := range nodes {
		id, _ := n["id"].(string)
		ref := n["ref"]
		if ref == c.FromID {
			ref = c.FromName
		}
		if ref == c.ToID {
			ref = c.ToName
		}
		identity := reviewJSON([]any{n["role"], n["symbol"]})
		if n["symbol"] == nil || n["symbol"] == "" {
			identity = reviewJSON([]any{n["role"], ref})
		}
		ids[id] = identity
		delete(n, "id")
		delete(n, "ref")
		n["identity"] = identity
	}
	for _, e := range edges {
		for _, key := range []string{"from", "to"} {
			if id, ok := e[key].(string); ok && ids[id] != "" {
				e[key] = ids[id]
			}
		}
	}
	nv, ev := []any{}, []any{}
	for _, n := range nodes {
		nv = append(nv, n)
	}
	for _, e := range edges {
		ev = append(ev, e)
	}
	return map[string]any{"nodes": sortedValues(nv), "edges": sortedValues(ev)}
}
