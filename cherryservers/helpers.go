package cherryservers

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/cherryservers/cherrygo/v3"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// projectIDSchema returns a standard schema for a project_id
func projectIDSchema() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Description:  "ID of the project you are working on",
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validation.NoZeroValues,
	}
}

func ServerIPSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the IP address",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the IP address",
			},
			"address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The address of the IP",
			},
			"address_family": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "IP address family",
			},
			"cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP address CIDR",
			},
		},
	}
}

func flattenServerIPs(ips []cherrygo.IPAddress) interface{} {
	if ips == nil {
		return nil
	}
	flattendIPs := []map[string]interface{}(nil)
	for _, ip := range ips {
		flattendIPs = append(flattendIPs, map[string]interface{}{
			"id":             ip.ID,
			"type":           ip.Type,
			"address":        ip.Address,
			"address_family": ip.AddressFamily,
			"cidr":           ip.Cidr,
		})
	}
	return flattendIPs
}

func IsBase64(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}

func ServerHostnameToID(hostname string, projectID int, ServerService cherrygo.ServersService) (int, error) {
	serversList, err := serverList(projectID, ServerService)
	for _, s := range serversList {
		if strings.EqualFold(hostname, s.Hostname) {
			return s.ID, err
		}
	}

	return 0, fmt.Errorf("Could not find server with `%s` hostname", hostname)
}

func serverList(projectID int, ServerService cherrygo.ServersService) ([]cherrygo.Server, error) {
	getOptions := cherrygo.GetOptions{
		Fields: []string{"id", "name", "hostname"},
	}
	serverList, _, err := ServerService.List(projectID, &getOptions)

	return serverList, err
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

// is404Error returns true if err is an HTTP 404 error
func is404Error(httpResponse *cherrygo.Response) bool {
	return httpResponse.StatusCode == 404
}
