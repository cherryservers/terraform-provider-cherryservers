package provider

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"testing"

	"github.com/cherryservers/cherrygo/v4"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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

func findPlanIndex(ctx context.Context, id int, client *cherrygo.Client) (int, error) {
	plans, _, err := client.Plans.List(ctx, 0, nil)
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

func TestGeneratePassword(t *testing.T) {
	const (
		cases  = 1000
		length = 20
	)
	passwords := make(map[string]struct{}, cases)

	for range cases {
		p, err := generatePassword()
		if err != nil {
			t.Fatal("failed to generate password")
		}

		if len(p) != length {
			t.Fatalf("password %q length %d, want %d", p, len(p), length)
		}

		if _, ok := passwords[p]; ok {
			t.Errorf("password %q repeated", p)
		}
		passwords[p] = struct{}{}

		hasLowercase := false
		hasNonFirstUppercase := false
		hasNonLastDigit := false
		allAlphaNums := true

		for i, c := range p {
			switch {
			case c >= 'a' && c <= 'z':
				hasLowercase = true
			case c >= 'A' && c <= 'Z':
				if i != 0 {
					hasNonFirstUppercase = true
				}
			case c >= '0' && c <= '9':
				if i != len(p)-1 {
					hasNonLastDigit = true
				}
			default:
				allAlphaNums = false
			}
		}

		if !hasLowercase || !hasNonFirstUppercase || !hasNonLastDigit || !allAlphaNums {
			t.Errorf("password %q didn't fit all constraints: hasLowercase: %t, "+
				"hasNonFirstUppercase: %t, hasNonLastDigit: %t, allAlphaNums: %t", p,
				hasLowercase, hasNonFirstUppercase, hasNonLastDigit, allAlphaNums)
		}

	}
}

type checkFuncWithContext func(context.Context, *terraform.State) error

func checkWithContext(ctx context.Context, f checkFuncWithContext) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return f(ctx, s)
	}
}