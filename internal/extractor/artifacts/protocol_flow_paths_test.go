package artifacts

import (
	"reflect"
	"testing"

	"github.com/mohammad-safakhou/diffmind/internal/extractor/model"
	"github.com/mohammad-safakhou/diffmind/protocol"
)

func TestProtocolFlowKeepsBranchesIndependentOfPathOrder(t *testing.T) {
	cond := model.Condition{Kind: "unconditional"}
	paths := []model.ConnectionPath{
		{Condition: cond, Steps: []model.ConnectionPathStep{{To: "Search.start", Condition: cond}, {To: "Flight.boost", Condition: cond}}},
		{Condition: cond, Steps: []model.ConnectionPathStep{{To: "Search.start", Condition: cond}, {To: "Topic.boost", Condition: cond}}},
	}
	build := func(paths []model.ConnectionPath) protocol.Flow {
		b := &protocolBuilder{doc: &protocol.Document{}, oldToNew: map[string]string{"entry": "http.get_search", "db": "dbq.read_facts"}, flowedFrom: map[string]bool{}}
		b.addFlow(model.Connection{FromExposureID: "entry", ToDependencyID: "db", Condition: cond, Paths: paths})
		return b.doc.Flows[0]
	}
	a := build(paths)
	b := build([]model.ConnectionPath{paths[1], paths[0], paths[1]})
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("path order/duplicates changed the flow:\n%+v\n%+v", a, b)
	}
	seen := map[string]bool{}
	for _, n := range a.Nodes {
		seen[n.Symbol] = true
	}
	if len(a.Nodes) != 6 || len(a.Edges) != 5 || !seen["Flight.boost"] || !seen["Topic.boost"] {
		t.Fatalf("lost a branch or duplicated the shared prefix: %+v", a)
	}
	changed := append([]model.ConnectionPath(nil), paths...)
	changed = append(changed, model.ConnectionPath{Condition: cond, Steps: []model.ConnectionPathStep{{To: "Other.write", Condition: cond}}})
	if reflect.DeepEqual(a, build(changed)) {
		t.Fatal("a real path addition must change the flow")
	}
}
