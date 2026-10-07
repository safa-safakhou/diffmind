package detectors

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/mohammad-safakhou/diffmind/internal/extractor/dependencies"
)

// Revision changes when extraction semantics or the validation matrix changes.
// It is recorded separately from application/library versions.
const Revision = "2026-10-07.1"

type Requirement struct {
	Ecosystem string `json:"ecosystem" yaml:"ecosystem"`
	Name      string `json:"name" yaml:"name"`
	Versions  string `json:"versions" yaml:"versions"`
}

func SummarizeCoverage(items []Coverage) []Coverage {
	seen := map[string]Coverage{}
	for _, c := range items {
		b, _ := json.Marshal(c)
		seen[string(b)] = c
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := make([]Coverage, 0, len(keys))
	for _, k := range keys {
		out = append(out, seen[k])
	}
	return out
}

type VersionRule struct {
	ID       string        `json:"id"`
	Requires []Requirement `json:"requires,omitempty"`
	// TestedVersions means static extraction fixtures, not runtime certification.
	TestedVersions map[string][]string `json:"tested_versions,omitempty"`
}
type Coverage struct {
	DetectorID   string              `json:"detector_id"`
	RuleID       string              `json:"rule_id"`
	Revision     string              `json:"revision"`
	Status       string              `json:"status"`
	Module       string              `json:"module"`
	Dependencies []dependencies.Fact `json:"dependencies,omitempty"`
	Reason       string              `json:"reason,omitempty"`
}

func (r Requirement) Validate() error {
	if strings.TrimSpace(r.Ecosystem) == "" || strings.TrimSpace(r.Name) == "" || strings.TrimSpace(r.Versions) == "" {
		return fmt.Errorf("ecosystem, name and versions are required")
	}
	if _, e := semver.NewConstraint(r.Versions); e != nil {
		return fmt.Errorf("invalid version range %q: %w", r.Versions, e)
	}
	return nil
}

// Evaluate keeps applicability separate from evidence quality. Future versions
// are still analyzed using source-backed fallback rules, but never called tested.
func Evaluate(id, file string, inventory dependencies.Inventory, rules []VersionRule) Coverage {
	out := Coverage{DetectorID: id, RuleID: id + ".source-fallback", Revision: Revision, Status: "unknown_version", Module: ".", Reason: "Dependency version could not be established; source evidence only"}
	if len(rules) == 0 {
		out.Status = "unversioned"
		out.Reason = "No version-specific coverage has been declared"
		return out
	}
	var best *Coverage
	for _, rule := range rules {
		c := Coverage{DetectorID: id, RuleID: rule.ID, Revision: Revision, Status: "validated", Module: "."}
		allKnown, allMatch, allTested := true, true, true
		for _, req := range rule.Requires {
			facts := selectFacts(inventory.ForFile(file), req)
			c.Dependencies = append(c.Dependencies, facts...)
			if len(facts) == 0 {
				allKnown = false
				continue
			}
			versions := map[string]bool{}
			for _, f := range facts {
				c.Module = f.Module
				if f.Version == "" || f.Conditions != "" {
					allKnown = false
				} else {
					versions[f.Version] = true
				}
			}
			if len(versions) > 1 {
				c.Status = "ambiguous_version"
				c.Reason = "Multiple effective versions apply in this module; no version-specific interpretation selected"
				return c
			}
			for v := range versions {
				constraint, _ := semver.NewConstraint(req.Versions)
				sv, e := semver.NewVersion(v)
				if e != nil || constraint == nil || !constraint.Check(sv) {
					allMatch = false
				}
				tested := false
				for _, tv := range rule.TestedVersions[req.Ecosystem+":"+req.Name] {
					if tv == v {
						tested = true
					}
				}
				if !tested {
					allTested = false
				}
			}
		}
		if !allKnown {
			c.Status = "unknown_version"
			c.RuleID = out.RuleID
			c.Reason = out.Reason
		} else if !allMatch {
			c.Status = "unsupported_version"
			c.RuleID = out.RuleID
			c.Reason = "Version is outside declared applicability; source evidence only"
		} else if !allTested {
			c.Status = "compatible_unverified"
			c.Reason = "Rule applies, but this exact version has no static extraction fixture"
		}
		if allKnown && allMatch {
			return c
		}
		if best == nil || c.Status == "unsupported_version" {
			copy := c
			best = &copy
		}
	}
	if best != nil {
		return *best
	}
	return out
}

func selectFacts(facts []dependencies.Fact, req Requirement) []dependencies.Fact {
	var out []dependencies.Fact
	rank := 0
	for _, f := range facts {
		if f.Ecosystem != req.Ecosystem || f.Name != req.Name {
			continue
		}
		if f.Scope == "parent" || f.Scope == "bom" || f.Scope == "peer" || f.Scope == "test" || f.Scope == "testImplementation" {
			continue
		}
		n := 1
		switch f.Resolution {
		case "lockfile":
			n = 5
		case "managed_exact":
			n = 3
		case "declared_exact":
			n = 3
		case "configured":
			n = 2
		}
		if f.Scope == "management" {
			n = 2
		}
		if n > rank {
			out = nil
			rank = n
		}
		if n == rank {
			out = append(out, f)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Version < out[j].Version })
	return out
}

func VersionRules(id string) []VersionRule {
	var eco, name, rangeSpec string
	switch id {
	case "javascript.http.express":
		return []VersionRule{
			{ID: id + ".v4", Requires: []Requirement{{"npm", "express", ">=4 <5"}}, TestedVersions: map[string][]string{"npm:express": {"4.21.2"}}},
			{ID: id + ".v5", Requires: []Requirement{{"npm", "express", ">=5 <6"}}, TestedVersions: map[string][]string{"npm:express": {"5.1.0"}}},
		}
	case "java.http.spring":
		eco, name, rangeSpec = "maven", "org.springframework:spring-web", ">=5 <8"
	case "java.activation.spring", "java.cache.spring":
		eco, name, rangeSpec = "maven", "org.springframework:spring-context", ">=5 <8"
	case "java.httpclient.feign":
		eco, name, rangeSpec = "maven", "org.springframework.cloud:spring-cloud-openfeign-core", ">=2 <6"
	case "java.httpclient.retrofit":
		eco, name, rangeSpec = "maven", "com.squareup.retrofit2:retrofit", ">=2 <4"
	case "java.queue.sqs":
		eco, name, rangeSpec = "maven", "io.awspring.cloud:spring-cloud-aws-sqs", ">=3 <5"
	case "java.queue.kafka":
		eco, name, rangeSpec = "maven", "org.springframework.kafka:spring-kafka", ">=2 <5"
	case "java.queue.rabbitmq":
		eco, name, rangeSpec = "maven", "org.springframework.amqp:spring-rabbit", ">=2 <5"
	case "java.queue.jms":
		eco, name, rangeSpec = "maven", "org.springframework:spring-jms", ">=5 <8"
	case "java.db.jpa":
		eco, name, rangeSpec = "maven", "org.springframework.data:spring-data-jpa", ">=2 <5"
	case "java.db.jdbc":
		eco, name, rangeSpec = "maven", "org.springframework:spring-jdbc", ">=5 <8"
	case "python.http.fastapi":
		eco, name, rangeSpec = "pypi", "fastapi", ">=0.60 <1"
	case "python.http.flask":
		eco, name, rangeSpec = "pypi", "flask", ">=1 <4"
	case "python.http.django", "python.db.django":
		eco, name, rangeSpec = "pypi", "django", ">=2 <7"
	case "python.cache.redis":
		eco, name, rangeSpec = "pypi", "redis", ">=3 <8"
	case "typescript.http.nestjs":
		eco, name, rangeSpec = "npm", "@nestjs/core", ">=7 <12"
	case "javascript.db.prisma":
		eco, name, rangeSpec = "npm", "@prisma/client", ">=2 <8"
	case "javascript.db.sequelize":
		eco, name, rangeSpec = "npm", "sequelize", ">=5 <8"
	case "golang.db.gorm":
		eco, name, rangeSpec = "go", "gorm.io/gorm", ">=1 <2"
	case "golang.db.bun":
		eco, name, rangeSpec = "go", "github.com/uptrace/bun", ">=1 <2"
	case "golang.http.fiber":
		eco, name, rangeSpec = "go", "github.com/gofiber/fiber/v2", ">=2 <3"
	case "golang.http.echo":
		eco, name, rangeSpec = "go", "github.com/labstack/echo/v4", ">=4 <5"
	case "golang.http.gin":
		eco, name, rangeSpec = "go", "github.com/gin-gonic/gin", ">=1 <2"
	case "golang.rpc.grpc":
		eco, name, rangeSpec = "go", "google.golang.org/grpc", ">=1 <2"
	case "ruby.http.rails":
		eco, name, rangeSpec = "gem", "rails", ">=5 <9"
	case "ruby.db.activerecord":
		eco, name, rangeSpec = "gem", "activerecord", ">=5 <9"
	default:
		return nil
	}
	// Applicability ranges are deliberately not a claim that every release in
	// the range is tested. Exact-version fixtures must be registered separately.
	tested := map[string][]string{
		"java.http.spring":    {"6.2.0", "7.0.0"},
		"python.http.fastapi": {"0.115.0"},
		"python.http.flask":   {"3.1.0"},
	}
	var matrix map[string][]string
	if versions := tested[id]; len(versions) > 0 {
		matrix = map[string][]string{eco + ":" + name: versions}
	}
	return []VersionRule{{ID: id + ".shared", Requires: []Requirement{{eco, name, rangeSpec}}, TestedVersions: matrix}}
}
