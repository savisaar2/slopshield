package registry

import (
	"context"
	"fmt"
	"net/http"

	"github.com/savisaar2/slopshield/internal/resilienthttp"
	"golang.org/x/time/rate"
)

type GoRegistry struct {
	client  *http.Client
	baseURL string
	limiter *rate.Limiter
}

func NewGoRegistry(baseURL string, limiter *rate.Limiter) *GoRegistry {
	if baseURL == "" {
		baseURL = "https://proxy.golang.org"
	}
	return &GoRegistry{
		client:  resilienthttp.NewClient(),
		baseURL: baseURL,
		limiter: limiter,
	}
}

func (r *GoRegistry) GetMetadata(ctx context.Context, name string) (*Metadata, error) {
	if r.limiter != nil {
		if err := r.limiter.Wait(ctx); err != nil {
			return nil, err
		}
	}

	// Go modules use the proxy.golang.org to verify existence
	url := fmt.Sprintf("%s/%s/@v/list", r.baseURL, name)
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
		return &Metadata{Exists: true}, nil
	}
	if resp.StatusCode == http.StatusNotFound {
		return &Metadata{Exists: false}, nil
	}
	return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}
