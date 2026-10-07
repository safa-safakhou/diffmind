package ui

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mohammad-safakhou/diffmind/internal/workspace/archgraph"
	"github.com/mohammad-safakhou/diffmind/internal/workspace/model"
	"github.com/mohammad-safakhou/diffmind/internal/workspace/store"
)

// Pull requests are intentionally a live projection rather than persisted
// project state. The architecture graph is the stable company model; GitHub is
// queried for the changing review queue and a selected PR is overlaid on that
// graph on demand.
type pullRequestSummary struct {
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	Draft     bool      `json:"draft"`
	Author    string    `json:"author"`
	HeadSHA   string    `json:"head_sha,omitempty"`
	Head      string    `json:"head"`
	Base      string    `json:"base"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Labels    []string  `json:"labels,omitempty"`
	RepoID    string    `json:"repo_id"`
	RepoName  string    `json:"repo_name"`
	Team      string    `json:"team,omitempty"`
}

type pullRequestRepo struct {
	RepoID       string               `json:"repo_id"`
	RepoName     string               `json:"repo_name"`
	Team         string               `json:"team,omitempty"`
	Provider     string               `json:"provider"`
	Status       string               `json:"status"`
	Message      string               `json:"message,omitempty"`
	Error        string               `json:"error,omitempty"`
	OpenCount    int                  `json:"open_count"`
	Truncated    bool                 `json:"truncated,omitempty"`
	PullRequests []pullRequestSummary `json:"pull_requests"`
}

type pullRequestsResponse struct {
	CheckedCount     int               `json:"checked_count"`
	UnavailableCount int               `json:"unavailable_count"`
	TotalOpen        int               `json:"total_open"`
	RepoCount        int               `json:"repo_count"`
	ErrorCount       int               `json:"error_count"`
	GeneratedAt      time.Time         `json:"generated_at"`
	Repositories     []pullRequestRepo `json:"repositories"`
}

type githubPull struct {
	Number  int    `json:"number"`
	Title   string `json:"title"`
	HTMLURL string `json:"html_url"`
	Draft   bool   `json:"draft"`
	User    struct {
		Login string `json:"login"`
	} `json:"user"`
	Head struct {
		Ref string `json:"ref"`
		SHA string `json:"sha"`
	} `json:"head"`
	Base struct {
		SHA string `json:"sha"`
		Ref string `json:"ref"`
	} `json:"base"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Labels    []struct {
		Name string `json:"name"`
	} `json:"labels"`
	Additions      int    `json:"additions"`
	Deletions      int    `json:"deletions"`
	ChangedFiles   int    `json:"changed_files"`
	Commits        int    `json:"commits"`
	Comments       int    `json:"comments"`
	ReviewComments int    `json:"review_comments"`
	MergeableState string `json:"mergeable_state"`
}

type githubPullFile struct {
	Filename  string `json:"filename"`
	Status    string `json:"status"`
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
	Changes   int    `json:"changes"`
	Patch     string `json:"patch"`
	BlobURL   string `json:"blob_url"`
}

type codeImpactFile struct {
	Path      string `json:"path"`
	Status    string `json:"status"`
	Category  string `json:"category"`
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
	URL       string `json:"url,omitempty"`
}

type changeCategory struct {
	ID        string `json:"id"`
	Label     string `json:"label"`
	Files     int    `json:"files"`
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
	Risk      int    `json:"risk"`
}

type semanticSignal struct {
	Kind     string   `json:"kind"`
	Label    string   `json:"label"`
	Severity string   `json:"severity"`
	Files    []string `json:"files"`
}

type codebaseImpact struct {
	ChangedFiles int              `json:"changed_files"`
	Additions    int              `json:"additions"`
	Deletions    int              `json:"deletions"`
	Commits      int              `json:"commits"`
	RiskScore    int              `json:"risk_score"`
	RiskLevel    string           `json:"risk_level"`
	RiskReasons  []string         `json:"risk_reasons"`
	Categories   []changeCategory `json:"categories"`
	Signals      []semanticSignal `json:"signals,omitempty"`
	Areas        []string         `json:"areas,omitempty"`
	Files        []codeImpactFile `json:"files"`
	Truncated    bool             `json:"truncated,omitempty"`
}

type impactedService struct {
	Name     string   `json:"name"`
	Team     string   `json:"team,omitempty"`
	Depth    int      `json:"depth"`
	Tier     string   `json:"tier"`
	Reason   string   `json:"reason"`
	Evidence []string `json:"evidence,omitempty"`
}

type changedEntrypoint struct {
	ID    string `json:"id"`
	Kind  string `json:"kind"`
	Name  string `json:"name"`
	File  string `json:"file"`
	Match string `json:"match"`
}

type entrypointDependency struct {
	Entrypoint   string `json:"entrypoint"`
	Target       string `json:"target"`
	Kind         string `json:"kind"`
	Reachability string `json:"reachability,omitempty"`
}

type graphRevision struct {
	Commit string `json:"commit,omitempty"`
	Branch string `json:"branch,omitempty"`
	Dirty  bool   `json:"dirty,omitempty"`
}

