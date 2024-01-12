package cherryservers

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/cherryservers/cherrygo/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCherryServersIP() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a CherryServers IP resource. This can be used to create, modify, and delete IP addresses",
		CreateContext: resourceCherryServersIPCreate,
		ReadContext:   resourceCherryServersIPRead,
		UpdateContext: resourceCherryServersIPUpdate,
		DeleteContext: resourceCherryServersIPDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"project_id": projectIDSchema(),
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Slug of the region. Example: eu_nord_1 [See List Regions](https://api.cherryservers.com/doc/#tag/Regions/operation/get-regions)",
			},
			"target_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the server to assign the created IP to",
			},
			"target_hostname": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The hostname of the server to assign the created IP to",
			},
			"target_ip_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Subnet or primary-ip type IP ID to route the created IP to",
			},
			"ddos_scrubbing": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If True, DDOS scrubbing protection will be applied in real-time",
			},
			"a_record": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Relative DNS name for the IP address. Resulting FQDN will be '<relative-dns-name>.cloud.cherryservers.net' and must be globally unique",
			},
			"ptr_record": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Reverse DNS name for the IP address",
			},
			"address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP address in canonical format used in the reverse DNS record",
			},
			"address_family": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "IP address family IPv4 or IPv6",
			},
			"cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CIDR block of the IP",
			},
			"gateway": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The gateway IP address",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of IP address",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Default:     nil,
				Description: "Key/value metadata for server tagging",
			},
		},
	}
}

func resourceCherryServersIPCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).cherrygoClient()
	request := &cherrygo.CreateIPAddress{
		Region:        d.Get("region").(string),
		DDoSScrubbing: d.Get("ddos_scrubbing").(bool),
	}
	projectID, _ := strconv.Atoi(d.Get("project_id").(string))

	if ptr, ptrOK := d.GetOk("ptr_record"); ptrOK {
		request.PtrRecord = ptr.(string)
	}

	if arecord, aOK := d.GetOk("a_record"); aOK {
		request.ARecord = arecord.(string)
	}

	tagsMap := make(map[string]string)
	if tags, tagsOK := d.GetOk("tags"); tagsOK {
		for k, v := range tags.(map[string]interface{}) {
			tagsMap[k] = v.(string)
		}
		request.Tags = &tagsMap
	}

	if targetHostname, hostnameOK := d.GetOk("target_hostname"); hostnameOK {
		srvID, err := ServerHostnameToID(targetHostname.(string), projectID, client.Servers)
		if err != nil {
			return diag.Errorf("%v", err)
		}
		request.TargetedTo = strconv.Itoa(srvID)
	} else if targetID, idOK := d.GetOk("target_id"); idOK {
		request.TargetedTo = strconv.Itoa(targetID.(int))
	} else if targetIPId, ipIdOK := d.GetOk("target_ip_id"); ipIdOK {
		request.TargetedTo = targetIPId.(string)
	}

	ipAddress, _, err := client.IPAddresses.Create(projectID, request)
	if err != nil {
		return diag.Errorf("error creating IP address: %v", err)
	}

	d.SetId(ipAddress.ID)
	log.Printf("[INFO] IP address ID: %s", d.Id())

	return resourceCherryServersIPRead(ctx, d, meta)
}

func resourceCherryServersIPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).cherrygoClient()
	ipAddress, _, err := client.IPAddresses.Get(d.Id(), nil)
	if err != nil {
		if strings.Contains(err.Error(), "IP object was not found") {
			log.Printf("[WARN] Removing IP address (%s) because it is gone", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting IP address (%s): %v", d.Id(), err)
	}

	_ = d.Set("address", ipAddress.Address)
	_ = d.Set("cidr", ipAddress.Cidr)
	_ = d.Set("gateway", ipAddress.Gateway)
	_ = d.Set("ptr_record", ipAddress.PtrRecord)
	_ = d.Set("a_record", ipAddress.ARecord)
	_ = d.Set("address_family", ipAddress.AddressFamily)
	_ = d.Set("ddos_scrubbing", ipAddress.DDoSScrubbing)
	_ = d.Set("type", ipAddress.Type)
	_ = d.Set("tags", ipAddress.Tags)
	_ = d.Set("region", ipAddress.Region.Slug)
	_ = d.Set("target_id", ipAddress.TargetedTo.ID)
	_ = d.Set("target_hostname", ipAddress.TargetedTo.Hostname)

	return nil
}

func resourceCherryServersIPUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).cherrygoClient()
	request := cherrygo.UpdateIPAddress{}
	projectID := d.Get("project_id").(int)

	if d.HasChange("tags") {
		tagsMap := make(map[string]string)
		if tags, tagsOK := d.GetOk("tags"); tagsOK {
			for k, v := range tags.(map[string]interface{}) {
				tagsMap[k] = v.(string)
			}
			request.Tags = &tagsMap
		}
	}
	if d.HasChange("ptr_record") {
		request.PtrRecord = d.Get("ptr_record").(string)
	}
	if d.HasChange("a_record") {
		request.PtrRecord = d.Get("a_record").(string)
	}

	if targetHostname, hostnameOK := d.GetOk("target_hostname"); hostnameOK && d.HasChange("target_hostname") {
		srvID, err := ServerHostnameToID(targetHostname.(string), projectID, client.Servers)
		if err != nil {
			return diag.Errorf("%v", err)
		}
		request.TargetedTo = strconv.Itoa(srvID)
	} else if targetID, idOK := d.GetOk("target_id"); idOK && d.HasChange("target_id") {
		request.TargetedTo = strconv.Itoa(targetID.(int))
	} else if targetIPId, ipIdOK := d.GetOk("target_ip_id"); ipIdOK && d.HasChange("target_id") {
		request.TargetedTo = targetIPId.(string)
	}

	log.Printf("[INFO] Updating IP address: %s", d.Id())
	if _, _, err := client.IPAddresses.Update(d.Id(), &request); err != nil {
		return diag.Errorf("error updating IP address (%s): %v", d.Id(), err)
	}

	return resourceCherryServersProjectRead(ctx, d, meta)
}

func resourceCherryServersIPDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).cherrygoClient()

	log.Printf("[INFO] Deleting IP address: %s", d.Id())
	ipAddress, _, err := client.IPAddresses.Get(d.Id(), nil)
	if err != nil {
		d.SetId("")
		return nil
	}

	if ipAddress.TargetedTo.ID != 0 {
		_, err = client.IPAddresses.Unassign(d.Id())
		if err != nil {
			return diag.Errorf("error failed to unassign IP address before deleting: %v", err)
		}
	}

	if _, err := client.IPAddresses.Remove(d.Id()); err != nil {
		return diag.Errorf("error deleting IP address (%s): %v", d.Id(), err)
	}
	return nil
}
