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
	inRequireBlock := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		if line == "require (" {
			inRequireBlock = true
			continue
		}
		if line == ")" && inRequireBlock {
			inRequireBlock = false
			continue
		}

		if inRequireBlock {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				deps = append(deps, Dependency{Name: parts[0], Version: parts[1], Source: "go.mod"})
			}
			continue
		}

		if strings.HasPrefix(line, "require") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				deps = append(deps, Dependency{Name: parts[1], Version: parts[2], Source: "go.mod"})
			}
		}
	}

	return deps, nil
}
