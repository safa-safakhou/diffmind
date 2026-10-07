package dependencies

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func write(t *testing.T, root, name, body string) {
	t.Helper()
	p := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(body), 0644); err != nil {
		t.Fatal(err)
	}
}
func inspect(t *testing.T, root string, opts Options) Inventory {
	t.Helper()
	i, e := Inspect(context.Background(), root, opts)
	if e != nil {
		t.Fatal(e)
	}
	return i
}
func fact(t *testing.T, i Inventory, name, m string) Fact {
	t.Helper()
	for _, f := range i.Dependencies {
		if f.Name == name && f.Module == m {
			return f
		}
	}
	t.Fatalf("missing %s in %s: %+v", name, m, i)
	return Fact{}
}

func TestMavenParentVersionIsNotFrameworkVersion(t *testing.T) {
	root, cache := t.TempDir(), t.TempDir()
	write(t, root, "pom.xml", `<project><parent><groupId>company</groupId><artifactId>parent</artifactId><version>7.0.1</version><relativePath/></parent><artifactId>app</artifactId><dependencies><dependency><groupId>org.springframework</groupId><artifactId>spring-web</artifactId></dependency></dependencies></project>`)
	write(t, cache, "company/parent/7.0.1/parent-7.0.1.pom", `<project><groupId>company</groupId><artifactId>parent</artifactId><version>7.0.1</version><properties><spring.version>6.2.0</spring.version></properties><dependencyManagement><dependencies><dependency><groupId>org.springframework</groupId><artifactId>spring-web</artifactId><version>${spring.version}</version></dependency></dependencies></dependencyManagement></project>`)
	i := inspect(t, root, Options{MavenCache: cache})
	f := fact(t, i, "org.springframework:spring-web", ".")
	if f.Version != "6.2.0" || f.Resolution != "managed_exact" {
		t.Fatalf("wrong framework version: %+v", f)
	}
	if f.Sources[len(f.Sources)-1] != "pom.xml" {
		t.Fatalf("missing declaration evidence: %+v", f)
	}
}

func TestLocalMavenParentPropertiesCanBeOverridden(t *testing.T) {
	root := t.TempDir()
	write(t, root, "pom.xml", `<project><groupId>company</groupId><artifactId>parent</artifactId><version>1.0.0</version><properties><framework.version>6.2.0</framework.version></properties><dependencyManagement><dependencies><dependency><groupId>org.springframework</groupId><artifactId>spring-web</artifactId><version>${framework.version}</version></dependency></dependencies></dependencyManagement></project>`)
	write(t, root, "api/pom.xml", `<project><parent><groupId>company</groupId><artifactId>parent</artifactId><version>1.0.0</version></parent><artifactId>api</artifactId><properties><framework.version>7.0.0</framework.version></properties><dependencies><dependency><groupId>org.springframework</groupId><artifactId>spring-web</artifactId></dependency></dependencies></project>`)
	i := inspect(t, root, Options{})
	if f := fact(t, i, "org.springframework:spring-web", "api"); f.Version != "7.0.0" {
		t.Fatalf("child override ignored: %+v", f)
	}
}

func TestMavenMissingParentDoesNotGuessLatestCachedVersion(t *testing.T) {
	root, cache := t.TempDir(), t.TempDir()
	write(t, root, "pom.xml", `<project><parent><groupId>company</groupId><artifactId>parent</artifactId><version>1.0.0</version></parent><dependencies><dependency><groupId>org.springframework</groupId><artifactId>spring-web</artifactId></dependency></dependencies></project>`)
	write(t, cache, "company/parent/9.0.0/parent-9.0.0.pom", `<project><dependencyManagement><dependencies><dependency><groupId>org.springframework</groupId><artifactId>spring-web</artifactId><version>7.0.0</version></dependency></dependencies></dependencyManagement></project>`)
	i := inspect(t, root, Options{MavenCache: cache})
	if f := fact(t, i, "org.springframework:spring-web", "."); f.Version != "" || f.Resolution != "unresolved" {
		t.Fatalf("guessed version: %+v", f)
	}
	if len(i.Limitations) == 0 {
		t.Fatal("missing parent was silent")
	}
}

