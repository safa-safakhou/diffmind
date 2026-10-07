package dependencies

import (
	"regexp"
	"strings"
)

func inspectGo(root, rel string, i *Inventory) error {
	b, e := readFile(root, rel)
	if e != nil {
		return e
	}
	inRequire := false
	replacements := map[string]bool{}
	inReplace := false
	for _, line := range strings.Split(string(b), "\n") {
		parts := strings.Fields(strings.Split(line, "//")[0])
		if len(parts) == 2 && parts[0] == "replace" && parts[1] == "(" {
			inReplace = true
			continue
		}
		if len(parts) == 1 && parts[0] == ")" {
			inReplace = false
			continue
		}
		if len(parts) > 0 && parts[0] == "replace" {
			parts = parts[1:]
		} else if !inReplace {
			continue
		}
		for n, p := range parts {
			if p == "=>" && n > 0 {
				key := parts[0]
				if n == 2 {
					key += "@" + parts[1]
				}
				replacements[key] = true
				break
			}
		}
	}
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(strings.Split(line, "//")[0])
		if line == "require (" {
			inRequire = true
			continue
		}
		if line == ")" {
			inRequire = false
			continue
		}
		parts := strings.Fields(line)
		if len(parts) == 3 && parts[0] == "require" {
			parts = parts[1:]
		} else if !inRequire {
			continue
		}
		if len(parts) != 2 {
			continue
		}
		f := Fact{Ecosystem: "go", Name: parts[0], Module: module(rel), Declared: parts[1], Sources: []string{rel}, Scope: "runtime"}
		if replacements[f.Name] || replacements[f.Name+"@"+f.Declared] {
			f.Resolution = "unresolved"
			i.Limitations = append(i.Limitations, rel+": replaced module "+f.Name+" requires replacement-aware interpretation")
		}
		i.Dependencies = append(i.Dependencies, f)
	}
	return nil
}

func inspectRuby(root, rel string, i *Inventory) error {
	b, e := readFile(root, rel)
	if e != nil {
		return e
	}
	if strings.HasSuffix(rel, ".lock") {
		re := regexp.MustCompile(`(?m)^    ([A-Za-z0-9_-]+) \(([^)]+)\)$`)
		for _, m := range re.FindAllStringSubmatch(string(b), -1) {
			i.Dependencies = append(i.Dependencies, Fact{Ecosystem: "gem", Name: m[1], Module: module(rel), Version: ExactVersion(m[2]), Resolution: "lockfile", Scope: "runtime", Sources: []string{rel}})
		}
	} else {
		re := regexp.MustCompile(`(?m)^\s*gem\s+["']([^"']+)["'](?:\s*,\s*["']([^"']+)["'])?`)
		for _, m := range re.FindAllStringSubmatch(string(b), -1) {
			i.Dependencies = append(i.Dependencies, Fact{Ecosystem: "gem", Name: m[1], Module: module(rel), Declared: m[2], Scope: "runtime", Sources: []string{rel}})
		}
	}
	return nil
}

func inspectGradle(root, rel string, i *Inventory) error {
	b, e := readFile(root, rel)
	if e != nil {
		return e
	}
	re := regexp.MustCompile(`(?m)^\s*(implementation|api|compileOnly|runtimeOnly|testImplementation|classpath)\s*\(?\s*["']([^:"']+):([^:"']+):([^"']+)["']`)
	for _, m := range re.FindAllStringSubmatch(string(b), -1) {
		i.Dependencies = append(i.Dependencies, Fact{Ecosystem: "maven", Name: m[2] + ":" + m[3], Module: module(rel), Declared: m[4], Scope: m[1], Sources: []string{rel}})
	}
	i.Limitations = append(i.Limitations, rel+": Gradle variable, catalog, plugin and transitive versions require additional resolution")
	return nil
}
