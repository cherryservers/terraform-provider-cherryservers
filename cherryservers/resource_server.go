package cherryservers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/cherryservers/cherrygo/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceCherryServersServer() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Cherry Servers server resource. This can be used to create, read, modify, and delete servers on your Cherry Servers account.",
		CreateContext: resourceCherryServersServerCreate,
		ReadContext:   resourceCherryServersServerRead,
		UpdateContext: resourceCherryServersServerUpdate,
		DeleteContext: resourceCherryServersServerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of the server",
			},
			"plan": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "Slug of the plan. Example: e5_1620v4. [See List Plans](https://api.cherryservers.com/doc/#tag/Plans/operation/get-plans)",
			},
			"project_id": projectIDSchema(),
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Slug of the region. Example: eu_nord_1. [See List Regions](https://api.cherryservers.com/doc/#tag/Regions/operation/get-regions)",
			},
			"hostname": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Hostname of the server",
			},
			"image": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Slug of the operating system. Example: ubuntu_22_04. [See List Images](https://api.cherryservers.com/doc/#tag/Images/operation/get-plan-images)",
			},
			"ssh_key_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of the SSH key IDs allowed to SSH to the servers",
			},
			"ip_addresses_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of the IP addresses IDs to be embed in to the Server",
			},
			"user_data": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Base64 encoded User-Data blob. It should be either bash or cloud-config script",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Default:     nil,
				Description: "Key/value metadata for server tagging",
			},
			"spot_instance": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If True, provisions the server as a spot instance",
			},
			"os_partition_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "OS partition size in GB",
			},

			"power_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The power state of the server, such as 'Powered off'",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The state of the server, such as 'pending'",
			},
			"ip_addresses": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        ServerIPSchema(),
				Description: "IP addresses attached to the server",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
		},
	}
}

func resourceCherryServersServerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).cherrygoClient()
	projectID, _ := strconv.Atoi(d.Get("project_id").(string))
	request := &cherrygo.CreateServer{
		ProjectID: projectID,
		Plan:      d.Get("plan").(string),
		Region:    d.Get("region").(string),
		Image:     d.Get("image").(string),
	}

	hostname, hostnameOk := d.GetOk("hostname")
	if hostnameOk {
		request.Hostname = hostname.(string)
	}

	if sshKeyIDs, sshKeyOK := d.GetOk("ssh_key_ids"); sshKeyOK {
		for _, v := range sshKeyIDs.([]interface{}) {
			request.SSHKeys = append(request.SSHKeys, v.(string))
		}
	}

	if ipAddressesIDs, ipAddressesOK := d.GetOk("ip_addresses_ids"); ipAddressesOK {
		for _, v := range ipAddressesIDs.([]interface{}) {
			request.IPAddresses = append(request.IPAddresses, v.(string))
		}
	}

	tagsMap := make(map[string]string)
	if tags, tagsOK := d.GetOk("tags"); tagsOK {
		for k, v := range tags.(map[string]interface{}) {
			tagsMap[k] = v.(string)
		}
		request.Tags = &tagsMap
	}

	spotInstance, spotOK := d.GetOk("spot_instance")
	if spotOK {
		request.SpotInstance = spotInstance.(bool)
	}

	userData, userDataExist := d.GetOk("user_data")
	if userDataExist {
		uData := userData.(string)
		err := IsBase64(uData)
		if err {
			return diag.Errorf("error creating server, user_data property must be base64 encoded value")
		}
		request.UserData = uData
	}

	osPartitionSize, partitionExist := d.GetOk("os_partition_size")
	if partitionExist {
		request.OSPartitionSize = osPartitionSize.(int)
	}

	log.Printf("[INFO] Creating server")
	var server *cherrygo.Server = nil

	// allow for retries on creation to handle retryable platform errors
	retryErr := retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate)-time.Minute, func() *retry.RetryError {
		s, _, err := client.Servers.Create(request)
		server = &s

		if err == nil {
			return nil
		}

		return retry.NonRetryableError(err)
	})

	if retryErr != nil {
		return diag.Errorf("error creating server: %v", retryErr)
	}

	d.SetId(strconv.Itoa(server.ID))

	if _, err := waitForServerAvailable(ctx, d, "active", []string{"pending", "provisioning"}, "state", meta); err != nil {
		return diag.Errorf("error while waiting for Server %s to be completed: %s", d.Id(), err)
	}

	return resourceCherryServersServerRead(ctx, d, meta)
}

func resourceCherryServersServerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).cherrygoClient()
	serverId, _ := strconv.Atoi(d.Id())
	server, _, err := client.Servers.Get(serverId, nil)
	if err != nil {
		if strings.Contains(err.Error(), "Server object was not found") {
			log.Printf("[WARN] Removing server (%s) because it is gone", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting server (%s): %v", d.Id(), err)
	}

	if len(server.IPAddresses) == 0 || server.State == "terminating" {
		log.Printf("[WARN] Removing server (%s), it might be in terminating mode or has no ip adresses", d.Id())
		d.SetId("")
		return nil
	}

	// var primaryIP, privateIP string
	// if len(server.IPAddresses) > 0 {
	// 	for _, ip := range server.IPAddresses {
	// 		if ip.Type == "primary-ip" {
	// 			primaryIP = ip.Address
	// 		}
	// 		if ip.Type == "private-ip" {
	// 			privateIP = ip.Address
	// 		}
	// 	}
	// }

	powerState, _, err := client.Servers.PowerState(serverId)
	if err != nil {
		log.Printf("Error while getting power state: %v", err)
	}

	_ = d.Set("name", server.Name)
	_ = d.Set("plan", server.Plan.Slug)
	_ = d.Set("hostname", server.Hostname)
	_ = d.Set("region", server.Region.Slug)
	_ = d.Set("power_state", powerState.Power)
	_ = d.Set("state", server.State)
	_ = d.Set("tags", server.Tags)
	_ = d.Set("spot_instance", server.SpotInstance)
	_ = d.Set("ip_addresses", flattenServerIPs(server.IPAddresses))

	var sshKeyIDs []string
	for _, sshKey := range server.SSHKeys {
		sshKeyID := strconv.Itoa(sshKey.ID)
		sshKeyIDs = append(sshKeyIDs, sshKeyID)
	}
	_ = d.Set("ssh_key_ids", sshKeyIDs)

	//d.Set("os_partition_size", unknown)

	return nil
}

func resourceCherryServersServerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).cherrygoClient()
	request := cherrygo.UpdateServer{}
	serverId, _ := strconv.Atoi(d.Id())

	if d.HasChange("tags") {
		tagsMap := make(map[string]string)
		if tags, tagsOK := d.GetOk("tags"); tagsOK {
			for k, v := range tags.(map[string]interface{}) {
				tagsMap[k] = v.(string)
			}
			request.Tags = &tagsMap
		}

		_, _, err := client.Servers.Update(serverId, &request)
		if err != nil {
			return diag.Errorf("error updating server (%s): %v", d.Id(), err)
		}
	}

	return resourceCherryServersServerRead(ctx, d, meta)
}

func resourceCherryServersServerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).cherrygoClient()
	serverId, _ := strconv.Atoi(d.Id())

	log.Printf("[INFO] Deleting server (%s)", d.Id())

	if _, _, err := client.Servers.Delete(serverId); err != nil {
		return diag.Errorf("error destroying server %s : %v", d.Id(), err)
	}
	return nil
}

func waitForServerAvailable(ctx context.Context, d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}) (interface{}, error) {
	log.Printf(
		"[INFO] Waiting for Server (%s) to have %s of %s",
		d.Id(), attribute, target)

	stateConf := &resource.StateChangeConf{ // nolint:all
		Pending:        pending,
		Target:         []string{target},
		Refresh:        newServerStateRefresh(ctx, d, meta, attribute),
		Timeout:        60 * time.Minute,
		Delay:          10 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 60,
	}

	return stateConf.WaitForStateContext(ctx)
}

func newServerStateRefresh(ctx context.Context, d *schema.ResourceData, meta interface{}, attr string) resource.StateRefreshFunc { // nolint:all
	client := meta.(*Client).cherrygoClient()
	return func() (interface{}, string, error) {

		log.Printf("[INFO] Creating Server")
		serverId, _ := strconv.Atoi(d.Id())
		options := cherrygo.GetOptions{Fields: []string{attr}}
		server, _, err := client.Servers.Get(serverId, &options)
		if err != nil {
			return nil, "", fmt.Errorf("error retrieving Server %s : %s", d.Id(), err)
		}

		if attr == "state" {
			log.Printf("[INFO] Server Status is %s", server.State)
			return server, server.State, nil
		} else if attr == "status" {
			//TODO: check if server is deployed (status field)
			return nil, "", nil
		} else {
			return nil, "", nil
		}
	}
}
