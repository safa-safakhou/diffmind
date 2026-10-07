package dependencies

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type mavenDependency struct {
	Group    string `xml:"groupId"`
	Artifact string `xml:"artifactId"`
	Version  string `xml:"version"`
	Scope    string `xml:"scope"`
	Type     string `xml:"type"`
}
type mavenParent struct {
	Group    string  `xml:"groupId"`
	Artifact string  `xml:"artifactId"`
	Version  string  `xml:"version"`
	Relative *string `xml:"relativePath"`
}
type mavenProperties struct{ Values map[string]string }

func (p *mavenProperties) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	p.Values = map[string]string{}
	for {
		t, e := d.Token()
		if e != nil {
			return e
		}
		switch v := t.(type) {
		case xml.StartElement:
			var value string
			if e = d.DecodeElement(&value, &v); e != nil {
				return e
			}
			p.Values[v.Name.Local] = strings.TrimSpace(value)
		case xml.EndElement:
			if v.Name == start.Name {
				return nil
			}
		}
	}
}

type mavenPOM struct {
	Group        string            `xml:"groupId"`
	Artifact     string            `xml:"artifactId"`
	Version      string            `xml:"version"`
	Parent       *mavenParent      `xml:"parent"`
	Properties   mavenProperties   `xml:"properties"`
	Dependencies []mavenDependency `xml:"dependencies>dependency"`
	Managed      []mavenDependency `xml:"dependencyManagement>dependencies>dependency"`
	Plugins      []mavenDependency `xml:"build>plugins>plugin"`
	Profiles     []struct {
		ID string `xml:"id"`
	} `xml:"profiles>profile"`
}
type pomFile struct {
	pom          mavenPOM
	file, source string
}
type managedVersion struct {
	value   string
	sources []string
}

var propertyRE = regexp.MustCompile(`\$\{([^}]+)\}`)

func substitute(value string, props map[string]string) string {
	for n := 0; n < 16; n++ {
		next := propertyRE.ReplaceAllStringFunc(value, func(p string) string {
			if v, ok := props[p[2:len(p)-1]]; ok {
				return v
			}
			return p
		})
		if next == value {
			break
		}
		value = next
	}
	return strings.TrimSpace(value)
}

