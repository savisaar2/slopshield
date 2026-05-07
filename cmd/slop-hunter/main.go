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

	reg, err := registry.GetRegistry(registry.Ecosystem(ecosystem))
	if err != nil {
		log.Fatalf("Unsupported ecosystem: %s", ecosystem)
	}
	registryFile := fmt.Sprintf("registry/%s.json", ecosystem)
	if ecosystem == "actions" { registryFile = "registry/actions.json" }

	hallucinations := make(map[string]bool)
	fmt.Printf("🎯 Hunting for hallucinations in %s...\n", ecosystem)

	for _, name := range names {
		name = strings.TrimSpace(name)
		meta, err := reg.GetMetadata(name)
		if err != nil {
			fmt.Printf("⚠️  Error checking %s: %v\n", name, err)
			continue
		}

		if !meta.Exists {
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
