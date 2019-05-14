package main

import (
	"strconv"

	"github.com/cherryservers/cherrygo"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSSHKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceSSHKeyCreate,
		Read:   resourceSSHKeyRead,
		Update: resourceSSHKeyUpdate,
		Delete: resourceSSHKeyDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"public_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"fingerprint": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"created": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSSHKeyCreate(d *schema.ResourceData, m interface{}) error {

	c, err := cherrygo.NewClient()
	if err != nil {
		return err
	}

	label := d.Get("name").(string)
	key := d.Get("public_key").(string)

	sshCreateRequest := cherrygo.CreateSSHKey{
		Label: label,
		Key:   key,
	}

	sshkey, _, err := c.SSHKey.Create(&sshCreateRequest)
	if err != nil {
		return err
	}

	keyIDString := strconv.Itoa(sshkey.ID)

	d.SetId(keyIDString)
	return resourceSSHKeyRead(d, m)
}

func resourceSSHKeyRead(d *schema.ResourceData, m interface{}) error {

	c, err := cherrygo.NewClient()
	if err != nil {
		return err
	}

	sshkey, _, err := c.SSHKey.List(d.Id())
	if err != nil {
		return err
	}

	d.Set("id", sshkey.ID)
	d.Set("name", sshkey.Label)
	d.Set("fingerprint", sshkey.Fingerprint)
	d.Set("created", sshkey.Created)
	d.Set("updated", sshkey.Updated)

	return nil
}

func resourceSSHKeyUpdate(d *schema.ResourceData, m interface{}) error {

	c, err := cherrygo.NewClient()
	if err != nil {
		return err
	}

	sshUpateRequest := cherrygo.UpdateSSHKey{}

	if d.HasChange("name") {
		keyLabel := d.Get("name").(string)
		sshUpateRequest.Label = keyLabel

	}

	if d.HasChange("public_key") {
		key := d.Get("public_key").(string)
		sshUpateRequest.Key = key
	}

	_, _, err = c.SSHKey.Update(d.Id(), &sshUpateRequest)
	if err != nil {
		return err
	}

	return resourceSSHKeyRead(d, m)
}

func resourceSSHKeyDelete(d *schema.ResourceData, m interface{}) error {

	c, err := cherrygo.NewClient()
	if err != nil {
		return err
	}

	sshDeleteRequest := cherrygo.DeleteSSHKey{ID: d.Id()}

	c.SSHKey.Delete(&sshDeleteRequest)

	d.SetId("")
	return nil
}
