package scanner

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// Rust
type RustScanner struct{}
func (s *RustScanner) Scan(path string) ([]Dependency, error) {
	data, err := os.ReadFile(filepath.Join(path, "Cargo.toml"))
	if err != nil { return nil, err }
	var cargo struct {
		Deps    map[string]interface{} `toml:"dependencies"`
		DevDeps map[string]interface{} `toml:"dev-dependencies"`
	}
	if err := toml.Unmarshal(data, &cargo); err != nil { return nil, err }
	var deps []Dependency
	for name := range cargo.Deps { deps = append(deps, Dependency{Name: name, Source: "Cargo.toml"}) }
	for name := range cargo.DevDeps { deps = append(deps, Dependency{Name: name, Source: "Cargo.toml"}) }
	return deps, nil
}

// PHP
type PHPScanner struct{}
func (s *PHPScanner) Scan(path string) ([]Dependency, error) {
	data, err := os.ReadFile(filepath.Join(path, "composer.json"))
	if err != nil { return nil, err }
	var composer struct {
		Require    map[string]string `json:"require"`
		RequireDev map[string]string `json:"require-dev"`
	}
	if err := json.Unmarshal(data, &composer); err != nil { return nil, err }
	var deps []Dependency
	for name := range composer.Require { deps = append(deps, Dependency{Name: name, Source: "composer.json"}) }
	for name := range composer.RequireDev { deps = append(deps, Dependency{Name: name, Source: "composer.json"}) }
	return deps, nil
}

// Ruby (Regex based for Gemfile)
type RubyScanner struct{}
func (s *RubyScanner) Scan(path string) ([]Dependency, error) {
	file, err := os.Open(filepath.Join(path, "Gemfile"))
	if err != nil { return nil, err }
	defer file.Close()
	var deps []Dependency
	re := regexp.MustCompile(`gem\s+['"]([^'"]+)['"]`)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		matches := re.FindStringSubmatch(scanner.Text())
		if len(matches) > 1 { deps = append(deps, Dependency{Name: matches[1], Source: "Gemfile"}) }
	}
	return deps, nil
}

// GitHub Actions (Scans .github/workflows/*.yml)
type ActionScanner struct{}
func (s *ActionScanner) Scan(path string) ([]Dependency, error) {
	workflowDir := filepath.Join(path, ".github", "workflows")
	files, err := os.ReadDir(workflowDir)
	if err != nil { return nil, err }
	var deps []Dependency
	re := regexp.MustCompile(`uses:\s+([^@\s]+)`)
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".yml") && !strings.HasSuffix(file.Name(), ".yaml") { continue }
		data, _ := os.ReadFile(filepath.Join(workflowDir, file.Name()))
		matches := re.FindAllStringSubmatch(string(data), -1)
		for _, m := range matches {
			if !strings.HasPrefix(m[1], "docker://") && !strings.HasPrefix(m[1], "./") {
				deps = append(deps, Dependency{Name: m[1], Source: file.Name()})
			}
		}
	}
	return deps, nil
}
