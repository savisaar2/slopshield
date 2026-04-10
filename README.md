# SlopShield 🛡️

AI Package Hallucination Scanner.

SlopShield helps prevent supply chain attacks by identifying packages in your project that might be hallucinations from LLMs. It cross-references your dependencies against official registries (NPM) and known hallucination lists.

## Features

- **NPM Support**: Scans `package.json` for non-existent dependencies.
- **Aggregated Hallucination Lists**: Fetches known hallucinated names from multiple community-driven sources.
- **GitHub Integration**: Generates SARIF reports for the GitHub Security tab.
- **Ignore List**: Support for `.slopignore` to skip internal or private packages.
- **CI/CD Ready**: Returns a non-zero exit code if hallucinations are found.

## Installation

```bash
go install github.com/savisaar2/slopshield/cmd/slopshield@latest
```

## Usage

### Scan a project
```bash
slopshield scan .
```

### Generate SARIF for GitHub
```bash
slopshield scan . --output sarif > results.sarif
```

### Ignore packages
Create a `.slopignore` file:
```text
my-private-package
test-hallucination-*
```

## How it works

1. **Extract**: Reads your `package.json`.
2. **Aggregate**: Fetches known hallucinated packages from multiple LLM-tracking repositories.
3. **Verify**: Checks the official registry (NPM) for any package that isn't on the known list.
4. **Report**: Outputs findings to the terminal or a SARIF file.
