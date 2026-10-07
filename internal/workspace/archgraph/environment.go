package archgraph

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/mohammad-safakhou/diffmind/internal/extractor/dependencies"
	"github.com/mohammad-safakhou/diffmind/internal/workspace/model"
)

type EnvironmentChange struct {
	Ecosystem string              `json:"ecosystem"`
	Name      string              `json:"name"`
	Module    string              `json:"module"`
	Scope     string              `json:"scope,omitempty"`
	Kind      string              `json:"kind"`
	Before    []dependencies.Fact `json:"before"`
	After     []dependencies.Fact `json:"after"`
}

// Environment changes are review context, not evidence that a flow changed.
func CompareEnvironments(a, b *model.RepoMetrics) []EnvironmentChange {
	left, right := environmentFacts(a), environmentFacts(b)
	keys := map[string]bool{}
	for k := range left {
		keys[k] = true
	}
	for k := range right {
		keys[k] = true
	}
	var out []EnvironmentChange
	for k := range keys {
		if environmentValue(left[k]) == environmentValue(right[k]) {
			continue
		}
		f := left[k]
		if len(right[k]) > 0 {
			f = right[k]
		}
		if len(f) == 0 {
			continue
		}
		kind := "version_changed"
		if len(left[k]) == 0 {
			kind = "added"
		} else if len(right[k]) == 0 {
			kind = "removed"
		} else {
			for _, f := range append(append([]dependencies.Fact{}, left[k]...), right[k]...) {
				if f.Version == "" {
					kind = "evidence_changed"
				}
			}
		}
		out = append(out, EnvironmentChange{Ecosystem: f[0].Ecosystem, Name: f[0].Name, Module: f[0].Module, Scope: f[0].Scope, Before: left[k], After: right[k], Kind: kind})
	}
	sort.Slice(out, func(i, j int) bool {
		a, b := out[i], out[j]
		return a.Module+":"+a.Ecosystem+":"+a.Name+":"+a.Scope < b.Module+":"+b.Ecosystem+":"+b.Name+":"+b.Scope
	})
	return out
}

func environmentFacts(m *model.RepoMetrics) map[string][]dependencies.Fact {
	out := map[string][]dependencies.Fact{}
	if m == nil || m.DependencyInventory == nil {
		return out
	}
	for _, f := range m.DependencyInventory.Dependencies {
		k := strings.Join([]string{f.Module, f.Ecosystem, f.Name, f.Scope}, "\x00")
		out[k] = append(out[k], f)
	}
	for _, t := range m.Toolchain {
		name := string(t.Language)
		if name == "javascript" || name == "typescript" {
			name = "node"
		}
		f := dependencies.Fact{Ecosystem: "runtime", Name: name, Module: ".", Version: t.Version, Declared: t.VersionConstraint, Sources: t.Sources}
		k := f.Ecosystem + ":" + f.Name
		out[k] = append(out[k], f)
		if t.BuildTool != "" {
			f = dependencies.Fact{Ecosystem: "tool", Name: t.BuildTool, Module: ".", Version: t.BuildToolVersion, Sources: t.Sources}
			k = f.Ecosystem + ":" + f.Name
			out[k] = append(out[k], f)
		}
	}
	return out
}

func environmentValue(facts []dependencies.Fact) string {
	values := make([]string, 0, len(facts))
	seen := map[string]bool{}
	for _, f := range facts {
		b, _ := json.Marshal([]string{f.Declared, f.Version, f.Conditions})
		if !seen[string(b)] {
			values = append(values, string(b))
			seen[string(b)] = true
		}
	}
	sort.Strings(values)
	b, _ := json.Marshal(values)
	return string(b)
}
