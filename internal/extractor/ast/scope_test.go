package ast

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestBuildHonorsConfiguredScopeForSourceAndConfig(t *testing.T) {
	root := t.TempDir()
	files := map[string]string{"diffmind-configuration.yaml": "schema: diffmind.config.v1\npaths:\n  include: [src/**]\n  exclude: [src/examples/**]\n", "src/main.go": "package main\nfunc main() {}\n", "src/application.yaml": "queue: production\n", "src/examples/demo.go": "package demo\nfunc Example() {}\n", "src/examples/application.yaml": "queue: fake\n", "other/main.go": "package other\n"}
	for name, body := range files {
		p := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0644); err != nil {
			t.Fatal(err)
		}
	}
	idx, err := Build(context.Background(), root, "go", 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(idx.Files) != 1 || idx.Files["src/main.go"] == nil || len(idx.Configs) != 1 || idx.Configs["src/application.yaml"] == nil {
		t.Fatalf("scope leak: files=%v configs=%v", idx.Files, idx.Configs)
	}
}
