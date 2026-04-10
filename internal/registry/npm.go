package registry

import (
	"fmt"
	"net/http"
	"time"
)

type NPMRegistry struct {
	client *http.Client
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
	resp, err := r.client.Head(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}
	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}
	return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}
