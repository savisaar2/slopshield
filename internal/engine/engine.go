package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/savisaar2/slopshield/internal/config"
	"github.com/savisaar2/slopshield/internal/registry"
	"github.com/savisaar2/slopshield/internal/scanner"
	"github.com/savisaar2/slopshield/internal/slopignore"
	"golang.org/x/sync/errgroup"
)

type Result struct {
	Dependency scanner.Dependency
	IsSlop     bool
	IsTyposquat bool
	Reason     string
}

type Engine struct {
	Config              *config.Config
	KnownHallucinations map[string]bool
	IgnoreList          *slopignore.IgnoreList
	Concurrency         int
	Cache               sync.Map // Cache for registry.Metadata
}

func NewEngine(path string, cfg *config.Config) (*Engine, error) {
	ignoreList, err := slopignore.Load(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load ignore list: %w", err)
	}

	if cfg == nil {
		var err error
		cfg, err = config.Load()
		if err != nil {
			return nil, err
		}
	}

	known := make(map[string]bool)
	files, _ := filepath.Glob("registry/*.json")
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err == nil {
			var m map[string]bool
			if err := json.Unmarshal(data, &m); err == nil {
				for k, v := range m {
					known[k] = v
				}
			}
		}
	}

	return &Engine{
		Config:              cfg,
		KnownHallucinations: known,
		IgnoreList:          ignoreList,
		Concurrency:         10,
	}, nil
}

func (e *Engine) Scan(ctx context.Context, path string) ([]Result, error) {
	scannerInfos, err := scanner.DetectScanners(path)
	if err != nil {
		return nil, err
	}

	var results []Result
	var mu sync.Mutex

	for _, info := range scannerInfos {
		deps, err := info.Scanner.Scan(path)
		if err != nil {
			slog.Error("Failed to scan manifest", "filename", info.Filename, "error", err)
			continue
		}

		baseURL := e.Config.PrivateRegistries[string(info.Ecosystem)]
		reg, err := registry.GetRegistry(info.Ecosystem, baseURL)
		if err != nil {
			slog.Warn("No registry found for ecosystem", "ecosystem", info.Ecosystem)
			continue
		}

		g, _ := errgroup.WithContext(ctx)
		g.SetLimit(e.Concurrency)

		for _, dep := range deps {
			dep := dep
			g.Go(func() error {
				if e.isIgnored(dep) {
					return nil
				}

				isSlop, isTyposquat, reason := e.evaluate(dep, reg, info.Ecosystem)
				if isSlop {
					mu.Lock()
					results = append(results, Result{Dependency: dep, IsSlop: true, IsTyposquat: isTyposquat, Reason: reason})
					mu.Unlock()
				}
				return nil
			})
		}

		if err := g.Wait(); err != nil {
			return nil, err
		}
	}

	return results, nil
}

func (e *Engine) isIgnored(dep scanner.Dependency) bool {
	// Built-in/local packages (specific to Flutter for now as in original code)
	if dep.Source == "pubspec.yaml" && (dep.Name == "flutter" || dep.Name == "flutter_test" || dep.Name == "flutter_localizations") {
		return true
	}
	return e.IgnoreList.IsIgnored(dep.Name)
}

func (e *Engine) evaluate(dep scanner.Dependency, reg registry.Registry, eco registry.Ecosystem) (bool, bool, string) {
	// 1. Check known slops
	if e.KnownHallucinations[dep.Name] {
		return true, false, "Known hallucination in local registry"
	}

	// 2. Typosquatting Check
	if isTypo, target := e.checkTyposquat(dep.Name, eco); isTypo {
		return true, true, fmt.Sprintf("Potential typosquatting of popular package '%s'", target)
	}

	// 3. Check Cache
	cacheKey := string(eco) + ":" + dep.Name
	if val, ok := e.Cache.Load(cacheKey); ok {
		meta := val.(*registry.Metadata)
		return e.checkMetadata(dep, meta)
	}

	// 3. Check upstream registry
	meta, err := reg.GetMetadata(dep.Name)
	if err != nil {
		slog.Warn("Error checking registry", "package", dep.Name, "error", err)
		return false, false, ""
	}

	// Store in cache
	e.Cache.Store(cacheKey, meta)

	return e.checkMetadata(dep, meta)
}

func (e *Engine) checkMetadata(dep scanner.Dependency, meta *registry.Metadata) (bool, bool, string) {
	if !meta.Exists {
		return true, false, "Package does not exist in official registry"
	}

	// 4. Reputation Check
	if e.Config.ReputationAgeDays > 0 && !meta.CreatedAt.IsZero() {
		if time.Since(meta.CreatedAt) < time.Duration(e.Config.ReputationAgeDays)*24*time.Hour {
			return true, false, fmt.Sprintf("Suspiciously new package (created %s)", meta.CreatedAt.Format("2006-01-02"))
		}
	}

	return false, false, ""
}
