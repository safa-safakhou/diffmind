package ast_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/mohammad-safakhou/diffmind/internal/extractor/ast"
)

func TestJavaCallbacksKeepReceiverAndDistinctCallSites(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "CampaignService.java", `class CampaignService {
    FlightService flights;
    void search(Stream stream) {
        stream.peek(this::assignBoost);
        stream.map(flights::assignBoost);
        combine(this::assignBoost, flights::assignBoost);
        assignBoost();
    }
    void assignBoost() { campaignDatabase.read(); }
    void combine(Object a, Object b) {}
}`)
	writeFile(t, dir, "FlightService.java", `class FlightService {
    void assignBoost() { flightDatabase.read(); }
}`)
	for run := 0; run < 8; run++ {
		idx, err := ast.Build(context.Background(), dir, "java", 2)
		if err != nil {
			t.Fatal(err)
		}
		calls := idx.CallGraph["CampaignService.search"]
		counts := map[string]int{}
		for _, call := range calls {
			counts[call.CalleeRaw]++
			want := []string{call.CalleeRaw}
			switch call.CalleeRaw {
			case "this.assignBoost", "assignBoost":
				want = []string{"CampaignService.assignBoost"}
			case "flights.assignBoost":
				want = []string{"FlightService.assignBoost"}
			case "combine":
				want = []string{"CampaignService.combine"}
			}
			if !reflect.DeepEqual(call.CalleeResolved, want) {
				t.Fatalf("run %d: %s resolved to %v, want %v", run, call.CalleeRaw, call.CalleeResolved, want)
			}
		}
		for raw, want := range map[string]int{"stream.peek": 1, "stream.map": 1, "this.assignBoost": 2, "flights.assignBoost": 2, "combine": 1, "assignBoost": 1} {
			if counts[raw] != want {
				t.Fatalf("run %d: counts = %v; %s want %d", run, counts, raw, want)
			}
		}
	}
}
