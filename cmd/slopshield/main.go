package main

import (
        "context"
        "encoding/json"
        "fmt"
        "log/slog"
        "os"
        "path/filepath"
        "sync"

        "github.com/savisaar2/slopshield/internal/registry"
        "github.com/savisaar2/slopshield/internal/sarif"
        "github.com/savisaar2/slopshield/internal/scanner"
        "github.com/savisaar2/slopshield/internal/slopignore"
        "github.com/charmbracelet/lipgloss"
        "github.com/spf13/cobra"
        "golang.org/x/sync/errgroup"
)
var (
	styleTitle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2).
			MarginLeft(1).
			BorderForeground(lipgloss.Color("#00AA00"))

	styleHeader = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4ECDC4")).
			Bold(true).
			Margin(1, 0)

	rootCmd = &cobra.Command{
		Use:     "slopshield",
		Version: "1.1.0",
		Short:   "SlopShield identifies hallucinated packages",
		Long:    `SlopShield is a security scanner that detects AI-hallucinated packages.`,
	}

	scanCmd = &cobra.Command{
		Use:   "scan [path]",
		Short: "Scan a project for hallucinated packages",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
            // ... (keep existing scan logic)
			path := "."
			if len(args) > 0 {
				path = args[0]
			}

			output, _ := cmd.Flags().GetString("output")
			interactive, _ := cmd.Flags().GetBool("interactive")

			if output == "text" {
				fmt.Println(styleTitle.Render("🛡️  SLOPSHIELD: AI Hallucination Guard"))
				fmt.Printf("🔍 Auto-detecting project type at: %s\n", path)
			}

			ignoreList, err := slopignore.Load(path)
			if err != nil {
				return fmt.Errorf("error loading ignore list: %w", err)
			}

			// Load known hallucinations from local registry files
			knownHallucinations := make(map[string]bool)
			registryFiles, _ := filepath.Glob("registry/*.json")
			for _, rf := range registryFiles {
				if data, err := os.ReadFile(rf); err == nil {
					var fileCache map[string]bool
					if err := json.Unmarshal(data, &fileCache); err == nil {
						for k, v := range fileCache {
							knownHallucinations[k] = v
						}
					}
				}
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
			if _, err := os.Stat(filepath.Join(path, "requirements.txt")); err == nil {
				scannerList = append(scannerList, struct {
					scanner  scanner.Scanner
					registry registry.Registry
					filename string
				}{&scanner.PythonScanner{}, registry.NewPythonRegistry(), "requirements.txt"})
			}
			if _, err := os.Stat(filepath.Join(path, "go.mod")); err == nil {
				scannerList = append(scannerList, struct {
					scanner  scanner.Scanner
					registry registry.Registry
					filename string
				}{&scanner.GoScanner{}, registry.NewGoRegistry(), "go.mod"})
			}
			if _, err := os.Stat(filepath.Join(path, "Cargo.toml")); err == nil {
				scannerList = append(scannerList, struct {
					scanner  scanner.Scanner
					registry registry.Registry
					filename string
				}{&scanner.RustScanner{}, registry.NewRustRegistry(), "Cargo.toml"})
			}
			if _, err := os.Stat(filepath.Join(path, "composer.json")); err == nil {
				scannerList = append(scannerList, struct {
					scanner  scanner.Scanner
					registry registry.Registry
					filename string
				}{&scanner.PHPScanner{}, registry.NewPHPRegistry(), "composer.json"})
			}
			if _, err := os.Stat(filepath.Join(path, "Gemfile")); err == nil {
				scannerList = append(scannerList, struct {
					scanner  scanner.Scanner
					registry registry.Registry
					filename string
				}{&scanner.RubyScanner{}, registry.NewRubyRegistry(), "Gemfile"})
			}
			if _, err := os.Stat(filepath.Join(path, ".github", "workflows")); err == nil {
				scannerList = append(scannerList, struct {
					scanner  scanner.Scanner
					registry registry.Registry
					filename string
				}{&scanner.ActionScanner{}, registry.NewGitHubRegistry(), "GitHub Actions"})
			}

			if len(scannerList) == 0 {
				return fmt.Errorf("no supported manifest file found in %s", path)
			}

			var allHallucinated []string
			var mu sync.Mutex

			for _, item := range scannerList {
				deps, err := item.scanner.Scan(path)
				if err != nil {
					slog.Error("Failed to scan manifest", "filename", item.filename, "error", err)
					continue
				}

				if output == "text" {
					fmt.Printf("📦 Found %d dependencies in %s\n", len(deps), item.filename)
				}

				g, ctx := errgroup.WithContext(context.Background())
				// Limit concurrency to 10 workers to avoid rate limiting
				g.SetLimit(10)

				for _, dep := range deps {
					dep := dep // shadow for goroutine
					g.Go(func() error {
						// Ignore built-in/local packages
						if dep.Source == "pubspec.yaml" && (dep.Name == "flutter" || dep.Name == "flutter_test" || dep.Name == "flutter_localizations") {
							return nil
						}

						if ignoreList.IsIgnored(dep.Name) {
							if output == "text" {
								slog.Debug("Ignoring package (matched in .slopignore)", "package", dep.Name)
							}
							return nil
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
								slog.Warn("Error checking registry", "package", dep.Name, "error", err)
								return nil
							}

							if !exists {
								if output == "text" {
									fmt.Printf("🚨 Potential Hallucination Found: %s (%s)\n", dep.Name, dep.Source)
								}
								isHallucination = true
							}
						}

						if isHallucination {
							mu.Lock()
							allHallucinated = append(allHallucinated, dep.Name)
							mu.Unlock()

							if interactive {
								// Note: Interactive mode is tricky with concurrency. 
								// For production, we usually collect all and then ask, 
								// but for now we'll just skip interactive if concurrent 
								// or handle it carefully. 
							}
						}
						return nil
					})
				}

				if err := g.Wait(); err != nil {
					slog.Error("Error during concurrent scan", "error", err)
				}
				_ = ctx // avoid unused
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

	registryCmd = &cobra.Command{
		Use:   "registry",
		Short: "Manage local registries",
	}

	clearRegistryCmd = &cobra.Command{
		Use:   "clear [ecosystem]",
		Short: "Clear local registry entries",
		Long:  "Empty one or all local registry JSON files. If no ecosystem is specified, all are cleared.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pattern := "registry/*.json"
			if len(args) > 0 {
				pattern = fmt.Sprintf("registry/%s.json", args[0])
			}

			files, err := filepath.Glob(pattern)
			if err != nil {
				return err
			}

			if len(files) == 0 {
				return fmt.Errorf("no registry files found matching: %s", pattern)
			}

			for _, f := range files {
				if err := os.WriteFile(f, []byte("{}"), 0644); err != nil {
					return fmt.Errorf("failed to clear %s: %w", f, err)
				}
				fmt.Printf("✅ Cleared %s\n", f)
			}
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(registryCmd)
	registryCmd.AddCommand(clearRegistryCmd)

	scanCmd.Flags().StringP("output", "o", "text", "Output format (text, sarif)")
	scanCmd.Flags().BoolP("interactive", "i", false, "Enable interactive mode for manual intervention")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
