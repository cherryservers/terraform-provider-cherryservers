package cherryservers

import (
	"fmt"
	"log"

	"github.com/cherryservers/cherrygo"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceIP() *schema.Resource {
	return &schema.Resource{
		Create: resourceIPCreate,
		Read:   resourceIPRead,
		Update: resourceIPUpdate,
		Delete: resourceIPDelete,

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"a_record": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"ptr_record": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"routed_to": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"routed_to_hostname": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"routed_to_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"cidr": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"gateway": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func getIDForServerIP(d *schema.ResourceData, m interface{}) (string, error) {

	projectID := d.Get("project_id").(string)
	routeToHostname := d.Get("routed_to_hostname").(string)
	routeToIP := d.Get("routed_to_ip").(string)

	c := m.(*cherrygo.Client)

	servers, _, err := c.Servers.List(projectID)
	if err != nil {
		log.Fatalf("Error while listing server: %v", err)
	}

	var routeTo string

	switch {
	case routeToHostname != "":
		for _, srv := range servers {
			if srv.Hostname == routeToHostname {
				if len(srv.IPAddresses) > 0 {
					for _, i := range srv.IPAddresses {
						if i.Type == "primary-ip" {
							routeTo = i.ID
						}
					}
				}
				fmt.Printf("PANICMODE: %v -> NUMBER: %v", srv.IPAddresses, routeTo)
			}
		}
	case routeToIP != "":
		for _, srv := range servers {
			if len(srv.IPAddresses) > 0 {
				for _, i := range srv.IPAddresses {
					if i.Type == "primary-ip" {
						if i.Address == routeToIP {
							routeTo = i.ID
						}
					}
				}
			}
		}
	}

	return routeTo, err
}

func resourceIPCreate(d *schema.ResourceData, m interface{}) error {

	c, _ := m.(*Config).Client()

	projectID := d.Get("project_id").(string)
	aRecord := d.Get("a_record").(string)
	ptrRecord := d.Get("ptr_record").(string)
	region := d.Get("region").(string)
	routedTo := d.Get("routed_to").(string)

	addIPRequest := cherrygo.CreateIPAddress{
		ARecord:   aRecord,
		PtrRecord: ptrRecord,
		Region:    region,
		RoutedTo:  routedTo,
	}

	ipAddress, _, err := c.client.IPAddress.Create(projectID, &addIPRequest)
	if err != nil {
		return err
	}

	//time.Sleep(5 * time.Second)

	d.SetId(ipAddress.ID)

	return resourceIPRead(d, m)
}

func resourceIPRead(d *schema.ResourceData, m interface{}) error {

	c, _ := m.(*Config).Client()

	projectID := d.Get("project_id").(string)

	ipAddress, _, err := c.client.IPAddress.List(projectID, d.Id())
	if err != nil {
		return err
	}

	d.Set("id", ipAddress.ID)
	d.Set("address", ipAddress.Address)
	d.Set("cidr", ipAddress.Cidr)
	d.Set("gateway", ipAddress.Gateway)
	d.Set("ptr", ipAddress.PtrRecord)
	d.Set("type", ipAddress.Type)
	d.Set("routed_to", ipAddress.RoutedTo)
	d.Set("region", ipAddress.Region)

	return nil
}

func resourceIPUpdate(d *schema.ResourceData, m interface{}) error {

	// c, err := cherrygo.NewClient()
	// if err != nil {
	// 	return err
	// }
	c, _ := m.(*Config).Client()

	projectID := d.Get("project_id").(string)

	updateIPRequest := cherrygo.UpdateIPAddress{}

	if d.HasChange("ptr_record") {
		ptrRecord := d.Get("ptr_record").(string)
		updateIPRequest.PtrRecord = ptrRecord
	}

	if d.HasChange("routed_to") {
		routedTo := d.Get("routed_to").(string)
		updateIPRequest.RoutedTo = routedTo
	}

	if d.HasChange("routed_to_hostname") {
		routedTo, err := getIDForServerIP(d, m)
		if err != nil {
			log.Fatalf("Error while gering IP address ID from hostname: %v", err)
		}
		updateIPRequest.RoutedTo = routedTo
	}

	c.client.IPAddress.Update(projectID, d.Id(), &updateIPRequest)

	return resourceIPRead(d, m)
}

func resourceIPDelete(d *schema.ResourceData, m interface{}) error {

	c, _ := m.(*Config).Client()

	projectID := d.Get("project_id").(string)

	ipDeleteRequest := cherrygo.RemoveIPAddress{ID: d.Id()}

	c.client.IPAddress.Remove(projectID, &ipDeleteRequest)

	d.SetId("")
	return nil
}
