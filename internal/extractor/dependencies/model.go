// Package dependencies records static dependency versions without executing a
// repository's build or installing its packages. A declaration is not evidence
// of what is deployed; resolution records the source of each version.
package dependencies

import (
	"path"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
)

type Fact struct {
	Ecosystem  string   `json:"ecosystem" yaml:"ecosystem"`
	Name       string   `json:"name" yaml:"name"`
	Module     string   `json:"module" yaml:"module"`
	Declared   string   `json:"declared,omitempty" yaml:"declared,omitempty"`
	Version    string   `json:"version,omitempty" yaml:"version,omitempty"`
	Resolution string   `json:"resolution" yaml:"resolution"`
	Scope      string   `json:"scope,omitempty" yaml:"scope,omitempty"`
	Conditions string   `json:"conditions,omitempty" yaml:"conditions,omitempty"`
	Sources    []string `json:"sources" yaml:"sources"`
}

type Inventory struct {
	Dependencies []Fact   `json:"dependencies"`
	Limitations  []string `json:"limitations,omitempty"`
}

// ForFile selects the nearest declaring module for each package independently.
// Workspace tools at the root must not overwrite a nested application's version.
func (i Inventory) ForFile(file string) []Fact {
	file = path.Clean(strings.ReplaceAll(file, "\\", "/"))
	nearest := map[string]int{}
	for _, f := range i.Dependencies {
		m := path.Clean(f.Module)
		if m != "." && file != m && !strings.HasPrefix(file, m+"/") {
			continue
		}
		k := f.Ecosystem + ":" + f.Name
		if n, ok := nearest[k]; !ok || len(m) > n {
			nearest[k] = len(m)
		}
	}
	var out []Fact
	for _, f := range i.Dependencies {
		m := path.Clean(f.Module)
		if (m == "." || file == m || strings.HasPrefix(file, m+"/")) && nearest[f.Ecosystem+":"+f.Name] == len(m) {
			out = append(out, f)
		}
	}
	return out
}

func ExactVersion(value string) string {
	value = strings.TrimSpace(value)
	if strings.ContainsAny(value, "<>=~^*| ,[]()${}") {
		return ""
	}
	if _, err := semver.NewVersion(value); err != nil {
		return ""
	}
	return value
}

func normalize(i *Inventory) {
	seen := map[string]bool{}
	out := i.Dependencies[:0]
	for _, f := range i.Dependencies {
		if f.Name == "" {
			continue
		}
		if f.Module == "" {
			f.Module = "."
		}
		if f.Resolution == "" {
			f.Version = ExactVersion(f.Declared)
			f.Resolution = "constraint"
			if f.Version != "" {
				f.Resolution = "declared_exact"
			} else if f.Declared == "" || strings.Contains(f.Declared, "${") {
				f.Resolution = "unresolved"
			}
		}
		if f.Ecosystem == "pypi" {
			f.Name = strings.NewReplacer("_", "-", ".", "-").Replace(strings.ToLower(f.Name))
		}
		sort.Strings(f.Sources)
		key := f.Ecosystem + "\x00" + f.Module + "\x00" + f.Name + "\x00" + f.Declared + "\x00" + f.Version + "\x00" + f.Resolution + "\x00" + f.Scope + "\x00" + f.Conditions + "\x00" + strings.Join(f.Sources, "\x00")
		if !seen[key] {
			seen[key] = true
			out = append(out, f)
		}
	}
	i.Dependencies = out
	sort.Slice(i.Dependencies, func(a, b int) bool {
		x, y := i.Dependencies[a], i.Dependencies[b]
		return x.Module+"\x00"+x.Ecosystem+"\x00"+x.Name+"\x00"+x.Scope+"\x00"+x.Version+"\x00"+strings.Join(x.Sources, ",") < y.Module+"\x00"+y.Ecosystem+"\x00"+y.Name+"\x00"+y.Scope+"\x00"+y.Version+"\x00"+strings.Join(y.Sources, ",")
	})
	seen = map[string]bool{}
	var limits []string
	for _, l := range i.Limitations {
		if !seen[l] {
			limits = append(limits, l)
			seen[l] = true
		}
	}
	sort.Strings(limits)
	i.Limitations = limits
}
