package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	RegistryURL  string `yaml:"registry_url"`
	RegistryPath string `yaml:"registry_path"` // Local directory for .json files
	Providers    struct {
		OpenAI    string `yaml:"openai_api_key"`
		Anthropic string `yaml:"anthropic_api_key"`
		Gemini    string `yaml:"gemini_api_key"`
		Ollama    struct {
			Enabled bool   `yaml:"enabled"`
			URL     string `yaml:"url"`
			Model   string `yaml:"model"`
		} `yaml:"ollama"`
	} `yaml:"providers"`
}

func Load() (*Config, error) {
	var cfg Config
	data, err := os.ReadFile("slopshield.yaml")
	if err == nil {
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, err
		}
	} else if !os.IsNotExist(err) {
		return nil, err
	}

	// Default Registry URL if not set
	if cfg.RegistryURL == "" {
		cfg.RegistryURL = "https://raw.githubusercontent.com/YOUR_USERNAME/slopshield/main/registry"
	}

	// Environment variable overrides
	if env := os.Getenv("OPENAI_API_KEY"); env != "" {
		cfg.Providers.OpenAI = env
	}
	if env := os.Getenv("ANTHROPIC_API_KEY"); env != "" {
		cfg.Providers.Anthropic = env
	}
	if env := os.Getenv("GEMINI_API_KEY"); env != "" {
		cfg.Providers.Gemini = env
	}
	if env := os.Getenv("OLLAMA_URL"); env != "" {
		cfg.Providers.Ollama.URL = env
	}
	if env := os.Getenv("OLLAMA_MODEL"); env != "" {
		cfg.Providers.Ollama.Model = env
	}

	return &cfg, nil
}
