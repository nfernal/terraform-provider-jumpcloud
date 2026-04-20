# Contributing to terraform-provider-jumpcloud

Thank you for your interest in contributing! This document provides guidelines and instructions for contributing.

## Development Environment

### Requirements

- [Go](https://golang.org/doc/install) >= 1.22
- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [golangci-lint](https://golangci-lint.run/welcome/install/)

### Building the Provider

```shell
make build
```

### Installing Locally

```shell
make install
```

### Running Tests

Run unit tests:

```shell
make test
```

Run acceptance tests (these create real resources in JumpCloud and require API credentials):

```shell
export JUMPCLOUD_API_KEY="your-api-key"
make testacc
```

### Linting

```shell
make lint
```

### Generating Documentation

Documentation is generated from the provider schema and example files using [terraform-plugin-docs](https://github.com/hashicorp/terraform-plugin-docs):

```shell
make generate
```

Always run `make generate` before submitting a PR to ensure documentation is up to date.

## Making Changes

1. Fork the repository and create a feature branch from `main`.
2. Make your changes, including tests for new functionality.
3. Run `make fmt lint test` to ensure code quality.
4. Run `make generate` to update documentation.
5. Commit your changes with a clear commit message.
6. Open a pull request against `main`.

## Commit Messages

Follow conventional commit style:

- `feat: add new resource for X`
- `fix: handle nil response in user creation`
- `docs: update provider configuration example`
- `test: add acceptance tests for user_group`
- `chore: update dependencies`

## Adding a New Resource or Data Source

1. Create the resource/data source implementation in `internal/provider/`.
2. Create the corresponding API client methods in `internal/client/`.
3. Register it in the provider's `Resources()` or `DataSources()` method in `internal/provider/provider.go`.
4. Add acceptance tests.
5. Add example configuration in `examples/`.
6. Run `make generate` to create documentation.

## Acceptance Tests

Acceptance tests create real infrastructure and are gated behind `TF_ACC=1`. Follow these conventions:

- Name test functions with `TestAcc` prefix: `TestAccUserResource_basic`
- Include a basic test, an update test, and an import test where applicable.
- Clean up resources in `CheckDestroy` functions.

## Reporting Issues

Use [GitHub Issues](https://github.com/nfernal/terraform-provider-jumpcloud/issues) to report bugs or request features. Please use the provided issue templates.

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.
