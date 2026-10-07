package ast

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mohammad-safakhou/diffmind/internal/extractor/detectors"
	"github.com/mohammad-safakhou/diffmind/internal/extractor/serviceconfig"
	"github.com/mohammad-safakhou/diffmind/internal/extractor/sourcefilter"
)

func defaultModule(m string) string {
	if strings.TrimSpace(m) == "" {
		return "."
	}
	return filepath.ToSlash(filepath.Clean(m))
}

// Custom patterns only match parsed calls or method annotations, never comments,
// arbitrary string literals or a whole-file regex hit with no executable owner.
func customBindings(idx *ProjectIndex, patterns []serviceconfig.CustomPattern) []FrameworkBinding {
	var out []FrameworkBinding
	for _, p := range patterns {
		re, e := regexp.Compile(p.Regex)
		if e != nil {
			continue
		}
		var includes []string
		if p.FileGlob != "" {
			includes = []string{p.FileGlob}
		}
		policy, e := sourcefilter.NewPolicy(includes, nil)
		if e != nil {
			continue
		}
		for file, fa := range idx.Files {
			language := p.Language
			if language == "golang" {
				language = "go"
			}
			if !policy.Allows(file) || (language != "" && language != fa.Language) {
				continue
			}
			src, e := os.ReadFile(filepath.Join(idx.RepoRoot, filepath.FromSlash(file)))
			if e != nil {
				continue
			}
			appendMatch := func(r Range, owner, text string) {
				match := re.FindStringSubmatchIndex(text)
				if len(match) == 0 || strings.TrimSpace(text[:match[0]]) != "" {
					return
				}
				fields := map[string]string{}
				for k, v := range p.Fields {
					fields[k] = string(re.ExpandString(nil, v, text, match))
					fields[k] = strings.TrimSpace(fields[k])
				}
				b := FrameworkBinding{Framework: "custom:" + p.ID, Kind: p.Kind, Symbol: owner, File: file, Range: r, TriggerSource: text, DetectorIDs: []string{"custom:" + p.ID}, ConfidenceReason: "configured_pattern_on_parsed_source"}
				switch p.Kind {
				case "http_endpoint", "http_route", "http_handler":
					b.Kind = "http_handler"
					b.Direction = "inbound"
					b.Trigger = strings.ToUpper(fields["method"]) + " " + fields["path"]
					if fields["method"] == "" || !strings.HasPrefix(fields["path"], "/") {
						b.RejectionReason = "Custom HTTP rule requires a method and literal path"
					}
				case "queue_consumer":
					b.Direction = "inbound"
					b.Trigger = fields["queue"]
					if b.Trigger == "" {
						b.Trigger = fields["trigger"]
					}
				case "scheduled_job", "scheduler":
					b.Kind = "scheduler"
					b.Direction = "inbound"
					b.Trigger = fields["schedule"]
				case "outbound_http", "http_client":
					b.Kind = "http_client"
					b.Direction = "outbound"
					b.Trigger = fields["url"]
				default:
					b.RejectionReason = "Custom pattern kind has no deterministic extraction mapping"
				}
				if handler := fields["handler"]; handler != "" {
					if target := customHandler(idx, file, handler); target != "" {
						b.Symbol = target
					} else {
						b.RejectionReason = "Custom handler reference is unresolved or ambiguous"
					}
				}
				if b.Symbol == "" || strings.TrimSpace(b.Trigger) == "" {
					b.RejectionReason = "Custom rule requires a resolved executable owner and trigger"
				}
				if len(p.Requires) > 0 {
					c := detectors.Evaluate("custom:"+p.ID, file, idx.DependencyInventory, []detectors.VersionRule{{ID: "custom:" + p.ID + ".configured", Requires: p.Requires}})
					b.VersionCoverage = []detectors.Coverage{c}
					if c.Status != "compatible_unverified" && c.Status != "validated" {
						b.RejectionReason = "Custom rule version requirements are not established: " + c.Status
					}
				} else {
					b.VersionCoverage = []detectors.Coverage{{DetectorID: "custom:" + p.ID, RuleID: "custom:" + p.ID + ".configured", Revision: detectors.Revision, Status: "unversioned", Module: ".", Reason: "Configured source pattern without version requirements"}}
				}
				out = append(out, b)
			}
			for _, call := range fa.Calls {
				// Call ranges identify the callee in the parser contract. Render
				// parsed argument expressions rather than guessing byte offsets.
				args := make([]string, 0, len(call.Arguments))
				for _, a := range call.Arguments {
					args = append(args, a.Source)
				}
				appendMatch(call.Range, call.Caller, call.CalleeRaw+"("+strings.Join(args, ", ")+")")
			}
			for _, sym := range fa.Symbols {
				if sym.Kind != SymbolKindFunction && sym.Kind != SymbolKindMethod {
					continue
				}
				for _, ann := range sym.Annotations {
					if int(ann.Range.EndByte) <= len(src) && ann.Range.StartByte < ann.Range.EndByte {
						appendMatch(ann.Range, sym.Qualified, string(src[ann.Range.StartByte:ann.Range.EndByte]))
					}
				}
			}
		}
	}
	return out
}

func customHandler(idx *ProjectIndex, file, ref string) string {
	var local, remote []string
	for name, defs := range idx.Symbols {
		if name != ref && !strings.HasSuffix(name, "."+ref) {
			continue
		}
		for _, def := range defs {
			if def.Kind != SymbolKindFunction && def.Kind != SymbolKindMethod {
				continue
			}
			if def.File == file {
				local = append(local, name)
			} else {
				remote = append(remote, name)
			}
		}
	}
	if len(local) == 1 {
		return local[0]
	}
	if len(local) == 0 && len(remote) == 1 {
		return remote[0]
	}
	return ""
}