func inspectMaven(root, rel, cache string, i *Inventory) error {
	chain, err := pomChain(root, filepath.Join(root, filepath.FromSlash(rel)), cache, map[string]bool{}, i)
	if err != nil {
		return err
	}
	if len(chain) == 0 {
		return nil
	}
	props := map[string]string{}
	for _, p := range chain {
		for k, v := range p.pom.Properties.Values {
			props[k] = v
		}
	}
	leaf := chain[len(chain)-1].pom
	group, version := leaf.Group, leaf.Version
	if leaf.Parent != nil {
		if group == "" {
			group = leaf.Parent.Group
		}
		if version == "" {
			version = leaf.Parent.Version
		}
	}
	props["project.groupId"] = group
	props["project.artifactId"] = leaf.Artifact
	props["project.version"] = version
	props["pom.version"] = version
	props["pom.groupId"] = group
	managed := map[string]managedVersion{}
	for _, p := range chain {
		layer := map[string]managedVersion{}
		for _, d := range p.pom.Managed {
			g, a, v := substitute(d.Group, props), substitute(d.Artifact, props), substitute(d.Version, props)
			if d.Type == "pom" && d.Scope == "import" {
				i.Dependencies = append(i.Dependencies, Fact{Ecosystem: "maven", Name: g + ":" + a, Module: module(rel), Declared: d.Version, Version: ExactVersion(v), Resolution: mavenResolution(v, "managed_exact"), Scope: "bom", Sources: []string{p.source}})
				addBOM(root, cache, g, a, v, props, layer, map[string]bool{}, i)
			}
		}
		for _, d := range p.pom.Managed {
			if d.Type == "pom" && d.Scope == "import" {
				continue
			}
			g, a, v := substitute(d.Group, props), substitute(d.Artifact, props), substitute(d.Version, props)
			if v != "" {
				layer[g+":"+a] = managedVersion{v, []string{p.source}}
			}
		}
		for k, v := range layer {
			managed[k] = v
		}
	}
	// A child declaration replaces the inherited dependency with the same key.
	deps := map[string]struct {
		dep    mavenDependency
		source string
	}{}
	for _, p := range chain {
		for _, d := range p.pom.Dependencies {
			deps[substitute(d.Group, props)+":"+substitute(d.Artifact, props)] = struct {
				dep    mavenDependency
				source string
			}{d, p.source}
		}
		if p.pom.Parent != nil {
			d := p.pom.Parent
			v := substitute(d.Version, props)
			i.Dependencies = append(i.Dependencies, Fact{Ecosystem: "maven", Name: d.Group + ":" + d.Artifact, Module: module(rel), Declared: d.Version, Version: ExactVersion(v), Resolution: mavenResolution(v, "declared_exact"), Scope: "parent", Sources: []string{p.source}})
		}
		for _, d := range p.pom.Plugins {
			g := d.Group
			if g == "" {
				g = "org.apache.maven.plugins"
			}
			v := substitute(d.Version, props)
			i.Dependencies = append(i.Dependencies, Fact{Ecosystem: "maven", Name: g + ":" + d.Artifact, Module: module(rel), Declared: d.Version, Version: ExactVersion(v), Resolution: mavenResolution(v, "declared_exact"), Scope: "tool", Sources: []string{p.source}})
		}
		if len(p.pom.Profiles) > 0 {
			i.Limitations = append(i.Limitations, rel+": Maven profile-dependent versions are not evaluated")
		}
	}
	for key, item := range deps {
		d := item.dep
		v := substitute(d.Version, props)
		sources := []string{item.source}
		resolution := "declared_exact"
		if v == "" {
			if m, ok := managed[key]; ok {
				v = m.value
				sources = append(sources, m.sources...)
				resolution = "managed_exact"
			}
		}
		i.Dependencies = append(i.Dependencies, Fact{Ecosystem: "maven", Name: key, Module: module(rel), Declared: d.Version, Version: ExactVersion(v), Resolution: mavenResolution(v, resolution), Scope: d.Scope, Sources: sources})
	}
	// Keep managed versions separately: dependency management does not prove
	// that the library is installed. Source/import evidence remains necessary.
	for key, m := range managed {
		if _, used := deps[key]; used {
			continue
		}
		if strings.HasPrefix(key, "org.springframework:") || strings.HasPrefix(key, "org.hibernate:") || strings.HasPrefix(key, "org.hibernate.orm:") || strings.HasPrefix(key, "io.awspring.cloud:") || strings.HasPrefix(key, "com.fasterxml.jackson") {
			i.Dependencies = append(i.Dependencies, Fact{Ecosystem: "maven", Name: key, Module: module(rel), Declared: m.value, Version: ExactVersion(m.value), Resolution: mavenResolution(m.value, "managed_exact"), Scope: "management", Sources: m.sources})
		}
	}
	return nil
}

func mavenResolution(v, exact string) string {
	if ExactVersion(v) != "" {
		return exact
	}
	if v == "" || strings.Contains(v, "${") {
		return "unresolved"
	}
	return "constraint"
}

func pomChain(root, file, cache string, seen map[string]bool, i *Inventory) ([]pomFile, error) {
	if seen[file] || len(seen) > 20 {
		return nil, fmt.Errorf("cyclic or deep Maven parent chain")
	}
	seen[file] = true
	p, e := readPOM(root, file, cache)
	if e != nil {
		return nil, e
	}
	var chain []pomFile
	if parent := p.pom.Parent; parent != nil {
		r := "../pom.xml"
		if parent.Relative != nil {
			r = *parent.Relative
		}
		candidate := filepath.Clean(filepath.Join(filepath.Dir(file), r))
		local := false
		if r != "" && within(root, candidate) {
			if lp, e := readPOM(root, candidate, cache); e == nil && lp.pom.Artifact == parent.Artifact && (lp.pom.Group == parent.Group || (lp.pom.Group == "" && lp.pom.Parent != nil && lp.pom.Parent.Group == parent.Group)) && (lp.pom.Version == parent.Version || (lp.pom.Version == "" && lp.pom.Parent != nil && lp.pom.Parent.Version == parent.Version)) {
				local = true
			}
		}
		if !local {
			candidate = cachedPOM(cache, parent.Group, parent.Artifact, substitute(parent.Version, p.pom.Properties.Values))
		}
		if candidate != "" {
			if parents, e := pomChain(root, candidate, cache, seen, i); e == nil {
				chain = parents
			} else {
				i.Limitations = append(i.Limitations, p.source+": parent metadata unavailable for "+parent.Group+":"+parent.Artifact+":"+parent.Version)
			}
		} else {
			i.Limitations = append(i.Limitations, p.source+": parent metadata unavailable for "+parent.Group+":"+parent.Artifact+":"+parent.Version)
		}
	}
	return append(chain, p), nil
}

