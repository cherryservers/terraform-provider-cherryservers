# cherryservers_storage (Data Source)

Fetch information about a specific CherryServers storage volume by ID.

## Example Usage

```hcl
data "cherryservers_storage" "example" {
  id = 12345
}

output "storage_vlan_ip" {
  value = data.cherryservers_storage.example.vlan_ip
}
```

## Argument Reference

- `id` - (Required) Storage volume ID to lookup.

## Attributes Reference

- `name` - Name of the storage volume.
- `project_id` - Project ID the storage belongs to.
- `region` - Region slug where the storage is located.
- `size` - Storage size in gigabytes.
- `description` - Storage description.
- `vlan_id` - iSCSI VLAN ID.
- `vlan_ip` - iSCSI VLAN IP address.
- `initiator` - iSCSI initiator name.
- `discovery_ip` - iSCSI discovery IP address.
- `allow_edit_size` - Whether the storage can be resized.
- `unit` - Storage unit (typically "GB").
- `attached_to` - Server ID if storage is attached, or null if unattached.
