# Contributing to SlopShield 🛡️

First off, thank you for considering contributing to SlopShield! As a security tool focused on accuracy and reliability, we maintain high standards for code quality and architectural integrity.

## Our Standards

We follow a **Senior Engineer** philosophy. This means:
- **Zero Laziness**: No temporary fixes or hacky workarounds. Find and resolve root causes.
- **Architectural Consistency**: Adhere strictly to the established **Factory Patterns** for scanners and registries.
- **Test-Driven**: Every bug fix must include a reproduction test case. Every new feature must have 80%+ test coverage.
- **Context-Aware**: All new registry or network-related code MUST support `context.Context` for cancellation and timeouts.

## How to Contribute

1.  **Fork and Clone**: Create your feature branch from `main`.
2.  **Implementation**: Follow the existing Go idioms and project structure.
3.  **Validate**: Run `go test ./...` and `go vet ./...` before submitting.
4.  **Documentation**: Update `README.md` or configuration examples if your change introduces new settings.
5.  **Pull Request**: Provide a clear description of the "why" behind your change, not just the "what".

## Reporting Hallucinations

If you find a new hallucinated package, please use the `slop-hunter` tool to verify it and then open an issue or submit a PR updating the relevant JSON file in the `registry/` directory.

---
*Secure your code. Stop the slop.*
