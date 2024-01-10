package cherryservers

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCherryServersServer() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := datasourceSchemaFromResourceSchema(resourceCherryServersServer().Schema)

	dsSchema["hostname"] = &schema.Schema{
		Type:          schema.TypeString,
		Optional:      true,
		Description:   "Hostname of the server",
		ConflictsWith: []string{"server_id"},
	}
	dsSchema["server_id"] = &schema.Schema{
		Type:          schema.TypeString,
		Optional:      true,
		Description:   "The ID of the server",
		ConflictsWith: []string{"hostname"},
	}
	dsSchema["project_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The ID of the project, required when searching server by hostname",
	}

	return &schema.Resource{
		ReadContext: dataSourceCherryServersServerRead,
		Schema:      dsSchema,
	}
}

func dataSourceCherryServersServerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).cherrygoClient()

	var serverID string

	if hostname, hostnameExists := d.GetOk("hostname"); hostnameExists {
		prjID, projectIDExists := d.GetOk("project_id")
		if !projectIDExists {
			return diag.Errorf("project_id is required argument when importing server by hostname")
		}

		projectID, _ := strconv.Atoi(prjID.(string))
		srvID, err := ServerHostnameToID(hostname.(string), projectID, client.Servers)
		if err != nil {
			return diag.Errorf("%v", err)
		}
		serverID = strconv.Itoa(srvID)
	} else if srvID, srvIDExists := d.GetOk("server_id"); srvIDExists {
		serverID = srvID.(string)
	}

	d.SetId(serverID)

	diags := resourceCherryServersServerRead(ctx, d, meta)
	if diags != nil {
		return append(diags, diag.Errorf("failed to read server")...)
	}

	if d.Id() == "" {
		return diag.Errorf("server (%s) not found", serverID)
	}

	return nil
}
