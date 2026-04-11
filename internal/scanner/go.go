package scanner

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type GoScanner struct{}

func (s *GoScanner) Scan(path string) ([]Dependency, error) {
	modPath := filepath.Join(path, "go.mod")
	file, err := os.Open(modPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var deps []Dependency
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Basic parsing for 'require package/name version'
		if strings.HasPrefix(line, "require") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				name := parts[1]
				version := ""
				if len(parts) >= 3 {
					version = parts[2]
				}
				deps = append(deps, Dependency{Name: name, Version: version, Source: "go.mod"})
			}
		}
	}

	return deps, nil
}
