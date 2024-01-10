package cherryservers

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceCherryServersIP() *schema.Resource {
	dsSchema := datasourceSchemaFromResourceSchema(resourceCherryServersIP().Schema)

	dsSchema["ip_address"] = &schema.Schema{
		Type:          schema.TypeString,
		Optional:      true,
		Description:   "The IPv4 address",
		ValidateFunc:  validation.IsIPv4Address,
		ConflictsWith: []string{"ip_id"},
	}
	dsSchema["ip_id"] = &schema.Schema{
		Type:          schema.TypeString,
		Optional:      true,
		Description:   "The IP address in canonical format used in the reverse DNS record",
		ConflictsWith: []string{"ip_address"},
	}

	dsSchema["project_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "ID of the project you are working on",
		Optional:    true,
	}

	return &schema.Resource{
		ReadContext: dataSourceCherryServersIPRead,
		Schema:      dsSchema,
	}
}

func dataSourceCherryServersIPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).cherrygoClient()

	ipID, ipIDExists := d.GetOk("ip_id")
	if !ipIDExists {
		prjID, projectIDExists := d.GetOk("project_id")
		if !projectIDExists {
			return diag.Errorf("project_id is required argument when importing IP by address")
		}

		projectID, _ := strconv.Atoi(prjID.(string))
		ipAddresses, _, err := client.IPAddresses.List(projectID, nil)
		if err != nil {
			return diag.FromErr(err)
		}

		for _, ip := range ipAddresses {
			if ip.Address == d.Get("ip_address").(string) {
				ipID = ip.ID
			}
		}
		if ipID == "" {
			return diag.Errorf("no IP found by given address %s", d.Get("ip_address"))
		}
	}

	d.SetId(ipID.(string))
	err := d.Set("ip_id", ipID)
	if err != nil {
		return diag.FromErr(err)
	}

	diags := resourceCherryServersIPRead(ctx, d, meta)
	if diags != nil {
		return append(diags, diag.Errorf("failed to read IP address")...)
	}

	if d.Id() == "" {
		return diag.Errorf("IP address (%s) not found", ipID)
	}

	return nil
}
