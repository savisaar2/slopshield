# SlopShield 🛡️

**AI Package Hallucination Scanner.**

SlopShield is a security tool designed to detect "AI hallucinations"—non-existent packages suggested by LLMs—that can lead to supply chain attacks. It cross-references your project's dependencies against official registries and a community-maintained hallucination database.

---

## 🚀 Getting Started

### Installation
Ensure you have Go installed, then clone the repository:
```bash
git clone https://github.com/YOUR_USERNAME/slopshield.git
cd slopshield
go build -o slopshield.exe cmd/slopshield/main.go
```

### Basic Usage
Scan your current project:
```bash
./slopshield.exe scan .
```
Scan a Flutter project:
```bash
./slopshield.exe scan path/to/flutter_project
```

---

## ⚙️ Configuration

SlopShield uses a `slopshield.yaml` file for advanced configuration, such as custom registry URLs and AI provider API keys.

1.  Copy the example config:
    ```bash
    cp slopshield.yaml.example slopshield.yaml
    ```
2.  Edit `slopshield.yaml` and add your API keys for OpenAI, Anthropic, or Gemini.

---

## 🚀 Usage

### Automatic Project Scan
SlopShield automatically detects the project type (Node.js, Flutter, Python, Go) in the target directory.
```bash
./slopshield.exe scan .
```

### Multi-Engine Hallucination Probing
To update your registry using all configured AI providers simultaneously:
```bash
./slop-prober.exe --ecosystem npm
```
*Note: This command requires at least one API key or a local Ollama instance configured in `slopshield.yaml`.*

*The prober will:*
1.  **Detect** available API keys.
2.  **Probe** OpenAI, Anthropic, Gemini, and Ollama in parallel.
3.  **Extract** potential package names.
4.  **Verify** each name against the official registry.
5.  **Merge** new findings into `registry/*.json`.

### 3. Commit and Push
Once your local registry files are updated, commit and push them to your repo:
```bash
git add registry/
git commit -m "feat: discover new hallucinations across multiple AI engines"
git push
```

---

## 🛠️ Features
- **Multi-Ecosystem**: Supports Node.js (`package.json`) and Flutter (`pubspec.yaml`).
- **Reputation-Aware**: Automatically flags packages less than 14 days old as "suspicious," even if they exist (to prevent attacker-registered hallucinations).
- **Interactive Resolution**: Interactively ignore false positives or report confirmed findings.
- **CI/CD Ready**: Generates **SARIF** reports for the GitHub Security tab.
- **Decentralized**: Points to any GitHub URL or static file server for its intelligence.

---

## 🤝 Contributing
Found a hallucination? 
1.  Verify it with `slop-hunter`.
2.  Open a Pull Request to the `/registry` folder.
3.  Help protect the community!

---
*SlopShield is an open-source security research project. Always verify dependencies before use.*
