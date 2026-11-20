# cherryservers_storages (Data Source)

List all CherryServers storage volumes in a project.

## Example Usage

```hcl
data "cherryservers_project" "main" {
}

data "cherryservers_storages" "all" {
  project_id = data.cherryservers_project.main.id
}

output "storage_ids" {
  value = [for storage in data.cherryservers_storages.all.storages : storage.id]
}
```

## Argument Reference

- `project_id` - (Required) Project ID to list storage volumes for.

## Attributes Reference

- `storages` - List of storage volumes. Each storage block contains:
  - `id` - Storage volume ID.
  - `name` - Name of the storage volume.
  - `region` - Region slug where the storage is located.
  - `size` - Storage size in gigabytes.
  - `description` - Storage description.
  - `vlan_id` - iSCSI VLAN ID.
  - `vlan_ip` - iSCSI VLAN IP address.
  - `initiator` - iSCSI initiator name.
  - `discovery_ip` - iSCSI discovery IP address.
  - `allow_edit_size` - Whether the storage can be resized.
  - `unit` - Storage unit (typically "GB").
  - `attached_to` - Server ID if attached, or null if unattached.
