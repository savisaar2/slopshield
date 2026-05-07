package registry

import (
	"fmt"
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

func GetRegistry(ecosystem Ecosystem, baseURL string) (Registry, error) {
	switch ecosystem {
	case EcosystemNPM:
		return NewNPMRegistry(baseURL), nil
	case EcosystemPub:
		return NewPubRegistry(baseURL), nil
	case EcosystemPython:
		return NewPythonRegistry(baseURL), nil
	case EcosystemGo:
		return NewGoRegistry(baseURL), nil
	case EcosystemRust:
		return NewRustRegistry(baseURL), nil
	case EcosystemPHP:
		return NewPHPRegistry(baseURL), nil
	case EcosystemRuby:
		return NewRubyRegistry(baseURL), nil
	case EcosystemGitHub:
		return NewGitHubRegistry(baseURL), nil
	case EcosystemMaven:
		return NewMavenRegistry(baseURL), nil
	case EcosystemNuGet:
		return NewNuGetRegistry(baseURL), nil
	default:
		return nil, fmt.Errorf("unsupported ecosystem: %s", ecosystem)
	}
}
