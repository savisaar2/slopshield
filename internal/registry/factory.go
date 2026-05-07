package registry

import (
	"fmt"

	"golang.org/x/time/rate"
)

type Ecosystem string

const (
	EcosystemNPM    Ecosystem = "npm"
	EcosystemPub    Ecosystem = "pub"
	EcosystemPython Ecosystem = "python"
	EcosystemGo     Ecosystem = "go"
	EcosystemRust   Ecosystem = "rust"
	EcosystemPHP    Ecosystem = "php"
	EcosystemRuby   Ecosystem = "ruby"
	EcosystemGitHub Ecosystem = "github"
	EcosystemMaven  Ecosystem = "maven"
	EcosystemNuGet  Ecosystem = "nuget"
)

func GetRegistry(ecosystem Ecosystem, baseURL string, limiter *rate.Limiter) (Registry, error) {
	switch ecosystem {
	case EcosystemNPM:
		return NewNPMRegistry(baseURL, limiter), nil
	case EcosystemPub:
		return NewPubRegistry(baseURL, limiter), nil
	case EcosystemPython:
		return NewPythonRegistry(baseURL, limiter), nil
	case EcosystemGo:
		return NewGoRegistry(baseURL, limiter), nil
	case EcosystemRust:
		return NewRustRegistry(baseURL, limiter), nil
	case EcosystemPHP:
		return NewPHPRegistry(baseURL, limiter), nil
	case EcosystemRuby:
		return NewRubyRegistry(baseURL, limiter), nil
	case EcosystemGitHub:
		return NewGitHubRegistry(baseURL, limiter), nil
	case EcosystemMaven:
		return NewMavenRegistry(baseURL, limiter), nil
	case EcosystemNuGet:
		return NewNuGetRegistry(baseURL, limiter), nil
	default:
		return nil, fmt.Errorf("unsupported ecosystem: %s", ecosystem)
	}
}
