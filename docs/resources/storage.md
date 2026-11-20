# cherryservers_storage

Provides a CherryServers Storage (EBS-like) resource. This can be used to create, modify, and delete elastic block storage volumes that can be attached to servers.

## Example Usage

### Basic Storage Volume

```hcl
data "cherryservers_project" "main" {
}

resource "cherryservers_storage" "example" {
  project_id  = data.cherryservers_project.main.id
  region      = "LT-Siauliai"
  size        = 100
  description = "My storage volume"
}
```

### With Server Attachment

```hcl
data "cherryservers_project" "main" {
}

resource "cherryservers_server" "web" {
  project_id = data.cherryservers_project.main.id
  plan       = "e5-1620v4"
  region     = "LT-Siauliai"
  hostname   = "web-server"
  image      = "ubuntu_22_04_x64"
}

resource "cherryservers_storage" "data" {
  project_id  = data.cherryservers_project.main.id
  region      = "LT-Siauliai"
  size        = 250
  attached_to = cherryservers_server.web.id
}
```

### Resize Storage (Creates New Volume)

```hcl
resource "cherryservers_storage" "resized" {
  project_id  = data.cherryservers_project.main.id
  region      = "LT-Siauliai"
  size        = 500  # Increased from 100 - creates new volume
  description = "Resized storage"
}
```

## Argument Reference

- `project_id` - (Required) CherryServers project ID. Cannot be changed after creation.
- `region` - (Required) Region slug where storage will be created (e.g., `LT-Siauliai`). Cannot be changed after creation.
- `size` - (Required) Storage size in gigabytes. Can only be increased. **Warning:** Changing the size creates a new storage volume with a new ID.
- `description` - (Optional) Storage description. Defaults to empty string.
- `attached_to` - (Optional) Server ID to attach the storage to. Can be changed to attach/detach or move between servers.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

- `id` - Storage volume ID.
- `name` - Auto-generated name of the storage volume.
- `vlan_id` - iSCSI VLAN ID for connecting to the storage.
- `vlan_ip` - iSCSI VLAN IP address.
- `initiator` - iSCSI initiator name (IQN).
- `discovery_ip` - iSCSI discovery IP address.
- `allow_edit_size` - Whether this storage volume can be resized (always true for EBS volumes).
- `unit` - Storage unit (typically "GB").
- `created_at` - Timestamp when the storage was created.

## Import

Storage volumes can be imported using their ID:

```bash
terraform import cherryservers_storage.example 12345
```

## Notes

### Storage Resizing

- Storage can only be resized **upward** (increased in size).
- **Important:** Resizing a storage volume creates a **new storage ID**. This means:
  - The old storage volume is destroyed
  - A new storage volume is created with the new size
  - iSCSI connection details (VLAN, IP, initiator) change
  - Any attached servers must reconnect with the new iSCSI details
- To resize, modify the `size` attribute and apply. Terraform will replace the resource.

### Storage Attachment

- Storage can exist unattached (when `attached_to` is not specified).
- To attach or move storage between servers, change the `attached_to` value.
- To detach storage, remove the `attached_to` attribute and apply.
- Storage must be detached before deletion if currently attached.

### iSCSI Connection

Once storage is created, you'll have iSCSI connection details available:

- **VLAN ID**: `vlan_id`
- **VLAN IP**: `vlan_ip`
- **Initiator**: `initiator`
- **Discovery IP**: `discovery_ip`

Use these details to connect the storage to servers via iSCSI protocol.
