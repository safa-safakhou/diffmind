package spring

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/mohammad-safakhou/diffmind/internal/extractor/ast"
	"github.com/mohammad-safakhou/diffmind/internal/extractor/detectors/languages/internal/frameworkutil"
	"github.com/mohammad-safakhou/diffmind/internal/extractor/serviceconfig"
	"gopkg.in/yaml.v3"
)

// Recover only operations with a unique concrete controller implementation of
// an API from an explicitly configured Spring interface generator.
func generatedSpringBindings(idx *ast.ProjectIndex, existing []ast.FrameworkBinding) []ast.FrameworkBinding {
	seen := map[string]bool{}
	for _, b := range existing {
		if b.Kind == "http_handler" {
			seen[b.Symbol] = true
		}
	}
	var out []ast.FrameworkBinding
	emitted := map[string]bool{}
	for _, spec := range serviceconfig.SpringOpenAPISpecs(idx.RepoRoot) {
		body, err := os.ReadFile(spec.Path)
		if err != nil || len(body) > 4<<20 {
			continue
		}
		var document yaml.Node
		if yaml.Unmarshal(body, &document) != nil {
			continue
		}
		var doc struct {
			OpenAPI string                    `yaml:"openapi"`
			Paths   map[string]map[string]any `yaml:"paths"`
		}
		if document.Decode(&doc) != nil || !strings.HasPrefix(doc.OpenAPI, "3.0.") {
			continue
		}
		for path, item := range doc.Paths {
			if !strings.HasPrefix(path, "/") {
				continue
			}
			for _, verb := range []string{"get", "post", "put", "patch", "delete", "head", "options"} {
				operation, _ := item[verb].(map[string]any)
				op, _ := operation["operationId"].(string)
				if op == "" {
					continue
				}
				var matches []ast.SymbolDef
				for _, fa := range idx.Files {
					if fa.Language != "java" {
						continue
					}
					classes := frameworkutil.ClassesByName(fa)
					for _, sym := range fa.Symbols {
						if sym.Kind != ast.SymbolKindMethod || sym.Name != op || seen[sym.Qualified] {
							continue
						}
						cls := frameworkutil.EnclosingClassForSymbol(fa, sym, classes)
						if cls == nil || !frameworkutil.HasAnyAnnotation(*cls, "RestController", "Controller") {
							continue
						}
						implements := false
						for iface, owners := range fa.Implements {
							for _, owner := range owners {
								if owner != cls.Name && owner != cls.Qualified {
									continue
								}
								for _, imp := range fa.Imports {
									if imp.Path == spec.APIPackage+"."+iface && strings.HasSuffix(iface, "Api") {
										implements = true
									}
								}
							}
						}
						if implements {
							matches = append(matches, sym)
						}
					}
				}
				if len(matches) != 1 {
					continue
				}
				sym := matches[0]
				rel, _ := filepath.Rel(idx.RepoRoot, spec.Path)
				contractLine := 0
				if node := yamlPathNode(&document, "paths", path, verb, "operationId"); node != nil {
					contractLine = node.Line
				}
				key := sym.Qualified + "\x00" + verb + "\x00" + path
				if emitted[key] {
					continue
				}
				emitted[key] = true
				out = append(out, ast.FrameworkBinding{ContractFile: filepath.ToSlash(rel), ContractLine: contractLine, Framework: "spring", Kind: "http_handler", Direction: "inbound", Symbol: sym.Qualified, Trigger: strings.ToUpper(verb) + " " + path, TriggerSource: filepath.ToSlash(rel) + " operationId=" + op, File: sym.File, Range: sym.Range, ConfidenceReason: "configured_spring_openapi_interface_controller"})
			}
		}
	}
	return out
}

func yamlPathNode(node *yaml.Node, keys ...string) *yaml.Node {
	if node.Kind == yaml.DocumentNode && len(node.Content) > 0 {
		node = node.Content[0]
	}
	for _, key := range keys {
		if node.Kind != yaml.MappingNode {
			return nil
		}
		var next *yaml.Node
		for i := 0; i+1 < len(node.Content); i += 2 {
			if node.Content[i].Value == key {
				next = node.Content[i+1]
				break
			}
		}
		if next == nil {
			return nil
		}
		node = next
	}
	return node
}
