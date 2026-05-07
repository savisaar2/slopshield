package registry

import (
	"fmt"
	"net/http"

	"github.com/savisaar2/slopshield/internal/resilienthttp"
)

type GoRegistry struct {
	client  *http.Client
	baseURL string
}

func NewGoRegistry(baseURL string) *GoRegistry {
	if baseURL == "" {
		baseURL = "https://proxy.golang.org"
	}
	return &GoRegistry{
		client:  resilienthttp.NewClient(),
		baseURL: baseURL,
	}
}

func (r *GoRegistry) GetMetadata(name string) (*Metadata, error) {
	// Go modules use the proxy.golang.org to verify existence
	url := fmt.Sprintf("%s/%s/@v/list", r.baseURL, name)
	resp, err := r.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return &Metadata{Exists: true}, nil
	}
	if resp.StatusCode == http.StatusNotFound {
		return &Metadata{Exists: false}, nil
	}
	return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}
