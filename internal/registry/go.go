package registry

import (
	"fmt"
	"net/http"
	"time"
)

type GoRegistry struct {
	client *http.Client
}

func NewGoRegistry() *GoRegistry {
	return &GoRegistry{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (r *GoRegistry) Exists(name string) (bool, error) {
	// Go modules use the proxy.golang.org to verify existence
	url := fmt.Sprintf("https://proxy.golang.org/%s/@v/list", name)
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
