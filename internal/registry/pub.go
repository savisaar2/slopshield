package registry

import (
	"fmt"
	"net/http"

	"github.com/savisaar2/slopshield/internal/resilienthttp"
)

type PubRegistry struct {
	client  *http.Client
	baseURL string
}

func NewPubRegistry(baseURL string) *PubRegistry {
	if baseURL == "" {
		baseURL = "https://pub.dev/packages"
	}
	return &PubRegistry{
		client:  resilienthttp.NewClient(),
		baseURL: baseURL,
	}
}

func (r *PubRegistry) GetMetadata(name string) (*Metadata, error) {
	url := fmt.Sprintf("%s/%s", r.baseURL, name)
	resp, err := r.client.Head(url)
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
