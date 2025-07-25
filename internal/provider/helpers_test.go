package provider

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"

	"github.com/cherryservers/cherrygo/v3"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccGetResourceIdInt(resourceName string, resourceType string, s *terraform.State) (int, error) {
	rs, ok := s.RootModule().Resources[resourceName]

	if !ok {
		return 0, fmt.Errorf("resource not found: %s", resourceName)
	}
	if rs.Primary.ID == "" {
		return 0, fmt.Errorf("%s ID is not set", resourceType)
	}
	ID, err := strconv.Atoi(rs.Primary.ID)
	if err != nil {
		return 0, fmt.Errorf("unable to convert %s ID", resourceType)
	}

	return ID, nil
}

// generateAlphaString generates a random lowercase alphabetic string.
func generateAlphaString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz"

	aRecord := make([]byte, length)
	for i := 0; i < length; i++ {
		aRecord[i] = charset[rand.Intn(len(charset))]
	}

	return string(aRecord)
}

var ipv4Regex = regexp.MustCompile(`^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`)

func findPlanIndex(id int, client *cherrygo.Client) (int, error) {
	plans, _, err := client.Plans.List(0, nil)
	if err != nil {
		return 0, err
	}

	for i, v := range plans {
		if v.ID == id {
			return i, nil
		}
	}

	return 0, errors.New("plan index not found")
}
