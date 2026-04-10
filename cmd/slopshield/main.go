package main

import (
	"fmt"
	"os"

	"github.com/savisaar2/slopshield/internal/aggregator"
	"github.com/savisaar2/slopshield/internal/registry"
	"github.com/savisaar2/slopshield/internal/sarif"
	"github.com/savisaar2/slopshield/internal/scanner"
	"github.com/savisaar2/slopshield/internal/slopignore"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "slopshield",
		Short: "SlopShield identifies hallucinated packages in your project",
		Long:  `SlopShield is a security scanner that detects AI-hallucinated packages by verifying them against official registries and known hallucination lists.`,
	}

	scanCmd = &cobra.Command{
		Use:   "scan [path]",
		Short: "Scan a project for hallucinated packages",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "."
			if len(args) > 0 {
				path = args[0]
			}

			output, _ := cmd.Flags().GetString("output")
			if output == "text" {
				fmt.Printf("🔍 Scanning project at: %s\n", path)
			}

			ignoreList, err := slopignore.Load(path)
			if err != nil {
				return fmt.Errorf("error loading ignore list: %w", err)
			}

			// Fetch known hallucinations
			agg := aggregator.NewAggregator()
			knownHallucinations, _ := agg.FetchAll()

			npmScanner := &scanner.NPMScanner{}
			deps, err := npmScanner.Scan(path)
			if err != nil {
				return fmt.Errorf("error scanning dependencies: %w", err)
			}

			if output == "text" {
				fmt.Printf("📦 Found %d dependencies\n", len(deps))
			}

			npmRegistry := registry.NewNPMRegistry()
			var hallucinatedDeps []string
			for _, dep := range deps {
				if ignoreList.IsIgnored(dep.Name) {
					if output == "text" {
						fmt.Printf("🙈 Ignoring %s (matched in .slopignore)\n", dep.Name)
					}
					continue
				}

				// Check known list first
				if knownHallucinations[dep.Name] {
					if output == "text" {
						fmt.Printf("🚨 KNOWN Hallucination Found: %s (%s)\n", dep.Name, dep.Source)
					}
					hallucinatedDeps = append(hallucinatedDeps, dep.Name)
					continue
				}

				exists, err := npmRegistry.Exists(dep.Name)
				if err != nil {
					if output == "text" {
						fmt.Printf("⚠️ Error checking %s: %v\n", dep.Name, err)
					}
					continue
				}

				if !exists {
					if output == "text" {
						fmt.Printf("🚨 Potential Hallucination Found: %s (%s)\n", dep.Name, dep.Source)
					}
					hallucinatedDeps = append(hallucinatedDeps, dep.Name)
				}
			}

			if output == "sarif" {
				if err := sarif.Generate(os.Stdout, hallucinatedDeps, "package.json"); err != nil {
					return fmt.Errorf("error generating SARIF: %w", err)
				}
				return nil
			}

			if len(hallucinatedDeps) == 0 {
				fmt.Println("✅ No hallucinated packages found!")
			} else {
				fmt.Printf("\n❌ Total hallucinated packages: %d\n", len(hallucinatedDeps))
				os.Exit(1)
			}

			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringP("output", "o", "text", "Output format (text, sarif)")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
