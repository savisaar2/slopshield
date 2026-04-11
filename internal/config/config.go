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
	data, err := os.ReadFile("slopshield.yaml")
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{RegistryURL: "https://raw.githubusercontent.com/YOUR_USERNAME/slopshield/main/registry"}, nil
		}
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
