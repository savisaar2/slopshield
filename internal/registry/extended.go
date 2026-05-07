package registry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/savisaar2/slopshield/internal/resilienthttp"
)

// Rust
type RustRegistry struct {
	client  *http.Client
	baseURL string
}

func NewRustRegistry(baseURL string) *RustRegistry {
	if baseURL == "" {
		baseURL = "https://crates.io/api/v1/crates"
	}
	return &RustRegistry{client: resilienthttp.NewClient(), baseURL: baseURL}
}
func (r *RustRegistry) GetMetadata(name string) (*Metadata, error) {
	resp, err := r.client.Get(fmt.Sprintf("%s/%s", r.baseURL, name))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return &Metadata{Exists: false}, nil
	}

	var meta struct {
		Crate struct {
			CreatedAt string `json:"created_at"`
		} `json:"crate"`
	}
	var createdAt time.Time
	if err := json.NewDecoder(resp.Body).Decode(&meta); err == nil {
		createdAt, _ = time.Parse(time.RFC3339, meta.Crate.CreatedAt)
	}
	return &Metadata{Exists: true, CreatedAt: createdAt}, nil
}

// PHP
type PHPRegistry struct {
	client  *http.Client
	baseURL string
}

func NewPHPRegistry(baseURL string) *PHPRegistry {
	if baseURL == "" {
		baseURL = "https://packagist.org/packages"
	}
	return &PHPRegistry{client: resilienthttp.NewClient(), baseURL: baseURL}
}
func (r *PHPRegistry) GetMetadata(name string) (*Metadata, error) {
	// Packagist requires vendor/package format. If not provided, it's definitely a slop
	if !strings.Contains(name, "/") {
		return &Metadata{Exists: false}, nil
	}
	resp, err := r.client.Get(fmt.Sprintf("%s/%s.json", r.baseURL, name))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return &Metadata{Exists: resp.StatusCode == http.StatusOK}, nil
}

// Ruby
type RubyRegistry struct {
	client  *http.Client
	baseURL string
}

func NewRubyRegistry(baseURL string) *RubyRegistry {
	if baseURL == "" {
		baseURL = "https://rubygems.org/api/v1/gems"
	}
	return &RubyRegistry{client: resilienthttp.NewClient(), baseURL: baseURL}
}
func (r *RubyRegistry) GetMetadata(name string) (*Metadata, error) {
	resp, err := r.client.Get(fmt.Sprintf("%s/%s.json", r.baseURL, name))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return &Metadata{Exists: resp.StatusCode == http.StatusOK}, nil
}

// NuGet (C#)
type NuGetRegistry struct {
	client  *http.Client
	baseURL string
}

func NewNuGetRegistry(baseURL string) *NuGetRegistry {
	if baseURL == "" {
		baseURL = "https://api.nuget.org/v3-flatcontainer"
	}
	return &NuGetRegistry{client: resilienthttp.NewClient(), baseURL: baseURL}
}
func (r *NuGetRegistry) GetMetadata(name string) (*Metadata, error) {
	resp, err := r.client.Get(fmt.Sprintf("%s/%s/index.json", r.baseURL, name))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return &Metadata{Exists: resp.StatusCode == http.StatusOK}, nil
}

// Maven (Java)
type MavenRegistry struct {
	client  *http.Client
	baseURL string
}

func NewMavenRegistry(baseURL string) *MavenRegistry {
	if baseURL == "" {
		baseURL = "https://search.maven.org/solrsearch/select"
	}
	return &MavenRegistry{client: resilienthttp.NewClient(), baseURL: baseURL}
}
func (r *MavenRegistry) GetMetadata(name string) (*Metadata, error) {
	// Maven search API
	url := fmt.Sprintf("%s?q=a:%s&rows=1&wt=json", r.baseURL, name)
	resp, err := r.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var res struct {
		Response struct{ NumFound int `json:"numFound"` } `json:"response"`
	}
	json.NewDecoder(resp.Body).Decode(&res)
	return &Metadata{Exists: res.Response.NumFound > 0}, nil
}

// GitHub Actions
type GitHubRegistry struct {
	client  *http.Client
	baseURL string
}

func NewGitHubRegistry(baseURL string) *GitHubRegistry {
	if baseURL == "" {
		baseURL = "https://github.com"
	}
	return &GitHubRegistry{client: resilienthttp.NewClient(), baseURL: baseURL}
}
func (r *GitHubRegistry) GetMetadata(name string) (*Metadata, error) {
	// name is usually "owner/repo"
	resp, err := r.client.Get(fmt.Sprintf("%s/%s", r.baseURL, name))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return &Metadata{Exists: resp.StatusCode == http.StatusOK}, nil
}
