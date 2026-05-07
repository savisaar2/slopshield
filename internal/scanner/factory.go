package scanner

import (
	"os"
	"path/filepath"

	"github.com/savisaar2/slopshield/internal/registry"
)

type ScannerInfo struct {
	Scanner   Scanner
	Ecosystem registry.Ecosystem
	Filename  string
}

func DetectScanners(path string) ([]ScannerInfo, error) {
	var list []ScannerInfo

	// NPM
	if _, err := os.Stat(filepath.Join(path, "package.json")); err == nil {
		list = append(list, ScannerInfo{&NPMScanner{}, registry.EcosystemNPM, "package.json"})
	}
	// Pub (Flutter/Dart)
	if _, err := os.Stat(filepath.Join(path, "pubspec.yaml")); err == nil {
		list = append(list, ScannerInfo{&PubScanner{}, registry.EcosystemPub, "pubspec.yaml"})
	}
	// Python
	if _, err := os.Stat(filepath.Join(path, "requirements.txt")); err == nil {
		list = append(list, ScannerInfo{&PythonScanner{}, registry.EcosystemPython, "requirements.txt"})
	}
	// Go
	if _, err := os.Stat(filepath.Join(path, "go.mod")); err == nil {
		list = append(list, ScannerInfo{&GoScanner{}, registry.EcosystemGo, "go.mod"})
	}
	// Rust
	if _, err := os.Stat(filepath.Join(path, "Cargo.toml")); err == nil {
		list = append(list, ScannerInfo{&RustScanner{}, registry.EcosystemRust, "Cargo.toml"})
	}
	// PHP
	if _, err := os.Stat(filepath.Join(path, "composer.json")); err == nil {
		list = append(list, ScannerInfo{&PHPScanner{}, registry.EcosystemPHP, "composer.json"})
	}
	// Ruby
	if _, err := os.Stat(filepath.Join(path, "Gemfile")); err == nil {
		list = append(list, ScannerInfo{&RubyScanner{}, registry.EcosystemRuby, "Gemfile"})
	}
	// GitHub Actions
	if _, err := os.Stat(filepath.Join(path, ".github", "workflows")); err == nil {
		list = append(list, ScannerInfo{&ActionScanner{}, registry.EcosystemGitHub, "GitHub Actions"})
	}
	// Java / Maven
	if _, err := os.Stat(filepath.Join(path, "pom.xml")); err == nil {
		list = append(list, ScannerInfo{&JavaScanner{}, registry.EcosystemMaven, "pom.xml"})
	}
	// C# / .NET
	csprojFiles, _ := filepath.Glob(filepath.Join(path, "*.csproj"))
	if len(csprojFiles) > 0 {
		list = append(list, ScannerInfo{&CSharpScanner{}, registry.EcosystemNuGet, ".csproj"})
	}

	return list, nil
}
