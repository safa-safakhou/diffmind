package ui

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/mohammad-safakhou/diffmind/internal/workspace/store"
)

// Explicit API endpoints are approved configuration. Redirects cannot extend
// that approval to another origin or forward its credentials there.
func githubHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{Timeout: timeout, CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) >= 10 {
			return fmt.Errorf("GitHub redirect limit exceeded")
		}
		first := via[0].URL
		if req.URL.Scheme != first.Scheme || !strings.EqualFold(req.URL.Host, first.Host) {
			return fmt.Errorf("GitHub redirected to another origin; check the approved API endpoint")
		}
		return nil
	}}
}

func validateGitHubAPIBase(raw string) error {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return fmt.Errorf("GitHub API base must be an absolute URL without credentials, query or fragment")
	}
	loopback := strings.EqualFold(u.Hostname(), "localhost")
	if ip := net.ParseIP(u.Hostname()); ip != nil {
		loopback = ip.IsLoopback()
	}
	if u.Scheme != "https" && !(u.Scheme == "http" && loopback) {
		return fmt.Errorf("GitHub API base requires HTTPS (HTTP is supported only on loopback)")
	}
	return nil
}

func gitRepositoryLocation(raw string) (host, owner, name string, ok bool) {
	raw = strings.TrimSpace(raw)
	if !strings.Contains(raw, "://") {
		if at := strings.Index(raw, "@"); at >= 0 {
			raw = raw[at+1:]
		}
		if colon := strings.Index(raw, ":"); colon > 0 {
			raw = "ssh://" + raw[:colon] + "/" + raw[colon+1:]
		} else {
			raw = "https://" + raw
		}
	}
	u, err := url.Parse(raw)
	if err != nil || u.Hostname() == "" {
		return
	}
	parts := strings.Split(strings.Trim(strings.TrimSuffix(u.Path, ".git"), "/"), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return
	}
	valid := func(s string) bool {
		for _, r := range s {
			if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || strings.ContainsRune("._-", r)) {
				return false
			}
		}
		return s != "." && s != ".."
	}
	if !valid(parts[0]) || !valid(parts[1]) {
		return
	}
	return strings.ToLower(u.Hostname()), parts[0], parts[1], true
}

func githubRepositoryEndpoint(ctx context.Context, repo store.Repo) (base, state, message string) {
	source := githubSourceForRepo(ctx, repo)
	host, owner, name, ok := gitRepositoryLocation(source)
	if !ok {
		if strings.TrimSpace(repo.GitURL) != "" {
			return "", "unsupported_provider", "This repository URL is not a supported GitHub source."
		}
		if repo.SourceType == "local" || repo.GitURL == "" && repo.Path != "" {
			return "", "local_only", "This local repository has no supported Git remote for PR queries."
		}
		return "", "missing_remote", "Add a supported GitHub repository URL before querying PRs."
	}
	if host != "github.com" && repo.GitProvider != "github" && repo.GitAPIBase == "" {
		return "", "unsupported_provider", "This repository provider is not supported for PR queries."
	}
	if host != "github.com" && repo.GitAPIBase == "" {
		return "", "provider_unavailable", "Configure the approved GitHub API endpoint for this repository host."
	}
	api := strings.TrimRight(strings.TrimSpace(repo.GitAPIBase), "/")
	if api == "" {
		if host == "github.com" {
			api = "https://api.github.com"
		} else {
			api = "https://" + host + "/api/v3"
		}
	}
	if err := validateGitHubAPIBase(api); err != nil {
		return "", "provider_unavailable", err.Error()
	}
	return api + "/repos/" + owner + "/" + name, "ready", ""
}

func githubResponseError(code int) error {
	switch code {
	case 401:
		return fmt.Errorf("GitHub authentication failed. Check server credentials for the approved API host.")
	case 403:
		return fmt.Errorf("GitHub denied this request. Check repository access and API rate limits.")
	case 404:
		return fmt.Errorf("GitHub repository or pull request is unavailable. Check its URL and your access.")
	case 429:
		return fmt.Errorf("GitHub rate limit reached. Wait before refreshing again.")
	default:
		return fmt.Errorf("GitHub request failed (HTTP %d). Check the approved API endpoint or retry later.", code)
	}
}

func githubConnectionError() error {
	return fmt.Errorf("GitHub request could not complete. Check the approved API endpoint, connection and redirect policy.")
}
