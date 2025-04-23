package tenantsettings_test

import (
	"testing"

	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model/settings"
	"github.com/CiscoDevnet/terraform-provider-sccfm/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testTenantSettingsDataSourceConfig = `data "sccfm_tenant_settings" "test" {}`

func TestAccTenantSettingsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: acctest.ProviderConfig() + testTenantSettingsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.sccfm_tenant_settings.test", "id", acctest.Env.TenantSettingsTenantUid()),
					resource.TestCheckResourceAttr("data.sccfm_tenant_settings.test", "change_request_support_enabled", "false"),
					resource.TestCheckResourceAttr("data.sccfm_tenant_settings.test", "auto_accept_device_changes_enabled", "false"),
					resource.TestCheckResourceAttr("data.sccfm_tenant_settings.test", "web_analytics_enabled", "false"),
					resource.TestCheckResourceAttr("data.sccfm_tenant_settings.test", "scheduled_deployments_enabled", "false"),
					resource.TestCheckResourceAttr("data.sccfm_tenant_settings.test", "deny_cisco_support_access_to_tenant_enabled", "false"),
					resource.TestCheckResourceAttr("data.sccfm_tenant_settings.test", "multi_cloud_defense_enabled", "false"),
					resource.TestCheckResourceAttr("data.sccfm_tenant_settings.test", "auto_discover_on_prem_fmcs_enabled", "false"),
					resource.TestCheckResourceAttr("data.sccfm_tenant_settings.test", "conflict_detection_interval", string(settings.ConflictDetectionIntervalEvery24Hours)),
				),
			},
		},
	})
}
