// Copyright (c) Nfernal
// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/nfernal/terraform-provider-jumpcloud/internal/client"
)

var (
	_ resource.Resource                = &UserResource{}
	_ resource.ResourceWithImportState = &UserResource{}
)

// NewUserResource returns a new user resource.
func NewUserResource() resource.Resource {
	return &UserResource{}
}

// UserResource defines the user resource implementation.
type UserResource struct {
	client *client.JumpCloudClient
}

// UserResourceModel describes the user resource data model.
type UserResourceModel struct {
	ID                          types.String `tfsdk:"id"`
	Username                    types.String `tfsdk:"username"`
	Email                       types.String `tfsdk:"email"`
	Firstname                   types.String `tfsdk:"firstname"`
	Lastname                    types.String `tfsdk:"lastname"`
	Password                    types.String `tfsdk:"password"`
	AccountLocked               types.Bool   `tfsdk:"account_locked"`
	Activated                   types.Bool   `tfsdk:"activated"`
	AllowPublicKey              types.Bool   `tfsdk:"allow_public_key"`
	Sudo                        types.Bool   `tfsdk:"sudo"`
	EnableManagedUID            types.Bool   `tfsdk:"enable_managed_uid"`
	UnixUID                     types.Int64  `tfsdk:"unix_uid"`
	UnixGUID                    types.Int64  `tfsdk:"unix_guid"`
	PasswordlessSudo            types.Bool   `tfsdk:"passwordless_sudo"`
	ExternallyManaged           types.Bool   `tfsdk:"externally_managed"`
	LdapBindingUser             types.Bool   `tfsdk:"ldap_binding_user"`
	EnableUserPortalMultifactor types.Bool   `tfsdk:"enable_user_portal_multifactor"`
	Suspended                   types.Bool   `tfsdk:"suspended"`
	State                       types.String `tfsdk:"state"`
	Description                 types.String `tfsdk:"description"`
	Department                  types.String `tfsdk:"department"`
	CostCenter                  types.String `tfsdk:"cost_center"`
	Company                     types.String `tfsdk:"company"`
	EmployeeType                types.String `tfsdk:"employee_type"`
	EmployeeIdentifier          types.String `tfsdk:"employee_identifier"`
	JobTitle                    types.String `tfsdk:"job_title"`
	Location                    types.String `tfsdk:"location"`
	Manager                     types.String `tfsdk:"manager"`
	MiddleName                  types.String `tfsdk:"middle_name"`
	DisplayName                 types.String `tfsdk:"display_name"`
	ExternalDN                  types.String `tfsdk:"external_dn"`
	ExternalSourceType          types.String `tfsdk:"external_source_type"`
	MFAConfigured               types.Bool   `tfsdk:"mfa_configured"`
	MFAExclusion                types.Bool   `tfsdk:"mfa_exclusion"`
}

