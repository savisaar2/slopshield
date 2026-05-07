package registry

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNPMRegistry_GetMetadata(t *testing.T) {
	tests := []struct {
		name           string
		packageName    string
		serverResponse string
		status         int
		wantExists     bool
		checkCreated   bool
	}{
		{
			name:           "package exists",
			packageName:    "lodash",
			status:         http.StatusOK,
			serverResponse: `{"time": {"created": "2012-04-23T18:25:43.511Z"}}`,
			wantExists:     true,
			checkCreated:   true,
		},
		{
			name:           "package does not exist",
			packageName:    "non-existent-slop",
			status:         http.StatusNotFound,
			serverResponse: `{}`,
			wantExists:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			r := &NPMRegistry{
				client:  &http.Client{Timeout: time.Second},
				baseURL: server.URL,
			}

			meta, err := r.GetMetadata(context.Background(), tt.packageName)
			if err != nil {
				t.Fatalf("GetMetadata() error = %v", err)
			}
			if meta.Exists != tt.wantExists {
				t.Errorf("GetMetadata() Exists = %v, want %v", meta.Exists, tt.wantExists)
			}
			if tt.checkCreated && meta.CreatedAt.IsZero() {
				t.Error("GetMetadata() expected non-zero CreatedAt")
			}
		})
	}
}
