package scanner

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type PubScanner struct{}

type Pubspec struct {
	Name         string            `yaml:"name"`
	Dependencies map[string]string `yaml:"dependencies"`
	DevDeps      map[string]string `yaml:"dev_dependencies"`
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
	for name, version := range pubspec.Dependencies {
		deps = append(deps, Dependency{Name: name, Version: version, Source: "pubspec.yaml"})
	}
	for name, version := range pubspec.DevDeps {
		deps = append(deps, Dependency{Name: name, Version: version, Source: "pubspec.yaml"})
	}

	return deps, nil
}
