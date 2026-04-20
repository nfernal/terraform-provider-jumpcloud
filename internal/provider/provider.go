// Copyright (c) Nfernal
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/nfernal/terraform-provider-jumpcloud/internal/client"
)

var _ provider.Provider = &JumpCloudProvider{}

// JumpCloudProvider defines the provider implementation.
type JumpCloudProvider struct {
	version string
}

// JumpCloudProviderModel describes the provider data model.
type JumpCloudProviderModel struct {
	APIKey types.String `tfsdk:"api_key"`
	OrgID  types.String `tfsdk:"org_id"`
	APIURL types.String `tfsdk:"api_url"`
}

// New returns a new provider factory function.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &JumpCloudProvider{
			version: version,
		}
	}
}

func (p *JumpCloudProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "jumpcloud"
	resp.Version = p.version
}

func (p *JumpCloudProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The JumpCloud provider allows you to manage JumpCloud resources such as users, user groups, system groups, and group memberships.",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "JumpCloud API key. Can also be set with the `JUMPCLOUD_API_KEY` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"org_id": schema.StringAttribute{
				MarkdownDescription: "JumpCloud organization ID for multi-tenant environments. Can also be set with the `JUMPCLOUD_ORG_ID` environment variable.",
				Optional:            true,
			},
			"api_url": schema.StringAttribute{
				MarkdownDescription: "JumpCloud API base URL. Defaults to `https://console.jumpcloud.com/api`. Can also be set with the `JUMPCLOUD_API_URL` environment variable.",
				Optional:            true,
			},
		},
	}
}

func (p *JumpCloudProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config JumpCloudProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check for unknown values
	if config.APIKey.IsUnknown() {
		resp.Diagnostics.AddError(
			"Unknown JumpCloud API Key",
			"The provider cannot create the JumpCloud API client as there is an unknown configuration value for the API key. "+
				"Set the value statically in the configuration or use the JUMPCLOUD_API_KEY environment variable.",
		)
	}
	if config.OrgID.IsUnknown() {
		resp.Diagnostics.AddError(
			"Unknown JumpCloud Org ID",
			"The provider cannot create the JumpCloud API client as there is an unknown configuration value for the org ID. "+
				"Set the value statically in the configuration or use the JUMPCLOUD_ORG_ID environment variable.",
		)
	}
	if config.APIURL.IsUnknown() {
		resp.Diagnostics.AddError(
			"Unknown JumpCloud API URL",
			"The provider cannot create the JumpCloud API client as there is an unknown configuration value for the API URL. "+
				"Set the value statically in the configuration or use the JUMPCLOUD_API_URL environment variable.",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Resolve with environment variable fallbacks
	apiKey := config.APIKey.ValueString()
	if apiKey == "" {
		apiKey = os.Getenv("JUMPCLOUD_API_KEY")
	}
	if apiKey == "" {
		resp.Diagnostics.AddError(
			"Missing JumpCloud API Key",
			"The provider requires a JumpCloud API key. Set it in the provider configuration or via the JUMPCLOUD_API_KEY environment variable.",
		)
		return
	}

	orgID := config.OrgID.ValueString()
	if orgID == "" {
		orgID = os.Getenv("JUMPCLOUD_ORG_ID")
	}

	apiURL := config.APIURL.ValueString()
	if apiURL == "" {
		apiURL = os.Getenv("JUMPCLOUD_API_URL")
	}

	jcClient := client.NewJumpCloudClient(apiKey, orgID, apiURL, p.version)

	resp.DataSourceData = jcClient
	resp.ResourceData = jcClient
}

func (p *JumpCloudProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewUserResource,
		NewUserGroupResource,
		NewUserGroupMembershipResource,
		NewSystemGroupResource,
	}
}

func (p *JumpCloudProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewUserDataSource,
		NewUserGroupDataSource,
	}
}
