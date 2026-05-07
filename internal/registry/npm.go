package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/savisaar2/slopshield/internal/resilienthttp"
	"golang.org/x/time/rate"
)

type NPMRegistry struct {
	client  *http.Client
	baseURL string
	limiter *rate.Limiter
}

type NPMMetadata struct {
	Time map[string]string `json:"time"`
}

func NewNPMRegistry(baseURL string, limiter *rate.Limiter) *NPMRegistry {
	if baseURL == "" {
		baseURL = "https://registry.npmjs.org"
	}
	return &NPMRegistry{
		client:  resilienthttp.NewClient(),
		baseURL: baseURL,
		limiter: limiter,
	}
}

func (r *NPMRegistry) GetMetadata(ctx context.Context, name string) (*Metadata, error) {
	if r.limiter != nil {
		if err := r.limiter.Wait(ctx); err != nil {
			return nil, err
		}
	}

	url := fmt.Sprintf("%s/%s", r.baseURL, name)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := r.client.Do(req)
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
