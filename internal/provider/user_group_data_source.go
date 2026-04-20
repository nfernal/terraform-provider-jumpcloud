// Copyright (c) Nfernal
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/nfernal/terraform-provider-jumpcloud/internal/client"
)

var _ datasource.DataSource = &UserGroupDataSource{}

// NewUserGroupDataSource returns a new user group data source.
func NewUserGroupDataSource() datasource.DataSource {
	return &UserGroupDataSource{}
}

// UserGroupDataSource defines the user group data source implementation.
type UserGroupDataSource struct {
	client *client.JumpCloudClient
}

// UserGroupDataSourceModel describes the user group data source data model.
type UserGroupDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

func (d *UserGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_group"
}

func (d *UserGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to look up an existing JumpCloud user group by name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the user group.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the user group to look up.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the user group.",
				Computed:            true,
			},
		},
	}
}

func (d *UserGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.JumpCloudClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.JumpCloudClient, got: %T", req.ProviderData),
		)
		return
	}

	d.client = c
}

func (d *UserGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config UserGroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	group, err := d.client.GetUserGroupByName(ctx, config.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading JumpCloud User Group",
			fmt.Sprintf("Could not find user group %q: %s", config.Name.ValueString(), err),
		)
		return
	}

	config.ID = types.StringValue(group.ID)
	config.Name = types.StringValue(group.Name)
	config.Description = types.StringValue(group.Description)

	tflog.Trace(ctx, "read JumpCloud user group data source", map[string]interface{}{
		"id":   group.ID,
		"name": group.Name,
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
