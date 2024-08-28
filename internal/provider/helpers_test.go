package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"math/rand"
	"strconv"
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

// generateAlphaString generates a random lowercase alphabetic string
func generateAlphaString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz"

	aRecord := make([]byte, length)
	for i := 0; i < length; i++ {
		aRecord[i] = charset[rand.Intn(len(charset))]
	}

	return string(aRecord)
}
