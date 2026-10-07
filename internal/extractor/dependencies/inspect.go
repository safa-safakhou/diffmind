package dependencies

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mohammad-safakhou/diffmind/internal/extractor/sourcefilter"
)

type Options struct {
	// MavenCache is optional and read-only. No remote artifact resolution occurs.
	MavenCache string
	Include    []string
	Exclude    []string
}

func DefaultOptions() Options {
	home, _ := os.UserHomeDir()
	if home == "" {
		return Options{}
	}
	return Options{MavenCache: filepath.Join(home, ".m2", "repository")}
}

func Inspect(ctx context.Context, root string, opts Options) (Inventory, error) {
	i := Inventory{Dependencies: []Fact{}}
	policy, err := sourcefilter.NewPolicy(opts.Include, opts.Exclude)
	if err != nil {
		return i, err
	}
	var manifests []string
	err = filepath.WalkDir(root, func(p string, d fs.DirEntry, e error) error {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if e != nil {
			i.Limitations = append(i.Limitations, "Cannot read dependency input: "+p)
			return nil
		}
		if d.IsDir() {
			if p != root && sourcefilter.SkipDirName(d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if d.Type()&os.ModeSymlink != 0 {
			return nil
		}
		r, _ := filepath.Rel(root, p)
		r = filepath.ToSlash(r)
		if !policy.Allows(r) {
			return nil
		}
		switch d.Name() {
		case "pom.xml", "package.json", "pyproject.toml", "requirements.txt", "go.mod", "Gemfile", "Gemfile.lock", "build.gradle", "build.gradle.kts":
			manifests = append(manifests, r)
		default:
			if strings.HasPrefix(d.Name(), "requirements") && strings.HasSuffix(d.Name(), ".txt") {
				manifests = append(manifests, r)
			}
		}
		return nil
	})
	if err != nil {
		return i, err
	}
	sort.Strings(manifests)
	nodeLocks := map[string]*nodeLock{}
	for _, rel := range manifests {
		if ctx.Err() != nil {
			return i, ctx.Err()
		}
		var e error
		switch filepath.Base(rel) {
		case "pom.xml":
			e = inspectMaven(root, rel, opts.MavenCache, &i)
		case "package.json":
			e = inspectNode(root, rel, &i, nodeLocks)
		case "pyproject.toml":
			e = inspectPythonProject(root, rel, &i)
		case "go.mod":
			e = inspectGo(root, rel, &i)
		case "Gemfile", "Gemfile.lock":
			e = inspectRuby(root, rel, &i)
		case "build.gradle", "build.gradle.kts":
			e = inspectGradle(root, rel, &i)
		default:
			e = inspectRequirements(root, rel, &i, map[string]bool{})
		}
		if e != nil {
			i.Limitations = append(i.Limitations, fmt.Sprintf("%s: %v", rel, e))
		}
	}
	normalize(&i)
	return i, nil
}

func readFile(root, rel string) ([]byte, error) {
	p := filepath.Join(root, filepath.FromSlash(rel))
	actual, e := filepath.EvalSymlinks(p)
	if e != nil {
		return nil, e
	}
	actualRoot, e := filepath.EvalSymlinks(root)
	if e != nil {
		return nil, e
	}
	if !within(actualRoot, actual) {
		return nil, fmt.Errorf("dependency input outside repository")
	}
	info, err := os.Lstat(p)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("not a regular file")
	}
	if info.Size() > 16<<20 {
		return nil, fmt.Errorf("dependency input exceeds 16 MiB")
	}
	return os.ReadFile(p)
}

func within(root, p string) bool {
	r, e := filepath.Rel(root, p)
	return e == nil && r != ".." && !strings.HasPrefix(r, ".."+string(filepath.Separator)) && !filepath.IsAbs(r)
}

func module(rel string) string { return filepath.ToSlash(filepath.Dir(rel)) }
