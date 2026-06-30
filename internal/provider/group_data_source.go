package provider

import (
	"context"
	"fmt"
	"terraform-provider-cis/internal/provider/person_api"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &GroupDataSource{}

func NewGroupDataSource() datasource.DataSource {
	return &GroupDataSource{}
}

// GroupDataSource defines the data source implementation.
type GroupDataSource struct {
	client *person_api.Client
}

// GroupDataSourceModel describes the data source data model.
type GroupDataSourceModel struct {
	LdapGroup types.String            `tfsdk:"ldap_group"`
	Members   []PeopleDataSourceModel `tfsdk:"members"`
}

func (d *GroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (d *GroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Group data source. Lists the members of an LDAP group.",

		Attributes: map[string]schema.Attribute{
			"ldap_group": schema.StringAttribute{
				MarkdownDescription: "LDAP group to list members of",
				Required:            true,
			},
			"members": schema.ListNestedAttribute{
				MarkdownDescription: "People who are members of the LDAP group",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: PeopleAttributes(true),
				},
			},
		},
	}
}

func (d *GroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *GroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data GroupDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	people, err := d.client.GetUsersByLDAPGroup(ctx, data.LdapGroup.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read group members, got error: %s", err.Error()))
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Read %d group members from API", len(people)), map[string]any{"ldap_group": data.LdapGroup.ValueString()})

	members := make([]PeopleDataSourceModel, 0, len(people))
	for i := range people {
		person := people[i]
		member, diags := personToModel(ctx, d.client, &person)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		member.Email = types.StringValue(person.PrimaryEmail.Value)
		members = append(members, member)
	}
	data.Members = members

	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