func cachedPOM(cache, g, a, v string) string {
	if cache == "" || g == "" || a == "" || v == "" {
		return ""
	}
	for _, s := range []string{g, a, v} {
		if strings.ContainsAny(s, "/\\${}") || s == ".." {
			return ""
		}
	}
	return filepath.Join(cache, strings.ReplaceAll(g, ".", string(filepath.Separator)), a, v, a+"-"+v+".pom")
}

func readPOM(root, file, cache string) (pomFile, error) {
	actual, e := filepath.EvalSymlinks(file)
	if e != nil {
		return pomFile{}, e
	}
	actualRoot, _ := filepath.EvalSymlinks(root)
	actualCache, _ := filepath.EvalSymlinks(cache)
	if !within(actualRoot, actual) && (actualCache == "" || !within(actualCache, actual)) {
		return pomFile{}, fmt.Errorf("POM outside allowed roots")
	}
	source := ""
	if within(root, file) {
		r, _ := filepath.Rel(root, file)
		source = filepath.ToSlash(r)
	} else if cache != "" && within(cache, file) {
		r, _ := filepath.Rel(cache, file)
		source = "maven-cache:" + filepath.ToSlash(r)
	} else {
		return pomFile{}, fmt.Errorf("POM outside allowed roots")
	}
	info, e := os.Lstat(file)
	if e != nil {
		return pomFile{}, e
	}
	if !info.Mode().IsRegular() || info.Size() > 4<<20 {
		return pomFile{}, fmt.Errorf("invalid POM file")
	}
	b, e := os.ReadFile(file)
	if e != nil {
		return pomFile{}, e
	}
	p := pomFile{file: file, source: source}
	e = xml.Unmarshal(b, &p.pom)
	return p, e
}

func addBOM(root, cache, g, a, v string, props map[string]string, managed map[string]managedVersion, seen map[string]bool, i *Inventory) {
	key := g + ":" + a + ":" + v
	if seen[key] || len(seen) > 20 {
		return
	}
	seen[key] = true
	file := cachedPOM(cache, g, a, v)
	if file == "" {
		i.Limitations = append(i.Limitations, "BOM metadata unavailable for "+key)
		return
	}
	chain, e := pomChain(root, file, cache, map[string]bool{}, i)
	if e != nil {
		i.Limitations = append(i.Limitations, "BOM metadata unavailable for "+key)
		return
	}
	bomProps := map[string]string{}
	_ = props // An imported BOM has its own property model.
	for _, p := range chain {
		for k, v := range p.pom.Properties.Values {
			bomProps[k] = v
		}
	}
	last := chain[len(chain)-1].pom
	bomProps["project.version"] = last.Version
	bomProps["project.groupId"] = last.Group
	bomProps["project.artifactId"] = last.Artifact
	effective := map[string]managedVersion{}
	// Imported BOM properties have their own model, not the importing project's.
	for _, p := range chain {
		layer := map[string]managedVersion{}
		for _, d := range p.pom.Managed {
			g, a, v := substitute(d.Group, bomProps), substitute(d.Artifact, bomProps), substitute(d.Version, bomProps)
			if d.Type == "pom" && d.Scope == "import" {
				addBOM(root, cache, g, a, v, bomProps, layer, seen, i)
			}
		}
		for _, d := range p.pom.Managed {
			if d.Type == "pom" && d.Scope == "import" {
				continue
			}
			g, a, v := substitute(d.Group, bomProps), substitute(d.Artifact, bomProps), substitute(d.Version, bomProps)
			if v != "" {
				layer[g+":"+a] = managedVersion{v, []string{p.source}}
			}
		}
		for k, v := range layer {
			effective[k] = v
		}
	}
	for k, v := range effective {
		if _, exists := managed[k]; !exists {
			managed[k] = v
		}
	}
}
