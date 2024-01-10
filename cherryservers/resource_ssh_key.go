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

func resourceCherryServersSSHKey() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a CherryServers SSH Key resource. This can be used to create, and delete SSH Keys associated with your Cherry account",
		CreateContext: resourceCherryServersSSHKeyCreate,
		ReadContext:   resourceCherryServersSSHKeyRead,
		UpdateContext: resourceCherryServersSSHKeyUpdate,
		DeleteContext: resourceCherryServersSSHKeyDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the SSH key",
			},
			"public_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The public SSH key",
			},
			"fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The fingerprint of your SSH Public key",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date when this Key was added",
			},
			"updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date when this Key was modified",
			},
		},
	}
}

func resourceCherryServersSSHKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).cherrygoClient()
	request := &cherrygo.CreateSSHKey{
		Label: d.Get("name").(string),
		Key:   strings.TrimSpace(d.Get("public_key").(string)),
	}

	key, _, err := client.SSHKeys.Create(request)
	if err != nil {
		return diag.Errorf("error creating SSH key: %v", err)
	}

	d.SetId(strconv.Itoa(key.ID))
	log.Printf("[INFO] SSH Key ID: %s", d.Id())

	return resourceCherryServersSSHKeyRead(ctx, d, meta)
}

func resourceCherryServersSSHKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).cherrygoClient()

	sskKeyID, _ := strconv.Atoi(d.Id())
	key, _, err := client.SSHKeys.Get(sskKeyID, nil)
	if err != nil {
		if strings.Contains(err.Error(), "Not Found") {
			tflog.Warn(ctx, fmt.Sprintf("Removing ssh key (%s) because it is gone", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting SSH keys: %v", err)
	}

	if err := d.Set("name", key.Label); err != nil {
		return diag.Errorf("unable to set resource ssh_key `name` read value: %v", err)
	}
	if err := d.Set("public_key", key.Key); err != nil {
		return diag.Errorf("unable to set resource ssh_key `public_key` read value: %v", err)
	}
	if err := d.Set("fingerprint", key.Fingerprint); err != nil {
		return diag.Errorf("unable to set resource ssh_key `fingerprint` read value: %v", err)
	}
	if err := d.Set("created", key.Created); err != nil {
		return diag.Errorf("unable to set resource ssh_key `created` read value: %v", err)
	}
	if err := d.Set("updated", key.Updated); err != nil {
		return diag.Errorf("unable to set resource ssh_key `updated` read value: %v", err)
	}

	return nil
}

func resourceCherryServersSSHKeyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).cherrygoClient()
	request := &cherrygo.UpdateSSHKey{}
	sskKeyID, _ := strconv.Atoi(d.Id())

	if d.HasChange("name") {
		label := d.Get("name").(string)
		request.Label = &label
	}

	if d.HasChange("public_key") {
		key := d.Get("ssh_key").(string)
		request.Key = &key
	}

	log.Printf("[INFO] Updating SSH Key: %s", d.Id())
	if _, _, err := client.SSHKeys.Update(sskKeyID, request); err != nil {
		return diag.Errorf("error updating SSH key (%s): %v", d.Id(), err)
	}

	return resourceCherryServersSSHKeyRead(ctx, d, meta)
}

func resourceCherryServersSSHKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).cherrygoClient()
	sskKeyID, _ := strconv.Atoi(d.Id())
	log.Printf("[INFO] Deleting SSH Key: %s", d.Id())

	if _, _, err := client.SSHKeys.Delete(sskKeyID); err != nil {
		return diag.Errorf("error deleting SSH key (%s): %v", d.Id(), err)
	}

	return nil
}
