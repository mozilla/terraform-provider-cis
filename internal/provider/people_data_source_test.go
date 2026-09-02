package provider

import (
	"context"
	"testing"

	"terraform-provider-cis/internal/provider/person_api"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestPersonToModelAdditionalAttributes(t *testing.T) {
	t.Parallel()

	person := &person_api.Person{
		Identities: person_api.IdentitiesAttributesValuesArray{
			BugzillaMozillaOrgID:           &person_api.StandardAttributeString{Value: "12345"},
			BugzillaMozillaOrgPrimaryEmail: &person_api.StandardAttributeString{Value: "person@example.com"},
			GithubPrimaryEmail:             &person_api.StandardAttributeString{Value: "person@users.noreply.github.com"},
		},
		Location: person_api.StandardAttributeString{Value: "Toronto, Canada"},
		StaffInformation: person_api.StaffInformationValuesArray{
			OfficeLocation: person_api.StandardAttributeString{Value: "Toronto"},
		},
		Timezone: person_api.StandardAttributeString{Value: "America/Toronto"},
	}

	model, diags := personToModel(context.Background(), person)
	if diags.HasError() {
		t.Fatalf("personToModel returned diagnostics: %v", diags)
	}

	want := map[string]struct {
		got  types.String
		want string
	}{
		"bugzilla_email":  {model.Bugzilla_Email, "person@example.com"},
		"bugzilla_id":     {model.Bugzilla_Id, "12345"},
		"github_email":    {model.GitHub_Email, "person@users.noreply.github.com"},
		"location":        {model.Location, "Toronto, Canada"},
		"office_location": {model.Office_Location, "Toronto"},
		"timezone":        {model.Timezone, "America/Toronto"},
	}
	for name, test := range want {
		if test.got.IsNull() || test.got.ValueString() != test.want {
			t.Errorf("%s = %q, want %q", name, test.got.ValueString(), test.want)
		}
	}
}

func TestPersonToModelMissingIdentityAttributesAreNull(t *testing.T) {
	t.Parallel()

	model, diags := personToModel(context.Background(), &person_api.Person{})
	if diags.HasError() {
		t.Fatalf("personToModel returned diagnostics: %v", diags)
	}

	for name, value := range map[string]types.String{
		"bugzilla_email": model.Bugzilla_Email,
		"bugzilla_id":    model.Bugzilla_Id,
		"github_email":   model.GitHub_Email,
	} {
		if !value.IsNull() {
			t.Errorf("%s = %q, want null", name, value.ValueString())
		}
	}
}

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
						tfjsonpath.New("bugzilla_email"),
						knownvalue.StringExact("jbuckley@mozilla.com"),
					),
					statecheck.ExpectKnownValue(
						"data.cis_people.test",
						tfjsonpath.New("bugzilla_id"),
						knownvalue.StringExact("567620"),
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
						tfjsonpath.New("github_email"),
						knownvalue.StringExact("jon@jbuckley.ca"),
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
						tfjsonpath.New("location"),
						knownvalue.StringExact("Toronto, ON"),
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
						tfjsonpath.New("office_location"),
						knownvalue.StringExact("Regional Toronto Office"),
					),
					statecheck.ExpectKnownValue(
						"data.cis_people.test",
						tfjsonpath.New("team"),
						knownvalue.StringExact("Site Reliability Engineering (Hamid Tahsildoost)"),
					),
					statecheck.ExpectKnownValue(
						"data.cis_people.test",
						tfjsonpath.New("timezone"),
						knownvalue.StringExact("America/Toronto"),
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
