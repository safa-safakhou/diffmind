package archgraph

import (
	"github.com/mohammad-safakhou/diffmind/internal/extractor/dependencies"
	"github.com/mohammad-safakhou/diffmind/internal/workspace/model"
	"testing"
)

func TestEnvironmentChangesDoNotInventFlowChanges(t *testing.T) {
	a := reviewFixture()
	b := cloneComparison(a)
	a.Services[0].HTTPRoutes[0].Details["detector_coverage"] = map[string]any{"revision": "old", "version": "6.2.0"}
	b.Services[0].HTTPRoutes[0].Details["detector_coverage"] = map[string]any{"revision": "new", "version": "7.0.0"}
	if r := reviewed(t, a, b); r.Change != "unchanged" {
		t.Fatalf("version evidence fabricated a flow diff: %+v", r)
	}
	before := &model.RepoMetrics{DependencyInventory: &dependencies.Inventory{Dependencies: []dependencies.Fact{{Ecosystem: "maven", Name: "spring-web", Module: "api", Version: "6.2.0", Sources: []string{"pom.xml"}}}}}
	after := &model.RepoMetrics{DependencyInventory: &dependencies.Inventory{Dependencies: []dependencies.Fact{{Ecosystem: "maven", Name: "spring-web", Module: "api", Version: "7.0.0", Sources: []string{"parent.xml"}}}}}
	if changes := CompareEnvironments(before, after); len(changes) != 1 || changes[0].Before[0].Version != "6.2.0" || changes[0].After[0].Version != "7.0.0" {
		t.Fatalf("missing separate version change: %+v", changes)
	}
	after.DependencyInventory.Dependencies[0].Version = "6.2.0"
	if changes := CompareEnvironments(before, after); len(changes) != 0 {
		t.Fatalf("source metadata appeared as dependency change: %+v", changes)
	}
}
