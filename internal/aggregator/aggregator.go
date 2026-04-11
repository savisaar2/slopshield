package aggregator

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type Aggregator struct {
	sources []string
	client  *http.Client
}

func NewAggregator() *Aggregator {
	baseURL := os.Getenv("SLOPSHIELD_REGISTRY_URL")
	if baseURL == "" {
		// Default to a placeholder that users will replace in their fork
		baseURL = "https://raw.githubusercontent.com/YOUR_USERNAME/slopshield/main/registry"
	} else {
		baseURL = strings.TrimSuffix(baseURL, "/")
	}

	return &Aggregator{
		sources: []string{
			fmt.Sprintf("%s/npm.json", baseURL),
			fmt.Sprintf("%s/pub.json", baseURL),
			fmt.Sprintf("%s/python.json", baseURL),
			fmt.Sprintf("%s/go.json", baseURL),
		},
		client: &http.Client{
			Timeout: 20 * time.Second,
		},
	}
}

func (a *Aggregator) Sync(cachePath string) (int, error) {
	hallucinated := make(map[string]bool)
	count := 0

	for _, source := range a.sources {
		resp, err := a.client.Get(source)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		// If the file is JSON (our new format)
		var remoteData map[string]bool
		if err := json.NewDecoder(resp.Body).Decode(&remoteData); err == nil {
			for name := range remoteData {
				if !hallucinated[name] {
					hallucinated[name] = true
					count++
				}
			}
			continue
		}
	}
    // ... rest of the existing WriteFile logic

	data, err := json.MarshalIndent(hallucinated, "", "  ")
	if err != nil {
		return 0, err
	}

	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		return 0, err
	}

	return count, nil
}

func (a *Aggregator) LoadCache(cachePath string) (map[string]bool, error) {
	data, err := os.ReadFile(cachePath)
	if err != nil {
		return nil, err
	}

	var cache map[string]bool
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, err
	}

	return cache, nil
}
