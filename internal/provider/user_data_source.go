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

var _ datasource.DataSource = &UserDataSource{}

// NewUserDataSource returns a new user data source.
func NewUserDataSource() datasource.DataSource {
	return &UserDataSource{}
}

// UserDataSource defines the user data source implementation.
type UserDataSource struct {
	client *client.JumpCloudClient
}

// UserDataSourceModel describes the user data source data model.
type UserDataSourceModel struct {
	ID                          types.String `tfsdk:"id"`
	Username                    types.String `tfsdk:"username"`
	Email                       types.String `tfsdk:"email"`
	Firstname                   types.String `tfsdk:"firstname"`
	Lastname                    types.String `tfsdk:"lastname"`
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

func (d *UserDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (d *UserDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to look up an existing JumpCloud user by email or username.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the user.",
				Computed:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "The username to look up. At least one of `username` or `email` must be specified.",
				Optional:            true,
				Computed:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "The email address to look up. At least one of `username` or `email` must be specified.",
				Optional:            true,
				Computed:            true,
			},
			"firstname": schema.StringAttribute{
				MarkdownDescription: "The first name of the user.",
				Computed:            true,
			},
			"lastname": schema.StringAttribute{
				MarkdownDescription: "The last name of the user.",
				Computed:            true,
			},
			"account_locked": schema.BoolAttribute{
				MarkdownDescription: "Whether the user account is locked.",
				Computed:            true,
			},
			"activated": schema.BoolAttribute{
				MarkdownDescription: "Whether the user account is activated.",
				Computed:            true,
			},
			"allow_public_key": schema.BoolAttribute{
				MarkdownDescription: "Whether the user is allowed to use public key authentication.",
				Computed:            true,
			},
			"sudo": schema.BoolAttribute{
				MarkdownDescription: "Whether the user has sudo privileges.",
				Computed:            true,
			},
			"enable_managed_uid": schema.BoolAttribute{
				MarkdownDescription: "Whether managed UID is enabled for the user.",
				Computed:            true,
			},
			"unix_uid": schema.Int64Attribute{
				MarkdownDescription: "The UNIX UID for the user.",
				Computed:            true,
			},
			"unix_guid": schema.Int64Attribute{
				MarkdownDescription: "The UNIX GID for the user.",
				Computed:            true,
			},
			"passwordless_sudo": schema.BoolAttribute{
				MarkdownDescription: "Whether the user can use sudo without a password.",
				Computed:            true,
			},
			"externally_managed": schema.BoolAttribute{
				MarkdownDescription: "Whether the user is externally managed.",
				Computed:            true,
			},
			"ldap_binding_user": schema.BoolAttribute{
				MarkdownDescription: "Whether the user is an LDAP binding user.",
				Computed:            true,
			},
			"enable_user_portal_multifactor": schema.BoolAttribute{
				MarkdownDescription: "Whether multi-factor authentication is enabled for the user portal.",
				Computed:            true,
			},
			"suspended": schema.BoolAttribute{
				MarkdownDescription: "Whether the user account is suspended.",
				Computed:            true,
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "The state of the user account.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the user.",
				Computed:            true,
			},
			"department": schema.StringAttribute{
				MarkdownDescription: "The department the user belongs to.",
				Computed:            true,
			},
			"cost_center": schema.StringAttribute{
				MarkdownDescription: "The cost center the user is associated with.",
				Computed:            true,
			},
			"company": schema.StringAttribute{
				MarkdownDescription: "The company the user is associated with.",
				Computed:            true,
			},
			"employee_type": schema.StringAttribute{
				MarkdownDescription: "The type of employee.",
				Computed:            true,
			},
			"employee_identifier": schema.StringAttribute{
				MarkdownDescription: "The employee identifier or number.",
				Computed:            true,
			},
			"job_title": schema.StringAttribute{
				MarkdownDescription: "The job title of the user.",
				Computed:            true,
			},
			"location": schema.StringAttribute{
				MarkdownDescription: "The location of the user.",
				Computed:            true,
			},
			"manager": schema.StringAttribute{
				MarkdownDescription: "The manager of the user.",
				Computed:            true,
			},
			"middle_name": schema.StringAttribute{
				MarkdownDescription: "The middle name of the user.",
				Computed:            true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the user.",
				Computed:            true,
			},
			"external_dn": schema.StringAttribute{
				MarkdownDescription: "The external distinguished name.",
				Computed:            true,
			},
			"external_source_type": schema.StringAttribute{
				MarkdownDescription: "The external source type.",
				Computed:            true,
			},
			"mfa_configured": schema.BoolAttribute{
				MarkdownDescription: "Whether MFA has been configured for the user.",
				Computed:            true,
			},
			"mfa_exclusion": schema.BoolAttribute{
				MarkdownDescription: "Whether the user is excluded from MFA requirements.",
				Computed:            true,
			},
		},
	}
}

func (d *UserDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *UserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config UserDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var user *client.SystemUser
	var err error

	if !config.Email.IsNull() && config.Email.ValueString() != "" {
		user, err = d.client.GetUserByEmail(ctx, config.Email.ValueString())
	} else if !config.Username.IsNull() && config.Username.ValueString() != "" {
		user, err = d.client.GetUserByUsername(ctx, config.Username.ValueString())
	} else {
		resp.Diagnostics.AddError(
			"Missing Search Criteria",
			"Either `email` or `username` must be specified to look up a user.",
		)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading JumpCloud User",
			fmt.Sprintf("Could not read user: %s", err),
		)
		return
	}

	state := userAPIToDataSourceModel(user)

	tflog.Trace(ctx, "read JumpCloud user data source", map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// userAPIToDataSourceModel converts a JumpCloud API response to a data source model.
func userAPIToDataSourceModel(user *client.SystemUser) UserDataSourceModel {
	model := UserDataSourceModel{
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

	// Boolean pointer fields
	model.AccountLocked = types.BoolValue(boolPtrValue(user.AccountLocked))
	model.AllowPublicKey = types.BoolValue(boolPtrValue(user.AllowPublicKey))
	model.Sudo = types.BoolValue(boolPtrValue(user.Sudo))
	model.EnableManagedUID = types.BoolValue(boolPtrValue(user.EnableManagedUID))
	model.PasswordlessSudo = types.BoolValue(boolPtrValue(user.PasswordlessSudo))
	model.ExternallyManaged = types.BoolValue(boolPtrValue(user.ExternallyManaged))
	model.LdapBindingUser = types.BoolValue(boolPtrValue(user.LdapBindingUser))
	model.EnableUserPortalMultifactor = types.BoolValue(boolPtrValue(user.EnableUserPortalMultifactor))
	model.Suspended = types.BoolValue(boolPtrValue(user.Suspended))

	// Integer pointer fields
	model.UnixUID = types.Int64Value(intPtrValue(user.UnixUID))
	model.UnixGUID = types.Int64Value(intPtrValue(user.UnixGUID))

	// MFA fields
	if user.MFA != nil {
		model.MFAConfigured = types.BoolValue(boolPtrValue(user.MFA.Configured))
		model.MFAExclusion = types.BoolValue(boolPtrValue(user.MFA.Exclusion))
	} else {
		model.MFAConfigured = types.BoolValue(false)
		model.MFAExclusion = types.BoolValue(false)
	}

	return model
}

func boolPtrValue(v *bool) bool {
	if v == nil {
		return false
	}
	return *v
}

func intPtrValue(v *int) int64 {
	if v == nil {
		return 0
	}
	return int64(*v)
}
