package slopignore

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type IgnoreList struct {
	patterns []string
}

func Load(path string) (*IgnoreList, error) {
	file, err := os.Open(filepath.Join(path, ".slopignore"))
	if err != nil {
		if os.IsNotExist(err) {
			return &IgnoreList{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var patterns []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		patterns = append(patterns, line)
	}
	return &IgnoreList{patterns: patterns}, nil
}

func (l *IgnoreList) IsIgnored(name string) bool {
	for _, pattern := range l.patterns {
		matched, _ := filepath.Match(pattern, name)
		if matched {
			return true
		}
	}
	return false
}
