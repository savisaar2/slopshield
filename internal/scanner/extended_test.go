package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRustScanner_Scan(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "slopshield-test-*")
	if err != nil { t.Fatalf("failed to create temp dir: %v", err) }
	defer os.RemoveAll(tmpDir)

	cargo := `[package]
name = "test"
version = "0.1.0"

[dependencies]
serde = "1.0"
tokio = { version = "1.0", features = ["full"] }

[dev-dependencies]
criterion = "0.4"`
	
	err = os.WriteFile(filepath.Join(tmpDir, "Cargo.toml"), []byte(cargo), 0644)
	if err != nil { t.Fatalf("failed to write Cargo.toml: %v", err) }

	scanner := &RustScanner{}
	deps, err := scanner.Scan(tmpDir)
	if err != nil { t.Fatalf("Scan() error = %v", err) }

	if len(deps) != 3 { t.Errorf("expected 3 dependencies, got %d", len(deps)) }
}

func TestPHPScanner_Scan(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "slopshield-test-*")
	if err != nil { t.Fatalf("failed to create temp dir: %v", err) }
	defer os.RemoveAll(tmpDir)

	composer := `{
		"require": { "monolog/monolog": "^3.0" },
		"require-dev": { "phpunit/phpunit": "^10.0" }
	}`
	err = os.WriteFile(filepath.Join(tmpDir, "composer.json"), []byte(composer), 0644)
	if err != nil { t.Fatalf("failed to write composer.json: %v", err) }

	scanner := &PHPScanner{}
	deps, err := scanner.Scan(tmpDir)
	if err != nil { t.Fatalf("Scan() error = %v", err) }
	if len(deps) != 2 { t.Errorf("expected 2 dependencies, got %d", len(deps)) }
}

func TestRubyScanner_Scan(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "slopshield-test-*")
	if err != nil { t.Fatalf("failed to create temp dir: %v", err) }
	defer os.RemoveAll(tmpDir)

	gemfile := `source 'https://rubygems.org'
gem 'rails', '~> 7.0.0'
gem "devise"`
	err = os.WriteFile(filepath.Join(tmpDir, "Gemfile"), []byte(gemfile), 0644)
	if err != nil { t.Fatalf("failed to write Gemfile: %v", err) }

	scanner := &RubyScanner{}
	deps, err := scanner.Scan(tmpDir)
	if err != nil { t.Fatalf("Scan() error = %v", err) }
	if len(deps) != 2 { t.Errorf("expected 2 dependencies, got %d", len(deps)) }
}

func TestActionScanner_Scan(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "slopshield-test-*")
	if err != nil { t.Fatalf("failed to create temp dir: %v", err) }
	defer os.RemoveAll(tmpDir)

	workflowDir := filepath.Join(tmpDir, ".github", "workflows")
	os.MkdirAll(workflowDir, 0755)

	workflow := `name: CI
jobs:
  test:
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'`
	
	err = os.WriteFile(filepath.Join(workflowDir, "ci.yml"), []byte(workflow), 0644)
	if err != nil { t.Fatalf("failed to write workflow file: %v", err) }

	scanner := &ActionScanner{}
	deps, err := scanner.Scan(tmpDir)
	if err != nil { t.Fatalf("Scan() error = %v", err) }
	if len(deps) != 2 { t.Errorf("expected 2 dependencies, got %d", len(deps)) }
}
