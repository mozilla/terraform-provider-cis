package provider

import (
	"context"
	"fmt"
	"terraform-provider-cis/internal/provider/person_api"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/sync/errgroup"
)

// githubLookupConcurrency bounds the number of in-flight whoami lookups when
// resolving GitHub usernames for a group's members.
const githubLookupConcurrency = 20

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
	Staff     types.Bool              `tfsdk:"staff"`
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
				Optional:            true,
			},
			"staff": schema.BoolAttribute{
				MarkdownDescription: "Filter members by whether they are Mozilla staff",
				Optional:            true,
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

func (d GroupDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.AtLeastOneOf(
			path.MatchRoot("ldap_group"),
			path.MatchRoot("staff"),
		),
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

	// staff is optional; only pass it as a filter when the caller set it.
	var staff *bool
	if !data.Staff.IsNull() && !data.Staff.IsUnknown() {
		staff = data.Staff.ValueBoolPointer()
	}

	people, err := d.client.GetUsersByAttributes(ctx, data.LdapGroup.ValueString(), staff)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read group members, got error: %s", err.Error()))
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Read %d group members from API", len(people)), map[string]any{"ldap_group": data.LdapGroup.ValueString()})

	// Map every member using only local data first; this is network-free.
	members := make([]PeopleDataSourceModel, len(people))
	for i := range people {
		member, diags := personToModel(ctx, &people[i])
		resp.Diagnostics.Append(diags...)
		member.Email = types.StringValue(people[i].PrimaryEmail.Value)
		members[i] = member
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Resolving GitHub usernames requires one whoami lookup per member, which
	// dominates the read time for large groups. Fan the lookups out with bounded
	// concurrency; each goroutine writes a distinct index, so no locking needed.
	group, groupCtx := errgroup.WithContext(ctx)
	group.SetLimit(githubLookupConcurrency)
	for i := range people {
		group.Go(func() error {
			members[i].GitHub_Username = types.StringValue(resolveGithubUsername(groupCtx, d.client, &people[i]))
			return nil
		})
	}
	// resolveGithubUsername never returns an error, so Wait cannot fail.
	_ = group.Wait()

	data.Members = members

	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
