package registry

import (
	"context"
	"fmt"
	"net/http"

	"github.com/savisaar2/slopshield/internal/resilienthttp"
	"golang.org/x/time/rate"
)

type PubRegistry struct {
	client  *http.Client
	baseURL string
	limiter *rate.Limiter
}

func NewPubRegistry(baseURL string, limiter *rate.Limiter) *PubRegistry {
	if baseURL == "" {
		baseURL = "https://pub.dev/packages"
	}
	return &PubRegistry{
		client:  resilienthttp.NewClient(),
		baseURL: baseURL,
		limiter: limiter,
	}
}

func (r *PubRegistry) GetMetadata(ctx context.Context, name string) (*Metadata, error) {
	if r.limiter != nil {
		if err := r.limiter.Wait(ctx); err != nil {
			return nil, err
		}
	}

	url := fmt.Sprintf("%s/%s", r.baseURL, name)
	req, err := http.NewRequestWithContext(ctx, "HEAD", url, nil)
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
