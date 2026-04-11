package registry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Rust
type RustRegistry struct{ client *http.Client }
func NewRustRegistry() *RustRegistry { return &RustRegistry{client: &http.Client{Timeout: 10 * time.Second}} }
func (r *RustRegistry) Exists(name string) (bool, error) {
	resp, err := r.client.Get(fmt.Sprintf("https://crates.io/api/v1/crates/%s", name))
	if err != nil { return false, err }
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK, nil
}

// PHP
type PHPRegistry struct{ client *http.Client }
func NewPHPRegistry() *PHPRegistry { return &PHPRegistry{client: &http.Client{Timeout: 10 * time.Second}} }
func (r *PHPRegistry) Exists(name string) (bool, error) {
	// Packagist requires vendor/package format. If not provided, it's definitely a slop
	if !strings.Contains(name, "/") { return false, nil }
	resp, err := r.client.Get(fmt.Sprintf("https://packagist.org/packages/%s.json", name))
	if err != nil { return false, err }
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK, nil
}

// Ruby
type RubyRegistry struct{ client *http.Client }
func NewRubyRegistry() *RubyRegistry { return &RubyRegistry{client: &http.Client{Timeout: 10 * time.Second}} }
func (r *RubyRegistry) Exists(name string) (bool, error) {
	resp, err := r.client.Get(fmt.Sprintf("https://rubygems.org/api/v1/gems/%s.json", name))
	if err != nil { return false, err }
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK, nil
}

// NuGet (C#)
type NuGetRegistry struct{ client *http.Client }
func NewNuGetRegistry() *NuGetRegistry { return &NuGetRegistry{client: &http.Client{Timeout: 10 * time.Second}} }
func (r *NuGetRegistry) Exists(name string) (bool, error) {
	resp, err := r.client.Get(fmt.Sprintf("https://api.nuget.org/v3-flatcontainer/%s/index.json", name))
	if err != nil { return false, err }
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK, nil
}

// Maven (Java)
type MavenRegistry struct{ client *http.Client }
func NewMavenRegistry() *MavenRegistry { return &MavenRegistry{client: &http.Client{Timeout: 10 * time.Second}} }
func (r *MavenRegistry) Exists(name string) (bool, error) {
	// Maven search API
	url := fmt.Sprintf("https://search.maven.org/solrsearch/select?q=a:%s&rows=1&wt=json", name)
	resp, err := r.client.Get(url)
	if err != nil { return false, err }
	defer resp.Body.Close()
	var res struct { Response struct { NumFound int `json:"numFound"` } `json:"response"` }
	json.NewDecoder(resp.Body).Decode(&res)
	return res.Response.NumFound > 0, nil
}

// GitHub Actions
type GitHubRegistry struct{ client *http.Client }
func NewGitHubRegistry() *GitHubRegistry { return &GitHubRegistry{client: &http.Client{Timeout: 10 * time.Second}} }
func (r *GitHubRegistry) Exists(name string) (bool, error) {
	// name is usually "owner/repo"
	resp, err := r.client.Get(fmt.Sprintf("https://github.com/%s", name))
	if err != nil { return false, err }
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK, nil
}
