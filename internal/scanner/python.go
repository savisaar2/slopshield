package scanner

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type PythonScanner struct{}

func (s *PythonScanner) Scan(path string) ([]Dependency, error) {
	reqPath := filepath.Join(path, "requirements.txt")
	file, err := os.Open(reqPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var deps []Dependency
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Basic parsing for name==version
		parts := strings.Split(line, "==")
		name := parts[0]
		version := ""
		if len(parts) > 1 {
			version = parts[1]
		}
		deps = append(deps, Dependency{Name: strings.ToLower(name), Version: version, Source: "requirements.txt"})
	}

	return deps, nil
}
