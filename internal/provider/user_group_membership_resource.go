// Copyright (c) Nfernal
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

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
	_ resource.Resource                = &UserGroupMembershipResource{}
	_ resource.ResourceWithImportState = &UserGroupMembershipResource{}
)

// NewUserGroupMembershipResource returns a new user group membership resource.
func NewUserGroupMembershipResource() resource.Resource {
	return &UserGroupMembershipResource{}
}

// UserGroupMembershipResource defines the user group membership resource implementation.
type UserGroupMembershipResource struct {
	client *client.JumpCloudClient
}

// UserGroupMembershipResourceModel describes the user group membership resource data model.
type UserGroupMembershipResourceModel struct {
	ID      types.String `tfsdk:"id"`
	GroupID types.String `tfsdk:"group_id"`
	UserID  types.String `tfsdk:"user_id"`
}

func (r *UserGroupMembershipResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_group_membership"
}

func (r *UserGroupMembershipResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the membership of a JumpCloud user in a user group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The composite identifier of the membership (`group_id/user_id`).",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"group_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the user group.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"user_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the user to add to the group.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *UserGroupMembershipResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *UserGroupMembershipResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan UserGroupMembershipResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	op := client.GroupMembershipOp{
		Op:   "add",
		Type: "user",
		ID:   plan.UserID.ValueString(),
	}

	err := r.client.ModifyUserGroupMembership(ctx, plan.GroupID.ValueString(), op)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating JumpCloud User Group Membership",
			fmt.Sprintf("Could not add user %q to group %q: %s",
				plan.UserID.ValueString(), plan.GroupID.ValueString(), err),
		)
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s/%s", plan.GroupID.ValueString(), plan.UserID.ValueString()))

	tflog.Trace(ctx, "created JumpCloud user group membership", map[string]interface{}{
		"group_id": plan.GroupID.ValueString(),
		"user_id":  plan.UserID.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *UserGroupMembershipResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state UserGroupMembershipResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	members, err := r.client.ListUserGroupMembers(ctx, state.GroupID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			tflog.Warn(ctx, "JumpCloud user group not found, removing membership from state", map[string]interface{}{
				"group_id": state.GroupID.ValueString(),
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading JumpCloud User Group Membership",
			fmt.Sprintf("Could not list members of group %q: %s", state.GroupID.ValueString(), err),
		)
		return
	}

	found := false
	for _, member := range members {
		if member.ID == state.UserID.ValueString() && member.Type == "user" {
			found = true
			break
		}
	}

	if !found {
		tflog.Warn(ctx, "JumpCloud user group membership not found, removing from state", map[string]interface{}{
			"group_id": state.GroupID.ValueString(),
			"user_id":  state.UserID.ValueString(),
		})
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *UserGroupMembershipResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Both attributes use RequiresReplace, so Update should never be called.
	resp.Diagnostics.AddError(
		"Unexpected Update",
		"User group membership does not support in-place updates. Both group_id and user_id require resource replacement.",
	)
}

func (r *UserGroupMembershipResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state UserGroupMembershipResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	op := client.GroupMembershipOp{
		Op:   "remove",
		Type: "user",
		ID:   state.UserID.ValueString(),
	}

	err := r.client.ModifyUserGroupMembership(ctx, state.GroupID.ValueString(), op)
	if err != nil {
		if client.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting JumpCloud User Group Membership",
			fmt.Sprintf("Could not remove user %q from group %q: %s",
				state.UserID.ValueString(), state.GroupID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted JumpCloud user group membership", map[string]interface{}{
		"group_id": state.GroupID.ValueString(),
		"user_id":  state.UserID.ValueString(),
	})
}

func (r *UserGroupMembershipResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID in the format 'group_id/user_id', got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("group_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("user_id"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
