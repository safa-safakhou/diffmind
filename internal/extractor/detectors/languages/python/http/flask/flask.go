package flask

import (
	"github.com/mohammad-safakhou/diffmind/internal/extractor/ast"
	"github.com/mohammad-safakhou/diffmind/internal/extractor/detectors/languages/internal/frameworkutil"
	"path/filepath"
	"sort"
	"strings"
)

func init() { ast.RegisterFrameworkDetector(&detector{}) }

// Flask (Python)

type detector struct{}

func (d *detector) Name() string { return "flask" }

func (d *detector) Detect(idx *ast.ProjectIndex) []ast.FrameworkBinding {
	var out []ast.FrameworkBinding
	prefixes := flaskBlueprintPrefixes(idx)
	for _, fa := range idx.Files {
		if fa.Language != "python" {
			continue
		}
		for _, sym := range fa.Symbols {
			for _, ann := range sym.Annotations {
				if b := flaskAnnotationToBinding(sym, ann, prefixes[fa.Path]); b != nil {
					route := strings.TrimPrefix(b.Trigger, "GET ")
					for _, method := range flaskRouteMethods(ann.Arguments) {
						binding := *b
						binding.Trigger = method + " " + route
						out = append(out, binding)
					}
				}
			}
		}
	}
	return out
}

func flaskAnnotationToBinding(sym ast.SymbolDef, ann ast.Annotation, prefixes map[string]string) *ast.FrameworkBinding {
	receiver, method := splitDecoratorName(ann.Name)
	if receiver == "" || !flaskRouteReceiver(receiver) {
		return nil
	}
	method = strings.ToLower(method)
	if method != "route" {
		return nil
	}
	path := frameworkutil.ExtractFirstStringArg(ann.Arguments)
	if path == "" {
		return nil
	}
	prefix := flaskBlueprintPrefix(prefixes, receiver)
	if prefix == unknownPrefix {
		return nil
	}
	if prefix != "" {
		path = frameworkutil.JoinPath(prefix, path)
	}
	httpMethod := "GET"
	reason := "flask_decorator_literal_path"
	if prefix != "" {
		reason = "flask_decorator_literal_path_blueprint_prefix"
	}
	return &ast.FrameworkBinding{
		Framework:     "flask",
		Kind:          "http_handler",
		Direction:     "inbound",
		Symbol:        sym.Qualified,
		Trigger:       httpMethod + " " + path,
		TriggerSource: "@" + ann.Name + "(" + ann.Arguments + ")",
		File:          sym.File,
		// Route identity is declared by this decorator, not the function body.
		Range:            ann.Range,
		ConfidenceReason: reason,
	}
}

const unknownPrefix = "\x00unknown"

func flaskBlueprintPrefixes(idx *ast.ProjectIndex) map[string]map[string]string {
	out := map[string]map[string]string{}
	if idx == nil {
		return out
	}
	type definition struct{ file, name string }
	var definitions []definition
	files := make([]string, 0, len(idx.Files))
	for file := range idx.Files {
		files = append(files, file)
	}
	sort.Strings(files)
	for _, file := range files {
		fa := idx.Files[file]
		if fa.Language != "python" {
			continue
		}
		out[fa.Path] = map[string]string{}
		for _, call := range fa.Calls {
			if (call.CalleeRaw != "Blueprint" && call.CalleeRaw != "flask.Blueprint") || call.AssignedTo == "" {
				continue
			}
			prefix, _ := flaskURLPrefix(call.Arguments)
			if _, duplicate := out[fa.Path][call.AssignedTo]; duplicate {
				prefix = unknownPrefix
			} else {
				definitions = append(definitions, definition{fa.Path, call.AssignedTo})
			}
			out[fa.Path][call.AssignedTo] = prefix
		}
	}
	// Registration overrides constructor defaults. Resolve to a specific
	// definition; reused names in different files must never share a prefix.
	registrations := map[definition]string{}
	for _, file := range files {
		fa := idx.Files[file]
		if fa.Language != "python" {
			continue
		}
		for _, call := range fa.Calls {
			if !strings.HasSuffix(call.CalleeRaw, ".register_blueprint") && call.CalleeRaw != "register_blueprint" {
				continue
			}
			if len(call.Arguments) == 0 {
				continue
			}
			prefix, present := flaskURLPrefix(call.Arguments)
			if !present {
				continue
			}
			receiver := strings.TrimSpace(call.Arguments[0].Source)
			var candidates []definition
			for _, d := range definitions {
				matches := receiver == d.name
				if dot := strings.LastIndex(receiver, "."); dot >= 0 {
					module := strings.ReplaceAll(receiver[:dot], ".", "/") + ".py"
					matches = receiver[dot+1:] == d.name && (filepath.ToSlash(d.file) == module || strings.HasSuffix(filepath.ToSlash(d.file), "/"+module))
				}
				if matches {
					candidates = append(candidates, d)
				}
			}
			if len(candidates) > 1 {
				var local []definition
				for _, d := range candidates {
					if d.file == fa.Path {
						local = append(local, d)
					}
				}
				candidates = local
			}
			if len(candidates) != 1 {
				continue
			}
			d := candidates[0]
			if previous, ok := registrations[d]; ok && previous != prefix {
				registrations[d] = unknownPrefix
			} else {
				registrations[d] = prefix
			}
		}
	}
	for d, prefix := range registrations {
		out[d.file][d.name] = prefix
	}
	return out
}

