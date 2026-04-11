package scanner

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type PubScanner struct{}

type Pubspec struct {
	Name         string                 `yaml:"name"`
	Dependencies map[string]interface{} `yaml:"dependencies"`
	DevDeps      map[string]interface{} `yaml:"dev_dependencies"`
}

func (s *PubScanner) Scan(path string) ([]Dependency, error) {
	pubspecPath := filepath.Join(path, "pubspec.yaml")
	data, err := os.ReadFile(pubspecPath)
	if err != nil {
		return nil, err
	}

	var pubspec Pubspec
	if err := yaml.Unmarshal(data, &pubspec); err != nil {
		return nil, err
	}

	var deps []Dependency
	for name, val := range pubspec.Dependencies {
		version := ""
		if v, ok := val.(string); ok {
			version = v
		}
		deps = append(deps, Dependency{Name: name, Version: version, Source: "pubspec.yaml"})
	}
	for name, val := range pubspec.DevDeps {
		version := ""
		if v, ok := val.(string); ok {
			version = v
		}
		deps = append(deps, Dependency{Name: name, Version: version, Source: "pubspec.yaml"})
	}

	return deps, nil
}
