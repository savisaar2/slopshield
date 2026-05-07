# SlopShield 🛡️ 

[![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/savisaar2/slopshield)](https://github.com/savisaar2/slopshield/tags)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/savisaar2/slopshield/ci.yml?branch=main)](https://github.com/savisaar2/slopshield/actions)
[![Go Version](https://img.shields.io/github/go-mod/go-version/savisaar2/slopshield)](https://github.com/savisaar2/slopshield)
[![License](https://img.shields.io/badge/license-MIT-green)](./LICENSE)
[![GitHub issues](https://img.shields.io/github/issues/savisaar2/slopshield)](https://github.com/savisaar2/slopshield/issues)

**The Universal AI Package Hallucination Scanner.**

SlopShield is a local-first security tool designed to protect developers from "AI Hallucinations"—non-existent or malicious packages suggested by Large Language Models (LLMs). By connecting SlopShield to your own LLM providers (OpenAI, Anthropic, Gemini, or Ollama), you can harvest, verify, and maintain your own private database of hallucinated packages.

---

## 🚀 Key Features

- **Typosquatting Detection**: Flags packages that are suspiciously close to popular libraries (e.g., `lodsh` vs `lodash`).
- **In-Memory Caching**: Ultra-fast scanning for monorepos by caching registry lookups.
- **Private Registry Support**: Point SlopShield to your internal Artifactory/Nexus to avoid false positives on private packages.
- **Multi-Ecosystem Support**: Auto-detects and scans Node.js, Flutter, Python, Go, Rust, PHP, Ruby, Java, C#, and GitHub Actions.
- **Multi-Engine Prober**: Automatically harvests new hallucinations across OpenAI, Anthropic, Gemini, and Ollama simultaneously to build your personal registry.
- **Reputation Analysis**: Flags suspiciously new packages (configurable threshold) even if they exist in the registry.
- **SARIF Integration**: Generates industry-standard reports for GitHub Security Tab and CI/CD pipelines.

---

## 🛠️ Quick Start

### One-Liner Install (Go 1.21+)

Install all SlopShield tools directly:
```bash
go install github.com/savisaar2/slopshield/cmd/...@latest
```
*Ensure your `$GOPATH/bin` (usually `~/go/bin`) is in your PATH.*

### Manual Installation

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/savisaar2/slopshield.git
    cd slopshield
    ```
2.  **Build the binaries**:
    ```bash
    go build -o slopshield cmd/slopshield/main.go
    go build -o slop-prober cmd/slop-prober/main.go
    go build -o slop-hunter cmd/slop-hunter/main.go
    ```

---

## 🔍 Usage

### 1. Scan a Project
Point SlopShield to any directory. It will automatically detect manifest files (e.g., `package.json`, `requirements.txt`, `go.mod`, `pom.xml`, etc.) and check them against your local registry and official sources.
```bash
slopshield scan .
```

### 2. Configure AI Providers
To use the Prober (automatic discovery), you can use environment variables (recommended for CI/CD) or a configuration file.

#### Option A: Environment Variables (Recommended)
```bash
export OPENAI_API_KEY="your-key"
export ANTHROPIC_API_KEY="your-key"
export GEMINI_API_KEY="your-key"
```

#### Option B: Configuration File
```bash
cp slopshield.yaml.example slopshield.yaml
```
Edit `slopshield.yaml` to include your API keys. This file is automatically ignored by git.

### 3. Harvest Hallucinations (The Prober)
Connect to your AI providers to discover niche hallucinations. SlopShield will verify them and save them to your local `registry/` folder.
```bash
slop-prober --ecosystem npm
```

### 4. Manual Verification (The Hunter)
If you find a specific potential hallucination, verify and merge it manually:
```bash
slop-hunter --update npm "fake-package-1,another-slop-pkg"
```

---

## 📦 Supported Manifests

| Ecosystem | Manifest File | Registry Source |
| :--- | :--- | :--- |
| **Node.js** | `package.json` | npmjs.org |
| **Python** | `requirements.txt` | pypi.org |
| **Go** | `go.mod` | proxy.golang.org |
| **Rust** | `Cargo.toml` | crates.io |
| **Java / Maven** | `pom.xml` | maven.org |
| **C# / .NET** | `.csproj` | nuget.org |
| **Flutter/Dart** | `pubspec.yaml` | pub.dev |
| **PHP** | `composer.json` | packagist.org |
| **Ruby** | `Gemfile` | rubygems.org |
| **GitHub Actions**| `.github/workflows/*.yml` | github.com |

---

## 🛡️ CI/CD Integration

SlopShield is designed to be integrated into CI/CD pipelines to prevent AI-hallucinated packages from being merged.

### GitHub Actions Example

```yaml
name: SlopShield Security Scan
on: [push, pull_request]

jobs:
  scan:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout Code
        uses: actions/checkout@v4

      - name: Install SlopShield
        run: go install github.com/savisaar2/slopshield/cmd/slopshield@latest

      - name: Run Scan
        run: slopshield scan --format sarif --output results.sarif .

      - name: Upload SARIF report
        uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: results.sarif
```

---

## 🐳 Docker Support

Run SlopShield without installing Go:
```bash
docker build -t slopshield .
docker run --rm -v $(pwd):/scan slopshield scan /scan
```

---

## 🧠 How it Works: Layered Detection

SlopShield uses a layered approach to verify every dependency:
1.  **Local Registry**: Fast-check against known "slops" you've previously identified.
2.  **Typosquatting Check**: Identifies packages that mimic popular libraries (e.g., `requesst` vs `requests`).
3.  **Official Registry (The Truth Check)**: Queries official sources (npmjs.org, pypi.org, etc.). If it returns a 404, it's flagged.
4.  **Reputation Check**: Flags packages that are suspiciously new (default is < 14 days, configurable in `slopshield.yaml`).

---

## ☕ Support the Project
If SlopShield helped secure your project, consider supporting the developer:

[!["Buy Me A Coffee"](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/savisaar2d)

---
*Built by savisaar2. Secure your code in the age of AI.*
