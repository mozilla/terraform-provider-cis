package provider

import (
	"context"
	"fmt"
	"terraform-provider-cis/internal/provider/person_api"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &PeopleDataSource{}

func NewPeopleDataSource() datasource.DataSource {
	return &PeopleDataSource{}
}

// PeopleDataSource defines the data source implementation.
type PeopleDataSource struct {
	client *person_api.Client
}

// PeopleDataSourceModel describes the data source data model.
type PeopleDataSourceModel struct {
	Email                types.String `tfsdk:"email"`
	GitHub_Id            types.String `tfsdk:"github_id"`
	GitHub_Node_Id       types.String `tfsdk:"github_node_id"`
	GitHub_Username      types.String `tfsdk:"github_username"`
	Id                   types.String `tfsdk:"id"`
	LDAP_Groups          types.List   `tfsdk:"ldap_groups"`
	Mozilliansorg_Groups types.List   `tfsdk:"mozilliansorg_groups"`
	Username             types.String `tfsdk:"username"`
}

func (d *PeopleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_people"
}

// PeopleAttributes returns the schema attributes describing a person. It is
// shared by the cis_people data source and the cis_group members list.
//
// When readOnly is false (the cis_people data source), email, id, and username
// are Optional because they identify which person to look up. When readOnly is
// true (embedded as cis_group members output), every attribute is Computed.
func PeopleAttributes(readOnly bool) map[string]schema.Attribute {
	// queryField is Optional when used as a lookup input, Computed when the whole
	// object is read-only output.
	queryField := func(description string) schema.StringAttribute {
		return schema.StringAttribute{
			MarkdownDescription: description,
			Optional:            !readOnly,
			Computed:            readOnly,
		}
	}

	return map[string]schema.Attribute{
		"email": queryField("People email address"),
		"github_id": schema.StringAttribute{
			MarkdownDescription: "GitHub ID",
			Computed:            true,
		},
		"github_node_id": schema.StringAttribute{
			MarkdownDescription: "GitHub node ID",
			Computed:            true,
		},
		"github_username": schema.StringAttribute{
			MarkdownDescription: "GitHub username",
			Computed:            true,
		},
		"id": queryField("People user identifier"),
		"ldap_groups": schema.ListAttribute{
			ElementType:         types.StringType,
			MarkdownDescription: "LDAP groups the user is in",
			Computed:            true,
		},
		"mozilliansorg_groups": schema.ListAttribute{
			ElementType:         types.StringType,
			MarkdownDescription: "Mozilliansorg groups the user is in",
			Computed:            true,
		},
		"username": queryField("People username"),
	}
}

func (d *PeopleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "People data source",

		Attributes: PeopleAttributes(false),
	}
}

func (d PeopleDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.AtLeastOneOf(
			path.MatchRoot("email"),
			path.MatchRoot("id"),
			path.MatchRoot("username"),
		),
	}
}

func (d *PeopleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*person_api.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *person_api.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *PeopleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PeopleDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := d.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read people, got error: %s", err))
	//     return
	// }

	tflog.Info(ctx, fmt.Sprintf("HTTP Request: %#v", d.client))

	var person *person_api.Person
	var err error
	var diags diag.Diagnostics

	if data.Email.ValueString() != "" {
		person, err = d.client.GetPersonByEmail(ctx, data.Email.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read person, got error: %s", err.Error()))
		}
	}

	tflog.Info(ctx, fmt.Sprintf("Read data from API %#v", person), map[string]any{"email": data.Email.ValueString()})

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	data.Id = types.StringValue(person.UserID.Value)

	if person.Identities.GithubIDV3 != nil {
		data.GitHub_Id = types.StringValue(person.Identities.GithubIDV3.Value)
	}
	if person.Identities.GithubIDV4 != nil {
		data.GitHub_Node_Id = types.StringValue(person.Identities.GithubIDV4.Value)
	}
	githubUsername := person.Usernames.Values.GitHubUsername
	if person.Identities.GithubIDV3 != nil && person.Identities.GithubIDV3.Value != "" {
		if username, err := d.client.GetGithubUsernameByNodeID(ctx, person.Identities.GithubIDV3.Value); err == nil {
			githubUsername = username
		}
	}
	data.GitHub_Username = types.StringValue(githubUsername)

	ldapGroups, ldapDiags := types.ListValueFrom(ctx, types.StringType, person.AccessInformation.LDAP.List)
	for _, d := range ldapDiags {
		resp.Diagnostics.Append(d)
	}
	data.LDAP_Groups = ldapGroups

	data.Mozilliansorg_Groups, diags = types.ListValueFrom(ctx, types.StringType, person.AccessInformation.Mozilliansorg.List)
	for _, d := range diags {
		resp.Diagnostics.Append(d)
	}
	data.Username = types.StringValue(person.PrimaryUsername.Value)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
