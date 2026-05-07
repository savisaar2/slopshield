package engine

import (
	"testing"
	"time"

	"github.com/savisaar2/slopshield/internal/config"
	"github.com/savisaar2/slopshield/internal/registry"
	"github.com/savisaar2/slopshield/internal/scanner"
	"github.com/savisaar2/slopshield/internal/slopignore"
)

type mockScanner struct {
	deps []scanner.Dependency
}

func (m *mockScanner) Scan(path string) ([]scanner.Dependency, error) {
	return m.deps, nil
}

type mockRegistry struct {
	meta map[string]*registry.Metadata
}

func (m *mockRegistry) GetMetadata(name string) (*registry.Metadata, error) {
	if meta, ok := m.meta[name]; ok {
		return meta, nil
	}
	return &registry.Metadata{Exists: false}, nil
}

func TestEngine_Evaluate(t *testing.T) {
	now := time.Now()
	e := &Engine{
		Config: &config.Config{
			ReputationAgeDays:   14,
			EnableTyposquatting: true,
		},
		KnownHallucinations: map[string]bool{"known-slop": true},
		IgnoreList:          &slopignore.IgnoreList{},
	}

	tests := []struct {
		name        string
		dep         scanner.Dependency
		meta        *registry.Metadata
		expected    bool
		isTyposquat bool
		reason      string
	}{
		{
			name:     "Known Slop",
			dep:      scanner.Dependency{Name: "known-slop"},
			expected: true,
			reason:   "Known hallucination in local registry",
		},
		{
			name:     "Non-existent Package",
			dep:      scanner.Dependency{Name: "non-existent"},
			meta:     &registry.Metadata{Exists: false},
			expected: true,
			reason:   "Package does not exist in official registry",
		},
		{
			name:     "Suspiciously New Package",
			dep:      scanner.Dependency{Name: "new-pkg"},
			meta:     &registry.Metadata{Exists: true, CreatedAt: now.Add(-24 * time.Hour)},
			expected: true,
			reason:   "Suspiciously new package",
		},
		{
			name:     "Reputable Package",
			dep:      scanner.Dependency{Name: "old-pkg"},
			meta:     &registry.Metadata{Exists: true, CreatedAt: now.Add(-365 * 24 * time.Hour)},
			expected: false,
		},
		{
			name:        "Typosquatting Package",
			dep:         scanner.Dependency{Name: "lodsh"}, // Close to lodash
			expected:    true,
			isTyposquat: true,
			reason:      "Potential typosquatting of popular package 'lodash'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := &mockRegistry{meta: map[string]*registry.Metadata{tt.dep.Name: tt.meta}}
			isSlop, isTypo, reason := e.evaluate(tt.dep, reg, registry.EcosystemNPM)
			if isSlop != tt.expected {
				t.Errorf("evaluate() isSlop = %v, expected %v", isSlop, tt.expected)
			}
			if isTypo != tt.isTyposquat {
				t.Errorf("evaluate() isTypo = %v, expected %v", isTypo, tt.isTyposquat)
			}
			if tt.expected && reason == "" {
				t.Error("evaluate() expected a reason for slop")
			}
		})
	}
}

func TestEngine_Cache(t *testing.T) {
	e := &Engine{
		Config: &config.Config{},
	}
	dep := scanner.Dependency{Name: "cached-pkg"}
	meta := &registry.Metadata{Exists: true}
	reg := &mockRegistry{meta: map[string]*registry.Metadata{dep.Name: meta}}

	// First call - should hit registry
	isSlop, _, _ := e.evaluate(dep, reg, registry.EcosystemNPM)
	if isSlop {
		t.Error("expected not slop")
	}

	// Verify it's in cache
	if _, ok := e.Cache.Load("npm:cached-pkg"); !ok {
		t.Error("expected metadata to be cached")
	}

	// Second call - should hit cache (registry doesn't matter)
	regEmpty := &mockRegistry{meta: make(map[string]*registry.Metadata)}
	isSlop, _, _ = e.evaluate(dep, regEmpty, registry.EcosystemNPM)
	if isSlop {
		t.Error("expected not slop from cache")
	}
}
