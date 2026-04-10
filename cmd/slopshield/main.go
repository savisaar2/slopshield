package main

import (
	"fmt"
	"os"

	"github.com/savisaar2/slopshield/internal/aggregator"
	"github.com/savisaar2/slopshield/internal/registry"
	"github.com/savisaar2/slopshield/internal/sarif"
	"github.com/savisaar2/slopshield/internal/scanner"
	"github.com/savisaar2/slopshield/internal/slopignore"
	"github.com/charmbracelet/huh"
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
			interactive, _ := cmd.Flags().GetBool("interactive")

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

				isHallucination := false
				// Check known list first
				if knownHallucinations[dep.Name] {
					if output == "text" {
						fmt.Printf("🚨 KNOWN Hallucination Found: %s (%s)\n", dep.Name, dep.Source)
					}
					isHallucination = true
				} else {
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
						isHallucination = true
					}
				}

				if isHallucination {
					if interactive {
						var action string
						form := huh.NewForm(
							huh.NewGroup(
								huh.NewNote().
									Title("Security Alert").
									Description(fmt.Sprintf("Package '%s' looks like a hallucination.", dep.Name)),
								huh.NewSelect[string]().
									Title("How do you want to handle this?").
									Options(
										huh.NewOption("Keep (ignore in future)", "ignore"),
										huh.NewOption("Report as confirmed hallucination", "report"),
										huh.NewOption("Do nothing", "none"),
									).
									Value(&action),
							),
						)

						if err := form.Run(); err != nil {
							return err
						}

						if action == "ignore" {
							if err := ignoreList.Add(path, dep.Name); err != nil {
								fmt.Printf("❌ Error adding to .slopignore: %v\n", err)
							} else {
								fmt.Printf("✅ Added %s to .slopignore\n", dep.Name)
							}
						}
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
	scanCmd.Flags().BoolP("interactive", "i", false, "Enable interactive mode for manual intervention")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
