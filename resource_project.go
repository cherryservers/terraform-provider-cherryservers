package main

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

	c, err := cherrygo.NewClient()
	if err != nil {
		return err
	}

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

	c, err := cherrygo.NewClient()
	if err != nil {
		return err
	}

	project, _, err := c.Project.List(d.Id())
	if err != nil {
		return err
	}

	d.Set("project_id", project.ID)
	d.Set("name", project.Name)

	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {

	c, err := cherrygo.NewClient()
	if err != nil {
		return err
	}

	projectUpateRequest := cherrygo.UpdateProject{}

	if d.HasChange("name") {
		projectName := d.Get("name").(string)
		projectUpateRequest.Name = projectName

	}

	_, _, err = c.Project.Update(d.Id(), &projectUpateRequest)
	if err != nil {
		return err
	}

	return resourceProjectRead(d, m)
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {

	c, err := cherrygo.NewClient()
	if err != nil {
		return err
	}

	projectDeleteRequest := cherrygo.DeleteProject{ID: d.Id()}

	c.Project.Delete(d.Id(), &projectDeleteRequest)

	d.SetId("")
	return nil
}
