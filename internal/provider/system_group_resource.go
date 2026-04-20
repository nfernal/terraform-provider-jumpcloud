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
	_ resource.Resource                = &SystemGroupResource{}
	_ resource.ResourceWithImportState = &SystemGroupResource{}
)

// NewSystemGroupResource returns a new system group resource.
func NewSystemGroupResource() resource.Resource {
	return &SystemGroupResource{}
}

// SystemGroupResource defines the system group resource implementation.
type SystemGroupResource struct {
	client *client.JumpCloudClient
}

// SystemGroupResourceModel describes the system group resource data model.
type SystemGroupResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

func (r *SystemGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system_group"
}

func (r *SystemGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a JumpCloud system group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the system group.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the system group.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the system group.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *SystemGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SystemGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SystemGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiGroup := client.SystemGroup{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	group, err := r.client.CreateSystemGroup(ctx, apiGroup)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating JumpCloud System Group",
			fmt.Sprintf("Could not create system group %q: %s", plan.Name.ValueString(), err),
		)
		return
	}

	plan.ID = types.StringValue(group.ID)
	plan.Name = types.StringValue(group.Name)
	plan.Description = types.StringValue(group.Description)

	tflog.Trace(ctx, "created JumpCloud system group", map[string]interface{}{
		"id":   group.ID,
		"name": group.Name,
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SystemGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SystemGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	group, err := r.client.GetSystemGroup(ctx, state.ID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			tflog.Warn(ctx, "JumpCloud system group not found, removing from state", map[string]interface{}{
				"id": state.ID.ValueString(),
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading JumpCloud System Group",
			fmt.Sprintf("Could not read system group %q: %s", state.ID.ValueString(), err),
		)
		return
	}

	state.ID = types.StringValue(group.ID)
	state.Name = types.StringValue(group.Name)
	state.Description = types.StringValue(group.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SystemGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SystemGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state SystemGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiGroup := client.SystemGroup{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	group, err := r.client.UpdateSystemGroup(ctx, state.ID.ValueString(), apiGroup)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating JumpCloud System Group",
			fmt.Sprintf("Could not update system group %q: %s", state.ID.ValueString(), err),
		)
		return
	}

	plan.ID = types.StringValue(group.ID)
	plan.Name = types.StringValue(group.Name)
	plan.Description = types.StringValue(group.Description)

	tflog.Trace(ctx, "updated JumpCloud system group", map[string]interface{}{
		"id": group.ID,
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SystemGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SystemGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSystemGroup(ctx, state.ID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting JumpCloud System Group",
			fmt.Sprintf("Could not delete system group %q: %s", state.ID.ValueString(), err),
		)
	}

	tflog.Trace(ctx, "deleted JumpCloud system group", map[string]interface{}{
		"id": state.ID.ValueString(),
	})
}

func (r *SystemGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
