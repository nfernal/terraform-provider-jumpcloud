# Terraform Provider for JumpCloud

[![CI](https://github.com/nfernal/terraform-provider-jumpcloud/actions/workflows/ci.yml/badge.svg)](https://github.com/nfernal/terraform-provider-jumpcloud/actions/workflows/ci.yml)
[![Release](https://github.com/nfernal/terraform-provider-jumpcloud/actions/workflows/release.yml/badge.svg)](https://github.com/nfernal/terraform-provider-jumpcloud/actions/workflows/release.yml)
[![Semantic Release](https://github.com/nfernal/terraform-provider-jumpcloud/actions/workflows/semantic-release.yml/badge.svg)](https://github.com/nfernal/terraform-provider-jumpcloud/actions/workflows/semantic-release.yml)
[![codecov](https://codecov.io/gh/nfernal/terraform-provider-jumpcloud/branch/main/graph/badge.svg)](https://codecov.io/gh/nfernal/terraform-provider-jumpcloud)
[![Go Version](https://img.shields.io/github/go-mod/go-version/nfernal/terraform-provider-jumpcloud)](https://github.com/nfernal/terraform-provider-jumpcloud/blob/main/go.mod)
[![GitHub Release](https://img.shields.io/github/v/release/nfernal/terraform-provider-jumpcloud)](https://github.com/nfernal/terraform-provider-jumpcloud/releases/latest)
[![GitHub Release Date](https://img.shields.io/github/release-date/nfernal/terraform-provider-jumpcloud)](https://github.com/nfernal/terraform-provider-jumpcloud/releases/latest)
[![GitHub Last Commit](https://img.shields.io/github/last-commit/nfernal/terraform-provider-jumpcloud)](https://github.com/nfernal/terraform-provider-jumpcloud/commits/main)
[![GitHub Languages](https://img.shields.io/github/languages/count/nfernal/terraform-provider-jumpcloud)](https://github.com/nfernal/terraform-provider-jumpcloud)
[![License](https://img.shields.io/github/license/nfernal/terraform-provider-jumpcloud)](https://github.com/nfernal/terraform-provider-jumpcloud/blob/main/LICENSE)
[![Terraform Registry](https://img.shields.io/badge/terraform-registry-blueviolet?logo=terraform)](https://registry.terraform.io/providers/nfernal/jumpcloud/latest)

The JumpCloud Terraform Provider allows you to manage [JumpCloud](https://jumpcloud.com/) resources using [Terraform](https://www.terraform.io/) or [OpenTofu](https://opentofu.org/).

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.22 (to build the provider)
- A [JumpCloud](https://jumpcloud.com/) account and API key

## Getting Started

```hcl
terraform {
  required_providers {
    jumpcloud = {
      source  = "nfernal/jumpcloud"
      version = "~> 0.1"
    }
  }
}

provider "jumpcloud" {
  # api_key = var.jumpcloud_api_key  # Or set JUMPCLOUD_API_KEY env var
  # org_id  = var.jumpcloud_org_id   # Or set JUMPCLOUD_ORG_ID env var
}
```

## Authentication

The provider requires a JumpCloud API key. You can configure it in one of two ways:

1. **Environment variable** (recommended): Set `JUMPCLOUD_API_KEY`
2. **Provider configuration**: Set the `api_key` attribute

For multi-tenant organizations, set `JUMPCLOUD_ORG_ID` or the `org_id` attribute.

## Resources

- `jumpcloud_user` - Manage JumpCloud users
- `jumpcloud_user_group` - Manage user groups
- `jumpcloud_user_group_membership` - Manage user group memberships
- `jumpcloud_system_group` - Manage system groups

## Data Sources

- `jumpcloud_user` - Look up a user by email
- `jumpcloud_user_group` - Look up a user group

## Example Usage

```hcl
resource "jumpcloud_user" "example" {
  username   = "john.doe"
  email      = "john.doe@example.com"
  firstname  = "John"
  lastname   = "Doe"
  department = "Engineering"
}

resource "jumpcloud_user_group" "engineering" {
  name        = "Engineering"
  description = "Engineering team"
}

resource "jumpcloud_user_group_membership" "example" {
  user_id  = jumpcloud_user.example.id
  group_id = jumpcloud_user_group.engineering.id
}
```

See the [`examples/`](examples/) directory for more complete examples.

## Developing the Provider

### Building

```shell
make build
```

### Running Tests

Unit tests:

```shell
make test
```

Acceptance tests (requires `JUMPCLOUD_API_KEY`):

```shell
make testacc
```

### Generating Documentation

```shell
make generate
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for full development guidelines.

## License

This project is licensed under the [MIT License](LICENSE).
