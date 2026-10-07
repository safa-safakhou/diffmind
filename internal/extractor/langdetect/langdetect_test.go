package langdetect

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

// TestInspectMultiLanguageRepo verifies the parser identifies every
// language a polyglot service tree contains, with the version each
// marker file declares.
func TestInspectMultiLanguageRepo(t *testing.T) {
	root := t.TempDir()
	write := func(rel, content string) {
		full := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	// Java via Maven 17
	write("backend/pom.xml", `<project>
  <modelVersion>4.0.0</modelVersion>
  <properties>
    <maven.compiler.source>17</maven.compiler.source>
    <maven.compiler.target>17</maven.compiler.target>
  </properties>
</project>`)

	// TypeScript via npm
	write("frontend/package.json", `{
  "name": "fe",
  "engines": { "node": "20.10.0" },
  "devDependencies": { "typescript": "^5.0.0" }
}`)

	// Go 1.22
	write("api/go.mod", `module example.com/api

go 1.22
`)

	// Python via pyproject
	write("ml/pyproject.toml", `[project]
requires-python = ">=3.11,<3.13"

[tool.poetry]
name = "ml"
`)

	// .NET 8
	write("dotnet/Foo.csproj", `<Project Sdk="Microsoft.NET.Sdk">
  <PropertyGroup>
    <TargetFramework>net8.0</TargetFramework>
  </PropertyGroup>
</Project>`)

	// Ruby with .ruby-version
	write("worker/.ruby-version", "3.2.2\n")

	facts, err := Inspect(context.Background(), root)
	if err != nil {
		t.Fatalf("Inspect: %v", err)
	}

	got := map[Language]Fact{}
	for _, f := range facts {
		got[f.Language] = f
	}

	checks := []struct {
		lang    Language
		version string
		tool    string
	}{
		{LangJava, "17", "maven"},
		{LangTypeScript, "20.10.0", "npm"},
		{LangGo, "1.22", "go"},
		{LangPython, "3.11", "poetry"},
		{LangCSharp, "8.0", "dotnet"},
		{LangRuby, "3.2.2", ""},
	}
	for _, c := range checks {
		f, ok := got[c.lang]
		if !ok {
			t.Errorf("missing %s", c.lang)
			continue
		}
		if f.Version != c.version {
			t.Errorf("%s version = %q, want %q", c.lang, f.Version, c.version)
		}
		if c.tool != "" && f.BuildTool != c.tool {
			t.Errorf("%s build_tool = %q, want %q", c.lang, f.BuildTool, c.tool)
		}
	}
}

// TestInspectSkipsNoiseDirs proves we don't pull markers from
// node_modules / vendor / etc. Those directories typically contain
// thousands of third-party manifests that don't represent the
// repo's own language inventory.
func TestInspectSkipsNoiseDirs(t *testing.T) {
	root := t.TempDir()
	write := func(rel, content string) {
		full := filepath.Join(root, rel)
		os.MkdirAll(filepath.Dir(full), 0o755)
		os.WriteFile(full, []byte(content), 0o644)
	}
	// The "real" service is Java.
	write("pom.xml", `<project>
  <properties><maven.compiler.source>21</maven.compiler.source></properties>
</project>`)
	// Noise: a vendored package.json that should NOT trigger TS/JS.
	write("node_modules/foo/package.json", `{"devDependencies":{"typescript":"5"}}`)
	write("vendor/bar/package.json", `{"engines":{"node":"18"}}`)

	facts, err := Inspect(context.Background(), root)
	if err != nil {
		t.Fatal(err)
	}
	langs := []string{}
	for _, f := range facts {
		langs = append(langs, string(f.Language))
	}
	sort.Strings(langs)
	if len(langs) != 1 || langs[0] != "java" {
		t.Errorf("expected only 'java'; got %v", langs)
	}
}

// TestNormalizeJava covers Maven's historical "1.8" form.
func TestNormalizeJava(t *testing.T) {
	cases := map[string]string{
		"1.8": "8",
		"11":  "11",
		"17":  "17",
		"21":  "21",
	}
	for in, want := range cases {
		if got := normalizeJava(in); got != want {
			t.Errorf("normalizeJava(%q) = %q, want %q", in, got, want)
		}
	}
}

// TestParseTFM tests the .NET target framework parser.
func TestParseTFM(t *testing.T) {
	cases := map[string]string{
		"net8.0":         "8.0",
		"net6.0":         "6.0",
		"netstandard2.1": "standard2.1",
		"4.7.2":          "4.7.2",
	}
	for in, want := range cases {
		if got := parseTFM(in); got != want {
			t.Errorf("parseTFM(%q) = %q, want %q", in, got, want)
		}
	}
}

// TestExtractMinPythonVersion exercises the constraint parser.
func TestExtractMinPythonVersion(t *testing.T) {
	cases := map[string]string{
		">=3.10":       "3.10",
		">=3.10,<3.13": "3.10",
		"3.11":         "3.11",
		"== 3.9":       "",
		"<4":           "",
		"":             "",
	}
	for in, want := range cases {
		if got := extractMinPythonVersion(in); got != want {
			t.Errorf("extractMinPythonVersion(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestBonialToolPinsAreCompleteAndSeparateFromConstraints(t *testing.T) {
	root := t.TempDir()
	os.WriteFile(filepath.Join(root, ".tool-versions"), []byte("java corretto-25.0.1\nmaven 3.9.12\nnodejs 24.12.0\npnpm 10.28.0\npython 3.14.5\nuv 0.11.19\nruby 3.3.4\n"), 0644)
	os.WriteFile(filepath.Join(root, "package.json"), []byte(`{"engines":{"node":">=24.10.0"},"packageManager":"pnpm@10.28.0","devDependencies":{"typescript":"^5.9.3"}}`), 0644)
	facts, err := Inspect(context.Background(), root)
	if err != nil {
		t.Fatal(err)
	}
	byLang := map[Language]Fact{}
	for _, f := range facts {
		byLang[f.Language] = f
	}
	for _, tc := range []struct {
		lang                       Language
		version, tool, toolVersion string
	}{{LangJava, "corretto-25.0.1", "maven", "3.9.12"}, {LangTypeScript, "24.12.0", "pnpm", "10.28.0"}, {LangPython, "3.14.5", "uv", "0.11.19"}, {LangRuby, "3.3.4", "", ""}} {
		f := byLang[tc.lang]
		if f.Version != tc.version || f.BuildTool != tc.tool || f.BuildToolVersion != tc.toolVersion {
			t.Errorf("%s: %+v", tc.lang, f)
		}
	}
	ts := byLang[LangTypeScript]
	if ts.VersionConstraint != ">=24.10.0" || len(ts.Declarations) != 2 {
		t.Fatalf("Node pin/constraint evidence lost: %+v", ts)
	}
}

func TestPoetryVersionAndJavaGradleKotlinDSL(t *testing.T) {
	root := t.TempDir()
	os.WriteFile(filepath.Join(root, "pyproject.toml"), []byte("[tool.poetry]\nname='app'\n[tool.poetry.dependencies]\npython='^3.14'\n"), 0644)
	os.WriteFile(filepath.Join(root, "build.gradle.kts"), []byte("plugins { java }\nsourceCompatibility = JavaVersion.VERSION_21\n"), 0644)
	os.MkdirAll(filepath.Join(root, "gradle/wrapper"), 0755)
	os.WriteFile(filepath.Join(root, "gradle/wrapper/gradle-wrapper.properties"), []byte("distributionUrl=https://services.gradle.org/distributions/gradle-8.14.3-bin.zip\n"), 0644)
	facts, err := Inspect(context.Background(), root)
	if err != nil {
		t.Fatal(err)
	}
	byLang := map[Language]Fact{}
	for _, f := range facts {
		byLang[f.Language] = f
	}
	if _, ok := byLang[LangKotlin]; ok {
		t.Fatal("Kotlin build DSL invented Kotlin application")
	}
	if f := byLang[LangJava]; f.Version != "21" || f.BuildToolVersion != "8.14.3" {
		t.Fatalf("Gradle JVM/tool version: %+v", f)
	}
	if f := byLang[LangPython]; f.Version != "3.14" || f.VersionConstraint != "^3.14" {
		t.Fatalf("Poetry version: %+v", f)
	}
}

func TestUpperBoundNodeConstraintIsNotAnInstalledVersion(t *testing.T) {
	f := detectPackageJSON("package.json", []byte(`{"engines":{"node":"<22"}}`))
	if f.Version != "" || f.VersionConstraint != "<22" {
		t.Fatalf("upper bound fabricated runtime: %+v", f)
	}
}

func TestConflictingRuntimeDeclarationsRemainVisible(t *testing.T) {
	root := t.TempDir()
	os.WriteFile(filepath.Join(root, ".tool-versions"), []byte("java corretto-21.0.6\n"), 0644)
	os.WriteFile(filepath.Join(root, "pom.xml"), []byte("<project><properties><java.version>25</java.version></properties></project>"), 0644)
	facts, err := Inspect(context.Background(), root)
	if err != nil {
		t.Fatal(err)
	}
	if len(facts) != 1 || !facts[0].VersionConflict || len(facts[0].Declarations) != 2 {
		t.Fatalf("conflicting pin/compiler declaration lost: %+v", facts)
	}
}
