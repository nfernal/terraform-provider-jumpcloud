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
