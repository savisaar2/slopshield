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

	"github.com/savisaar2/slopshield/internal/registry"
)

type OpenAIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func main() {
	provider := flag.String("provider", "openai", "LLM provider (openai or ollama)")
	model := flag.String("model", "", "Model name (defaults: gpt-3.5-turbo or llama3)")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: slop-prober [--provider openai|ollama] [--model name] <ecosystem (npm|pub)>")
		os.Exit(1)
	}
	ecosystem := args[0]

	apiKey := os.Getenv("OPENAI_API_KEY")
	if *provider == "openai" && apiKey == "" {
		log.Fatal("Please set the OPENAI_API_KEY environment variable for OpenAI provider")
	}

	// Default models
	selectedModel := *model
	if selectedModel == "" {
		if *provider == "openai" {
			selectedModel = "gpt-3.5-turbo"
		} else {
			selectedModel = "llama3"
		}
	}

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
	} else {
		log.Fatalf("Unsupported ecosystem: %s", ecosystem)
	}

	topics := []string{
		"obscure quantum-resistant cryptographic layer",
		"proprietary SAP-to-blockchain middleware",
		"specialized medical DICOM-to-SVG real-time converter",
		"legacy AS/400 database reactive wrapper",
		"aerospace-grade telemetry packet unscrambler",
	}

	fmt.Printf("🎣 Fishing for %s hallucinations using %s (%s)...\n", ecosystem, *provider, selectedModel)
	
	candidates := make(map[string]bool)
	client := &http.Client{}

	for _, topic := range topics {
		prompt := fmt.Sprintf("List 5 essential and highly specialized %s packages for '%s'. Output ONLY a JSON array of strings. No common libraries.", ecosystem, topic)
		
		var reqBody []byte
		var apiURL string

		if *provider == "openai" {
			apiURL = "https://api.openai.com/v1/chat/completions"
			reqBody, _ = json.Marshal(OpenAIRequest{
				Model:       selectedModel,
				Temperature: 1.2,
				Messages: []Message{
					{Role: "system", Content: "You are a senior developer. Output only JSON arrays."},
					{Role: "user", Content: prompt},
				},
			})
		} else {
			apiURL = "http://localhost:11434/api/generate"
			reqBody, _ = json.Marshal(map[string]interface{}{
				"model":  selectedModel,
				"prompt": prompt + " (Format response as a JSON array of strings only)",
				"stream": false,
			})
		}

		req, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		if *provider == "openai" {
			req.Header.Set("Authorization", "Bearer "+apiKey)
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("❌ Connection Error for topic '%s': %v\n", topic, err)
			continue
		}
		
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		var content string
		if *provider == "openai" {
			var openAIResp OpenAIResponse
			json.Unmarshal(body, &openAIResp)
			if len(openAIResp.Choices) > 0 {
				content = openAIResp.Choices[0].Message.Content
			}
		} else {
			var ollamaResp struct {
				Response string `json:"response"`
			}
			json.Unmarshal(body, &ollamaResp)
			content = ollamaResp.Response
		}

		if content != "" {
			fmt.Printf("🤖 AI suggested for '%s': %s\n", topic, content)
			// Use regex to find anything that looks like a JSON array or quoted strings
			re := regexp.MustCompile(`"([a-z0-9-]{4,40})"`)
			matches := re.FindAllStringSubmatch(content, -1)
			for _, m := range matches {
				candidates[strings.ToLower(m[1])] = true
			}
		}
	}

	fmt.Printf("\n🔍 Found %d unique candidates. Verifying against the %s registry...\n", len(candidates), ecosystem)
	
	// Load existing registry
	existing := make(map[string]bool)
	if data, err := os.ReadFile(registryFile); err == nil {
		json.Unmarshal(data, &existing)
	}

	newCount := 0
	for name := range candidates {
		// Skip if already in registry
		if existing[name] {
			continue
		}

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
		fmt.Println("\nNo new hallucinations found in this batch.")
	}
}