func TestNpmWorkspaceVersionsAreScopedAndPinned(t *testing.T) {
	root := t.TempDir()
	write(t, root, "package.json", `{"dependencies":{"express":"^5.0.0"}}`)
	write(t, root, "apps/legacy/package.json", `{"dependencies":{"express":"^4.0.0"}}`)
	write(t, root, "pnpm-lock.yaml", `lockfileVersion: '9.0'
importers:
  .:
    dependencies:
      express: {specifier: ^5.0.0, version: 5.1.0}
  apps/legacy:
    dependencies:
      express: {specifier: ^4.0.0, version: 4.21.2}
`)
	i := inspect(t, root, Options{})
	if f := fact(t, i, "express", "."); f.Version != "5.1.0" || f.Resolution != "lockfile" {
		t.Fatalf("root: %+v", f)
	}
	if f := fact(t, i, "express", "apps/legacy"); f.Version != "4.21.2" {
		t.Fatalf("workspace: %+v", f)
	}
	fs := i.ForFile("apps/legacy/src/server.ts")
	if len(fs) != 1 || fs[0].Version != "4.21.2" {
		t.Fatalf("cross-module version leakage: %+v", fs)
	}
}

func TestUnknownLockfileAndStaleSpecifierAreNotResolved(t *testing.T) {
	for _, body := range []string{`lockfileVersion: '99.0'
importers: {'.': {dependencies: {express: {specifier: ^5, version: 5.1.0}}}}`, `lockfileVersion: '9.0'
importers: {'.': {dependencies: {express: {specifier: ^4, version: 4.21.2}}}}`} {
		root := t.TempDir()
		write(t, root, "package.json", `{"dependencies":{"express":"^5"}}`)
		write(t, root, "pnpm-lock.yaml", body)
		i := inspect(t, root, Options{})
		if f := fact(t, i, "express", "."); f.Version != "" {
			t.Fatalf("unsafe lock resolution: %+v", f)
		}
		if len(i.Limitations) == 0 {
			t.Fatal("missing limitation")
		}
	}
}

func TestNpmNestedResolutionAndYarnSelectors(t *testing.T) {
	root := t.TempDir()
	write(t, root, "apps/api/package.json", `{"dependencies":{"express":"^4.0.0"}}`)
	write(t, root, "package-lock.json", `{"lockfileVersion":3,"packages":{"node_modules/express":{"version":"5.1.0"},"apps/api/node_modules/express":{"version":"4.21.2"}}}`)
	if f := fact(t, inspect(t, root, Options{}), "express", "apps/api"); f.Version != "4.21.2" {
		t.Fatalf("wrong npm scope: %+v", f)
	}
	root = t.TempDir()
	write(t, root, "package.json", `{"dependencies":{"express":"^4.0.0"}}`)
	write(t, root, "yarn.lock", `express@^4.0.0:
  version "4.21.2"
express@^5.0.0:
  version "5.1.0"
`)
	if f := fact(t, inspect(t, root, Options{}), "express", "."); f.Version != "4.21.2" {
		t.Fatalf("wrong yarn selector: %+v", f)
	}
}

func TestPythonPinsConstraintsAndConditionalVersions(t *testing.T) {
	root := t.TempDir()
	write(t, root, "requirements.txt", "FastAPI==0.115.0\npydantic>=1,<3\nFlask==3.1.0; python_version >= '3.9'\n-r api/requirements.txt\n")
	write(t, root, "api/requirements.txt", "Django==5.2.0\n")
	i := inspect(t, root, Options{})
	if f := fact(t, i, "fastapi", "."); f.Version != "0.115.0" {
		t.Fatalf("exact pin lost: %+v", f)
	}
	if f := fact(t, i, "pydantic", "."); f.Version != "" || f.Resolution != "constraint" {
		t.Fatalf("constraint treated as version: %+v", f)
	}
	if f := fact(t, i, "flask", "."); f.Conditions == "" {
		t.Fatal("environment marker discarded")
	}
	if f := fact(t, i, "django", "api"); f.Version != "5.2.0" {
		t.Fatal("included requirements missing")
	}
}

