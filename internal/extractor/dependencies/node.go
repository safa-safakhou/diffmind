package dependencies

import (
	"encoding/json"
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Masterminds/semver/v3"
	"gopkg.in/yaml.v3"
)

type nodePackage struct {
	Name                 string            `json:"name"`
	Version              string            `json:"version"`
	Dependencies         map[string]string `json:"dependencies"`
	DevDependencies      map[string]string `json:"devDependencies"`
	OptionalDependencies map[string]string `json:"optionalDependencies"`
	PeerDependencies     map[string]string `json:"peerDependencies"`
}
type nodeLock struct {
	kind     string
	packages map[string]nodePackage
	legacy   map[string]struct {
		Version string `json:"version"`
	}
	pnpm map[string]struct {
		Dependencies         map[string]pnpmDependency `yaml:"dependencies"`
		DevDependencies      map[string]pnpmDependency `yaml:"devDependencies"`
		OptionalDependencies map[string]pnpmDependency `yaml:"optionalDependencies"`
	}
	yarn map[string]string
	err  error
}
type pnpmDependency struct {
	Specifier string `yaml:"specifier"`
	Version   string `yaml:"version"`
}

func inspectNode(root, rel string, i *Inventory, cache map[string]*nodeLock) error {
	b, e := readFile(root, rel)
	if e != nil {
		return e
	}
	var p nodePackage
	if e = json.Unmarshal(b, &p); e != nil {
		return e
	}
	m := module(rel)
	for scope, deps := range map[string]map[string]string{"runtime": p.Dependencies, "dev": p.DevDependencies, "optional": p.OptionalDependencies, "peer": p.PeerDependencies} {
		for name, decl := range deps {
			f := Fact{Ecosystem: "npm", Name: name, Module: m, Declared: decl, Scope: scope, Sources: []string{rel}}
			if scope != "peer" {
				if v, lock := nodeLockedVersion(root, m, name, decl, scope, cache, i); v != "" {
					f.Version = v
					f.Resolution = "lockfile"
					f.Sources = append(f.Sources, lock)
				}
			}
			i.Dependencies = append(i.Dependencies, f)
		}
	}
	return nil
}

func nodeLockedVersion(root, m, name, decl, scope string, cache map[string]*nodeLock, i *Inventory) (string, string) {
	for dir := m; ; dir = path.Dir(dir) {
		for _, filename := range []string{"npm-shrinkwrap.json", "package-lock.json", "pnpm-lock.yaml", "yarn.lock"} {
			rel := path.Join(dir, filename)
			l, ok := cache[rel]
			if !ok {
				b, e := readFile(root, rel)
				if e != nil {
					cache[rel] = nil
					continue
				}
				l = parseNodeLock(filename, b)
				cache[rel] = l
				if l.err != nil {
					i.Limitations = append(i.Limitations, fmt.Sprintf("%s: unsupported or malformed lockfile: %v", rel, l.err))
				}
			}
			if l == nil || l.err != nil {
				continue
			}
			var v string
			relModule, _ := filepath.Rel(filepath.FromSlash(dir), filepath.FromSlash(m))
			relModule = filepath.ToSlash(relModule)
			switch l.kind {
			case "pnpm":
				importer, exists := l.pnpm[relModule]
				if !exists {
					return "", ""
				}
				deps := importer.Dependencies
				if scope == "dev" {
					deps = importer.DevDependencies
				}
				if scope == "optional" {
					deps = importer.OptionalDependencies
				}
				d := deps[name]
				if d.Specifier != decl {
					i.Limitations = append(i.Limitations, fmt.Sprintf("%s: lockfile specifier for %s differs from %s", rel, name, path.Join(m, "package.json")))
					return "", ""
				}
				v = strings.Split(d.Version, "(")[0]
			case "npm":
				base := relModule
				if base == "." {
					base = ""
				}
				for dir2 := base; ; dir2 = path.Dir(dir2) {
					key := path.Join(dir2, "node_modules", name)
					v = l.packages[key].Version
					if v != "" || dir2 == "" || dir2 == "." {
						break
					}
				}
				if v == "" && base == "" {
					v = l.legacy[name].Version
				}
			case "yarn":
				v = l.yarn[name+"@"+decl]
				if v == "" {
					v = l.yarn[name+"@npm:"+decl]
				}
			}
			if ExactVersion(v) == "" {
				return "", ""
			}
			// Overrides may intentionally differ from the declared range. PNPM
			// records the matching specifier explicitly; npm/yarn must satisfy it.
			if l.kind != "pnpm" {
				c, e := semver.NewConstraint(decl)
				sv, _ := semver.NewVersion(v)
				if e != nil || !c.Check(sv) {
					return "", ""
				}
			}
			return v, rel
		}
		if dir == "." || dir == "/" {
			break
		}
	}
	return "", ""
}

func parseNodeLock(filename string, b []byte) *nodeLock {
	l := &nodeLock{}
	switch filename {
	case "npm-shrinkwrap.json", "package-lock.json":
		var doc struct {
			LockfileVersion int                    `json:"lockfileVersion"`
			Packages        map[string]nodePackage `json:"packages"`
			Dependencies    map[string]struct {
				Version string `json:"version"`
			} `json:"dependencies"`
		}
		l.err = json.Unmarshal(b, &doc)
		l.kind = "npm"
		l.packages = doc.Packages
		l.legacy = doc.Dependencies
		if l.err == nil && (doc.LockfileVersion < 1 || doc.LockfileVersion > 3) {
			l.err = fmt.Errorf("npm lockfileVersion %d", doc.LockfileVersion)
		}
	case "pnpm-lock.yaml":
		var doc struct {
			LockfileVersion string `yaml:"lockfileVersion"`
			Importers       map[string]struct {
				Dependencies         map[string]pnpmDependency `yaml:"dependencies"`
				DevDependencies      map[string]pnpmDependency `yaml:"devDependencies"`
				OptionalDependencies map[string]pnpmDependency `yaml:"optionalDependencies"`
			} `yaml:"importers"`
		}
		l.err = yaml.Unmarshal(b, &doc)
		l.kind = "pnpm"
		l.pnpm = doc.Importers
		if l.err == nil && doc.LockfileVersion != "9.0" && doc.LockfileVersion != "6.0" {
			l.err = fmt.Errorf("pnpm lockfileVersion %s", doc.LockfileVersion)
		}
	case "yarn.lock":
		l.kind = "yarn"
		l.yarn = map[string]string{}
		if strings.Contains(string(b), "__metadata:") {
			var doc map[string]struct {
				Version string `yaml:"version"`
			}
			l.err = yaml.Unmarshal(b, &doc)
			for key, val := range doc {
				for _, selector := range strings.Split(key, ", ") {
					l.yarn[selector] = val.Version
				}
			}
		} else {
			var selectors []string
			vRE := regexp.MustCompile(`^\s+version\s+"([^"]+)"`)
			for _, line := range strings.Split(string(b), "\n") {
				if line != "" && line[0] != ' ' && strings.HasSuffix(line, ":") && !strings.HasPrefix(line, "#") {
					selectors = nil
					for _, s := range strings.Split(strings.TrimSuffix(line, ":"), ", ") {
						selectors = append(selectors, strings.Trim(s, `"`))
					}
				}
				if match := vRE.FindStringSubmatch(line); len(match) > 1 {
					for _, s := range selectors {
						l.yarn[s] = match[1]
					}
				}
			}
		}
	}
	return l
}
