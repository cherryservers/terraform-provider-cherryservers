package cherryservers

import (
	"strconv"

	"github.com/cherryservers/cherrygo"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,

		Schema: map[string]*schema.Schema{
			"team_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {

	c := m.(*cherrygo.Client)

	projectName := d.Get("name").(string)
	teamID := d.Get("team_id").(string)

	addProjectRequest := cherrygo.CreateProject{
		Name: projectName,
	}

	intTeamID, err := strconv.Atoi(teamID)

	project, _, err := c.Project.Create(intTeamID, &addProjectRequest)
	if err != nil {
		return err
	}

	projectIDString := strconv.Itoa(project.ID)

	d.SetId(projectIDString)
	return resourceProjectRead(d, m)
}

func resourceProjectRead(d *schema.ResourceData, m interface{}) error {

	c := m.(*cherrygo.Client)

	project, _, err := c.Project.List(d.Id())
	if err != nil {
		return err
	}

	d.Set("project_id", project.ID)
	d.Set("name", project.Name)

	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {

	c := m.(*cherrygo.Client)

	projectUpateRequest := cherrygo.UpdateProject{}

	if d.HasChange("name") {
		projectName := d.Get("name").(string)
		projectUpateRequest.Name = projectName

	}

	c.Project.Update(d.Id(), &projectUpateRequest)

	return resourceProjectRead(d, m)
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {

	c := m.(*cherrygo.Client)

	projectDeleteRequest := cherrygo.DeleteProject{ID: d.Id()}

	c.Project.Delete(d.Id(), &projectDeleteRequest)

	d.SetId("")
	return nil
}
