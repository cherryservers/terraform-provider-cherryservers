package cherryservers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCherryServersSSHKey() *schema.Resource {
	dsSchema := datasourceSchemaFromResourceSchema(resourceCherryServersSSHKey().Schema)
	addOptionalFieldsToSchema(dsSchema, "name")

	dsSchema["name"].ConflictsWith = []string{"ssh_key_id"}
	dsSchema["ssh_key_id"] = &schema.Schema{
		Type:          schema.TypeString,
		Optional:      true,
		Description:   "The ID of the SSH key",
		ConflictsWith: []string{"name"},
	}
	dsSchema["project_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The ID of the project",
	}

	return &schema.Resource{
		ReadContext: dataSourceCherrySSHKeyRead,
		Schema:      dsSchema,
	}
}

func dataSourceCherrySSHKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).cherrygoClient()

	sshKeyID, sshKeyIDExists := d.GetOk("ssh_key_id")
	if !sshKeyIDExists {
		projectID, _ := strconv.Atoi(d.Get("project_id").(string))
		sshKeys, _, err := client.Projects.ListSSHKeys(projectID, nil)
		if err != nil {
			return diag.FromErr(err)
		}
		for _, sshKey := range sshKeys {
			if sshKey.Label == d.Get("name").(string) {
				if sshKeyID != "" {
					return diag.FromErr(fmt.Errorf("more than 1 SSH Key found with the same name %s", d.Get("name")))
				}
				sshKeyID = strconv.Itoa(sshKey.ID)
			}
		}
		if sshKeyID == "" {
			return diag.FromErr(fmt.Errorf("no SSH Key found with the name %s", d.Get("name")))
		}
	}

	d.SetId(sshKeyID.(string))

	diags := resourceCherryServersSSHKeyRead(ctx, d, meta)
	if diags != nil {
		return append(diags, diag.Errorf("failed to read SSH key")...)
	}

	if d.Id() == "" {
		return diag.Errorf("SSH key (%s) not found", sshKeyID)
	}

	return nil
}