type companyImpact struct {
	Available              bool                   `json:"available"`
	RootService            string                 `json:"root_service,omitempty"`
	RunID                  string                 `json:"run_id,omitempty"`
	DirectServices         int                    `json:"direct_services"`
	IndirectServices       int                    `json:"indirect_services"`
	Teams                  []string               `json:"teams,omitempty"`
	Resources              []string               `json:"resources,omitempty"`
	Services               []impactedService      `json:"services,omitempty"`
	PotentialServices      []impactedService      `json:"potential_services,omitempty"`
	PotentialResources     []string               `json:"potential_resources,omitempty"`
	ChangedEntrypoints     []changedEntrypoint    `json:"changed_entrypoints,omitempty"`
	EntrypointDependencies []entrypointDependency `json:"entrypoint_dependencies,omitempty"`
	Flow                   *archgraph.FlowView    `json:"flow,omitempty"`
	Confidence             string                 `json:"confidence"`
	Freshness              string                 `json:"freshness"`
	GraphRevision          graphRevision          `json:"graph_revision,omitempty"`
	ScoreEligible          bool                   `json:"score_eligible"`
	Notes                  []string               `json:"notes,omitempty"`
	NextAction             string                 `json:"next_action"`
	Limitations            []string               `json:"limitations"`
}

type pullRequestImpactResponse struct {
	PullRequest  pullRequestSummary `json:"pull_request"`
	Codebase     codebaseImpact     `json:"codebase"`
	Company      companyImpact      `json:"company"`
	RiskScore    int                `json:"risk_score"`
	RiskLevel    string             `json:"risk_level"`
	GeneratedAt  time.Time          `json:"generated_at"`
	ScoreMeaning string             `json:"score_meaning"`
	Delivery     string             `json:"delivery"`
}

func (s *Server) handlePullRequests(w http.ResponseWriter, r *http.Request) {
	pid := r.PathValue("pid")
	repos, err := s.workspaceRepos(pid)
	if err != nil {
		s.writeStoreErr(w, err)
		return
	}
	results := make([]pullRequestRepo, len(repos))
	sem := make(chan struct{}, 6)
	var wg sync.WaitGroup
	for i := range repos {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			results[i] = githubOpenPullRequests(r.Context(), repos[i])
		}()
	}
	wg.Wait()
	sort.Slice(results, func(i, j int) bool {
		if results[i].OpenCount != results[j].OpenCount {
			return results[i].OpenCount > results[j].OpenCount
		}
		return results[i].RepoName < results[j].RepoName
	})
	response := pullRequestsResponse{RepoCount: len(results), GeneratedAt: time.Now().UTC(), Repositories: results}
	for _, result := range results {
		response.TotalOpen += result.OpenCount
		if result.Status == "ok" {
			response.CheckedCount++
		} else if result.Status != "error" {
			response.UnavailableCount++
		}
		if result.Status == "error" {
			response.ErrorCount++
		}
	}
	writeJSON(w, http.StatusOK, response)
}

func githubOpenPullRequests(ctx context.Context, repo workspaceRepo) pullRequestRepo {
	result := pullRequestRepo{
		RepoID: repo.ID, RepoName: repo.Name, Team: repo.EffectiveTeam,
		Provider: firstNonEmpty(repo.GitProvider, repo.SourceType, "git"), Status: "unavailable",
		PullRequests: []pullRequestSummary{},
	}
	base, state, message := githubRepositoryEndpoint(ctx, repo.Repo)
	if state != "ready" {
		result.Status = state
		result.Message = message
		return result
	}
	result.Provider = "github"
	client := githubHTTPClient(20 * time.Second)
	token := githubAPIToken(ctx, base)
	for page := 1; page <= 10; page++ {
		endpoint := fmt.Sprintf("%s/pulls?state=open&per_page=100&page=%d&sort=updated&direction=desc", base, page)
		var pulls []githubPull
		if err := githubJSON(ctx, client, token, endpoint, &pulls); err != nil {
			result.Status = "error"
			result.Error = err.Error()
			return result
		}
		for _, pull := range pulls {
			summary := pullSummary(pull, repo.Repo)
			summary.Team = repo.EffectiveTeam
			result.PullRequests = append(result.PullRequests, summary)
		}
		if len(pulls) < 100 {
			break
		}
		if page == 10 {
			result.Truncated = true
		}
	}
	result.OpenCount = len(result.PullRequests)
	result.Status = "ok"
	return result
}

