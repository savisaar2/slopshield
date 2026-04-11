package registry

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNPMRegistry_Exists(t *testing.T) {
	tests := []struct {
		name           string
		packageName    string
		serverResponse string
		status         int
		want           bool
	}{
		{
			name:           "package exists and is old",
			packageName:    "lodash",
			status:         http.StatusOK,
			serverResponse: `{"time": {"created": "2012-04-23T18:25:43.511Z"}}`,
			want:           true,
		},
		{
			name:           "package exists but is brand new",
			packageName:    "suspicious-new-pkg",
			status:         http.StatusOK,
			serverResponse: `{"time": {"created": "` + time.Now().Format(time.RFC3339) + `"}}`,
			want:           false,
		},
		{
			name:           "package does not exist",
			packageName:    "non-existent-slop",
			status:         http.StatusNotFound,
			serverResponse: `{}`,
			want:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			// We need to inject the mock server URL into the registry check
			// For testing, we'll temporarily modify the registry check or use a helper
			r := &NPMRegistry{
				client:  &http.Client{Timeout: time.Second},
				baseURL: server.URL,
			}

			got, err := r.Exists(tt.packageName)
			if err != nil {
				t.Fatalf("Exists() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("Exists() got = %v, want %v", got, tt.want)
			}
		})
	}
}
