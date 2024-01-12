package cherryservers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/cherryservers/cherrygo/v3"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCherryServersProject() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a CherryServers Project resource. This can be used to create, modify, and delete projects",
		CreateContext: resourceCherryServersProjectCreate,
		ReadContext:   resourceCherryServersProjectRead,
		UpdateContext: resourceCherryServersProjectUpdate,
		DeleteContext: resourceCherryServersProjectDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the project",
			},
			"team_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of team that owns the Project",
			},
		},
	}
}

func resourceCherryServersProjectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).cherrygoClient()
	teamID, _ := strconv.Atoi(d.Get("team_id").(string))
	request := &cherrygo.CreateProject{
		Name: d.Get("name").(string),
	}

	project, _, err := client.Projects.Create(teamID, request)
	if err != nil {
		return diag.Errorf("error creating project: %v", err)
	}

	d.SetId(strconv.Itoa(project.ID))
	log.Printf("[INFO] Project ID: %s", d.Id())

	return resourceCherryServersProjectRead(ctx, d, meta)
}

func resourceCherryServersProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).cherrygoClient()
	projectID, _ := strconv.Atoi(d.Id())
	project, _, err := client.Projects.Get(projectID, nil)
	if err != nil {
		if strings.Contains(err.Error(), "Not Found") {
			tflog.Warn(ctx, fmt.Sprintf("Removing project (%s) because it is gone", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting project: %v", err)
	}

	_ = d.Set("name", project.Name)

	return nil
}

func resourceCherryServersProjectUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).cherrygoClient()
	request := &cherrygo.UpdateProject{}
	projectID, _ := strconv.Atoi(d.Id())

	if d.HasChange("name") {
		label := d.Get("name").(string)
		request.Name = &label
	}

	log.Printf("[INFO] Updating project: %s", d.Id())
	if _, _, err := client.Projects.Update(projectID, request); err != nil {
		return diag.Errorf("Error updating project (%s): %v", d.Id(), err)
	}

	return resourceCherryServersProjectRead(ctx, d, meta)
}

func resourceCherryServersProjectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).cherrygoClient()
	projectID, _ := strconv.Atoi(d.Id())
	log.Printf("[INFO] Deleting project: %s", d.Id())

	if _, err := client.Projects.Delete(projectID); err != nil {
		return diag.Errorf("error deleting Project (%s): %v", d.Id(), err)
	}

	return nil
}
