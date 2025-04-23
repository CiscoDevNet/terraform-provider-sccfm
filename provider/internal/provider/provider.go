// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-sccfm/internal/device/ftd/ftdversion"
	"github.com/CiscoDevnet/terraform-provider-sccfm/internal/msp/msp_tenant"
	"github.com/CiscoDevnet/terraform-provider-sccfm/internal/msp/msp_tenant_user_api_token"
	"github.com/CiscoDevnet/terraform-provider-sccfm/internal/msp/msp_tenant_users"
	"os"

	"github.com/CiscoDevnet/terraform-provider-sccfm/internal/connector"
	"github.com/CiscoDevnet/terraform-provider-sccfm/internal/connector/connectoronboarding"
	"github.com/CiscoDevnet/terraform-provider-sccfm/internal/connector/sec"
	"github.com/CiscoDevnet/terraform-provider-sccfm/internal/connector/sec/seconboarding"
	"github.com/CiscoDevnet/terraform-provider-sccfm/internal/service/duoadminpanel"
	"github.com/CiscoDevnet/terraform-provider-sccfm/internal/tenantsettings"

	"github.com/CiscoDevnet/terraform-provider-sccfm/internal/cdfmc"
	"github.com/CiscoDevnet/terraform-provider-sccfm/internal/device/ftd/ftdonboarding"

	"github.com/CiscoDevnet/terraform-provider-sccfm/internal/device/ftd"
	"github.com/CiscoDevnet/terraform-provider-sccfm/internal/tenant"

	"github.com/CiscoDevnet/terraform-provider-sccfm/internal/user"
	"github.com/CiscoDevnet/terraform-provider-sccfm/internal/user_api_token"

	"github.com/CiscoDevnet/terraform-provider-sccfm/internal/device/ios"

	sccFwMgrClient "github.com/CiscoDevnet/terraform-provider-sccfm/go-client"
	"github.com/CiscoDevnet/terraform-provider-sccfm/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/CiscoDevnet/terraform-provider-sccfm/internal/device/asa"
)

var _ provider.Provider = &SccFirewallManagerProvider{}

type SccFirewallManagerProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// SccFirewallManagerProviderModel describes the provider data model.
type SccFirewallManagerProviderModel struct {
	ApiToken types.String `tfsdk:"api_token"`
	BaseURL  types.String `tfsdk:"base_url"`
}

func (p *SccFirewallManagerProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sccfm"
	resp.Version = p.version
}

func (p *SccFirewallManagerProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use the Cisco Security Cloud Control Firewall Manager (SCC Firewall Manager) provider to onboard and manage the many devices and other resources supported by SCC Firewall Manager. You must configure the provider with the proper credentials and region before you can use it.",
		Attributes: map[string]schema.Attribute{
			"api_token": schema.StringAttribute{
				MarkdownDescription: "The API token used to authenticate with SCC Firewall Manager. [See here](https://docs.manage.security.cisco.com/c_api-tokens.html#!t-generatean-api-token.html) to learn how to generate an API token.",
				Optional:            true,
				Sensitive:           true,
				Validators: []validator.String{
					validators.OneOfRoles("ROLE_SUPER_ADMIN", "ROLE_ADMIN"),
				},
			},
			"base_url": schema.StringAttribute{
				MarkdownDescription: "The base SCC Firewall Manager URL. This is the URL you enter when logging into your SCC Firewall Manager account.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("https://www.defenseorchestrator.com", "https://us.manage.security.cisco.com", "https://www.defenseorchestrator.eu", "https://eu.manage.security.cisco.com", "https://apj.cdo.cisco.com", "https://apj.manage.security.cisco.com", "https://staging.dev.lockhart.io", "https://staging.manage.security.cisco.com", "https://ci.dev.lockhart.io", "https://ci.manage.security.cisco.com", "https://scale.dev.lockhart.io", "https://scale.manage.security.cisco.com", "http://localhost:9000", "https://aus.cdo.cisco.com", "https://aus.manage.security.cisco.com", "https://in.cdo.cisco.com", "https://aus.manage.security.cisco.com"),
				},
			},
		},
	}
}

func (p *SccFirewallManagerProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data SccFirewallManagerProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.ApiToken.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			"Unknown Cisco SCC Firewall Manager Token",
			"The provider cannot create the Cisco SCC Firewall Manager client as there is an unknown configuration value for the Cisco SCC Firewall Manager token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the CISCO_SCC Firewall Manager_API_TOKEN environment variable.",
		)
	}

	if data.BaseURL.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("base_url"),
			"Unknown Cisco SCC Firewall Manager Base URL",
			"The provider cannot create the Cisco SCC Firewall Manager client as there is an unknown configuration value for the Cisco SCC Firewall Manager Base URL. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the CISCO_SCC Firewall Manager_BASE_URL environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	apiToken := os.Getenv("CISCO_SCC Firewall Manager_API_TOKEN")
	if !data.ApiToken.IsNull() {
		apiToken = data.ApiToken.ValueString()
	}

	if apiToken == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			"Missing Cisco SCC Firewall Manager Token",
			"The provider cannot create the Cisco SCC Firewall Manager client as there is a missing or empty value for the Cisco SCC Firewall Manager token. "+
				"Set the API token value in the configuration or use the CISCO_SCC Firewall Manager_API_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	baseURL := os.Getenv("CISCO_SCC Firewall Manager_BASE_URL")
	if !data.BaseURL.IsNull() {
		baseURL = data.BaseURL.ValueString()
	}
	if baseURL == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("base_url"),
			"Missing Cisco SCC Firewall Manager Base URL",
			"The provider cannot create the Cisco SCC Firewall Manager client as there is a missing or empty value for the Cisco SCC Firewall Manager base URL. "+
				"Set the API token value in the configuration or use the CISCO_SCC Firewall Manager_BASE_URL environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := sccFwMgrClient.New(baseURL, apiToken)
	if err != nil {
		resp.Diagnostics.AddError("Error while trying to create SCC Firewall Manager client", fmt.Sprintf("cause=%s", err.Error()))
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *SccFirewallManagerProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		connector.NewResource,
		asa.NewAsaDeviceResource,
		ios.NewIosDeviceResource,
		ftd.NewResource,
		user.NewResource,
		user_api_token.NewResource,
		ftdonboarding.NewResource,
		connectoronboarding.NewResource,
		cdfmc.NewResource,
		sec.NewResource,
		seconboarding.NewResource,
		duoadminpanel.NewResource,
		tenantsettings.NewTenantSettingsResource,
		msp_tenant.NewTenantResource,
		msp_tenant_users.NewMspManagedTenantUsersResource,
		msp_tenant_user_api_token.NewMspManagedTenantUserApiTokenResource,
		ftdversion.NewResource,
	}
}

func (p *SccFirewallManagerProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		connector.NewDataSource,
		asa.NewAsaDataSource,
		ios.NewIosDataSource,
		ftd.NewDataSource,
		user.NewDataSource,
		tenant.NewDataSource,
		cdfmc.NewDataSource,
		tenantsettings.NewTenantSettingsDataSource,
		msp_tenant.NewTenantDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SccFirewallManagerProvider{
			version: version,
		}
	}
}
