# SlopShield 🛡️ (v1.1.0)

**The Universal AI Package Hallucination Scanner.**

SlopShield is a local-first security tool designed to protect developers from "AI Hallucinations"—non-existent or malicious packages suggested by LLMs. By connecting SlopShield to your own LLM providers (OpenAI, Anthropic, Gemini, or Ollama), you can harvest, verify, and maintain your own private database of hallucinated packages.

---

## 🚀 Key Features

- **Multi-Ecosystem Support**: Auto-detects and scans Node.js, Flutter, Python, Go, Rust, PHP, Ruby, Java, C#, and GitHub Actions.
- **Multi-Engine Prober**: Automatically harvests new hallucinations across OpenAI, Anthropic, Gemini, and Ollama simultaneously to build your personal registry.
- **Reputation Analysis**: Flags suspiciously new packages (less than 14 days old) even if they exist in the registry.
- **SARIF Integration**: Generates industry-standard reports for GitHub Security Tab and CI/CD pipelines.
- **Interactive TUI**: A beautiful terminal interface for manual intervention and ignore-list management.
- **Local Intelligence**: Your hallucination database stays on your machine, updated by your own AI probing.

---

## 🛠️ Installation

```bash
git clone https://github.com/savisaar2/slopshield.git
cd slopshield
go build -o slopshield.exe cmd/slopshield/main.go
go build -o slop-prober.exe cmd/slop-prober/main.go
```

---

## ⚙️ Configuration

SlopShield uses `slopshield.yaml` for API keys and local settings.

1.  **Initialize Config**:
    ```bash
    cp slopshield.yaml.example slopshield.yaml
    ```
2.  **Add API Keys**: Edit `slopshield.yaml` to include keys for OpenAI, Anthropic, or Gemini to enable your personal Prober.

---

## 🔍 Usage

### Scan a Project
Simply point SlopShield to any directory. It will automatically detect the manifest files and check them against your local registry and the official ecosystem registries.
```bash
./slopshield.exe scan .
```

---

## 🛡️ CI/CD Integration

SlopShield is designed to be easily integrated into your existing CI/CD pipelines to prevent AI-hallucinated packages from being merged into your codebase.

### How it Works in CI
Even if your local `registry/*.json` files are empty, SlopShield will still protect you. It uses a **Tiered Verification** approach:
1.  **Local Registry**: Fast-check against known "slops" you've already found.
2.  **Official Registry (The Truth Check)**: If a package is not in your registry, SlopShield queries the official source (e.g., npmjs.org). If the package doesn't exist (404), it's flagged as a hallucination. It also flags packages less than 14 days old as suspicious.

This means you get immediate protection out-of-the-box, while your personal registry grows more intelligent over time as you use the Prober.

### GitHub Actions
You can use SlopShield to scan your repository on every push or pull request. The following example shows how to run SlopShield and upload the results as a SARIF report to GitHub's Security Tab.

```yaml
name: SlopShield Security Scan
on: [push, pull_request]

jobs:
  scan:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout Code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'

      - name: Install SlopShield
        run: |
          git clone https://github.com/savisaar2/slopshield.git
          cd slopshield
          go build -o ../slopshield cmd/slopshield/main.go

      - name: Run SlopShield Scan
        run: |
          ./slopshield scan --format sarif --output results.sarif .

      - name: Upload SARIF report
        uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: results.sarif
```

### Other CI Environments
For other CI systems (GitLab, Jenkins, etc.), you can run the scanner and check the exit code. SlopShield will exit with a non-zero code if hallucinations are detected.

```bash
# Example for a generic CI script
go build -o slopshield cmd/slopshield/main.go
./slopshield scan .
```

---

## 🎯 Maintainer Tools (Personal Intelligence)

Keep your local registry fresh using the built-in AI harvester.

### The Prober (Automatic Discovery & Verification)
Connect to your AI providers to solve niche tasks. SlopShield will extract the non-existent packages they suggest, verify them, and save them to your local `registry/` folder automatically.
```bash
# Run for a specific ecosystem
./slop-prober.exe --ecosystem npm
```

### The Hunter (Manual Entry)
If you find a specific hallucination, verify and merge it manually:
```bash
go run cmd/slop-hunter/main.go --update npm "obscure-pkg-1,non-existent-pkg-2"
```

---

## 📦 Supported Manifests

| Ecosystem | Manifest File | Registry |
| :--- | :--- | :--- |
| **Node.js** | `package.json` | npmjs.org |
| **Flutter/Dart** | `pubspec.yaml` | pub.dev |
| **Python** | `requirements.txt` | pypi.org |
| **Go** | `go.mod` | proxy.golang.org |
| **Rust** | `Cargo.toml` | crates.io |
| **PHP** | `composer.json` | packagist.org |
| **Ruby** | `Gemfile` | rubygems.org |
| **C# / .NET** | `.csproj` | nuget.org |
| **Java / Maven** | `pom.xml` | maven.org |
| **GitHub Actions**| `.github/workflows/*.yml` | github.com |

---

## 🧠 How it Works: The Operational Flow

1.  **Auto-Discovery**: The scanner identifies project manifests (e.g., `Cargo.toml`).
2.  **Local Intelligence**: It loads your `registry/*.json` files—packages you've previously identified as slops using the Prober.
3.  **Dependency Extraction**: Specific ecosystem parsers extract all dependencies.
4.  **The Filter**:
    -   **Tier 1**: Checks your `.slopignore`.
    -   **Tier 2**: Checks your local personal registry.
    -   **Tier 3 (The Truth Check)**: Asks the official registry. If 404, it's a slop. If < 14 days old, it's suspicious.
5.  **Report**: Beautiful UI or SARIF export.

---

## ☕ Support the Project
If SlopShield helped secure your project, consider supporting the developer:

[!["Buy Me A Coffee"](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/savisaar2d)

---
*Built by savisaar2. Secure your code in the age of AI.*
