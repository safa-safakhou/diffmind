package ast_test

import (
	"github.com/mohammad-safakhou/diffmind/internal/extractor/ast"
	"testing"
)

func TestFrameworkVersionFixtureMatrix(t *testing.T) {
	for _, tt := range []struct{ framework, manifest, file, body string }{
		{"spring", `<project><dependencies><dependency><groupId>org.springframework</groupId><artifactId>spring-web</artifactId><version>6.2.0</version></dependency></dependencies></project>`, "Controller.java", `import org.springframework.web.bind.annotation.*; @RestController public class Controller { @GetMapping("/orders") public String orders() { return "ok"; } }`},
		{"spring", `<project><dependencies><dependency><groupId>org.springframework</groupId><artifactId>spring-web</artifactId><version>7.0.0</version></dependency></dependencies></project>`, "Controller.java", `import org.springframework.web.bind.annotation.*; @RestController public class Controller { @GetMapping("/orders") public String orders() { return "ok"; } }`},
		{"fastapi", "fastapi==0.115.0\n", "app.py", "from fastapi import FastAPI\napp = FastAPI()\n@app.get('/orders')\ndef orders():\n    return []\n"},
		{"flask", "Flask==3.1.0\n", "app.py", "from flask import Flask\napp = Flask(__name__)\n@app.route('/orders', methods=['GET'])\ndef orders():\n    return []\n"},
	} {
		t.Run(tt.framework+tt.manifest, func(t *testing.T) {
			dir := t.TempDir()
			name := "requirements.txt"
			if tt.framework == "spring" {
				name = "pom.xml"
			}
			writeFile(t, dir, name, tt.manifest)
			writeFile(t, dir, tt.file, tt.body)
			idx := buildIndex(t, dir)
			found := false
			for _, b := range idx.Frameworks {
				if b.Framework == tt.framework && b.Trigger == "GET /orders" {
					found = true
					if len(b.VersionCoverage) != 1 || b.VersionCoverage[0].Status != "validated" {
						t.Fatalf("fixture validation absent: %+v", b)
					}
				}
			}
			if !found {
				t.Fatalf("fixture route not extracted: %+v", idx.Frameworks)
			}
		})
	}
}

func TestExpressVersionChangesRouteInterpretation(t *testing.T) {
	for _, tt := range []struct {
		version, path string
		accepted      bool
		status        string
	}{
		{"4.21.2", "/:id?", true, "validated"},
		{"5.1.0", "/:id?", false, "validated"},
		{"5.1.0", "/*", false, "validated"},
		{"5.1.0", "/*rest", true, "validated"},
		{"5.1.0", "/{*rest}", true, "validated"},
		{"5.2.99", "/orders", true, "compatible_unverified"},
		{"6.0.0", "/orders", true, "unsupported_version"},
	} {
		t.Run(tt.version+tt.path, func(t *testing.T) {
			dir := t.TempDir()
			writeFile(t, dir, "package.json", `{"dependencies":{"express":"`+tt.version+`"}}`)
			writeFile(t, dir, "server.js", `import express from 'express'; const app = express(); app.get('`+tt.path+`', handler); function handler(req,res) {res.send('ok')}`)
			idx := buildIndex(t, dir)
			bindings := idx.Frameworks
			if !tt.accepted {
				bindings = idx.RejectedFrameworks
			}
			var found *ast.FrameworkBinding
			for n := range bindings {
				if bindings[n].Framework == "express" {
					found = &bindings[n]
				}
			}
			if found == nil {
				t.Fatalf("route acceptance=%v wrong: %+v / %+v", tt.accepted, idx.Frameworks, idx.RejectedFrameworks)
			}
			if len(found.VersionCoverage) != 1 || found.VersionCoverage[0].Status != tt.status {
				t.Fatalf("wrong version coverage: %+v", found)
			}
		})
	}
}

func TestCustomPatternVersionGateAndExecutableEvidence(t *testing.T) {
	for _, tt := range []struct {
		version  string
		accepted bool
	}{{"2.4.0", true}, {"3.0.0", false}, {"", false}} {
		t.Run(tt.version, func(t *testing.T) {
			dir := t.TempDir()
			writeFile(t, dir, "routes.go", `package service
// RegisterWebhook("/comment", Handler)
func Handler() {}
func init() { RegisterWebhook("/orders", Handler); log.Info("RegisterWebhook(\"/text\", Handler)") }
`)
			hint := ""
			if tt.version != "" {
				hint = "dependency_versions:\n  - ecosystem: maven\n    name: company:router\n    version: " + tt.version + "\n"
			}
			writeFile(t, dir, "diffmind-configuration.yaml", hint+`patterns:
  - id: company-router-v2
    kind: http_endpoint
    language: go
    regex: 'RegisterWebhook\("(?P<path>[^"]+)",\s*(?P<handler>[A-Za-z0-9_.]+)\)'
    fields: {method: POST, path: "$path", handler: "$handler"}
    requires:
      - {ecosystem: maven, name: 'company:router', versions: '>=2 <3'}
`)
			idx := buildIndex(t, dir)
			bindings := idx.Frameworks
			if !tt.accepted {
				bindings = idx.RejectedFrameworks
			}
			count := 0
			for _, b := range bindings {
				if b.Framework == "custom:company-router-v2" {
					count++
					if b.Trigger != "POST /orders" || b.Symbol == "" {
						t.Fatalf("bad custom route: %+v", b)
					}
				}
			}
			if count != 1 {
				t.Fatalf("version gating or comment/string exclusion failed: %+v / %+v; calls=%+v", idx.Frameworks, idx.RejectedFrameworks, idx.Files["routes.go"].Calls)
			}
		})
	}
}

func TestCustomPatternDisabledAndAmbiguousHandler(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "router.go", `package app
func init() { RegisterWebhook("/orders", Handler) }
`)
	writeFile(t, dir, "one.go", `package one;func Handler() {}`)
	writeFile(t, dir, "two.go", `package two;func Handler() {}`)
	writeFile(t, dir, "diffmind-configuration.yaml", `detectors:
  disabled: ["custom:company-router"]
patterns:
  - id: company-router
    kind: http_endpoint
    language: go
    regex: 'RegisterWebhook\("(?P<path>[^"]+)",\s*(?P<handler>[A-Za-z0-9_.]+)\)'
    fields: {method: POST, path: "$path", handler: "$handler"}
`)
	idx := buildIndex(t, dir)
	if len(idx.Frameworks) != 0 || len(idx.RejectedFrameworks) != 1 {
		t.Fatalf("custom disable ignored: %+v / %+v", idx.Frameworks, idx.RejectedFrameworks)
	}
}
