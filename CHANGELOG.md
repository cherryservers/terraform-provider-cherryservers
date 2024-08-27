# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Server re-installation functionality. Updating server resource attributes `image`, `os_partition_size`, `ssh_key_ids` or `user_data` will now trigger a
  server re-install, if `allow_reinstall` is set to true.

### Deprecated

- Server resource attribute `ip_addresses_ids` is now deprecated. Use `extra_ip_addresses_ids`.

### Changed

- Completed migration from SDKv2 to Framework.
- Terraform 1.8 is now the minimum required version.
- Go 1.21 is now the minimum required version, for building the provider.
- In all data sources, the identifying attribute is now called `id`, instead of `{resource-name}_id` or something
  similar.
- IP resource and data source attribute `target_ip_id` is now called `route_ip_id`.
- IP resource and data source attributes `ptr_record` and `a_record` now have derived attributes `ptr_record_actual` and
  `a_record_actual` that are used to keep the actual resource state, while the old ones are used for configuration.
- SSH key resource and data source attribute `name` is now called `label`.
- Server resource and data source attribute `project_id` type changed from `string` to `int64`.
- Server resource and data source attribute `ssh_key_ids` type changed from `list` to `set`.

[unreleased]: https://github.com/cherryservers/terraform-provider-cherryservers/compare/v0.0.6...0.1.0
