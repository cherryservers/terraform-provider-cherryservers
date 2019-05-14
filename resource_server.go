package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/cherryservers/cherrygo"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"hostname": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"image": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"plan_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"power_state": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"primary_ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssh_keys_ids": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ip_addresses_ids": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceServerCreate(d *schema.ResourceData, m interface{}) error {

	c, err := cherrygo.NewClient()
	if err != nil {
		return err
	}

	projectID := d.Get("project_id").(string)
	hostname := d.Get("hostname").(string)
	image := d.Get("image").(string)
	region := d.Get("region").(string)
	planID := d.Get("plan_id").(string)
	sshKeys1 := d.Get("ssh_keys_ids").([]interface{})
	ipAddresses := d.Get("ip_addresses_ids").([]interface{})

	//var sshKeysArr []string
	var sshKeysArr = make([]string, 0)
	for _, v := range sshKeys1 {
		if v == nil {
			continue
		}
		sshKeysArr = append(sshKeysArr, v.(string))

	}

	// below is nil default value
	//var ipAddressesArr []string
	var ipAddressesArr = make([]string, 0)
	for _, v := range ipAddresses {
		if v == nil {
			continue
		}
		ipAddressesArr = append(ipAddressesArr, v.(string))
	}

	addServerRequest := cherrygo.CreateServer{
		ProjectID:   projectID,
		Hostname:    hostname,
		Image:       image,
		Region:      region,
		SSHKeys:     sshKeysArr,
		IPAddresses: ipAddressesArr,
		PlanID:      planID,
	}

	server, _, err := c.Server.Create(projectID, &addServerRequest)
	if err != nil {
		log.Fatal("Error while creating new server: ", err)
	}

	serverID := strconv.Itoa(server.ID)

	d.SetId(serverID)

	err = waitForNetwork(d, m)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	return resourceServerRead(d, m)
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {

	c, err := cherrygo.NewClient()
	if err != nil {
		return err
	}

	server, _, err := c.Server.List(d.Id())
	if err != nil {
		log.Fatalf("Error while listing server: %v", err)
	}

	// If the server doesn't have any IPs or it's in terminating
	// state, assume it's gone.
	if len(server.IPAddresses) == 0 || server.State == "terminating" {
		d.SetId("")
		return nil
	}

	var primaryIP, privateIP string

	if len(server.IPAddresses) > 0 {
		for _, ip := range server.IPAddresses {
			if ip.Type == "primary-ip" {
				primaryIP = ip.Address
			}
			if ip.Type == "private-ip" {
				privateIP = ip.Address
			}
		}
	}

	srvPower, _, err := c.Server.PowerState(d.Id())
	if err != nil {
		log.Fatalf("Error while getting power sstate: %v", err)
	}

	d.Set("name", server.Name)
	d.Set("hostname", server.Hostname)
	d.Set("image", server.Image)
	d.Set("price", server.Pricing.Price)
	d.Set("region", server.Region.Name)
	d.Set("power_state", srvPower.Power)
	d.Set("state", server.State)
	d.Set("primary_ip", primaryIP)
	d.Set("private_ip", privateIP)

	return nil
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceServerRead(d, m)
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {

	c, err := cherrygo.NewClient()
	if err != nil {
		return err
	}

	serverDeleteRequest := cherrygo.DeleteServer{ID: d.Id()}

	c.Server.Delete(&serverDeleteRequest)

	d.SetId("")
	return nil
}

func waitForNetwork(d *schema.ResourceData, m interface{}) error {

	c, err := cherrygo.NewClient()
	if err != nil {
		return err
	}

	for i := 1; i < 10; i++ {

		time.Sleep(time.Second * 6)

		server, _, err := c.Server.List(d.Id())
		if err != nil {
			log.Fatalf("Error while listing server: %v", err)
		}

		for _, ip := range server.IPAddresses {
			if ip.Type == "primary-ip" {
				if ip.Address != "" {
					return nil
				}
			}
		}
	}

	err = fmt.Errorf("timed out waiting for active device: %v", d.Id())

	return err
}
