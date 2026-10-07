package dependencies

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/pelletier/go-toml/v2"
)

var pythonRequirementRE = regexp.MustCompile(`^([A-Za-z0-9][A-Za-z0-9_.-]*)(?:\[[^]]+\])?\s*(.*)$`)

func pythonFact(req, m, source, scope string) (Fact, bool) {
	req = strings.TrimSpace(req)
	if req == "" || strings.HasPrefix(req, "-") || strings.HasPrefix(req, "#") {
		return Fact{}, false
	}
	req = strings.Split(req, " #")[0]
	parts := strings.SplitN(req, ";", 2)
	match := pythonRequirementRE.FindStringSubmatch(strings.TrimSpace(parts[0]))
	if len(match) < 3 {
		return Fact{}, false
	}
	f := Fact{Ecosystem: "pypi", Name: match[1], Module: m, Declared: strings.TrimSpace(match[2]), Scope: scope, Sources: []string{source}}
	if len(parts) > 1 {
		f.Conditions = strings.TrimSpace(parts[1])
	}
	if strings.HasPrefix(f.Declared, "===") {
		f.Version = ExactVersion(strings.TrimPrefix(f.Declared, "==="))
	} else if strings.HasPrefix(f.Declared, "==") {
		f.Version = ExactVersion(strings.TrimPrefix(f.Declared, "=="))
	}
	if f.Version != "" {
		f.Resolution = "declared_exact"
	}
	return f, true
}

func inspectRequirements(root, rel string, i *Inventory, seen map[string]bool) error {
	if seen[rel] {
		return nil
	}
	seen[rel] = true
	b, e := readFile(root, rel)
	if e != nil {
		return e
	}
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "-r ") || strings.HasPrefix(line, "--requirement ") {
			parts := strings.Fields(line)
			r := path.Clean(path.Join(path.Dir(rel), parts[1]))
			if r == ".." || strings.HasPrefix(r, "../") || path.IsAbs(parts[1]) {
				i.Limitations = append(i.Limitations, rel+": requirement include outside repository ignored")
				continue
			}
			if e = inspectRequirements(root, r, i, seen); e != nil {
				i.Limitations = append(i.Limitations, fmt.Sprintf("%s: %v", r, e))
			}
			continue
		}
		if strings.HasPrefix(line, "-c ") || strings.HasPrefix(line, "--constraint ") {
			i.Limitations = append(i.Limitations, rel+": pip constraint include not resolved")
			continue
		}
		if f, ok := pythonFact(line, module(rel), rel, "runtime"); ok {
			i.Dependencies = append(i.Dependencies, f)
		}
	}
	return nil
}

func mapAt(doc map[string]any, keys ...string) map[string]any {
	for _, k := range keys {
		next, ok := doc[k].(map[string]any)
		if !ok {
			return nil
		}
		doc = next
	}
	return doc
}

func inspectPythonProject(root, rel string, i *Inventory) error {
	b, e := readFile(root, rel)
	if e != nil {
		return e
	}
	var doc map[string]any
	if e = toml.Unmarshal(b, &doc); e != nil {
		return e
	}
	m := module(rel)
	var facts []Fact
	addArray := func(raw any, scope string) {
		if arr, ok := raw.([]any); ok {
			for _, item := range arr {
				if s, ok := item.(string); ok {
					if f, ok := pythonFact(s, m, rel, scope); ok {
						facts = append(facts, f)
					}
				}
			}
		}
	}
	project := mapAt(doc, "project")
	addArray(project["dependencies"], "runtime")
	if extras, ok := project["optional-dependencies"].(map[string]any); ok {
		for group, deps := range extras {
			addArray(deps, "optional:"+group)
		}
	}
	for group, deps := range mapAt(doc, "dependency-groups") {
		addArray(deps, "dev:"+group)
	}
	poetry := mapAt(doc, "tool", "poetry", "dependencies")
	for name, raw := range poetry {
		if name == "python" {
			continue
		}
		var version, conditions string
		switch v := raw.(type) {
		case string:
			version = v
		case map[string]any:
			version, _ = v["version"].(string)
			conditions, _ = v["markers"].(string)
			if v["path"] != nil || v["git"] != nil {
				version = ""
			}
		}
		facts = append(facts, Fact{Ecosystem: "pypi", Name: name, Module: m, Declared: version, Conditions: conditions, Scope: "runtime", Sources: []string{rel}})
	}
	for _, f := range facts {
		if v, lock := pythonLockedVersion(root, m, f); v != "" {
			f.Version = v
			f.Resolution = "lockfile"
			f.Sources = append(f.Sources, lock)
		}
		i.Dependencies = append(i.Dependencies, f)
	}
	return nil
}

func pythonLockedVersion(root, m string, f Fact) (string, string) {
	if f.Conditions != "" {
		return "", ""
	}
	for _, filename := range []string{"uv.lock", "poetry.lock"} {
		rel := path.Join(m, filename)
		b, e := readFile(root, rel)
		if e != nil {
			continue
		}
		var lock struct {
			Version  int `toml:"version"`
			Metadata struct {
				LockVersion string `toml:"lock-version"`
			} `toml:"metadata"`
			Packages []struct {
				Name    string   `toml:"name"`
				Version string   `toml:"version"`
				Marker  string   `toml:"marker"`
				Markers []string `toml:"resolution-markers"`
			} `toml:"package"`
		}
		if toml.Unmarshal(b, &lock) != nil {
			return "", ""
		}
		if filename == "uv.lock" && lock.Version != 1 || filename == "poetry.lock" && lock.Metadata.LockVersion != "1.1" && lock.Metadata.LockVersion != "2.0" && lock.Metadata.LockVersion != "2.1" {
			return "", ""
		}
		versions := map[string]bool{}
		for _, p := range lock.Packages {
			if strings.EqualFold(strings.ReplaceAll(p.Name, "_", "-"), strings.ReplaceAll(f.Name, "_", "-")) && p.Marker == "" && len(p.Markers) == 0 && ExactVersion(p.Version) != "" {
				versions[p.Version] = true
			}
		}
		if len(versions) == 1 {
			for v := range versions {
				// A changed manifest invalidates a contradictory exact lock pin.
				exact := f.Version
				if exact == "" {
					exact = ExactVersion(f.Declared)
				}
				if exact != "" && exact != v || !pythonLockSatisfies(f.Declared, v) {
					return "", ""
				}
				return v, rel
			}
		}
		return "", ""
	}
	return "", ""
}

// Only resolve the conventional numeric subset shared with semver. PEP 440
// epochs, local/post/dev versions and compatible-release clauses stay unknown.
// A lock pin must not override a changed manifest's constraint.
func pythonLockSatisfies(declared, version string) bool {
	declared = strings.TrimSpace(declared)
	if declared == "" || declared == "*" {
		return true
	}
	if strings.Contains(declared, "~=") || strings.Contains(declared, "===") || strings.ContainsAny(strings.ReplaceAll(declared, "!=", ""), "!+;@") {
		return false
	}
	if strings.IndexFunc(declared, func(r rune) bool { return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' }) >= 0 {
		return false
	}
	constraint, err := semver.NewConstraint(declared)
	if err != nil {
		return false
	}
	v, err := semver.NewVersion(version)
	return err == nil && constraint.Check(v)
}
