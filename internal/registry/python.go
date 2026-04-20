package registry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type PythonRegistry struct {
	client *http.Client
}

func NewPythonRegistry() *PythonRegistry {
	return &PythonRegistry{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (r *PythonRegistry) Exists(name string) (bool, error) {
	url := fmt.Sprintf("https://pypi.org/pypi/%s/json", name)
	resp, err := r.client.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var meta struct {
			Info struct {
				Created string `json:"created"` // PyPI actually puts this in releases or info sometimes
			} `json:"info"`
			Releases map[string][]struct {
				UploadTime string `json:"upload_time"`
			} `json:"releases"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&meta); err == nil {
			// Find the earliest upload time
			var earliest time.Time
			for _, releases := range meta.Releases {
				for _, release := range releases {
					t, _ := time.Parse("2006-01-02T15:04:05", release.UploadTime)
					if earliest.IsZero() || t.Before(earliest) {
						earliest = t
					}
				}
			}
			if !earliest.IsZero() && time.Since(earliest) < 14*24*time.Hour {
				return false, nil // Suspiciously new
			}
		}
		return true, nil
	}
	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}
	return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}
