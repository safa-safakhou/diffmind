package query

import "testing"

func TestEvidenceNeverTurnsMissingOrUnknownCoverageIntoReady(t *testing.T) {
	for _, tc := range []struct {
		run              string
		input            []string
		graph, freshness string
	}{
		{"", nil, "missing", "unknown"}, {"saved", []string{"fresh"}, "saved", "fresh"}, {"saved", []string{"fresh", "unknown"}, "saved", "unknown"}, {"saved", []string{"fresh", "stale"}, "saved", "needs_update"}, {"saved", []string{"dirty"}, "saved", "needs_update"},
	} {
		got := DescribeEvidence(tc.run, tc.input)
		if got.GraphState != tc.graph || got.AnalysisFreshness != tc.freshness || got.Coverage != "unverified" {
			t.Fatalf("%+v => %+v", tc, got)
		}
	}
}
