# Contributing to Pandora Core

Thank you for your interest in contributing to **Pandora Core** — we appreciate your support in improving this open source project!

This document outlines the basic process for contributing code, reporting issues, and participating in discussions.

## :raising_hand: Ways to Contribute

We welcome contributions in many forms:

* Reporting bugs
* Requesting features or enhancements
* Submitting pull requests with code improvements
* Contributing to documentation (code comments, README, DEVELOPMENT.md, etc.)
* Sharing IDE setups or developer tips in [DEVELOPMENT.md](./DEVELOPMENT.md)

## :bug: Reporting Issues

If you find a bug or unexpected behavior, please [open an issue](https://github.com/PandoraSuite/pandora-core/issues) with the following information:

* A descriptive title
* Steps to reproduce the issue
* Expected vs actual behavior
* Version or commit hash (if applicable)

> Note: At the moment we don't use issue templates — just be as clear and concise as possible.

## :package: Submitting Pull Requests

We welcome pull requests to fix bugs, improve code quality, or add new functionality. Before submitting a PR:

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature-name`
3. Make your changes and write tests if applicable
4. Run tests: `go test ./...`
5. Commit using a clear message (e.g., `fix: prevent quota underflow`)
6. Push your changes and open a pull request to the `master` branch

### :white_check_mark: Pull Request Guidelines

* Clearly describe the purpose of the PR
* Include references to related issues if applicable (e.g., `Fixes #42`)
* Small, focused PRs are preferred over large or mixed changes
* Ensure the code passes checks

> Note: We currently don't use PR templates, but you're encouraged to explain your changes in detail.

## :test_tube: Code Style & Testing

* Go code should follow standard formatting `go fmt ./...`
* Validate your changes with `go vet ./...` and `go test ./...`
* New features should include relevant unit tests where appropriate

## :lock: Security Disclosure

If you discover a security vulnerability, **do not open an issue or PR**. Instead, please follow the secure disclosure process described in [`SECURITY.md`](./SECURITY.md).

## :handshake: Contributor Expectations

All contributors are expected to respect others, collaborate in good faith, and follow community norms. Inappropriate behavior will not be tolerated.

Thank you for being part of the Pandora Core community!
