package registry

import (
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
		return true, nil
	}
	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}
	return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}
