package registry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/savisaar2/slopshield/internal/resilienthttp"
)

type NPMRegistry struct {
	client  *http.Client
	baseURL string
}

type NPMMetadata struct {
	Time map[string]string `json:"time"`
}

func NewNPMRegistry(baseURL string) *NPMRegistry {
	if baseURL == "" {
		baseURL = "https://registry.npmjs.org"
	}
	return &NPMRegistry{
		client:  resilienthttp.NewClient(),
		baseURL: baseURL,
	}
}

func (r *NPMRegistry) GetMetadata(name string) (*Metadata, error) {
	url := fmt.Sprintf("%s/%s", r.baseURL, name)
	resp, err := r.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return &Metadata{Exists: false}, nil
	}

	if resp.StatusCode == http.StatusOK {
		var meta NPMMetadata
		var createdAt time.Time
		if err := json.NewDecoder(resp.Body).Decode(&meta); err == nil {
			if created, ok := meta.Time["created"]; ok {
				createdAt, _ = time.Parse(time.RFC3339, created)
			}
		}
		return &Metadata{Exists: true, CreatedAt: createdAt}, nil
	}
	return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}
