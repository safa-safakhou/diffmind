package query

import (
	"errors"
	"time"

	"github.com/mohammad-safakhou/diffmind/internal/workspace/store"
)

// Readiness is an observation, not an authorization token or a coverage claim.
// Work and source freshness are independent of the saved graph's validity.
type Readiness struct {
	ProjectID         string             `json:"project_id"`
	ObservedAt        time.Time          `json:"observed_at"`
	Runtime           string             `json:"runtime"`
	Maintenance       string             `json:"maintenance"`
	ConnectionMode    string             `json:"connection_mode"`
	Access            string             `json:"access"`
	Graph             string             `json:"graph"`
	SavedRunID        string             `json:"saved_run_id,omitempty"`
	SavedAt           *time.Time         `json:"saved_at,omitempty"`
	Sources           []ReadinessSource  `json:"sources"`
	Inputs            []store.RunRepoRef `json:"inputs"`
	Evidence          EvidenceState      `json:"evidence"`
	PRHeadEligibility string             `json:"pr_head_eligibility"`
	Work              WorkState          `json:"work"`
	Actions           ReadinessActions   `json:"actions"`
	NextAction        string             `json:"next_action"`
	Limits            []string           `json:"limits"`
}

// Revisions are the last stored observations; a readiness read never fetches
// provider heads. Analysis IDs and immutable graph inputs remain independent.
type ReadinessSource struct {
	RepositoryID           string `json:"repository_id"`
	LastAnalysisID         string `json:"last_analysis_id,omitempty"`
	LastObservedHead       string `json:"last_observed_head,omitempty"`
	LastObservedRemoteHead string `json:"last_observed_remote_head,omitempty"`
	RevisionBasis          string `json:"revision_basis"`
	AnalysisFreshness      string `json:"analysis_freshness"`
}

// UnavailableReadiness is a safe transport/access projection. It contains no
// project identifiers or saved provenance after authority can no longer be checked.
func UnavailableReadiness(runtime, access string) *Readiness {
	next := "reconnect"
	if access == "denied" {
		next = "request_access"
	}
	evidence := DescribeEvidence("", nil)
	evidence.GraphState = "unknown"
	return &Readiness{ObservedAt: time.Now().UTC(), Runtime: runtime, Maintenance: "unknown", ConnectionMode: "unknown", Access: access, Graph: "unknown", Sources: []ReadinessSource{}, Inputs: []store.RunRepoRef{}, Evidence: evidence, PRHeadEligibility: "unknown", Work: WorkState{Kind: "unknown", Status: "unknown", Phase: "unknown"}, NextAction: next, Limits: []string{"Recheck connection and access before querying or changing context."}}
}

