package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccGroupDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing with the ldap_group and staff filters set
			{
				Config: testAccGroupDataSourceConfigGroupAndStaff,
				ConfigStateChecks: []statecheck.StateCheck{
					// Filtering to staff should still return at least one member
					// with a populated user identifier.
					statecheck.ExpectKnownValue(
						"data.cis_group.test",
						tfjsonpath.New("members").AtSliceIndex(0).AtMapKey("id"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

const testAccGroupDataSourceConfigGroupAndStaff = `
data "cis_group" "test" {
  ldap_group = "team_moco"
  staff      = true
}
`
