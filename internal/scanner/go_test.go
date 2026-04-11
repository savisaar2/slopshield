package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGoScanner_Scan(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "slopshield-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	goMod := `module example.com/test

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/stretchr/testify v1.8.4
)

require github.com/spf13/cobra v1.8.0
`
	err = os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644)
	if err != nil {
		t.Fatalf("failed to write go.mod: %v", err)
	}

	scanner := &GoScanner{}
	deps, err := scanner.Scan(tmpDir)
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	if len(deps) != 3 {
		t.Errorf("expected 3 dependencies, got %d", len(deps))
	}

	expected := map[string]bool{
		"github.com/gin-gonic/gin":     true,
		"github.com/stretchr/testify":  true,
		"github.com/spf13/cobra":      true,
	}

	for _, d := range deps {
		if !expected[d.Name] {
			t.Errorf("unexpected dependency found: %s", d.Name)
		}
		delete(expected, d.Name)
	}

	if len(expected) != 0 {
		t.Errorf("some expected dependencies were not found: %v", expected)
	}
}
