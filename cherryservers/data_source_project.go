package cherryservers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCherryServersProject() *schema.Resource {
	dsSchema := datasourceSchemaFromResourceSchema(resourceCherryServersProject().Schema)
	dsSchema["project_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
		Optional: true,
		//Required:    true,
		Description: "The ID of the Project",
	}

	return &schema.Resource{
		ReadContext: dataSourceCherryServersProjectRead,
		Schema:      dsSchema,
	}
}

func dataSourceCherryServersProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectID, projectIDExists := d.GetOk("project_id")
	if !projectIDExists {
		return diag.Errorf("project_id is required argument")
	}

	d.SetId(projectID.(string))
	_ = d.Set("project_id", projectID)

	diags := resourceCherryServersProjectRead(ctx, d, meta)
	if diags != nil {
		return append(diags, diag.Errorf("failed to read project")...)
	}

	if d.Id() == "" {
		return diag.Errorf("project (%s) not found", projectID)
	}

	return nil
}
