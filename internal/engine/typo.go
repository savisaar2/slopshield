package engine

import (
	"github.com/savisaar2/slopshield/internal/registry"
)

var popularPackages = map[registry.Ecosystem][]string{
	registry.EcosystemNPM: {
		"lodash", "react", "vue", "express", "moment", "axios", "chalk", "commander", "debug", "request",
	},
	registry.EcosystemPython: {
		"requests", "numpy", "pandas", "django", "flask", "boto3", "cryptography", "scikit-learn", "urllib3", "pip",
	},
}

// levenshteinDistance calculates the Levenshtein distance between two strings.
func levenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	column := make([]int, len(s1)+1)
	for i := 1; i <= len(s1); i++ {
		column[i] = i
	}

	for j := 1; j <= len(s2); j++ {
		column[0] = j
		lastDiagonal := j - 1
		for i := 1; i <= len(s1); i++ {
			oldColumn := column[i]
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}
			column[i] = min(column[i]+1, column[i-1]+1, lastDiagonal+cost)
			lastDiagonal = oldColumn
		}
	}

	return column[len(s1)]
}

func min(a, b, c int) int {
	if a < b && a < c {
		return a
	}
	if b < c {
		return b
	}
	return c
}

func (e *Engine) checkTyposquat(name string, eco registry.Ecosystem) (bool, string) {
	if !e.Config.EnableTyposquatting {
		return false, ""
	}

	targets := popularPackages[eco]
	if custom := e.Config.TyposquattingTargets[string(eco)]; len(custom) > 0 {
		targets = append(targets, custom...)
	}

	for _, p := range targets {
		if name == p {
			return false, ""
		}
		dist := levenshteinDistance(name, p)
		if dist > 0 && dist <= 2 {
			return true, p
		}
	}
	return false, ""
}
