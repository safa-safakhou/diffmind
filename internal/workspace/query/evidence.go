package query

// EvidenceState distinguishes a saved graph from the freshness of repository
// analyses. Fresh analyses do not imply that the selected graph incorporates them.
type EvidenceState struct {
	FreshnessReference string `json:"freshness_reference"`
	FreshnessBasis     string `json:"freshness_basis"`
	GraphState         string `json:"graph_state"`
	SavedRunID         string `json:"saved_run_id,omitempty"`
	RepositoryCount    int    `json:"repository_count"`
	AnalysisFreshness  string `json:"analysis_freshness"`
	Fresh              int    `json:"fresh"`
	Stale              int    `json:"stale"`
	Dirty              int    `json:"dirty"`
	Unknown            int    `json:"unknown"`
	Basis              string `json:"basis"`
	Coverage           string `json:"coverage"`
}

func DescribeEvidence(runID string, freshness []string) EvidenceState {
	state := EvidenceState{FreshnessReference: "latest_repository_analysis", GraphState: "missing", FreshnessBasis: "stored_repository_status", SavedRunID: runID, RepositoryCount: len(freshness), AnalysisFreshness: "unknown", Basis: "static_source", Coverage: "unverified"}
	if runID != "" {
		state.GraphState = "saved"
	}
	for _, value := range freshness {
		switch value {
		case "fresh":
			state.Fresh++
		case "stale":
			state.Stale++
		case "dirty":
			state.Dirty++
		default:
			state.Unknown++
		}
	}
	if state.Stale > 0 || state.Dirty > 0 {
		state.AnalysisFreshness = "needs_update"
	} else if state.Unknown == 0 && state.Fresh > 0 {
		state.AnalysisFreshness = "fresh"
	}
	return state
}
