package ui

import "testing"

func TestAnalyzerUsesConfiguredExecutableAndExplicitOverride(t *testing.T) {
	s := &Server{}
	t.Setenv("DIFFMIND_BINARY", "")
	s.SetAnalyzerBinary("/installed outside PATH/diffmind")
	if s.analyzerExecutable() != "/installed outside PATH/diffmind" {
		t.Fatal("worker lost executable identity")
	}
	t.Setenv("DIFFMIND_BINARY", "/explicit/analyzer")
	if s.analyzerExecutable() != "/explicit/analyzer" {
		t.Fatal("override ignored")
	}
}
