package aggregator

import (
	"bufio"
	"net/http"
	"strings"
	"time"
)

type Aggregator struct {
	sources []string
	client  *http.Client
}

func NewAggregator() *Aggregator {
	return &Aggregator{
		sources: []string{
			"https://raw.githubusercontent.com/ai-security/hallucinated-packages/main/npm.txt",
		},
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (a *Aggregator) FetchAll() (map[string]bool, error) {
	hallucinated := make(map[string]bool)
	for _, source := range a.sources {
		resp, err := a.client.Get(source)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			name := strings.TrimSpace(scanner.Text())
			if name != "" && !strings.HasPrefix(name, "#") {
				hallucinated[name] = true
			}
		}
	}
	return hallucinated, nil
}