func (s *Server) handlePullRequestImpact(w http.ResponseWriter, r *http.Request) {
	pid, repoID := r.PathValue("pid"), r.PathValue("repo_id")
	number, err := strconv.Atoi(r.PathValue("number"))
	if err != nil || number <= 0 {
		writeErr(w, http.StatusBadRequest, errors.New("invalid pull request number"))
		return
	}
	repo, err := s.store.GetRepo(pid, repoID)
	if err != nil {
		s.writeStoreErr(w, err)
		return
	}
	endpoint, state, message := githubRepositoryEndpoint(r.Context(), *repo)
	if state != "ready" {
		writeErr(w, http.StatusBadRequest, errors.New(message))
		return
	}
	client := githubHTTPClient(30 * time.Second)
	token := githubAPIToken(r.Context(), endpoint)
	base := fmt.Sprintf("%s/pulls/%d", endpoint, number)
	var pull githubPull
	if err := githubJSON(r.Context(), client, token, base, &pull); err != nil {
		writeErr(w, http.StatusBadGateway, err)
		return
	}
	files, truncated, err := githubPullFiles(r.Context(), client, token, base)
	if err != nil {
		writeErr(w, http.StatusBadGateway, err)
		return
	}
	codebase := analyzeCodebaseImpact(pull, files, truncated)
	company := s.pullRequestCompanyImpact(pid, r.URL.Query().Get("run_id"), *repo, pull, files)
	companyScore := companyImpactScore(company)
	overall := minInt(100, int(math.Round(float64(codebase.RiskScore)*0.62+float64(companyScore)*0.38)))
	if !company.Available || !company.ScoreEligible {
		overall = codebase.RiskScore
	}
	writeJSON(w, http.StatusOK, pullRequestImpactResponse{
		PullRequest: pullSummary(pull, *repo), Codebase: codebase, Company: company,
		RiskScore: overall, RiskLevel: riskLevel(overall), GeneratedAt: time.Now().UTC(),
		ScoreMeaning: "Uncalibrated attention heuristic from file categories and eligible changed-surface evidence; not a probability, merge recommendation or proof of safety. Ineligible company context is excluded.",
		Delivery:     "On-demand inspection only. No automatic code-host comments, checks or merge decisions are delivered.",
	})
}

func githubPullFiles(ctx context.Context, client *http.Client, token, base string) ([]githubPullFile, bool, error) {
	var out []githubPullFile
	for page := 1; page <= 30; page++ {
		var files []githubPullFile
		endpoint := fmt.Sprintf("%s/files?per_page=100&page=%d", base, page)
		if err := githubJSON(ctx, client, token, endpoint, &files); err != nil {
			return nil, false, err
		}
		out = append(out, files...)
		if len(files) < 100 {
			return out, false, nil
		}
	}
	return out, true, nil
}

