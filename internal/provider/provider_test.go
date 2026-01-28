package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories is used to instantiate a provider during acceptance testing.
// The factory function is called for each Terraform CLI command to create a provider
// server that the CLI can connect to and interact with.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"cis": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
	if auth0_client_id := os.Getenv("AUTH0_CLIENT_ID"); auth0_client_id == "" {
		t.Fatal("AUTH0_CLIENT_ID must be set for acceptance tests")
	}
	if auth0_client_secret := os.Getenv("AUTH0_CLIENT_SECRET"); auth0_client_secret == "" {
		t.Fatal("AUTH0_CLIENT_SECRET must be set for acceptance tests")
	}
}
