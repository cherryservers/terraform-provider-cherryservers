package provider

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/cherryservers/cherrygo/v4"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	testProjectNamePrefix = "terraform_test_project_"
	defaultTestImage      = "ubuntu_26_04_64bit"
)

var (
	testCherryGoClient *cherrygo.Client
	testTeam           int
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"cherryservers": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
}

func setTestTeam() error {
	const teamIDVar = "CHERRY_TEST_TEAM_ID"

	team, ok := os.LookupEnv(teamIDVar)
	if !ok {
		return fmt.Errorf("%s must be set for acceptance tests", teamIDVar)
	}
	id, err := strconv.Atoi(team)
	if err != nil {
		return fmt.Errorf("%s must be an integer: %s", teamIDVar, err.Error())
	}

	testTeam = id
	return nil
}

func setupClient() (*cherrygo.Client, error) {
	apiKey := os.Getenv(apiKeyVar)
	if apiKey == "" {
		return nil, fmt.Errorf("%s must be set for acceptance tests", apiKeyVar)
	}

	userAgent := "terraform-provider/cherryservers/test terraform/dev"
	args := []cherrygo.ClientOpt{cherrygo.WithAPIKey(apiKey), cherrygo.WithUserAgent(userAgent)}
	client, err := cherrygo.NewClient(args...)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func TestMain(m *testing.M) {
	// Skip setup on unit tests.
	if acc := os.Getenv(resource.EnvTfAcc); acc != "" {
		var err error
		testCherryGoClient, err = setupClient()
		if err != nil {
			log.Fatalf("failed to initialize api client: %s", err.Error())
		}

		err = setTestTeam()
		if err != nil {
			log.Fatalf("failed to get test team: %s", err.Error())
		}

	}
	resource.TestMain(m)
}

func TestAPITokenConflictsWithAPIKey(t *testing.T) {
	resource.ParallelTest(
		t, resource.TestCase{
			IsUnitTest: true,
			Steps: []resource.TestStep{
				{
					// A resource is required, otherwise the provider is skipped.
					Config: `provider "cherryservers" {
				api_key = "abc"
				api_token = "abc"
				}

				resource "cherryservers_server" "srv" {
				region = "test"
				plan = "test"
				project = 123
				}`,
					ExpectError:              regexp.MustCompile(`Attribute "api_token" cannot be specified when "api_key" is specified`),
					ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				},
			},
		},
	)
}

func assertDiags(t *testing.T, want, got diag.Diagnostics) {
	t.Helper()

	for _, d := range want {
		if !got.Contains(d) {
			t.Errorf("missing diagnostic with type %T, severity %q, summary %q and detail %q", d, d.Severity(), d.Summary(), d.Detail())
		}
	}

	for _, d := range got {
		if !want.Contains(d) {
			t.Errorf("unexpected diagnostic with type %T, severity %q, summary %q and detail %q", d, d.Severity(), d.Summary(), d.Detail())
		}
	}
}

func TestAPISecretHierarchy(t *testing.T) {
	cases := []struct {
		name         string
		authTokenVar string
		authKeyVar   string
		apiKeyVar    string
		apiTokenTF   types.String
		apiKeyTF     types.String
		wantDiags    []diag.Diagnostic
		want         string
	}{
		{
			name:         "configure with CHERRY_AUTH_TOKEN env var",
			authTokenVar: "test-secret",
			want:         "test-secret",
			wantDiags: []diag.Diagnostic{
				diag.NewWarningDiagnostic("CHERRY_AUTH_TOKEN is deprecated",
					"CHERRY_AUTH_TOKEN is deprecated and will be removed in the next major version of the provider, please use CHERRY_API_KEY instead."),
			},
		},
		{
			name:         "CHERRY_AUTH_KEY beats CHERRY_AUTH_TOKEN",
			authTokenVar: "bad",
			authKeyVar:   "test-secret",
			wantDiags: []diag.Diagnostic{
				diag.NewWarningDiagnostic("CHERRY_AUTH_KEY is deprecated",
					"CHERRY_AUTH_KEY is deprecated and will be removed in the next major version of the provider, please use CHERRY_API_KEY instead."),
			},
			want: "test-secret",
		},
		{
			name:       "CHERRY_API_KEY beats CHERRY_AUTH_KEY",
			authKeyVar: "bad",
			apiKeyVar:  "test-secret",
			want:       "test-secret",
		},
		{
			name:         "api_token tf attribute beats all env vars",
			authTokenVar: "bad",
			authKeyVar:   "bad",
			apiKeyVar:    "bad",
			apiTokenTF:   types.StringValue("test-secret"),
			want:         "test-secret",
		},
		{
			name:         "api_key tf attribute beats all env vars",
			authTokenVar: "bad",
			authKeyVar:   "bad",
			apiKeyVar:    "bad",
			apiKeyTF:     types.StringValue("test-secret"),
			want:         "test-secret",
		},
		{
			name: "error when no api secret supplied",
			wantDiags: []diag.Diagnostic{
				diag.NewAttributeErrorDiagnostic(path.Root("api_key"), "Missing CherryServers API key",
					"The provider cannot create the CherryServers API client as there is a missing or empty value for the CherryServers API key. Set the API key value in the configuration or use the CHERRY_API_KEY environment variable."),
			},
		},
		{
			name:       "error with unknown APIToken tf attribute",
			apiTokenTF: types.StringUnknown(),
			wantDiags: []diag.Diagnostic{
				diag.NewAttributeErrorDiagnostic(path.Root("api_token"), "Unknown CherryServers API token",
					"The provider cannot create the CherryServers API client as there is an unknown configuration value for the CherryServers API token. Either target apply the source of the value first, set the value statically in the configuration, or use the CHERRY_API_KEY environment variable."),
			},
		},
		{
			name:     "error with unknown APIKey tf attribute",
			apiKeyTF: types.StringUnknown(),
			wantDiags: []diag.Diagnostic{
				diag.NewAttributeErrorDiagnostic(path.Root("api_key"), "Unknown CherryServers API key",
					"The provider cannot create the CherryServers API client as there is an unknown configuration value for the CherryServers API key. Either target apply the source of the value first, set the value statically in the configuration, or use the CHERRY_API_KEY environment variable."),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				diags, wantDiags diag.Diagnostics
				cfg              CherryServersProviderModel
			)

			wantDiags = append(wantDiags, tc.wantDiags...)
			t.Setenv("CHERRY_AUTH_TOKEN", tc.authTokenVar)
			t.Setenv("CHERRY_AUTH_KEY", tc.authKeyVar)
			t.Setenv("CHERRY_API_KEY", tc.apiKeyVar)
			cfg.APIToken = tc.apiTokenTF
			cfg.APIKey = tc.apiKeyTF

			got := apiKey(&diags, cfg)

			if got != tc.want {
				t.Errorf("want %q api secret, got %q", tc.want, got)
			}

			assertDiags(t, wantDiags, diags)
		})
	}
}
