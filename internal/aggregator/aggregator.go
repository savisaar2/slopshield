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
	return &Aggregator{
		sources: []string{
			"https://raw.githubusercontent.com/ai-security/hallucinated-packages/main/npm.txt",
			"https://raw.githubusercontent.com/LassoSecurity/hallucinated-packages/main/hallucinated_packages.txt",
		},
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (a *Aggregator) Sync(cachePath string) (int, error) {
	hallucinated := make(map[string]bool)
	count := 0

	for _, source := range a.sources {
		resp, err := a.client.Get(source)
		if err != nil {
			fmt.Printf("⚠️  Failed to fetch from %s: %v\n", source, err)
			continue
		}
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			name := strings.TrimSpace(scanner.Text())
			if name != "" && !strings.HasPrefix(name, "#") {
				if !hallucinated[name] {
					hallucinated[name] = true
					count++
				}
			}
		}
	}

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
