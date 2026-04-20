package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/savisaar2/slopshield/internal/registry"
)

func main() {
	update := flag.Bool("update", false, "Update the local registry files directly")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		fmt.Println("Usage: slop-hunter [--update] <ecosystem> <name1,name2,...>")
		fmt.Println("Example: slop-hunter --update npm express-gpt,react-ai-core")
		os.Exit(1)
	}

	ecosystem := args[0]
	names := strings.Split(args[1], ",")

	var reg registry.Registry
	var registryFile string

	if ecosystem == "npm" {
		reg = registry.NewNPMRegistry()
		registryFile = "registry/npm.json"
	} else if ecosystem == "pub" {
		reg = registry.NewPubRegistry()
		registryFile = "registry/pub.json"
	} else if ecosystem == "python" {
		reg = registry.NewPythonRegistry()
		registryFile = "registry/python.json"
	} else if ecosystem == "go" {
		reg = registry.NewGoRegistry()
		registryFile = "registry/go.json"
	} else if ecosystem == "rust" {
		reg = registry.NewRustRegistry()
		registryFile = "registry/rust.json"
	} else if ecosystem == "php" {
		reg = registry.NewPHPRegistry()
		registryFile = "registry/php.json"
	} else if ecosystem == "ruby" {
	        reg = registry.NewRubyRegistry()
	        registryFile = "registry/ruby.json"
	} else if ecosystem == "actions" {
	        reg = registry.NewGitHubRegistry()
	        registryFile = "registry/actions.json"
	} else if ecosystem == "maven" {
	        reg = registry.NewMavenRegistry()
	        registryFile = "registry/maven.json"
	} else if ecosystem == "nuget" {
	        reg = registry.NewNuGetRegistry()
	        registryFile = "registry/nuget.json"
	} else {
	        log.Fatalf("Unsupported ecosystem: %s", ecosystem)
	}
	hallucinations := make(map[string]bool)
	fmt.Printf("🎯 Hunting for hallucinations in %s...\n", ecosystem)

	for _, name := range names {
		name = strings.TrimSpace(name)
		exists, err := reg.Exists(name)
		if err != nil {
			fmt.Printf("⚠️  Error checking %s: %v\n", name, err)
			continue
		}

		if !exists {
			fmt.Printf("🚨 CONFIRMED HALLUCINATION: %s\n", name)
			hallucinations[name] = true
		} else {
			fmt.Printf("✅ Real package: %s\n", name)
		}
	}

	if len(hallucinations) == 0 {
		fmt.Println("\nNo new hallucinations found.")
		return
	}

	if *update {
		// Load existing
		existing := make(map[string]bool)
		data, err := os.ReadFile(registryFile)
		if err == nil {
			json.Unmarshal(data, &existing)
		}

		// Merge
		newCount := 0
		for name := range hallucinations {
			if !existing[name] {
				existing[name] = true
				newCount++
			}
		}

		// Save
		updatedData, _ := json.MarshalIndent(existing, "", "  ")
		if err := os.WriteFile(registryFile, updatedData, 0644); err != nil {
			log.Fatalf("Failed to write to %s: %v", registryFile, err)
		}
		fmt.Printf("\n✅ Successfully added %d new hallucinations to %s!\n", newCount, registryFile)
	} else {
		data, _ := json.MarshalIndent(hallucinations, "", "  ")
		fmt.Println("\n--- NEW FINDINGS ---")
		fmt.Println(string(data))
		fmt.Println("\nRun with --update to save these to the registry.")
	}
}
