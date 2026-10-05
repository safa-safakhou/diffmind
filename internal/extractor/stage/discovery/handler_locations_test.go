package discovery

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	astpkg "github.com/mohammad-safakhou/diffmind/internal/extractor/ast"
	_ "github.com/mohammad-safakhou/diffmind/internal/extractor/detectors/register"
)

func TestFlaskBindingRetainsDecoratorAndOwnHandlerBody(t *testing.T) {
	root := t.TempDir()
	source := "from flask import Flask, request\napp = Flask(__name__)\n\n@app.route('/login', methods=['GET', 'POST'])\ndef login():\n    name = request.form['account_name']\n    return name\n\n@app.route('/other')\ndef other():\n    return 'other'\n"
	if err := os.WriteFile(filepath.Join(root, "app.py"), []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	idx, err := astpkg.Build(context.Background(), root, "python", 1)
	if err != nil {
		t.Fatal(err)
	}
	var count int
	for _, binding := range idx.Frameworks {
		if binding.Kind != "http_handler" {
			continue
		}
		entity, ok := EntityFromFrameworkBinding(idx, objectiveByType(t, "http_route"), binding)
		if !ok {
			t.Fatal("route rejected")
		}
		if len(entity.Locations) != 2 {
			t.Fatalf("%s locations = %+v", entity.Name, entity.Locations)
		}
		count++
		decl, body := entity.Locations[0], entity.Locations[1]
		if entity.Name == "GET /other" {
			if decl.StartLine != 9 || body.StartLine != 10 || body.EndLine != 11 {
				t.Fatalf("other route includes unrelated source: %+v", entity.Locations)
			}
		} else if decl.StartLine != 4 || decl.EndLine != 4 || body.StartLine != 5 || body.EndLine != 7 {
			t.Fatalf("login source range = %+v", entity.Locations)
		}
	}
	if count != 3 {
		t.Fatalf("routes = %d, want GET/POST login and GET other", count)
	}
}

func TestBindingHandlerLocationRejectsUnprovenOwnership(t *testing.T) {
	sym := astpkg.SymbolDef{Name: "login", Qualified: "login", Kind: astpkg.SymbolKindFunction, File: "app.py",
		Range:       astpkg.Range{StartLine: 4, EndLine: 8},
		Annotations: []astpkg.Annotation{{Range: astpkg.Range{StartLine: 3, EndLine: 3}}}}
	for _, tc := range []struct {
		name   string
		symbol string
		line   uint32
		defs   []astpkg.SymbolDef
		want   bool
	}{
		{"owned decorator", "login", 3, []astpkg.SymbolDef{sym}, true},
		{"owned body", "login", 4, []astpkg.SymbolDef{sym}, true},
		{"unrelated setup", "login", 20, []astpkg.SymbolDef{sym}, false},
		{"wrong receiver", "Other.login", 3, []astpkg.SymbolDef{sym}, false},
		{"missing symbol", "", 3, []astpkg.SymbolDef{sym}, false},
		{"ambiguous", "login", 3, []astpkg.SymbolDef{sym, sym}, false},
		{"unindexed file", "login", 3, nil, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			idx := &astpkg.ProjectIndex{Files: map[string]*astpkg.FileAST{"app.py": {Symbols: tc.defs}}}
			_, ok := bindingHandlerLocation(idx, astpkg.FrameworkBinding{File: "app.py", Symbol: tc.symbol, Range: astpkg.Range{StartLine: tc.line, EndLine: tc.line}})
			if ok != tc.want {
				t.Fatalf("resolved=%v want=%v", ok, tc.want)
			}
		})
	}
}
