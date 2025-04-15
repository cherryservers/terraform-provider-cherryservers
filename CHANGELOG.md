# Changelog

Changelog moved to Release Notes in [Github Releases](https://github.com/cherryservers/terraform-provider-cherryservers/releases)

## [1.0.1] - 2024-11-29

### Removed

- `ddos_scrubbing` attribute from IP resources and data sources.
This attribute was made unusable by upstream API changes.

## [1.0.0] - 2024-08-29

### Added

- Server re-installation functionality. Updating server resource attributes `image`, `os_partition_size`, `ssh_key_ids`
  or `user_data` will now trigger a
  server re-install, if `allow_reinstall` is set to true.
- Server resource attributes `hostname` and `name` can now be updated.
- BGP attribute to project resource and data source.

### Changed

- Completed migration from SDKv2 to Framework.
- Terraform 1.8 is now the minimum required version.
- Go 1.21 is now the minimum required version for building the provider.
- IP resource and data source attributes `ptr_record_effective` and `a_record_effective` now keep the actual resource
  state, while `ptr_record` and `a_record` are used for configuration.
- Attribute `project_id` type changed from `string` to `int64` for all resources and data sources.
- Server resource and data source attribute `ssh_key_ids` type changed from `list` to `set`.
- Project resource and data source attribute `team_id` type changed from `string` to `int64`

### Deprecated

- Server resource attribute `ip_addresses_ids` is now deprecated. Use `extra_ip_addresses_ids` instead.

### Removed

- Server data source attribute `server_id` has been removed. Use `id` instead.
- IP data source attributes `ip_id` and `ip_address` have been removed. Use `id` and `address` instead.
- SSH key data source attribute `ssh_key_id` has been removed. Use `id` instead.

[1.0.1]: https://github.com/cherryservers/terraform-provider-cherryservers/compare/v1.0.0...1.0.1
[1.0.0]: https://github.com/cherryservers/terraform-provider-cherryservers/compare/v0.0.6...1.0.0
