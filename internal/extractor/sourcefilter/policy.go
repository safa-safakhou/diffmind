package sourcefilter

import (
	"fmt"
	"path"
	"strings"
)

// Policy applies repository-relative globs. ** matches zero or more complete
// path segments; exclusion takes precedence. Traversal safety rules still apply.
type Policy struct{ include, exclude []string }

func NewPolicy(include, exclude []string) (Policy, error) {
	for _, patterns := range [][]string{include, exclude} {
		for _, pattern := range patterns {
			if pattern == "" || strings.HasPrefix(pattern, "/") || strings.Contains(pattern, "\\") {
				return Policy{}, fmt.Errorf("invalid repository-relative path glob %q", pattern)
			}
			for _, segment := range strings.Split(pattern, "/") {
				if segment == ".." {
					return Policy{}, fmt.Errorf("path glob escapes repository: %q", pattern)
				}
				if segment != "**" {
					if _, err := path.Match(segment, ""); err != nil {
						return Policy{}, fmt.Errorf("invalid path glob %q: %w", pattern, err)
					}
				}
			}
		}
	}
	return Policy{append([]string(nil), include...), append([]string(nil), exclude...)}, nil
}

func (p Policy) Allows(relative string) bool {
	relative = strings.TrimPrefix(relative, "./")
	if SkipTestPath(relative) {
		return false
	}
	parts := strings.Split(relative, "/")
	for _, segment := range parts {
		if segment == ".." || SkipDirName(segment) {
			return false
		}
	}
	if SkipFileName(relative) {
		return false
	}
	matches := func(patterns []string) bool {
		for _, pattern := range patterns {
			if matchSegments(strings.Split(pattern, "/"), parts) {
				return true
			}
		}
		return false
	}
	return (len(p.include) == 0 || matches(p.include)) && !matches(p.exclude)
}

func matchSegments(pattern, parts []string) bool {
	if len(pattern) == 0 {
		return len(parts) == 0
	}
	if pattern[0] == "**" {
		for i := 0; i <= len(parts); i++ {
			if matchSegments(pattern[1:], parts[i:]) {
				return true
			}
		}
		return false
	}
	if len(parts) == 0 {
		return false
	}
	matched, _ := path.Match(pattern[0], parts[0])
	return matched && matchSegments(pattern[1:], parts[1:])
}

// SkipTestPath excludes test sources and configurations from production facts.
func SkipTestPath(relative string) bool {
	value := "/" + strings.ToLower(strings.TrimPrefix(relative, "./"))
	for _, segment := range []string{"/test/", "/tests/", "/__tests__/", "/fixtures/", "/fixture/"} {
		if strings.Contains(value, segment) {
			return true
		}
	}
	base := path.Base(value)
	return strings.Contains(base, "_test.") || strings.Contains(base, ".test.") || strings.Contains(base, ".spec.") || strings.HasSuffix(base, "test.java") || strings.HasSuffix(base, "tests.java")
}
