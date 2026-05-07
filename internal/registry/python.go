package registry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/savisaar2/slopshield/internal/resilienthttp"
)

type PythonRegistry struct {
	client  *http.Client
	baseURL string
}

func NewPythonRegistry(baseURL string) *PythonRegistry {
	if baseURL == "" {
		baseURL = "https://pypi.org/pypi"
	}
	return &PythonRegistry{
		client:  resilienthttp.NewClient(),
		baseURL: baseURL,
	}
}

func (r *PythonRegistry) GetMetadata(name string) (*Metadata, error) {
	url := fmt.Sprintf("%s/%s/json", r.baseURL, name)
	resp, err := r.client.Get(url)
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
