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

type PythonRegistry struct {
	client  *http.Client
	baseURL string
	limiter *rate.Limiter
}

func NewPythonRegistry(baseURL string, limiter *rate.Limiter) *PythonRegistry {
	if baseURL == "" {
		baseURL = "https://pypi.org/pypi"
	}
	return &PythonRegistry{
		client:  resilienthttp.NewClient(),
		baseURL: baseURL,
		limiter: limiter,
	}
}

func (r *PythonRegistry) GetMetadata(ctx context.Context, name string) (*Metadata, error) {
	if r.limiter != nil {
		if err := r.limiter.Wait(ctx); err != nil {
			return nil, err
		}
	}

	url := fmt.Sprintf("%s/%s/json", r.baseURL, name)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var meta struct {
			Releases map[string][]struct {
				UploadTime string `json:"upload_time"`
			} `json:"releases"`
		}
		var earliest time.Time
		if err := json.NewDecoder(resp.Body).Decode(&meta); err == nil {
			for _, releases := range meta.Releases {
				for _, release := range releases {
					t, _ := time.Parse("2006-01-02T15:04:05", release.UploadTime)
					if earliest.IsZero() || t.Before(earliest) {
						earliest = t
					}
				}
			}
		}
		return &Metadata{Exists: true, CreatedAt: earliest}, nil
	}
	if resp.StatusCode == http.StatusNotFound {
		return &Metadata{Exists: false}, nil
	}
	return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}
