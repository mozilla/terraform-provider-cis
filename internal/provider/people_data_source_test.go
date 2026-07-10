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
						tfjsonpath.New("active"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"data.cis_people.test",
						tfjsonpath.New("cost_center"),
						knownvalue.Int64Exact(14150),
					),
					statecheck.ExpectKnownValue(
						"data.cis_people.test",
						tfjsonpath.New("email"),
						knownvalue.StringExact("jbuckley@mozilla.com"),
					),
					statecheck.ExpectKnownValue(
						"data.cis_people.test",
						tfjsonpath.New("first_name"),
						knownvalue.StringExact("Jon"),
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
						tfjsonpath.New("is_director"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"data.cis_people.test",
						tfjsonpath.New("is_manager"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"data.cis_people.test",
						tfjsonpath.New("is_staff"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"data.cis_people.test",
						tfjsonpath.New("last_name"),
						knownvalue.StringExact("Buckley"),
					),
					statecheck.ExpectKnownValue(
						"data.cis_people.test",
						tfjsonpath.New("ldap_groups"),
						knownvalue.ListPartial(map[int]knownvalue.Check{
							0: knownvalue.StringExact("gh_access_mozilla"),
						}),
					),
					statecheck.ExpectKnownValue(
						"data.cis_people.test",
						tfjsonpath.New("manager_email"),
						knownvalue.StringExact("htahsildoost@mozilla.com"),
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
						tfjsonpath.New("team"),
						knownvalue.StringExact("Site Reliability Engineering (Hamid Tahsildoost)"),
					),
					statecheck.ExpectKnownValue(
						"data.cis_people.test",
						tfjsonpath.New("title"),
						knownvalue.StringExact("Senior Staff Software Engineer"),
					),
					statecheck.ExpectKnownValue(
						"data.cis_people.test",
						tfjsonpath.New("username"),
						knownvalue.StringExact("jbuck"),
					),
					statecheck.ExpectKnownValue(
						"data.cis_people.test",
						tfjsonpath.New("worker_type"),
						knownvalue.StringExact("Employee"),
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
