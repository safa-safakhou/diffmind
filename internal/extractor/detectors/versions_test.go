package detectors

import (
	"github.com/mohammad-safakhou/diffmind/internal/extractor/dependencies"
	"testing"
)

func TestVersionSelectionStatesAndWorkspaceScope(t *testing.T) {
	for _, tt := range []struct{ version, resolution, conditions, want, rule string }{
		{"4.21.2", "lockfile", "", "validated", "javascript.http.express.v4"},
		{"5.1.0", "declared_exact", "", "validated", "javascript.http.express.v5"},
		{"5.2.99", "lockfile", "", "compatible_unverified", "javascript.http.express.v5"},
		{"6.0.0", "lockfile", "", "unsupported_version", "javascript.http.express.source-fallback"},
		{"", "constraint", "", "unknown_version", "javascript.http.express.source-fallback"},
		{"5.1.0", "declared_exact", "python_version > '3.9'", "unknown_version", "javascript.http.express.source-fallback"},
	} {
		t.Run(tt.want+tt.version, func(t *testing.T) {
			i := dependencies.Inventory{Dependencies: []dependencies.Fact{{Ecosystem: "npm", Name: "express", Module: ".", Version: "4.21.2", Resolution: "lockfile"}, {Ecosystem: "npm", Name: "express", Module: "api", Version: tt.version, Resolution: tt.resolution, Conditions: tt.conditions}}}
			c := Evaluate("javascript.http.express", "api/server.js", i, VersionRules("javascript.http.express"))
			if c.Status != tt.want || c.RuleID != tt.rule {
				t.Fatalf("unexpected selection: %+v", c)
			}
		})
	}
}

func TestVersionSelectionPreservesAmbiguityAndPrefersLockEvidence(t *testing.T) {
	i := dependencies.Inventory{Dependencies: []dependencies.Fact{{Ecosystem: "npm", Name: "express", Module: ".", Version: "4.21.2", Resolution: "declared_exact"}, {Ecosystem: "npm", Name: "express", Module: ".", Version: "5.1.0", Resolution: "lockfile"}}}
	if c := Evaluate("javascript.http.express", "server.js", i, VersionRules("javascript.http.express")); c.Status != "validated" || c.Dependencies[0].Version != "5.1.0" {
		t.Fatalf("lock ignored: %+v", c)
	}
	i.Dependencies[0].Resolution = "lockfile"
	if c := Evaluate("javascript.http.express", "server.js", i, VersionRules("javascript.http.express")); c.Status != "ambiguous_version" {
		t.Fatalf("ambiguous versions silently selected: %+v", c)
	}
}

func TestAllDeclaredRangesAreValidAndValidationIsExact(t *testing.T) {
	for _, d := range All() {
		for _, r := range d.VersionRules {
			for _, req := range r.Requires {
				if e := req.Validate(); e != nil {
					t.Fatalf("%s: %v", d.ID, e)
				}
			}
		}
	}
	for _, bad := range []Requirement{{"", "foo", ">=2"}, {"npm", "foo", "banana"}} {
		if bad.Validate() == nil {
			t.Fatal("invalid requirement accepted")
		}
	}
}
