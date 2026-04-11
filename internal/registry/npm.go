package registry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type NPMRegistry struct {
	client *http.Client
}

type NPMMetadata struct {
	Time map[string]string `json:"time"`
}

func NewNPMRegistry() *NPMRegistry {
	return &NPMRegistry{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (r *NPMRegistry) Exists(name string) (bool, error) {
	url := fmt.Sprintf("https://registry.npmjs.org/%s", name)
	resp, err := r.client.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	if resp.StatusCode == http.StatusOK {
		var meta NPMMetadata
		if err := json.NewDecoder(resp.Body).Decode(&meta); err == nil {
			if created, ok := meta.Time["created"]; ok {
				t, _ := time.Parse(time.RFC3339, created)
				// If the package is less than 14 days old, we still consider it "suspicious" 
				// even if it exists, as it might be an attacker-registered hallucination.
				if time.Since(t) < 14*24*time.Hour {
					return false, nil // Treat as "not reputable"
				}
			}
		}
		return true, nil
	}
	return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}
