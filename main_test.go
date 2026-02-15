package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIsValidDomain(t *testing.T) {
	t.Parallel()

	valid := []string{"example.com", "api.dev.example.com", "sub.test.co.uk"}
	invalid := []string{"", ".", "not_a_domain", "http://"}

	for _, domain := range valid {
		if !isValidDomain(domain) {
			t.Fatalf("expected valid domain: %s", domain)
		}
	}

	for _, domain := range invalid {
		if isValidDomain(domain) {
			t.Fatalf("expected invalid domain: %s", domain)
		}
	}
}

func TestUpdateDomainTree(t *testing.T) {
	t.Parallel()

	tree := make(DomainTree)
	updateDomainTree(tree, "api.dev.example.com")

	comNode, ok := tree["com"]
	if !ok {
		t.Fatal("expected top-level tld node 'com'")
	}

	exampleNode, ok := comNode["example.com"]
	if !ok {
		t.Fatal("expected second-level node 'example.com'")
	}

	devNode, ok := exampleNode["dev.example.com"]
	if !ok {
		t.Fatal("expected third-level node 'dev.example.com'")
	}

	if _, ok := devNode["api.dev.example.com"]; !ok {
		t.Fatal("expected leaf node 'api.dev.example.com'")
	}
}

func TestProcessInputGeneratesReport(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}

	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir temp dir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(wd)
	})

	fixturePath := filepath.Join(wd, "test_data", "h1_subs.txt")
	f, err := os.Open(fixturePath)
	if err != nil {
		t.Fatalf("open fixture: %v", err)
	}
	defer f.Close()

	processInput(f)

	reportDir := filepath.Join(tmpDir, ".report")
	entries, err := os.ReadDir(reportDir)
	if err != nil {
		t.Fatalf("read report dir: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected one report file, got %d", len(entries))
	}
	if !strings.HasPrefix(entries[0].Name(), "assetviz_report_") || !strings.HasSuffix(entries[0].Name(), ".html") {
		t.Fatalf("unexpected report file name: %s", entries[0].Name())
	}
}
