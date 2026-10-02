package sourcefilter

import "testing"

func TestPolicyGlobScopeAndSafety(t *testing.T) {
	p, err := NewPolicy([]string{"src/**", "openapi.y?ml"}, []string{"**/examples/**", "src/generated/**"})
	if err != nil {
		t.Fatal(err)
	}
	for file, want := range map[string]bool{"src/main.go": true, "src/a/b/main.ts": true, "openapi.yaml": true, "docs/main.go": false, "src/examples/x.go": false, "src/generated/x.go": false, "src/node_modules/x.js": false, "src/../secret.go": false, "src/test/resources/openapi.yaml": false, "src/handler_test.go": false} {
		if got := p.Allows(file); got != want {
			t.Errorf("Allows(%q)=%v, want %v", file, got, want)
		}
	}
	for _, pattern := range []string{"[", "../**", "/outside/**"} {
		if _, err := NewPolicy([]string{pattern}, nil); err == nil {
			t.Errorf("accepted %q", pattern)
		}
	}
	p, _ = NewPolicy([]string{"**/*.go"}, nil)
	if !p.Allows("main.go") || !p.Allows("a/main.go") {
		t.Fatal("** must match zero or more segments")
	}
}
