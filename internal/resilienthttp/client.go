package resilienthttp

import (
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

// NewClient returns a pre-configured HTTP client with retries and exponential backoff.
func NewClient() *http.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	retryClient.RetryWaitMin = 1 * time.Second
	retryClient.RetryWaitMax = 5 * time.Second
	retryClient.Logger = nil // Disable verbose logging for CLI production use

	return retryClient.StandardClient()
}

// DoWithRetry is a helper for simple GET requests if needed, though standard client is preferred.
func DoWithRetry(req *http.Request) (*http.Response, error) {
	client := NewClient()
	return client.Do(req)
}
