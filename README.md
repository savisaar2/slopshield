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

## 🔗 Connecting Your Own Registry (Decentralization)

By default, SlopShield points to its own GitHub repository for hallucination data. You can "own" your registry by pointing it to your own fork or a different repository.

### 1. Host Your Own Registry
1.  Fork this repository or create a new one.
2.  Ensure there is a `/registry` folder with `npm.json` and `pub.json`.
3.  Set the `SLOPSHIELD_REGISTRY_URL` environment variable:
    ```bash
    # Windows (PowerShell)
    $env:SLOPSHIELD_REGISTRY_URL="https://raw.githubusercontent.com/YOUR_USERNAME/YOUR_REPO/main/registry"
    
    # Unix/Linux
    export SLOPSHIELD_REGISTRY_URL="https://raw.githubusercontent.com/YOUR_USERNAME/YOUR_REPO/main/registry"
    ```

### 2. Manual Sync
Update your local cache from your remote registry at any time:
```bash
./slopshield.exe sync
```

---

## 🎯 Maintaining the Registry (The Hunter-Gatherer Workflow)

To keep your registry up-to-date, use the built-in "Bait and Catch" toolset.

### Phase 1: The Prober (Harvesting Bait)
The `slop-prober` uses an LLM (like OpenAI) to intentionally "bait" hallucinations by asking for niche code snippets.
```bash
$env:OPENAI_API_KEY="your_key"
go run cmd/slop-prober/main.go
```
*Output: `express-gpt-parser, flutter-zeus-iot-sdk`*

### Phase 2: The Hunter (Verifying & Merging)
Use `slop-hunter` to verify those names against official registries. If it gets a **404**, it's a confirmed hallucination. The `--update` flag will automatically merge it into your local `registry/*.json` files.
```bash
go run cmd/slop-hunter/main.go --update npm "express-gpt-parser,flutter-zeus-iot-sdk"
```

### Phase 3: Committing
Once your local registry files are updated, commit and push them to your repo:
```bash
git add registry/
git commit -m "feat: discover 2 new hallucinations"
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
