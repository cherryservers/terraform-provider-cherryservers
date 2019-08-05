package cherryservers

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

	c, _ := m.(*Config).Client()

	label := d.Get("name").(string)
	key := d.Get("public_key").(string)

	sshCreateRequest := cherrygo.CreateSSHKey{
		Label: label,
		Key:   key,
	}

	sshkey, _, err := c.client.SSHKey.Create(&sshCreateRequest)
	if err != nil {
		return err
	}

	keyIDString := strconv.Itoa(sshkey.ID)

	d.SetId(keyIDString)
	return resourceSSHKeyRead(d, m)
}

func resourceSSHKeyRead(d *schema.ResourceData, m interface{}) error {

	c, _ := m.(*Config).Client()

	sshkey, _, err := c.client.SSHKey.List(d.Id())
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

	c, _ := m.(*Config).Client()

	sshUpateRequest := cherrygo.UpdateSSHKey{}

	if d.HasChange("name") {
		keyLabel := d.Get("name").(string)
		sshUpateRequest.Label = keyLabel

	}

	if d.HasChange("public_key") {
		key := d.Get("public_key").(string)
		sshUpateRequest.Key = key
	}

	c.client.SSHKey.Update(d.Id(), &sshUpateRequest)

	return resourceSSHKeyRead(d, m)
}

func resourceSSHKeyDelete(d *schema.ResourceData, m interface{}) error {
	c, _ := m.(*Config).Client()

	sshDeleteRequest := cherrygo.DeleteSSHKey{ID: d.Id()}

	c.client.SSHKey.Delete(&sshDeleteRequest)

	d.SetId("")
	return nil
}