func (r *UserResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *UserResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a JumpCloud system user.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the user.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "The username for the user. Must be unique.",
				Required:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "The email address of the user. Must be unique.",
				Required:            true,
			},
			"firstname": schema.StringAttribute{
				MarkdownDescription: "The first name of the user.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"lastname": schema.StringAttribute{
				MarkdownDescription: "The last name of the user.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The password for the user. Write-only; not returned by the API.",
				Optional:            true,
				Sensitive:           true,
			},
			"account_locked": schema.BoolAttribute{
				MarkdownDescription: "Whether the user account is locked.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"activated": schema.BoolAttribute{
				MarkdownDescription: "Whether the user account is activated.",
				Computed:            true,
			},
			"allow_public_key": schema.BoolAttribute{
				MarkdownDescription: "Whether the user is allowed to use public key authentication.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"sudo": schema.BoolAttribute{
				MarkdownDescription: "Whether the user has sudo privileges.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"enable_managed_uid": schema.BoolAttribute{
				MarkdownDescription: "Whether managed UID is enabled for the user.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"unix_uid": schema.Int64Attribute{
				MarkdownDescription: "The UNIX UID for the user.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"unix_guid": schema.Int64Attribute{
				MarkdownDescription: "The UNIX GID for the user.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"passwordless_sudo": schema.BoolAttribute{
				MarkdownDescription: "Whether the user can use sudo without a password.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"externally_managed": schema.BoolAttribute{
				MarkdownDescription: "Whether the user is externally managed.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"ldap_binding_user": schema.BoolAttribute{
				MarkdownDescription: "Whether the user is an LDAP binding user.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"enable_user_portal_multifactor": schema.BoolAttribute{
				MarkdownDescription: "Whether multi-factor authentication is enabled for the user portal.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"suspended": schema.BoolAttribute{
				MarkdownDescription: "Whether the user account is suspended.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "The state of the user account (e.g., `ACTIVATED`, `STAGED`).",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the user.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"department": schema.StringAttribute{
				MarkdownDescription: "The department the user belongs to.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cost_center": schema.StringAttribute{
				MarkdownDescription: "The cost center the user is associated with.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"company": schema.StringAttribute{
				MarkdownDescription: "The company the user is associated with.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"employee_type": schema.StringAttribute{
				MarkdownDescription: "The type of employee (e.g., full-time, contractor).",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"employee_identifier": schema.StringAttribute{
				MarkdownDescription: "The employee identifier or number.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"job_title": schema.StringAttribute{
				MarkdownDescription: "The job title of the user.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"location": schema.StringAttribute{
				MarkdownDescription: "The location of the user.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"manager": schema.StringAttribute{
				MarkdownDescription: "The manager of the user (user ID or email).",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"middle_name": schema.StringAttribute{
				MarkdownDescription: "The middle name of the user.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the user.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"external_dn": schema.StringAttribute{
				MarkdownDescription: "The external distinguished name, set by directory sync.",
				Computed:            true,
			},
			"external_source_type": schema.StringAttribute{
				MarkdownDescription: "The external source type, set by directory sync.",
				Computed:            true,
			},
			"mfa_configured": schema.BoolAttribute{
				MarkdownDescription: "Whether MFA has been configured for the user.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"mfa_exclusion": schema.BoolAttribute{
				MarkdownDescription: "Whether the user is excluded from MFA requirements.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *UserResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan UserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiUser := userModelToAPI(plan)

	user, err := r.client.CreateUser(ctx, apiUser)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating JumpCloud User",
			fmt.Sprintf("Could not create user %q: %s", plan.Username.ValueString(), err),
		)
		return
	}

	state := userAPIToModel(user, plan)

	tflog.Trace(ctx, "created JumpCloud user", map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state UserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.client.GetUser(ctx, state.ID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			tflog.Warn(ctx, "JumpCloud user not found, removing from state", map[string]interface{}{
				"id": state.ID.ValueString(),
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading JumpCloud User",
			fmt.Sprintf("Could not read user %q: %s", state.ID.ValueString(), err),
		)
		return
	}

	newState := userAPIToModel(user, state)

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan UserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state UserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiUser := userModelToAPI(plan)

	user, err := r.client.UpdateUser(ctx, state.ID.ValueString(), apiUser)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating JumpCloud User",
			fmt.Sprintf("Could not update user %q: %s", state.ID.ValueString(), err),
		)
		return
	}

	newState := userAPIToModel(user, plan)

	tflog.Trace(ctx, "updated JumpCloud user", map[string]interface{}{
		"id": user.ID,
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state UserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteUser(ctx, state.ID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting JumpCloud User",
			fmt.Sprintf("Could not delete user %q: %s", state.ID.ValueString(), err),
		)
	}

	tflog.Trace(ctx, "deleted JumpCloud user", map[string]interface{}{
		"id": state.ID.ValueString(),
	})
}

func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// userModelToAPI converts a Terraform resource model to a JumpCloud API struct.
func userModelToAPI(model UserResourceModel) client.SystemUser {
	user := client.SystemUser{
		Username: model.Username.ValueString(),
		Email:    model.Email.ValueString(),
	}

	if !model.Firstname.IsNull() && !model.Firstname.IsUnknown() {
		user.Firstname = model.Firstname.ValueString()
	}
	if !model.Lastname.IsNull() && !model.Lastname.IsUnknown() {
		user.Lastname = model.Lastname.ValueString()
	}
	if !model.Password.IsNull() && !model.Password.IsUnknown() {
		user.Password = model.Password.ValueString()
	}
	if !model.Description.IsNull() && !model.Description.IsUnknown() {
		user.Description = model.Description.ValueString()
	}
	if !model.Department.IsNull() && !model.Department.IsUnknown() {
		user.Department = model.Department.ValueString()
	}
	if !model.CostCenter.IsNull() && !model.CostCenter.IsUnknown() {
		user.CostCenter = model.CostCenter.ValueString()
	}
	if !model.Company.IsNull() && !model.Company.IsUnknown() {
		user.Company = model.Company.ValueString()
	}
	if !model.EmployeeType.IsNull() && !model.EmployeeType.IsUnknown() {
		user.EmployeeType = model.EmployeeType.ValueString()
	}
	if !model.EmployeeIdentifier.IsNull() && !model.EmployeeIdentifier.IsUnknown() {
		user.EmployeeIdentifier = model.EmployeeIdentifier.ValueString()
	}
	if !model.JobTitle.IsNull() && !model.JobTitle.IsUnknown() {
		user.JobTitle = model.JobTitle.ValueString()
	}
	if !model.Location.IsNull() && !model.Location.IsUnknown() {
		user.Location = model.Location.ValueString()
	}
	if !model.Manager.IsNull() && !model.Manager.IsUnknown() {
		user.Manager = model.Manager.ValueString()
	}
	if !model.MiddleName.IsNull() && !model.MiddleName.IsUnknown() {
		user.MiddleName = model.MiddleName.ValueString()
	}
	if !model.DisplayName.IsNull() && !model.DisplayName.IsUnknown() {
		user.DisplayName = model.DisplayName.ValueString()
	}

	// Boolean fields
	if !model.AccountLocked.IsNull() && !model.AccountLocked.IsUnknown() {
		v := model.AccountLocked.ValueBool()
		user.AccountLocked = &v
	}
	if !model.AllowPublicKey.IsNull() && !model.AllowPublicKey.IsUnknown() {
		v := model.AllowPublicKey.ValueBool()
		user.AllowPublicKey = &v
	}
	if !model.Sudo.IsNull() && !model.Sudo.IsUnknown() {
		v := model.Sudo.ValueBool()
		user.Sudo = &v
	}
	if !model.EnableManagedUID.IsNull() && !model.EnableManagedUID.IsUnknown() {
		v := model.EnableManagedUID.ValueBool()
		user.EnableManagedUID = &v
	}
	if !model.PasswordlessSudo.IsNull() && !model.PasswordlessSudo.IsUnknown() {
		v := model.PasswordlessSudo.ValueBool()
		user.PasswordlessSudo = &v
	}
	if !model.ExternallyManaged.IsNull() && !model.ExternallyManaged.IsUnknown() {
		v := model.ExternallyManaged.ValueBool()
		user.ExternallyManaged = &v
	}
	if !model.LdapBindingUser.IsNull() && !model.LdapBindingUser.IsUnknown() {
		v := model.LdapBindingUser.ValueBool()
		user.LdapBindingUser = &v
	}
	if !model.EnableUserPortalMultifactor.IsNull() && !model.EnableUserPortalMultifactor.IsUnknown() {
		v := model.EnableUserPortalMultifactor.ValueBool()
		user.EnableUserPortalMultifactor = &v
	}
	if !model.Suspended.IsNull() && !model.Suspended.IsUnknown() {
		v := model.Suspended.ValueBool()
		user.Suspended = &v
	}

	// Integer fields
	if !model.UnixUID.IsNull() && !model.UnixUID.IsUnknown() {
		v := int(model.UnixUID.ValueInt64())
		user.UnixUID = &v
	}
	if !model.UnixGUID.IsNull() && !model.UnixGUID.IsUnknown() {
		v := int(model.UnixGUID.ValueInt64())
		user.UnixGUID = &v
	}

	// MFA
	if (!model.MFAConfigured.IsNull() && !model.MFAConfigured.IsUnknown()) ||
		(!model.MFAExclusion.IsNull() && !model.MFAExclusion.IsUnknown()) {
		user.MFA = &client.MFA{}
		if !model.MFAConfigured.IsNull() && !model.MFAConfigured.IsUnknown() {
			v := model.MFAConfigured.ValueBool()
			user.MFA.Configured = &v
		}
		if !model.MFAExclusion.IsNull() && !model.MFAExclusion.IsUnknown() {
			v := model.MFAExclusion.ValueBool()
			user.MFA.Exclusion = &v
		}
	}

	return user
}

// userAPIToModel converts a JumpCloud API response to a Terraform resource model.
// The priorState is used to preserve write-only fields (like password) that the API does not return.
func userAPIToModel(user *client.SystemUser, priorState UserResourceModel) UserResourceModel {
	model := UserResourceModel{
		ID:                 types.StringValue(user.ID),
		Username:           types.StringValue(user.Username),
		Email:              types.StringValue(user.Email),
		Firstname:          types.StringValue(user.Firstname),
		Lastname:           types.StringValue(user.Lastname),
		Activated:          types.BoolValue(user.Activated),
		State:              types.StringValue(user.State),
		Description:        types.StringValue(user.Description),
		Department:         types.StringValue(user.Department),
		CostCenter:         types.StringValue(user.CostCenter),
		Company:            types.StringValue(user.Company),
		EmployeeType:       types.StringValue(user.EmployeeType),
		EmployeeIdentifier: types.StringValue(user.EmployeeIdentifier),
		JobTitle:           types.StringValue(user.JobTitle),
		Location:           types.StringValue(user.Location),
		Manager:            types.StringValue(user.Manager),
		MiddleName:         types.StringValue(user.MiddleName),
		DisplayName:        types.StringValue(user.DisplayName),
		ExternalDN:         types.StringValue(user.ExternalDN),
		ExternalSourceType: types.StringValue(user.ExternalSourceType),
	}

	// Preserve password from prior state (API never returns it)
	model.Password = priorState.Password

	// Boolean pointer fields
	if user.AccountLocked != nil {
		model.AccountLocked = types.BoolValue(*user.AccountLocked)
	} else {
		model.AccountLocked = types.BoolValue(false)
	}
	if user.AllowPublicKey != nil {
		model.AllowPublicKey = types.BoolValue(*user.AllowPublicKey)
	} else {
		model.AllowPublicKey = types.BoolValue(false)
	}
	if user.Sudo != nil {
		model.Sudo = types.BoolValue(*user.Sudo)
	} else {
		model.Sudo = types.BoolValue(false)
	}
	if user.EnableManagedUID != nil {
		model.EnableManagedUID = types.BoolValue(*user.EnableManagedUID)
	} else {
		model.EnableManagedUID = types.BoolValue(false)
	}
	if user.PasswordlessSudo != nil {
		model.PasswordlessSudo = types.BoolValue(*user.PasswordlessSudo)
	} else {
		model.PasswordlessSudo = types.BoolValue(false)
	}
	if user.ExternallyManaged != nil {
		model.ExternallyManaged = types.BoolValue(*user.ExternallyManaged)
	} else {
		model.ExternallyManaged = types.BoolValue(false)
	}
	if user.LdapBindingUser != nil {
		model.LdapBindingUser = types.BoolValue(*user.LdapBindingUser)
	} else {
		model.LdapBindingUser = types.BoolValue(false)
	}
	if user.EnableUserPortalMultifactor != nil {
		model.EnableUserPortalMultifactor = types.BoolValue(*user.EnableUserPortalMultifactor)
	} else {
		model.EnableUserPortalMultifactor = types.BoolValue(false)
	}
	if user.Suspended != nil {
		model.Suspended = types.BoolValue(*user.Suspended)
	} else {
		model.Suspended = types.BoolValue(false)
	}

	// Integer pointer fields
	if user.UnixUID != nil {
		model.UnixUID = types.Int64Value(int64(*user.UnixUID))
	} else {
		model.UnixUID = types.Int64Value(0)
	}
	if user.UnixGUID != nil {
		model.UnixGUID = types.Int64Value(int64(*user.UnixGUID))
	} else {
		model.UnixGUID = types.Int64Value(0)
	}

	// MFA fields
	if user.MFA != nil {
		if user.MFA.Configured != nil {
			model.MFAConfigured = types.BoolValue(*user.MFA.Configured)
		} else {
			model.MFAConfigured = types.BoolValue(false)
		}
		if user.MFA.Exclusion != nil {
			model.MFAExclusion = types.BoolValue(*user.MFA.Exclusion)
		} else {
			model.MFAExclusion = types.BoolValue(false)
		}
	} else {
		model.MFAConfigured = types.BoolValue(false)
		model.MFAExclusion = types.BoolValue(false)
	}

	return model
}
