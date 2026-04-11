package main

import (
	"fmt"
	"os"
	"path/filepath"

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

			// Load known hallucinations from cache
			agg := aggregator.NewAggregator()
			knownHallucinations, err := agg.LoadCache(".slop_cache")
			if err != nil && output == "text" {
				fmt.Println("⚠️  Warning: No local cache found. Run 'slopshield sync' for better detection.")
				knownHallucinations = make(map[string]bool)
			}

			// Detect ecosystems
			var scannerList []struct {
				scanner  scanner.Scanner
				registry registry.Registry
				filename string
			}

			if _, err := os.Stat(filepath.Join(path, "package.json")); err == nil {
				scannerList = append(scannerList, struct {
					scanner  scanner.Scanner
					registry registry.Registry
					filename string
				}{&scanner.NPMScanner{}, registry.NewNPMRegistry(), "package.json"})
			}
			if _, err := os.Stat(filepath.Join(path, "pubspec.yaml")); err == nil {
				scannerList = append(scannerList, struct {
					scanner  scanner.Scanner
					registry registry.Registry
					filename string
				}{&scanner.PubScanner{}, registry.NewPubRegistry(), "pubspec.yaml"})
			}

			if len(scannerList) == 0 {
				return fmt.Errorf("no supported manifest file (package.json, pubspec.yaml) found in %s", path)
			}

			var allHallucinated []string
			for _, item := range scannerList {
				deps, err := item.scanner.Scan(path)
				if err != nil {
					if output == "text" {
						fmt.Printf("⚠️ Error scanning %s: %v\n", item.filename, err)
					}
					continue
				}

				if output == "text" {
					fmt.Printf("📦 Found %d dependencies in %s\n", len(deps), item.filename)
				}

				for _, dep := range deps {
					// Ignore built-in/local packages
					if dep.Source == "pubspec.yaml" && (dep.Name == "flutter" || dep.Name == "flutter_test" || dep.Name == "flutter_localizations") {
						continue
					}

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
						exists, err := item.registry.Exists(dep.Name)
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
						allHallucinated = append(allHallucinated, dep.Name)
					}
				}
			}

			if output == "sarif" {
				if err := sarif.Generate(os.Stdout, allHallucinated, "manifest"); err != nil {
					return fmt.Errorf("error generating SARIF: %w", err)
				}
				return nil
			}

			if len(allHallucinated) == 0 {
				fmt.Println("✅ No hallucinated packages found!")
			} else {
				fmt.Printf("\n❌ Total hallucinated packages: %d\n", len(allHallucinated))
				os.Exit(1)
			}

			return nil
		},
	}

	syncCmd = &cobra.Command{
		Use:   "sync",
		Short: "Fetch and update the local hallucination registry",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("🔄 Syncing hallucination registries...")
			agg := aggregator.NewAggregator()
			count, err := agg.Sync(".slop_cache")
			if err != nil {
				return fmt.Errorf("sync failed: %w", err)
			}
			fmt.Printf("✅ Successfully synced %d known hallucinated packages to .slop_cache\n", count)
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(syncCmd)
	scanCmd.Flags().StringP("output", "o", "text", "Output format (text, sarif)")
	scanCmd.Flags().BoolP("interactive", "i", false, "Enable interactive mode for manual intervention")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
