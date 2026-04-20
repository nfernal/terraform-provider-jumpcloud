// Copyright (c) Nfernal
// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/nfernal/terraform-provider-jumpcloud/internal/client"
)

var (
	_ resource.Resource                = &UserGroupResource{}
	_ resource.ResourceWithImportState = &UserGroupResource{}
)

// NewUserGroupResource returns a new user group resource.
func NewUserGroupResource() resource.Resource {
	return &UserGroupResource{}
}

// UserGroupResource defines the user group resource implementation.
type UserGroupResource struct {
	client *client.JumpCloudClient
}

// UserGroupResourceModel describes the user group resource data model.
type UserGroupResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

func (r *UserGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_group"
}

func (r *UserGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a JumpCloud user group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the user group.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the user group.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the user group.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *UserGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.JumpCloudClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.JumpCloudClient, got: %T", req.ProviderData),
		)
		return
	}

	r.client = c
}

func (r *UserGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan UserGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiGroup := client.UserGroup{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	group, err := r.client.CreateUserGroup(ctx, apiGroup)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating JumpCloud User Group",
			fmt.Sprintf("Could not create user group %q: %s", plan.Name.ValueString(), err),
		)
		return
	}

	plan.ID = types.StringValue(group.ID)
	plan.Name = types.StringValue(group.Name)
	plan.Description = types.StringValue(group.Description)

	tflog.Trace(ctx, "created JumpCloud user group", map[string]interface{}{
		"id":   group.ID,
		"name": group.Name,
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *UserGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state UserGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	group, err := r.client.GetUserGroup(ctx, state.ID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			tflog.Warn(ctx, "JumpCloud user group not found, removing from state", map[string]interface{}{
				"id": state.ID.ValueString(),
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading JumpCloud User Group",
			fmt.Sprintf("Could not read user group %q: %s", state.ID.ValueString(), err),
		)
		return
	}

	state.ID = types.StringValue(group.ID)
	state.Name = types.StringValue(group.Name)
	state.Description = types.StringValue(group.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *UserGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan UserGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state UserGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiGroup := client.UserGroup{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	group, err := r.client.UpdateUserGroup(ctx, state.ID.ValueString(), apiGroup)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating JumpCloud User Group",
			fmt.Sprintf("Could not update user group %q: %s", state.ID.ValueString(), err),
		)
		return
	}

	plan.ID = types.StringValue(group.ID)
	plan.Name = types.StringValue(group.Name)
	plan.Description = types.StringValue(group.Description)

	tflog.Trace(ctx, "updated JumpCloud user group", map[string]interface{}{
		"id": group.ID,
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *UserGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state UserGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteUserGroup(ctx, state.ID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting JumpCloud User Group",
			fmt.Sprintf("Could not delete user group %q: %s", state.ID.ValueString(), err),
		)
	}

	tflog.Trace(ctx, "deleted JumpCloud user group", map[string]interface{}{
		"id": state.ID.ValueString(),
	})
}

func (r *UserGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
