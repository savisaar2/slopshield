package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/savisaar2/slopshield/internal/config"
	"github.com/savisaar2/slopshield/internal/engine"
	"github.com/savisaar2/slopshield/internal/sarif"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	styleTitle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2).
			MarginLeft(1).
			BorderForeground(lipgloss.Color("#00AA00"))

	rootCmd = &cobra.Command{
		Use:     "slopshield",
		Version: "1.2.0",
		Short:   "SlopShield identifies hallucinated packages",
		Long:    `SlopShield is a security scanner that detects AI-hallucinated packages.`,
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
				fmt.Println(styleTitle.Render("🛡️  SLOPSHIELD: AI Hallucination Guard"))
				fmt.Printf("🔍 Scanning project at: %s\n", path)
			}

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			eng, err := engine.NewEngine(path, cfg)
			if err != nil {
				return fmt.Errorf("failed to initialize engine: %w", err)
			}

			results, err := eng.Scan(context.Background(), path)
			if err != nil {
				return fmt.Errorf("scan failed: %w", err)
			}

			if output == "sarif" {
				var names []string
				for _, r := range results {
					names = append(names, r.Dependency.Name)
				}
				if err := sarif.Generate(os.Stdout, names, "manifest"); err != nil {
					return fmt.Errorf("error generating SARIF: %w", err)
				}
				return nil
			}

			if len(results) == 0 {
				fmt.Println("✅ No hallucinated packages found!")
				return nil
			}

			fmt.Printf("\n❌ Found %d potential hallucination(s):\n", len(results))
			for _, r := range results {
				fmt.Printf("  - %s (%s): %s\n", r.Dependency.Name, r.Dependency.Source, r.Reason)
			}

			os.Exit(1)
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
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