func TestPythonProjectUvLockAndPoetry(t *testing.T) {
	root := t.TempDir()
	write(t, root, "pyproject.toml", `[project]
name = "api"
dependencies = ["fastapi>=0.100", "pydantic>=1,<3"]
[tool.poetry.dependencies]
python = ">=3.10"
flask = "^3.0"
`)
	write(t, root, "uv.lock", `version = 1
[[package]]
name = "fastapi"
version = "0.115.0"
[[package]]
name = "pydantic"
version = "1.10.0"
[[package]]
name = "pydantic"
version = "2.11.0"
`)
	i := inspect(t, root, Options{})
	if f := fact(t, i, "fastapi", "."); f.Version != "0.115.0" || f.Resolution != "lockfile" {
		t.Fatalf("uv pin missing: %+v", f)
	}
	if f := fact(t, i, "pydantic", "."); f.Version != "" {
		t.Fatal("ambiguous uv version guessed")
	}
	if f := fact(t, i, "flask", "."); f.Declared != "^3.0" {
		t.Fatal("Poetry dependency lost")
	}
}

func TestDependencyInspectorDoesNotFollowSymlinks(t *testing.T) {
	root, outside := t.TempDir(), t.TempDir()
	write(t, outside, "package.json", `{"dependencies":{"express":"5.1.0"}}`)
	if e := os.Symlink(filepath.Join(outside, "package.json"), filepath.Join(root, "package.json")); e != nil {
		t.Fatal(e)
	}
	if i := inspect(t, root, Options{}); len(i.Dependencies) != 0 {
		t.Fatalf("followed outside source: %+v", i)
	}
}

func TestGoReplacementFormsRemainUnresolved(t *testing.T) {
	for _, replace := range []string{
		"replace gorm.io/gorm => ../local",
		"replace gorm.io/gorm v1.25.0 => example.com/fork v1.26.0",
		"replace (\n gorm.io/gorm => ../local\n)",
	} {
		root := t.TempDir()
		write(t, root, "go.mod", "module example.com/api\nrequire gorm.io/gorm v1.25.0\n"+replace+"\n")
		f := fact(t, inspect(t, root, Options{}), "gorm.io/gorm", ".")
		if f.Version != "" || f.Resolution != "unresolved" {
			t.Fatalf("replacement treated as upstream version: %+v", f)
		}
	}
}

func TestPythonChangedConstraintAndFutureLockStayUnknown(t *testing.T) {
	for _, lock := range []string{
		"version = 1\n[[package]]\nname = 'fastapi'\nversion = '0.99.0'\n",
		"version = 99\n[[package]]\nname = 'fastapi'\nversion = '0.115.0'\n",
	} {
		root := t.TempDir()
		write(t, root, "pyproject.toml", "[project]\ndependencies = ['fastapi>=0.100,<1']\n")
		write(t, root, "uv.lock", lock)
		if f := fact(t, inspect(t, root, Options{}), "fastapi", "."); f.Version != "" {
			t.Fatalf("unsafe Python lock pin: %+v", f)
		}
	}
}

func TestDifferentPythonEnvironmentMarkersAreRetained(t *testing.T) {
	root := t.TempDir()
	write(t, root, "requirements.txt", "flask==3.1.0; python_version >= '3.10'\nflask==3.1.0; python_version < '3.10'\n")
	i := inspect(t, root, Options{})
	if len(i.Dependencies) != 2 {
		t.Fatalf("conditional evidence collapsed: %+v", i.Dependencies)
	}
}
