package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNPMScanner_Scan(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "slopshield-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	pkgJSON := `{
		"dependencies": {
			"lodash": "^4.17.21"
		},
		"devDependencies": {
			"jest": "^29.0.0"
		}
	}`
	err = os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(pkgJSON), 0644)
	if err != nil {
		t.Fatalf("failed to write package.json: %v", err)
	}

	scanner := &NPMScanner{}
	deps, err := scanner.Scan(tmpDir)
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	if len(deps) != 2 {
		t.Errorf("expected 2 dependencies, got %d", len(deps))
	}

	foundLodash := false
	foundJest := false
	for _, d := range deps {
		if d.Name == "lodash" {
			foundLodash = true
		}
		if d.Name == "jest" {
			foundJest = true
		}
	}

	if !foundLodash || !foundJest {
		t.Errorf("did not find expected dependencies: lodash=%v, jest=%v", foundLodash, foundJest)
	}
}