func flaskURLPrefix(args []ast.ArgumentExpr) (string, bool) {
	for _, arg := range args {
		pair := strings.SplitN(strings.TrimSpace(arg.Source), "=", 2)
		if len(pair) != 2 || strings.TrimSpace(pair[0]) != "url_prefix" {
			continue
		}
		prefix, ok := flaskLiteralString(pair[1])
		if !ok || (prefix != "" && !strings.HasPrefix(prefix, "/")) {
			return unknownPrefix, true
		}
		return prefix, true
	}
	return "", false
}

func flaskLiteralString(value string) (string, bool) {
	value = strings.TrimSpace(value)
	if len(value) < 2 || (value[0] != 39 && value[0] != 34) || value[len(value)-1] != value[0] {
		return "", false
	}
	body := value[1 : len(value)-1]
	if strings.ContainsAny(body, "\\\"'\r\n") {
		return "", false
	}
	return body, true
}

func flaskBlueprintPrefix(prefixes map[string]string, receiver string) string {
	return prefixes[strings.TrimSpace(receiver)]
}

func splitDecoratorName(name string) (receiver, method string) {
	name = strings.TrimSpace(name)
	i := strings.LastIndex(name, ".")
	if i < 0 {
		return "", name
	}
	return strings.TrimSpace(name[:i]), strings.TrimSpace(name[i+1:])
}

func flaskRouteReceiver(receiver string) bool {
	receiver = strings.ToLower(strings.TrimSpace(receiver))
	if receiver == "app" || receiver == "application" {
		return true
	}
	if strings.Contains(receiver, "blueprint") || strings.HasSuffix(receiver, "_bp") || strings.HasSuffix(receiver, "bp") {
		return true
	}
	for _, r := range receiver {
		if !(r == '_' || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')) {
			return false
		}
	}
	return receiver != ""
}

func flaskRouteMethods(args string) []string {
	named, _, _ := frameworkutil.ParseAnnotationArgs(args)
	raw, present := named["methods"]
	if !present {
		return []string{"GET"}
	}
	raw = strings.TrimSpace(raw)
	if len(raw) < 2 || !((raw[0] == '[' && raw[len(raw)-1] == ']') || (raw[0] == '(' && raw[len(raw)-1] == ')')) {
		return nil
	}
	parts := strings.Split(raw[1:len(raw)-1], ",")
	var methods []string
	seen := map[string]bool{}
	for i, part := range parts {
		if strings.TrimSpace(part) == "" && i == len(parts)-1 {
			continue
		}
		method, ok := flaskLiteralString(part)
		if !ok || method == "" {
			return nil
		}
		method = strings.ToUpper(method)
		for _, c := range method {
			if c < 'A' || c > 'Z' {
				return nil
			}
		}
		if !seen[method] {
			methods = append(methods, method)
			seen[method] = true
		}
	}
	return methods
}
