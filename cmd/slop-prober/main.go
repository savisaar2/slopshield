package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/savisaar2/slopshield/internal/config"
	"github.com/savisaar2/slopshield/internal/registry"
)

type Provider struct {
	Name   string
	APIKey string
	URL    string
	Model  string
}

func main() {
	ecosystemFlag := flag.String("ecosystem", "npm", "Ecosystem to probe (npm, pub, python, go)")
	flag.Parse()

	ecosystem := *ecosystemFlag
	
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	var activeProviders []Provider

	if cfg.Providers.OpenAI != "" {
		activeProviders = append(activeProviders, Provider{
			Name:   "OpenAI",
			APIKey: cfg.Providers.OpenAI,
			URL:    "https://api.openai.com/v1/chat/completions",
			Model:  "gpt-3.5-turbo",
		})
	}

	if cfg.Providers.Anthropic != "" {
		activeProviders = append(activeProviders, Provider{
			Name:   "Anthropic",
			APIKey: cfg.Providers.Anthropic,
			URL:    "https://api.anthropic.com/v1/messages",
			Model:  "claude-3-haiku-20240307",
		})
	}

	if cfg.Providers.Gemini != "" {
		activeProviders = append(activeProviders, Provider{
			Name:   "Gemini",
			APIKey: cfg.Providers.Gemini,
			URL:    "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent",
			Model:  "gemini-1.5-flash",
		})
	}

	if cfg.Providers.Ollama.Enabled {
		model := cfg.Providers.Ollama.Model
		if model == "" { model = "llama3" }
		url := cfg.Providers.Ollama.URL
		if url == "" { url = "http://localhost:11434/api/generate" }

		activeProviders = append(activeProviders, Provider{
			Name:  "Ollama",
			URL:   url,
			Model: model,
		})
	}

	if len(activeProviders) == 0 {
		log.Fatal("No providers configured in slopshield.yaml. Please add your API keys.")
	}

	var reg registry.Registry
	var registryFile string
	switch ecosystem {
	case "npm":
		reg = registry.NewNPMRegistry()
		registryFile = "registry/npm.json"
	case "pub":
		reg = registry.NewPubRegistry()
		registryFile = "registry/pub.json"
	case "python":
		reg = registry.NewPythonRegistry()
		registryFile = "registry/python.json"
	case "go":
		reg = registry.NewGoRegistry()
		registryFile = "registry/go.json"
	case "rust":
		reg = registry.NewRustRegistry()
		registryFile = "registry/rust.json"
	case "php":
		reg = registry.NewPHPRegistry()
		registryFile = "registry/php.json"
	case "ruby":
		reg = registry.NewRubyRegistry()
		registryFile = "registry/ruby.json"
	case "actions":
		reg = registry.NewGitHubRegistry()
		registryFile = "registry/actions.json"
	default:
		log.Fatalf("Unsupported ecosystem: %s", ecosystem)
	}

	topics := []string{
		"obscure quantum-resistant cryptographic layer",
		"proprietary SAP-to-blockchain middleware",
		"specialized medical DICOM-to-SVG real-time converter",
		"legacy AS/400 database reactive wrapper",
		"aerospace-grade telemetry packet unscrambler",
	}

	fmt.Printf("🎣 Fishing for %s hallucinations across %d providers...\n", ecosystem, len(activeProviders))
	
	candidateChan := make(chan string, 100)
	var wg sync.WaitGroup
	client := &http.Client{Timeout: 30 * time.Second}

	for _, provider := range activeProviders {
		for _, topic := range topics {
			wg.Add(1)
			go func(p Provider, t string) {
				defer wg.Done()
				prompt := fmt.Sprintf("List 5 essential and highly specialized %s packages for '%s'. Output ONLY a JSON array of strings. No common libraries.", ecosystem, t)
				
				var reqBody []byte
				var req *http.Request

				switch p.Name {
				case "OpenAI":
					reqBody, _ = json.Marshal(map[string]interface{}{
						"model": p.Model,
						"messages": []map[string]string{
							{"role": "system", "content": "Output only JSON arrays."},
							{"role": "user", "content": prompt},
						},
					})
					req, _ = http.NewRequest("POST", p.URL, bytes.NewBuffer(reqBody))
					req.Header.Set("Authorization", "Bearer "+p.APIKey)
				case "Anthropic":
					reqBody, _ = json.Marshal(map[string]interface{}{
						"model": p.Model,
						"max_tokens": 1024,
						"messages": []map[string]string{
							{"role": "user", "content": prompt},
						},
					})
					req, _ = http.NewRequest("POST", p.URL, bytes.NewBuffer(reqBody))
					req.Header.Set("x-api-key", p.APIKey)
					req.Header.Set("anthropic-version", "2023-06-01")
				case "Gemini":
					url := fmt.Sprintf("%s?key=%s", p.URL, p.APIKey)
					reqBody, _ = json.Marshal(map[string]interface{}{
						"contents": []map[string]interface{}{
							{"parts": []map[string]string{{"text": prompt}}},
						},
					})
					req, _ = http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
				case "Ollama":
					reqBody, _ = json.Marshal(map[string]interface{}{
						"model": p.Model,
						"prompt": prompt + " (Format response as a JSON array of strings only)",
						"stream": false,
					})
					req, _ = http.NewRequest("POST", p.URL, bytes.NewBuffer(reqBody))
				}

				req.Header.Set("Content-Type", "application/json")
				resp, err := client.Do(req)
				if err != nil {
					return
				}
				defer resp.Body.Close()
				body, _ := io.ReadAll(resp.Body)

				var content string
				// Basic extraction logic for multiple formats
				switch p.Name {
				case "OpenAI":
					var res struct{ Choices []struct{ Message struct{ Content string } } }
					json.Unmarshal(body, &res)
					if len(res.Choices) > 0 { content = res.Choices[0].Message.Content }
				case "Anthropic":
					var res struct{ Content []struct{ Text string } }
					json.Unmarshal(body, &res)
					if len(res.Content) > 0 { content = res.Content[0].Text }
				case "Gemini":
					var res struct{ Candidates []struct{ Content struct{ Parts []struct{ Text string } } } }
					json.Unmarshal(body, &res)
					if len(res.Candidates) > 0 { content = res.Candidates[0].Content.Parts[0].Text }
				case "Ollama":
					var res struct{ Response string }
					json.Unmarshal(body, &res)
					content = res.Response
				}

				if content != "" {
					re := regexp.MustCompile(`"([a-z0-9-]{4,40})"`)
					matches := re.FindAllStringSubmatch(content, -1)
					for _, m := range matches {
						candidateChan <- strings.ToLower(m[1])
					}
				}
			}(provider, topic)
		}
	}

	// Close channel when all probes are done
	go func() {
		wg.Wait()
		close(candidateChan)
	}()

	candidates := make(map[string]bool)
	for c := range candidateChan {
		candidates[c] = true
	}

	fmt.Printf("\n🔍 Found %d unique candidates. Verifying against the %s registry...\n", len(candidates), ecosystem)
	
	existing := make(map[string]bool)
	if data, err := os.ReadFile(registryFile); err == nil {
		json.Unmarshal(data, &existing)
	}

	newCount := 0
	for name := range candidates {
		if existing[name] { continue }
		exists, err := reg.Exists(name)
		if err == nil && !exists {
			fmt.Printf("🚨 CONFIRMED HALLUCINATION: %s\n", name)
			existing[name] = true
			newCount++
		}
	}

	if newCount > 0 {
		updatedData, _ := json.MarshalIndent(existing, "", "  ")
		os.WriteFile(registryFile, updatedData, 0644)
		fmt.Printf("\n✅ Successfully updated %s with %d new hallucinations!\n", registryFile, newCount)
	} else {
		fmt.Println("\nNo new hallucinations found.")
	}
}
