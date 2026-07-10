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
	Active               types.Bool   `tfsdk:"active"`
	Email                types.String `tfsdk:"email"`
	First_Name           types.String `tfsdk:"first_name"`
	GitHub_Id            types.String `tfsdk:"github_id"`
	GitHub_Node_Id       types.String `tfsdk:"github_node_id"`
	GitHub_Username      types.String `tfsdk:"github_username"`
	Id                   types.String `tfsdk:"id"`
	Is_Director          types.Bool   `tfsdk:"is_director"`
	Is_Manager           types.Bool   `tfsdk:"is_manager"`
	Is_Staff             types.Bool   `tfsdk:"is_staff"`
	Last_Name            types.String `tfsdk:"last_name"`
	LDAP_Groups          types.List   `tfsdk:"ldap_groups"`
	Manager_Email        types.String `tfsdk:"manager_email"`
	Mozilliansorg_Groups types.List   `tfsdk:"mozilliansorg_groups"`
	Team                 types.String `tfsdk:"team"`
	Title                types.String `tfsdk:"title"`
	Username             types.String `tfsdk:"username"`
	Worker_Type          types.String `tfsdk:"worker_type"`
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
		"active": schema.BoolAttribute{
			MarkdownDescription: "Whether the person's account is active",
			Computed:            true,
		},
		"email": queryField("People email address"),
		"first_name": schema.StringAttribute{
			MarkdownDescription: "First name",
			Computed:            true,
		},
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
		"is_director": schema.BoolAttribute{
			MarkdownDescription: "Whether the person is a director",
			Computed:            true,
		},
		"is_manager": schema.BoolAttribute{
			MarkdownDescription: "Whether the person is a manager",
			Computed:            true,
		},
		"is_staff": schema.BoolAttribute{
			MarkdownDescription: "Whether the person is Mozilla staff",
			Computed:            true,
		},
		"last_name": schema.StringAttribute{
			MarkdownDescription: "Last name",
			Computed:            true,
		},
		"ldap_groups": schema.ListAttribute{
			ElementType:         types.StringType,
			MarkdownDescription: "LDAP groups the user is in",
			Computed:            true,
		},
		"manager_email": schema.StringAttribute{
			MarkdownDescription: "Primary work email of the person's manager",
			Computed:            true,
		},
		"mozilliansorg_groups": schema.ListAttribute{
			ElementType:         types.StringType,
			MarkdownDescription: "Mozilliansorg groups the user is in",
			Computed:            true,
		},
		"team": schema.StringAttribute{
			MarkdownDescription: "Team",
			Computed:            true,
		},
		"title": schema.StringAttribute{
			MarkdownDescription: "Job title",
			Computed:            true,
		},
		"username": queryField("People username"),
		"worker_type": schema.StringAttribute{
			MarkdownDescription: "Worker type (e.g. Employee)",
			Computed:            true,
		},
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

	if data.Email.ValueString() != "" {
		person, err = d.client.GetPersonByEmail(ctx, data.Email.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read person, got error: %s", err.Error()))
		}
	}

	tflog.Info(ctx, fmt.Sprintf("Read data from API %#v", person), map[string]any{"email": data.Email.ValueString()})

	model, diags := personToModel(ctx, person)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	model.GitHub_Username = types.StringValue(resolveGithubUsername(ctx, d.client, person))

	// Preserve the email the caller queried with; the rest comes from the API.
	model.Email = data.Email
	data = model

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

// personToModel maps a person_api.Person into a PeopleDataSourceModel using only
// locally-available data. It is shared by the cis_people and cis_group data
// sources. GitHub_Username is set to the value CIS provides; callers wanting the
// freshest value should overwrite it with resolveGithubUsername (a network call).
// The Email field is left unset, as it is not consistently present on profiles
// returned in bulk; the caller can populate it from the source query when known.
func personToModel(ctx context.Context, person *person_api.Person) (PeopleDataSourceModel, diag.Diagnostics) {
	var data PeopleDataSourceModel
	var diags diag.Diagnostics

	data.Id = types.StringValue(person.UserID.Value)
	data.Active = types.BoolValue(person.Active.Value)
	data.First_Name = types.StringValue(person.FirstName.Value)
	data.Last_Name = types.StringValue(person.LastName.Value)

	if person.Identities.GithubIDV3 != nil {
		data.GitHub_Id = types.StringValue(person.Identities.GithubIDV3.Value)
	}
	if person.Identities.GithubIDV4 != nil {
		data.GitHub_Node_Id = types.StringValue(person.Identities.GithubIDV4.Value)
	}
	data.GitHub_Username = types.StringValue(person.Usernames.Values.GitHubUsername)

	ldapGroups, ldapDiags := types.ListValueFrom(ctx, types.StringType, person.AccessInformation.LDAP.List)
	diags.Append(ldapDiags...)
	data.LDAP_Groups = ldapGroups

	groups, groupDiags := types.ListValueFrom(ctx, types.StringType, person.AccessInformation.Mozilliansorg.List)
	diags.Append(groupDiags...)
	data.Mozilliansorg_Groups = groups

	// managers_primary_work_email lives in the HRIS values map, which is untyped.
	if managerEmail, ok := person.AccessInformation.Hris.Values["managers_primary_work_email"].(string); ok {
		data.Manager_Email = types.StringValue(managerEmail)
	}

	data.Is_Director = types.BoolValue(person.StaffInformation.Director.Value)
	data.Is_Manager = types.BoolValue(person.StaffInformation.Manager.Value)
	data.Is_Staff = types.BoolValue(person.StaffInformation.Staff.Value)
	data.Team = types.StringValue(person.StaffInformation.Team.Value)
	data.Title = types.StringValue(person.StaffInformation.Title.Value)
	data.Worker_Type = types.StringValue(person.StaffInformation.WorkerType.Value)

	data.Username = types.StringValue(person.PrimaryUsername.Value)

	return data, diags
}

// resolveGithubUsername returns the freshest GitHub username for the person,
// performing the dino-park whoami lookup when a v3 id is available. On any
// lookup failure it falls back to the CIS-provided value. This makes a network
// call and is safe to invoke concurrently for different people.
func resolveGithubUsername(ctx context.Context, client *person_api.Client, person *person_api.Person) string {
	githubUsername := person.Usernames.Values.GitHubUsername
	if person.Identities.GithubIDV3 != nil && person.Identities.GithubIDV3.Value != "" {
		if username, err := client.GetGithubUsernameByNodeID(ctx, person.Identities.GithubIDV3.Value); err == nil {
			githubUsername = username
		}
	}
	return githubUsername
}
