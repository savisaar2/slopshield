package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPubScanner_Scan(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "slopshield-test-*")
	if err != nil { t.Fatalf("failed to create temp dir: %v", err) }
	defer os.RemoveAll(tmpDir)

	pubspec := `name: test_project
dependencies:
  flutter:
    sdk: flutter
  http: ^1.1.0
dev_dependencies:
  test: ^1.24.0`
	
	err = os.WriteFile(filepath.Join(tmpDir, "pubspec.yaml"), []byte(pubspec), 0644)
	if err != nil { t.Fatalf("failed to write pubspec.yaml: %v", err) }

	scanner := &PubScanner{}
	deps, err := scanner.Scan(tmpDir)
	if err != nil { t.Fatalf("Scan() error = %v", err) }

	if len(deps) != 3 { t.Errorf("expected 3 dependencies, got %d", len(deps)) }
}