type WorkState struct {
	Kind      string    `json:"kind"`
	ID        string    `json:"id,omitempty"`
	Status    string    `json:"status"`
	Phase     string    `json:"phase"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ReadinessActions struct {
	Query       bool `json:"query"`
	Refresh     bool `json:"refresh"`
	Configure   bool `json:"configure"`
	InspectWork bool `json:"inspect_work"`
}

// WithReadinessPermissions returns a request-scoped copy; default query-only
// connections never advertise mutation authority, even on a trusted local store.
func (s *Service) WithReadinessPermissions(resolve func(string) (ReadinessActions, error)) *Service {
	copy := *s
	copy.readinessPermissions = resolve
	return &copy
}

// WithReadinessGraphLoader reuses the host's validated immutable artifact cache.
// The loader must return only a completed run with a readable, matching graph.
func (s *Service) WithReadinessGraphLoader(load func(string) (*store.RunManifest, error)) *Service {
	copy := *s
	copy.readinessGraph = load
	return &copy
}

func (s *Service) Readiness(projectID string) (*Readiness, error) {
	p, err := s.ResolveProject(projectID)
	if err != nil {
		return nil, err
	}
	out := &Readiness{ProjectID: p.ID, ObservedAt: time.Now().UTC(), Runtime: "available", Maintenance: "unknown", ConnectionMode: "query_only", Access: "allowed", Graph: "missing", Inputs: []store.RunRepoRef{}, Sources: []ReadinessSource{}, PRHeadEligibility: "unknown", Work: WorkState{Kind: "none", Status: "empty", Phase: "idle"}, Limits: []string{"Static source evidence; dependency coverage remains unverified.", "Checkout freshness compares the latest repository analysis, not the saved graph or a PR head."}}
	permissions := ReadinessActions{}
	if s.readinessPermissions != nil {
		permissions, err = s.readinessPermissions(p.ID)
		if err != nil {
			return nil, err
		}
		out.Maintenance, out.ConnectionMode = "available", "managed"
	}
	repos, err := s.currentRepositories(p.ID)
	if err != nil {
		return nil, err
	}
	freshness := make([]string, 0, len(repos))
	for _, repo := range repos {
		freshness = append(freshness, repo.DiffMindFreshness)
		out.Sources = append(out.Sources, ReadinessSource{RepositoryID: repo.ID, LastAnalysisID: repo.LastDiffMindRunID, LastObservedHead: repo.HeadSHA, LastObservedRemoteHead: repo.RemoteHeadSHA, RevisionBasis: "stored_repository_status", AnalysisFreshness: repo.DiffMindFreshness})
	}
	var run *store.RunManifest
	var graphErr error
	if s.readinessGraph != nil {
		run, graphErr = s.readinessGraph(p.ID)
	} else {
		run, _, graphErr = s.loadGraph(p.ID, "")
	}
	if graphErr == nil {
		out.Graph, out.SavedRunID, out.Inputs = "queryable", run.ID, run.Repos
		if out.Inputs == nil {
			out.Inputs = []store.RunRepoRef{}
		}
		if !run.FinishedAt.IsZero() {
			stamp := run.FinishedAt
			out.SavedAt = &stamp
		}
	} else if !errors.Is(graphErr, ErrNoCompletedGraph) {
		// A damaged or inaccessible artifact is not an empty workspace.
		out.Graph = "unavailable"
		out.Limits = append(out.Limits, "Saved graph could not be read; inspect the graph run and restore or rebuild it.")
	}
	out.Evidence = DescribeEvidence(out.SavedRunID, freshness)
	out.Evidence.FreshnessBasis = "live_checkout_status"
	if out.Graph == "unavailable" {
		out.Evidence.GraphState = "unavailable"
	}
	ingestion, err := s.store.GetIngestion(p.ID)
	if err != nil && !errors.Is(err, store.ErrNotFound) {
		return nil, err
	}
	if ingestion != nil {
		out.Work = WorkState{Kind: "ingestion", ID: ingestion.ID, Status: ingestion.Status, Phase: ingestion.Phase, UpdatedAt: ingestion.UpdatedAt}
	}
	runs, err := s.store.ListRuns(p.ID)
	if err != nil {
		return nil, err
	}
	for _, r := range runs {
		stamp := r.FinishedAt
		if stamp.IsZero() {
			stamp = r.StartedAt
		}
		candidate := WorkState{Kind: "graph", ID: r.ID, Status: r.Status, Phase: "graph", UpdatedAt: stamp}
		out.Work = selectWork(out.Work, candidate)
	}
	jobs, err := s.store.ListJobs(p.ID)
	if err != nil {
		return nil, err
	}
	for _, job := range jobs {
		status := job.Status
		if status == "succeeded" {
			status = "completed"
		}
		candidate := WorkState{Kind: "refresh_job", ID: job.ID, Status: status, Phase: job.Status, UpdatedAt: job.UpdatedAt}
		// A job's current ingestion has more useful phase detail than its envelope.
		if ingestion != nil && ingestion.JobID == job.ID {
			candidate.Phase = ingestion.Phase
			if job.Status == "succeeded" && ingestion.Status == "partial" {
				candidate.Status = "partial"
			}
		}
		out.Work = selectWork(out.Work, candidate)
	}
	for _, repo := range repos {
		if repo.SyncStatus == "diffmind_running" {
			out.Work = selectWork(out.Work, WorkState{Kind: "repository_analysis", ID: repo.ID, Status: "running", Phase: "analyzing", UpdatedAt: repo.UpdatedAt})
		}
	}
	active := workActive(out.Work.Status)
	out.Actions = ReadinessActions{Query: out.Graph == "queryable", InspectWork: out.Work.Kind != "none", Refresh: permissions.Refresh && len(repos) > 0 && !active, Configure: permissions.Configure && !active}
	switch {
	case active:
		out.NextAction = "inspect_work"
	case out.Work.Status == "partial" || out.Work.Status == "failed" || out.Work.Status == "cancelled" || out.Work.Status == "interrupted":
		out.NextAction = "inspect_work"
	case out.Graph == "unavailable":
		out.NextAction = "inspect_work"
	case out.Graph == "queryable":
		out.NextAction = "query_saved_graph"
	case out.Actions.Refresh:
		out.NextAction = "refresh"
	case out.Actions.Configure:
		out.NextAction = "review_import"
	default:
		out.NextAction = "request_editor"
	}
	return out, nil
}

func workActive(status string) bool {
	return status == "running" || status == "queued" || status == "cancelling"
}

func selectWork(current, candidate WorkState) WorkState {
	// Active work takes priority over any completed/failed history. Among active
	// work, executing work takes priority over waiting work; time breaks ties.
	rank := func(w WorkState) int {
		if w.Status == "running" || w.Status == "cancelling" {
			return 2
		}
		if w.Status == "queued" {
			return 1
		}
		return 0
	}
	if rank(candidate) > rank(current) || (rank(candidate) == rank(current) && candidate.UpdatedAt.After(current.UpdatedAt)) {
		return candidate
	}
	return current
}
