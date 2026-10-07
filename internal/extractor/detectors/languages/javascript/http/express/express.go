package express

import (
	"github.com/mohammad-safakhou/diffmind/internal/extractor/ast"
	"github.com/mohammad-safakhou/diffmind/internal/extractor/detectors"
	"github.com/mohammad-safakhou/diffmind/internal/extractor/detectors/languages/internal/frameworkutil"
	"regexp"
	"strings"
)

func init() { ast.RegisterFrameworkDetector(&detector{}) }

// Express / Node.js

type detector struct{}

func (d *detector) Name() string { return "express" }

func (d *detector) Detect(idx *ast.ProjectIndex) []ast.FrameworkBinding {
	var out []ast.FrameworkBinding
	for _, fa := range idx.Files {
		if fa.Language != "javascript" && fa.Language != "typescript" && fa.Language != "tsx" && fa.Language != "jsx" {
			continue
		}
		for _, call := range fa.Calls {
			b := expressCallToBinding(call)
			if b != nil {
				coverage := detectors.Evaluate("javascript.http.express", b.File, idx.DependencyInventory, detectors.VersionRules("javascript.http.express"))
				if coverage.RuleID == "javascript.http.express.v5" && !validV5Path(strings.SplitN(b.Trigger, " ", 2)[1]) {
					b.RejectionReason = "Express 5 path uses removed optional/regexp syntax or an unnamed wildcard"
				}
				out = append(out, *b)
			}
		}
	}
	return out
}

// Express 5 uses path-to-regexp's named wildcards and braces for optional
// segments. Keep the verbatim pattern; do not pretend to expand it into routes.
func validV5Path(p string) bool {
	for n := 0; n < len(p); n++ {
		if p[n] == '\\' {
			n++
			continue
		}
		if strings.ContainsRune("?+()[]!", rune(p[n])) {
			return false
		}
		if p[n] == '*' && (n+1 == len(p) || !regexp.MustCompile(`[A-Za-z_"$]`).MatchString(p[n+1:n+2])) {
			return false
		}
	}
	return true
}

func expressCallToBinding(call ast.CallSite) *ast.FrameworkBinding {
	raw := call.CalleeRaw
	parts := strings.Split(raw, ".")
	if len(parts) != 2 {
		return nil
	}
	receiver, verb := parts[0], parts[1]
	if receiver != "app" && receiver != "router" {
		return nil
	}
	methods := map[string]string{
		"get": "GET", "post": "POST", "put": "PUT", "patch": "PATCH", "delete": "DELETE",
	}
	method, ok := methods[verb]
	if !ok || len(call.Arguments) < 2 || !frameworkutil.IsLiteralPathArg(call.Arguments, 0) {
		return nil
	}
	path := frameworkutil.LiteralPathArg(call.Arguments, 0)
	return &ast.FrameworkBinding{
		Framework:        "express",
		Kind:             "http_handler",
		Direction:        "inbound",
		Symbol:           call.Caller,
		Trigger:          method + " " + path,
		TriggerSource:    raw + "(" + path + ", ...)",
		File:             call.File,
		Range:            call.Range,
		ConfidenceReason: "express_receiver_literal_path_handler",
	}
}
