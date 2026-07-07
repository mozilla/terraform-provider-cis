package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccExampleDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccExampleDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.cis_people.test",
						tfjsonpath.New("email"),
						knownvalue.StringExact("jbuckley@mozilla.com"),
					),
					statecheck.ExpectKnownValue(
						"data.cis_people.test",
						tfjsonpath.New("github_id"),
						knownvalue.StringExact("578466"),
					),
					statecheck.ExpectKnownValue(
						"data.cis_people.test",
						tfjsonpath.New("github_node_id"),
						knownvalue.StringExact("MDQ6VXNlcjU3ODQ2Ng=="),
					),
					statecheck.ExpectKnownValue(
						"data.cis_people.test",
						tfjsonpath.New("github_username"),
						knownvalue.StringExact("jbuck"),
					),
					statecheck.ExpectKnownValue(
						"data.cis_people.test",
						tfjsonpath.New("id"),
						knownvalue.StringExact("ad|Mozilla-LDAP|jbuckley"),
					),
					statecheck.ExpectKnownValue(
						"data.cis_people.test",
						tfjsonpath.New("ldap_groups"),
						knownvalue.ListPartial(map[int]knownvalue.Check{
							0: knownvalue.StringExact("inventory"),
						}),
					),
					statecheck.ExpectKnownValue(
						"data.cis_people.test",
						tfjsonpath.New("mozilliansorg_groups"),
						knownvalue.ListPartial(map[int]knownvalue.Check{
							0: knownvalue.StringExact("aws_billing_access"),
						}),
					),
					statecheck.ExpectKnownValue(
						"data.cis_people.test",
						tfjsonpath.New("username"),
						knownvalue.StringExact("jbuck"),
					),
				},
			},
		},
	})
}

const testAccExampleDataSourceConfig = `
data "cis_people" "test" {
  email = "jbuckley@mozilla.com"
}
`
