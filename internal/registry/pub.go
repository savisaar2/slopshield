package registry

import (
	"fmt"
	"net/http"
	"time"
)

type PubRegistry struct {
	client *http.Client
}

func NewPubRegistry() *PubRegistry {
	return &PubRegistry{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (r *PubRegistry) Exists(name string) (bool, error) {
	url := fmt.Sprintf("https://pub.dev/packages/%s", name)
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