func githubJSON(ctx context.Context, client *http.Client, token, endpoint string, dst any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := client.Do(req)
	if err != nil {
		return githubConnectionError()
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return githubResponseError(resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(dst)
}

func githubSourceForRepo(ctx context.Context, repo store.Repo) string {
	if source := strings.TrimSpace(repo.GitURL); source != "" {
		return source
	}
	path := firstNonEmpty(repo.Path, repo.ClonePath)
	if path != "" {
		if remote := gitOutput(ctx, path, "remote", "get-url", "origin"); remote != "" {
			return remote
		}
	}
	return path
}

func pullSummary(p githubPull, repo store.Repo) pullRequestSummary {
	labels := make([]string, 0, len(p.Labels))
	for _, label := range p.Labels {
		if strings.TrimSpace(label.Name) != "" {
			labels = append(labels, label.Name)
		}
	}
	sort.Strings(labels)
	return pullRequestSummary{
		Number: p.Number, Title: p.Title, URL: p.HTMLURL, Draft: p.Draft, Author: p.User.Login,
		Head: p.Head.Ref, HeadSHA: p.Head.SHA, Base: p.Base.Ref, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
		Labels: labels, RepoID: repo.ID, RepoName: repo.Name, Team: repo.Team,
	}
}

var changeCategoryMeta = map[string]struct {
	label string
	risk  int
}{
	"api":            {"API & contracts", 20},
	"data":           {"Data & migrations", 25},
	"security":       {"Security & identity", 25},
	"infrastructure": {"Infrastructure", 18},
	"dependencies":   {"Dependencies", 15},
	"configuration":  {"Configuration", 10},
	"tests":          {"Tests", 0},
	"documentation":  {"Documentation", 0},
	"code":           {"Application code", 8},
}

func analyzeCodebaseImpact(pull githubPull, files []githubPullFile, truncated bool) codebaseImpact {
	impact := codebaseImpact{
		ChangedFiles: pull.ChangedFiles, Additions: pull.Additions, Deletions: pull.Deletions, Commits: pull.Commits,
		Files: make([]codeImpactFile, 0, len(files)), Truncated: truncated,
	}
	if impact.ChangedFiles == 0 {
		impact.ChangedFiles = len(files)
	}
	if impact.Additions == 0 && impact.Deletions == 0 {
		for _, f := range files {
			impact.Additions += f.Additions
			impact.Deletions += f.Deletions
		}
	}
	byCategory := map[string]*changeCategory{}
	areaCounts := map[string]int{}
	signalFiles := map[string]map[string]bool{}
	hasSource, hasTests := false, false
	for _, file := range files {
		category := classifyChangedFile(file.Filename)
		meta := changeCategoryMeta[category]
		cat := byCategory[category]
		if cat == nil {
			cat = &changeCategory{ID: category, Label: meta.label, Risk: meta.risk}
			byCategory[category] = cat
		}
		cat.Files++
		cat.Additions += file.Additions
		cat.Deletions += file.Deletions
		impact.Files = append(impact.Files, codeImpactFile{Path: file.Filename, Status: file.Status, Category: category, Additions: file.Additions, Deletions: file.Deletions, URL: file.BlobURL})
		areaCounts[changeArea(file.Filename)]++
		if category == "tests" {
			hasTests = true
		}
		if category == "code" || category == "api" || category == "data" {
			hasSource = true
		}
		for _, signal := range fileSignals(file) {
			if signalFiles[signal] == nil {
				signalFiles[signal] = map[string]bool{}
			}
			signalFiles[signal][file.Filename] = true
		}
	}
	for _, cat := range byCategory {
		impact.Categories = append(impact.Categories, *cat)
	}
	sort.Slice(impact.Categories, func(i, j int) bool {
		if impact.Categories[i].Risk != impact.Categories[j].Risk {
			return impact.Categories[i].Risk > impact.Categories[j].Risk
		}
		return impact.Categories[i].Files > impact.Categories[j].Files
	})
	for area := range areaCounts {
		impact.Areas = append(impact.Areas, area)
	}
	sort.Slice(impact.Areas, func(i, j int) bool {
		if areaCounts[impact.Areas[i]] != areaCounts[impact.Areas[j]] {
			return areaCounts[impact.Areas[i]] > areaCounts[impact.Areas[j]]
		}
		return impact.Areas[i] < impact.Areas[j]
	})
	if len(impact.Areas) > 8 {
		impact.Areas = impact.Areas[:8]
	}
	for kind, paths := range signalFiles {
		files := sortedKeys(paths)
		label, severity := signalLabel(kind)
		impact.Signals = append(impact.Signals, semanticSignal{Kind: kind, Label: label, Severity: severity, Files: files})
	}
	sort.Slice(impact.Signals, func(i, j int) bool {
		return signalRank(impact.Signals[i].Severity) > signalRank(impact.Signals[j].Severity)
	})

	score := 5
	lines := impact.Additions + impact.Deletions
	switch {
	case lines > 1500:
		score += 24
	case lines > 500:
		score += 16
	case lines > 100:
		score += 8
	}
	switch {
	case impact.ChangedFiles > 50:
		score += 20
	case impact.ChangedFiles > 20:
		score += 12
	case impact.ChangedFiles > 8:
		score += 5
	}
	for _, cat := range impact.Categories {
		if cat.Risk > 0 {
			score += minInt(cat.Risk, 8+cat.Files*3)
		}
	}
	if hasSource && !hasTests {
		score += 10
		impact.RiskReasons = append(impact.RiskReasons, "production code changed without test-file changes")
	}
	if impact.Deletions > impact.Additions && impact.Deletions > 100 {
		score += 6
		impact.RiskReasons = append(impact.RiskReasons, "large removal surface")
	}
	for _, signal := range impact.Signals {
		if signal.Severity == "high" {
			score += 8
		}
		impact.RiskReasons = append(impact.RiskReasons, signal.Label)
	}
	if lines > 500 {
		impact.RiskReasons = append(impact.RiskReasons, fmt.Sprintf("large diff: %d changed lines", lines))
	}
	if impact.ChangedFiles > 20 {
		impact.RiskReasons = append(impact.RiskReasons, fmt.Sprintf("broad change across %d files", impact.ChangedFiles))
	}
	impact.RiskScore = minInt(100, score)
	impact.RiskLevel = riskLevel(impact.RiskScore)
	if len(impact.RiskReasons) == 0 {
		impact.RiskReasons = []string{"localized change with no high-risk file patterns detected"}
	}
	sort.Slice(impact.Files, func(i, j int) bool {
		return impact.Files[i].Additions+impact.Files[i].Deletions > impact.Files[j].Additions+impact.Files[j].Deletions
	})
	return impact
}

func classifyChangedFile(path string) string {
	lower := strings.ToLower(filepath.ToSlash(path))
	base := filepath.Base(lower)
	ext := strings.ToLower(filepath.Ext(base))
	contains := func(parts ...string) bool {
		for _, p := range parts {
			if strings.Contains(lower, p) {
				return true
			}
		}
		return false
	}
	if contains("/test/", "/tests/", "_test.", ".spec.", ".test.", "__tests__") {
		return "tests"
	}
	if ext == ".md" || ext == ".rst" || contains("/docs/") {
		return "documentation"
	}
	switch base {
	case ".gitleaksignore", ".gitleaks.toml", "gitleaks.toml":
		return "configuration"
	case "go.mod", "go.sum", "package.json", "package-lock.json", "pnpm-lock.yaml", "yarn.lock", "pom.xml", "build.gradle", "build.gradle.kts", "requirements.txt", "pyproject.toml", "poetry.lock", "uv.lock", "pipfile", "pipfile.lock", "cargo.toml", "cargo.lock", ".tool-versions", ".python-version", ".node-version", ".nvmrc", "gradle-wrapper.properties", "maven-wrapper.properties":
		return "dependencies"
	}
	// Deployment settings for a migration component are configuration, not SQL changes.
	if base == "configuration.yaml" || base == "configuration.yml" {
		return "configuration"
	}
	if contains("migration", "schema.sql", "liquibase", "flyway", "prisma/schema", "db/schema") {
		return "data"
	}
	if contains("openapi", "swagger", ".proto", "graphql", "/routes/", "/controller", "/api/") {
		return "api"
	}
	if contains("auth", "security", "permission", "policy", "iam", "rbac", "oauth", "secret", "credential", "certificate", "tls") {
		return "security"
	}
	if ext == ".tf" || contains("terraform", "helm/", "charts/", "k8s/", "kubernetes", "dockerfile", "docker-compose", ".github/workflows") {
		return "infrastructure"
	}
	if ext == ".yaml" || ext == ".yml" || ext == ".json" || ext == ".toml" || ext == ".ini" || ext == ".conf" || ext == ".properties" || strings.HasPrefix(base, ".env") {
		return "configuration"
	}
	return "code"
}

func fileSignals(file githubPullFile) []string {
	category := classifyChangedFile(file.Filename)
	set := map[string]bool{}
	if category == "api" {
		set["contract_change"] = true
	}
	if category == "data" {
		set["data_migration"] = true
	}
	if category == "security" {
		set["security_boundary"] = true
	}
	if category == "infrastructure" {
		set["deployment_change"] = true
	}
	if category == "dependencies" {
		set["dependency_change"] = true
	}
	patch := strings.ToLower(file.Patch)
	if strings.Contains(patch, "-func ") || strings.Contains(patch, "-public ") || strings.Contains(patch, "-export ") || strings.Contains(patch, "-interface ") {
		set["public_surface_removal"] = true
	}
	if strings.Contains(patch, "drop table") || strings.Contains(patch, "drop column") || strings.Contains(patch, "alter table") {
		set["destructive_schema"] = true
	}
	return sortedKeys(set)
}

func signalLabel(kind string) (string, string) {
	switch kind {
	case "contract_change":
		return "API or contract surface changed", "high"
	case "data_migration":
		return "Database schema or migration changed", "high"
	case "destructive_schema":
		return "Potentially destructive schema operation", "high"
	case "security_boundary":
		return "Security or authorization boundary changed", "high"
	case "public_surface_removal":
		return "Public symbol removal candidate", "high"
	case "deployment_change":
		return "Deployment or infrastructure changed", "medium"
	case "dependency_change":
		return "Dependency or tool version inputs changed", "medium"
	default:
		return kind, "low"
	}
}

func (s *Server) pullRequestCompanyImpact(pid, requestedRun string, repo store.Repo, pull githubPull, files []githubPullFile) companyImpact {
	result := companyImpact{Confidence: "unavailable", Freshness: "unknown", Notes: []string{}, NextAction: "Inspect available PR files; select saved evidence from a matching clean PR-head revision for exact caller checks.", Limitations: []string{
		"Removed entrypoints may be absent from the head graph. No matching merge-base evidence is used here, so deleted-surface callers remain unproven.",
		"Internal, transitive and configuration changes may affect unchanged routes without intersecting their source locations; file-scope matches are candidates only.",
		"No exact matches is incomplete evidence, not proof that a PR is safe to merge. Extracted request-field compatibility is a separate saved-snapshot comparison.",
	}}
	for _, file := range files {
		if classifyChangedFile(file.Filename) == "dependencies" {
			result.Limitations = append(result.Limitations, "Dependency or tool versions changed. Unchanged extracted flows do not verify framework compatibility, serialization, or runtime behavior; inherited and transitive versions may be unresolved.")
			break
		}
	}
	runID := strings.TrimSpace(requestedRun)
	if runID == "" {
		if run := s.latestCompletedWorkspaceRun(pid); run != nil {
			runID = run.ID
		}
	}
	if runID == "" {
		result.Notes = append(result.Notes, "build a company graph to calculate cross-service impact")
		return result
	}
	graph, err := s.fullArchGraphForRun(pid, runID)
	if err != nil {
		result.Notes = append(result.Notes, "company graph is unavailable: "+err.Error())
		return result
	}

	if requestedRun == "" {
		if capture, _, headGraph, captureErr := s.loadPRCapture(pid, repo.ID, "", pull.Head.SHA); captureErr == nil {
			graph = headGraph
			runID = capture.After.ID
			result.Notes = append(result.Notes, "PR-head evidence comes from an isolated prepared snapshot; company context is held at "+capture.ContextRun)
		}
	}
	root := graphServiceForRepo(graph, repo)
	if root == "" {
		result.RunID = runID
		result.Notes = append(result.Notes, "repository is not represented as a service in this graph run")
		return result
	}
	flow, ok := archgraph.BuildImpactView(graph, root, archgraph.FlowOptions{Depth: 8, MaxNodes: 750})
	if !ok {
		return result
	}
	rootNode := graphService(graph, root)
	result.Available, result.RootService, result.RunID, result.Flow = true, root, runID, flow
	result.GraphRevision = serviceGraphRevision(rootNode)
	result.Freshness = pullRequestGraphFreshness(result.GraphRevision, pull.Head.SHA)
	result.ChangedEntrypoints = changedEntrypoints(rootNode, files, pull.Head.SHA)
	result.EntrypointDependencies = dependenciesForChangedEntrypoints(rootNode, result.ChangedEntrypoints)
	exactCallers := exactChangedSurfaceCallers(graph, root, result.ChangedEntrypoints)
	teamSet, potentialTeamSet, resourceSet := map[string]bool{}, map[string]bool{}, map[string]bool{}
	for _, service := range flow.Services {
		if service.Name == root {
			continue
		}
		if evidence := exactCallers[service.Name]; len(evidence) > 0 {
			if service.Team != "" {
				teamSet[service.Team] = true
			}
			result.DirectServices++
			result.Services = append(result.Services, impactedService{
				Name: service.Name, Team: service.Team, Depth: service.Depth, Tier: "exact_endpoint_caller",
				Reason: "caller matches a changed entrypoint", Evidence: evidence,
			})
			continue
		}
		if service.Team != "" {
			potentialTeamSet[service.Team] = true
		}
		result.PotentialServices = append(result.PotentialServices, impactedService{
			Name: service.Name, Team: service.Team, Depth: service.Depth, Tier: "potential_service_level_caller",
			Reason: "dependency path exists, but no changed endpoint or resource was matched",
		})
	}
	for _, node := range flow.Nodes {
		if node.Kind != "service" && node.Kind != "external" && node.Label != "" {
			resourceSet[node.Kind+": "+node.Label] = true
		}
	}
	result.Teams, result.PotentialResources = sortedKeys(teamSet), sortedKeys(resourceSet)
	result.Resources = []string{}
	result.ScoreEligible = result.Freshness == "fresh" && len(result.Services) > 0
	switch {
	case result.Freshness != "fresh":
		result.Confidence = "stale_graph_estimate"
		result.Notes = append(result.Notes, fmt.Sprintf("graph snapshot is %s relative to PR head; updating the same default branch does not guarantee a PR-head match", result.Freshness))
		result.NextAction = "Prepare PR flows above to analyze the exact PR revisions in isolated checkouts, then refresh impact evidence."
	case len(result.Services) > 0:
		result.Confidence = "changed_surface_evidence"
	default:
		result.Confidence = "service_level_candidates"
	}
	if len(result.Services) == 0 {
		result.Notes = append(result.Notes, "no exact caller of a changed entrypoint or resource is proven by the graph")
	}
	if len(result.PotentialServices) > 0 {
		result.Notes = append(result.Notes, fmt.Sprintf("%d service-level dependency candidate(s) are shown separately and are not proven PR impacts", len(result.PotentialServices)))
	}
	if len(potentialTeamSet) > 0 {
		result.Notes = append(result.Notes, fmt.Sprintf("candidate paths span %d additional team(s)", len(potentialTeamSet)))
	}
	result.Notes = append(result.Notes, "the technical graph is repository-wide context; only exact changed-surface matches count toward PR risk")
	return result
}

var diffHunkPattern = regexp.MustCompile(`^@@ -([0-9]+)(?:,[0-9]+)? \+([0-9]+)(?:,[0-9]+)? @@`)
var routeParameterPattern = regexp.MustCompile(`\$?\{[^}]+\}|:[A-Za-z_][A-Za-z0-9_]*`)

type changedLines struct {
	old map[int]bool
	new map[int]bool
	ok  bool
}

func parseChangedLines(patch string) changedLines {
	result := changedLines{old: map[int]bool{}, new: map[int]bool{}}
	oldLine, newLine := 0, 0
	for _, line := range strings.Split(patch, "\n") {
		if match := diffHunkPattern.FindStringSubmatch(line); len(match) == 3 {
			oldLine, _ = strconv.Atoi(match[1])
			newLine, _ = strconv.Atoi(match[2])
			result.ok = true
			continue
		}
		if !result.ok || strings.HasPrefix(line, "\\ No newline") {
			continue
		}
		switch {
		case strings.HasPrefix(line, "+"):
			result.new[newLine] = true
			newLine++
		case strings.HasPrefix(line, "-"):
			result.old[oldLine] = true
			oldLine++
		default:
			oldLine++
			newLine++
		}
	}
	return result
}

func graphService(graph *ArchGraph, name string) *archgraph.ServiceNode {
	if graph == nil {
		return nil
	}
	for _, service := range graph.Services {
		if service != nil && service.Name == name {
			return service
		}
	}
	return nil
}

func changedEntrypoints(service *archgraph.ServiceNode, files []githubPullFile, headSHA string) []changedEntrypoint {
	if service == nil {
		return nil
	}
	byPath := map[string]changedLines{}
	for _, file := range files {
		byPath[filepath.ToSlash(file.Filename)] = parseChangedLines(file.Patch)
	}
	var result []changedEntrypoint
	collections := [][]archgraph.EntitySummary{service.HTTPRoutes, service.RPCEndpoints, service.QueueConsumers, service.ScheduledJobs, service.Webhooks, service.CLICommands}
	for _, collection := range collections {
		for _, entity := range collection {
			// PR patches use head coordinates for added lines. A graph from an
			// older revision (including the base branch) cannot use those lines.
			// GitHub's PR diff starts at the merge base, which need not be Base.SHA.
			revision := entityGraphRevision(entity)
			atHead := pullRequestGraphFreshness(revision, headSHA) == "fresh" && archgraph.DescribeRelationship("", entity).Class == "source_extracted"
			var best *changedEntrypoint
			for _, location := range entitySourceLocations(entity) {
				lines, changed := byPath[filepath.ToSlash(location.File)]
				if !changed {
					continue
				}
				match := "file_scope"
				if lines.ok && atHead && location.StartLine > 0 {
					if !locationIntersectsChangedLines(location, lines.new) {
						continue
					}
					match = "changed_line"
				}
				entry := changedEntrypoint{ID: entity.ID, Kind: entity.Kind, Name: entity.Name, File: location.File, Match: match}
				if best == nil || match == "changed_line" {
					best = &entry
				}
				if match == "changed_line" {
					break
				}
			}
			if best != nil {
				result = append(result, *best)
			}
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	return result
}

func entitySourceLocations(entity archgraph.EntitySummary) []model.Location {
	value := entity.Details["source_locations"]
	switch locations := value.(type) {
	case []model.Location:
		return locations
	case []any:
		result := make([]model.Location, 0, len(locations))
		for _, raw := range locations {
			if item, ok := raw.(map[string]any); ok {
				result = append(result, model.Location{File: stringValue(item["file"]), StartLine: intValue(item["start_line"]), EndLine: intValue(item["end_line"])})
			}
		}
		return result
	default:
		return nil
	}
}

func locationIntersectsChangedLines(location model.Location, lines map[int]bool) bool {
	end := location.EndLine
	if end < location.StartLine {
		end = location.StartLine
	}
	for line := location.StartLine; line <= end; line++ {
		if lines[line] {
			return true
		}
	}
	return false
}

func exactChangedSurfaceCallers(graph *ArchGraph, root string, changed []changedEntrypoint) map[string][]string {
	operations, ids := map[string]string{}, map[string]string{}
	for _, entrypoint := range changed {
		if entrypoint.Match != "changed_line" {
			continue
		}
		// Names and IDs are not proof across protocols. Unknown surface kinds
		// remain candidates instead of borrowing another protocol's identity.
		protocol := exactSurfaceProtocol(entrypoint.Kind)
		if protocol == "" {
			continue
		}
		if protocol == "http" {
			if key := normalizedOperation(entrypoint.Name); key != "" {
				operations[key] = entrypoint.Name
			}
		}
		if entrypoint.ID != "" {
			ids[protocol+"\x00"+entrypoint.ID] = entrypoint.Name
		}
	}
	result := map[string][]string{}
	for _, edge := range graph.Edges {
		if edge == nil || edge.To != root || graphService(graph, edge.From) == nil {
			continue
		}
		caller := graphService(graph, edge.From)
		status := caller.AnalysisStatus
		if status == nil || status.Dirty || status.State != "analyzed_clean" || status.AnalyzedRevision == "" {
			continue
		}
		for _, detail := range edge.Details {
			if archgraph.DescribeRelationship(graph.RunID, detail).Class != "source_extracted" || pullRequestGraphFreshness(entityGraphRevision(detail), status.AnalyzedRevision) != "fresh" {
				continue
			}
			matched := ""
			if detail.ID != "" {
				matched = ids[edge.Type+"\x00"+detail.ID]
			}
			if matched == "" && edge.Type == "http" {
				matched = operations[httpCallerOperation(detail)]
			}
			if matched != "" {
				result[edge.From] = appendUnique(result[edge.From], matched)
			}
		}
	}
	return result
}

// Display names may contain a service prefix and a full URL. Use the
// source fact's HTTP fields after the graph has resolved its destination.
func httpCallerOperation(detail archgraph.EntitySummary) string {
	fields := detail.Details
	method, _ := fields["method"].(string)
	address, _ := fields["path"].(string)
	if address == "" {
		address, _ = fields["url_template"].(string)
	}
	if method == "" && address == "" {
		metadata, _ := fields["metadata"].(map[string]any)
		nested, _ := metadata["details"].(map[string]any)
		method, _ = nested["method"].(string)
		address, _ = nested["path"].(string)
		if address == "" {
			address, _ = nested["url_template"].(string)
		}
	}
	if method == "" && address == "" {
		return normalizedOperation(detail.Name)
	}
	if method == "" || address == "" || len(strings.Fields(method)) != 1 || strings.ContainsAny(address, " \t\r\n") {
		return ""
	}
	parsed, err := url.Parse(address)
	if err != nil {
		return ""
	}
	if parsed.IsAbs() {
		if (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
			return ""
		}
	} else if !strings.HasPrefix(address, "/") || strings.HasPrefix(address, "//") {
		return ""
	}
	route := parsed.EscapedPath()
	if route == "" {
		route = "/"
	}
	// Preserve literal escaping/case/trailing slash; normalize only the existing
	// recognized route-parameter notation.
	if !strings.Contains(address, "%") {
		route = parsed.Path
		if route == "" {
			route = "/"
		}
	}
	return normalizedOperation(method + " " + route)
}

func exactSurfaceProtocol(kind string) string {
	switch kind {
	case "http_route", "http_endpoint", "webhook":
		return "http"
	case "rpc_endpoint":
		return "rpc"
	default:
		return ""
	}
}

func normalizedOperation(value string) string {
	fields := strings.Fields(strings.TrimSpace(value))
	if len(fields) < 2 {
		return ""
	}
	method := strings.ToUpper(fields[0])
	path := fields[1]
	path = routeParameterPattern.ReplaceAllString(path, "{}")
	return method + " " + path
}

func appendUnique(values []string, value string) []string {
	for _, existing := range values {
		if existing == value {
			return values
		}
	}
	return append(values, value)
}

func dependenciesForChangedEntrypoints(service *archgraph.ServiceNode, changed []changedEntrypoint) []entrypointDependency {
	if service == nil {
		return nil
	}
	names := map[string]string{}
	for _, entrypoint := range changed {
		if entrypoint.Match == "changed_line" && entrypoint.ID != "" {
			names[entrypoint.ID] = entrypoint.Name
		}
	}
	var result []entrypointDependency
	for _, connection := range service.Connections {
		entrypoint := names[connection.FromID]
		if entrypoint == "" {
			continue
		}
		result = append(result, entrypointDependency{Entrypoint: entrypoint, Target: connection.ToName, Kind: connection.Kind, Reachability: connection.Reachability})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Entrypoint != result[j].Entrypoint {
			return result[i].Entrypoint < result[j].Entrypoint
		}
		return result[i].Target < result[j].Target
	})
	return result
}

func serviceGraphRevision(service *archgraph.ServiceNode) graphRevision {
	if service == nil {
		return graphRevision{}
	}
	collections := [][]archgraph.EntitySummary{service.HTTPRoutes, service.RPCEndpoints, service.QueueConsumers, service.ScheduledJobs, service.Webhooks, service.CLICommands, service.Dependencies}
	var result graphRevision
	if status := service.AnalysisStatus; status != nil {
		if status.Dirty || status.State == "analyzed_dirty" {
			return graphRevision{Dirty: true}
		}
		if status.State == "analyzed_clean" {
			result = graphRevision{Commit: status.AnalyzedRevision, Branch: status.Branch}
		}
	}
	for _, collection := range collections {
		for _, entity := range collection {
			revision := entityGraphRevision(entity)
			if revision.Dirty {
				return graphRevision{Dirty: true}
			}
			if revision.Commit == "" {
				continue
			}
			if result.Commit != "" && result.Commit != revision.Commit {
				return graphRevision{}
			}
			result = revision
		}
	}
	return result
}

func entityGraphRevision(entity archgraph.EntitySummary) graphRevision {
	revision, _ := entity.Details["repository_revision"].(map[string]any)
	return graphRevision{Commit: stringValue(revision["commit"]), Branch: stringValue(revision["branch"]), Dirty: boolValue(revision["dirty"])}
}

func pullRequestGraphFreshness(revision graphRevision, headSHA string) string {
	if revision.Dirty {
		return "dirty"
	}
	if revision.Commit != "" && headSHA != "" {
		if revision.Commit == headSHA {
			return "fresh"
		}
		return "stale"
	}
	return "unknown"
}

func stringValue(value any) string {
	valueString, _ := value.(string)
	return valueString
}

func intValue(value any) int {
	switch number := value.(type) {
	case int:
		return number
	case float64:
		return int(number)
	default:
		return 0
	}
}

func boolValue(value any) bool {
	valueBool, _ := value.(bool)
	return valueBool
}

func graphServiceForRepo(graph *ArchGraph, repo store.Repo) string {
	wantPath := cleanPath(firstNonEmpty(repo.Path, repo.ClonePath))
	wantName := normalizedRepoName(repo.Name)
	for _, service := range graph.Services {
		if service == nil {
			continue
		}
		if repo.ID != "" && service.RepoID == repo.ID {
			return service.Name
		}
		if wantPath != "" && cleanPath(service.RepoPath) == wantPath {
			return service.Name
		}
	}
	for _, service := range graph.Services {
		if service != nil && normalizedRepoName(service.Name) == wantName {
			return service.Name
		}
	}
	return ""
}

func cleanPath(path string) string {
	if strings.TrimSpace(path) == "" {
		return ""
	}
	return filepath.Clean(path)
}
func normalizedRepoName(name string) string {
	return strings.Trim(strings.ToLower(strings.NewReplacer("_", "-", ".", "-").Replace(name)), "-")
}
func changeArea(path string) string {
	parts := strings.Split(filepath.ToSlash(path), "/")
	if len(parts) >= 2 {
		return parts[0] + "/" + parts[1]
	}
	if len(parts) == 1 {
		return parts[0]
	}
	return "root"
}
func riskLevel(score int) string {
	switch {
	case score >= 80:
		return "critical"
	case score >= 55:
		return "high"
	case score >= 30:
		return "moderate"
	default:
		return "low"
	}
}
func companyImpactScore(c companyImpact) int {
	if !c.Available || !c.ScoreEligible || (c.DirectServices == 0 && c.IndirectServices == 0 && len(c.Resources) == 0) {
		return 0
	}
	return minInt(100, 8+c.DirectServices*12+c.IndirectServices*5+len(c.Teams)*8+len(c.Resources)*2)
}
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func sortedKeys[T ~bool](in map[string]T) []string {
	out := make([]string, 0, len(in))
	for key := range in {
		out = append(out, key)
	}
	sort.Strings(out)
	return out
}
func signalRank(severity string) int {
	switch severity {
	case "high":
		return 3
	case "medium":
		return 2
	default:
		return 1
	}
}
