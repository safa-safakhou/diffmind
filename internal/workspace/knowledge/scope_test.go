package knowledge

import (
	"context"
	"github.com/mohammad-safakhou/diffmind/internal/workspace/util"
	"os"
	"path/filepath"
	"testing"
)

func TestPackExtractionAndDetectionRespectServiceScope(t *testing.T) {
	root := t.TempDir()
	for name, body := range map[string]string{"diffmind-configuration.yaml": "paths:\n  include: [production/**]\n", "production/app.yaml": "queue: orders\n", "examples/app.yaml": "queue: demo\n"} {
		p := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0644); err != nil {
			t.Fatal(err)
		}
	}
	pack := &Pack{ID: "scoped", Extractions: []Extraction{{Name: "queues", Source: ExtractionSource{Glob: "**/app.yaml"}, Extract: []ExtractField{{Field: "queue", MapsTo: "queue"}}}}}
	results := NewEngine(util.NewLogger(util.LevelInfo)).Run(pack, root)
	if len(results) != 1 || results[0].SourceFile != "production/app.yaml" {
		t.Fatalf("extraction scope leak: %+v", results)
	}
	files, err := detectorFiles(context.Background(), root, "**/app.yaml", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 || files[0] != filepath.Join(root, "production", "app.yaml") {
		t.Fatalf("detector scope leak: %v", files)
	}
}
