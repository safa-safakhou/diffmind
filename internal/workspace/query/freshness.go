package query

import (
	"context"
	"os/exec"
	"strings"
	"time"

	"github.com/mohammad-safakhou/diffmind/internal/workspace/artifacts"
	"github.com/mohammad-safakhou/diffmind/internal/workspace/store"
)

// RepositoryFreshness compares saved analysis with a current checkout without
// fetching, pulling, registering repositories or changing saved evidence.
func RepositoryFreshness(repo store.Repo, latest *artifacts.DiffMindRunInfo) string {
	if latest == nil || latest.RepoGitSHA == "" {
		return "unknown"
	}
	current := repo.RemoteHeadSHA
	if current == "" {
		current = repo.HeadSHA
	}
	if repo.SourceType == "local" || (repo.GitURL == "" && repo.Path != "") {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, "git", "-C", repo.Path, "status", "--porcelain=v2", "--branch")
		output, err := cmd.Output()
		if err != nil {
			return "unknown"
		}
		current = ""
		dirty := false
		for _, line := range strings.Split(string(output), "\n") {
			if strings.HasPrefix(line, "# branch.oid ") {
				current = strings.TrimSpace(strings.TrimPrefix(line, "# branch.oid "))
			} else if line != "" && !strings.HasPrefix(line, "# ") {
				dirty = true
			}
		}
		if current == "" || current == "(initial)" {
			return "unknown"
		}
		if dirty {
			return "dirty"
		}
	}
	if current == "" {
		return "unknown"
	}
	if latest.RepoGitSHA == current {
		return "fresh"
	}
	return "stale"
}

func (s *Service) currentRepositories(pid string) ([]store.Repo, error) {
	repos, err := s.store.ListRepos(pid)
	if err != nil {
		return nil, err
	}
	groups, err := artifacts.DiscoverDiffMindRunsByRepo(s.runsDir)
	if err != nil {
		// A saved graph remains usable when analysis artifacts are unavailable.
		for i := range repos {
			repos[i].DiffMindFreshness = "unknown"
		}
		return repos, nil
	}
	for i := range repos {
		repo := &repos[i]
		var latest *artifacts.DiffMindRunInfo
		if runs := groups[repo.Path]; len(runs) > 0 {
			latest = &runs[0]
		} else if repo.LastDiffMindRunID != "" {
			if info, ok := artifacts.DiffMindRunByID(s.runsDir, repo.LastDiffMindRunID); ok && artifacts.RunMatchesRepo(info, repo.Name, repo.ID, repo.Path) {
				latest = &info
			}
		}
		repo.DiffMindFreshness = RepositoryFreshness(*repo, latest)
	}
	return repos, nil
}
