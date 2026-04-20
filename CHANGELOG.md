## [1.0.1](https://github.com/nfernal/terraform-provider-jumpcloud/compare/v1.0.0...v1.0.1) (2026-04-20)

### Bug Fixes

* install Terraform in CI generate step to avoid expired GPG key ([8b5e39e](https://github.com/nfernal/terraform-provider-jumpcloud/commit/8b5e39eb3f5d9cebfc1718fae4c4dc76e2a5de53))

## 1.0.0 (2026-04-20)

### Bug Fixes

* handle remaining unchecked error returns in test files ([88d50d2](https://github.com/nfernal/terraform-provider-jumpcloud/commit/88d50d2ec1876a5d04ff92b0130eedf748aec9b0))
* handle unchecked error returns flagged by errcheck linter ([eb71ba6](https://github.com/nfernal/terraform-provider-jumpcloud/commit/eb71ba6264f9e04bf3b1c311d2d6e753bd251d61))
* pin golangci-lint version and add missing tools/go.sum ([1d51996](https://github.com/nfernal/terraform-provider-jumpcloud/commit/1d51996937af636e6a93d097abd6267f1b5fbaf5))
* remove gosimple linter, merged into staticcheck in v2 ([bc36441](https://github.com/nfernal/terraform-provider-jumpcloud/commit/bc364416fa262486d9b91aba8933f64236f8cb01))
* resolve all remaining errcheck lint violations ([42f0525](https://github.com/nfernal/terraform-provider-jumpcloud/commit/42f052504b180fec0ffc9c6750d8da90938e344d))
* update Go version to 1.25.0 and fix semantic-release auth ([18a7f43](https://github.com/nfernal/terraform-provider-jumpcloud/commit/18a7f434bcdef21094bb8a58d15d7aea840d46b6))
* update golangci-lint to v2.11.4 for action v9 compatibility ([ef35033](https://github.com/nfernal/terraform-provider-jumpcloud/commit/ef3503350a708adb2cb93949566a0bbd5bad4bcf))
* upgrade terraform-plugin-docs and add missing semantic-release dep ([acfb668](https://github.com/nfernal/terraform-provider-jumpcloud/commit/acfb668c152d3928384e9b8725dbfa57087d04a3))
* use -coverprofile instead of -cover to avoid missing covdata tool ([ac40ebe](https://github.com/nfernal/terraform-provider-jumpcloud/commit/ac40ebe44bd00ec68b6c092cefb15715270e4a26))

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Initial provider implementation
- **Resources:** `jumpcloud_user`, `jumpcloud_user_group`, `jumpcloud_user_group_membership`, `jumpcloud_system_group`
- **Data Sources:** `jumpcloud_user`, `jumpcloud_user_group`
- Provider configuration with `api_key`, `org_id`, and `api_url` attributes
- Environment variable support for `JUMPCLOUD_API_KEY`, `JUMPCLOUD_ORG_ID`, `JUMPCLOUD_API_URL`
