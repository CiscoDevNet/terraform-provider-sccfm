package msp_tenant_test

import (
	"github.com/CiscoDevnet/terraform-provider-scc-firewall-manager/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"strings"
	"testing"
)

var testMspTenantResource = struct {
	Name        string
	DisplayName string
	Id          string
	Region      string
	ApiToken    string
}{
	ApiToken:    acctest.Env.AddedMspManagedTenantApiToken(),
	Name:        acctest.Env.AddedMspManagedTenantName(),
	DisplayName: acctest.Env.AddedMspManagedTenantDisplayName(),
	Id:          acctest.Env.AddedMspManagedTenantId(),
	Region:      strings.ToUpper(acctest.Env.MspTenantRegion()),
}

const testMspTenantResourceTemplate = `
resource "sccfm_msp_managed_tenant" "test" {
	api_token = "{{.ApiToken}}"
}`

var testMspTenantResourceConfig = acctest.MustParseTemplate(testMspTenantResourceTemplate, testMspTenantResource)

func TestAccMspTenantResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: acctest.MspProviderConfig() + testMspTenantResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_msp_managed_tenant.test", "name", testMspTenantResource.Name),
					resource.TestCheckResourceAttr("sccfm_msp_managed_tenant.test", "display_name", testMspTenantResource.DisplayName),
					resource.TestCheckResourceAttr("sccfm_msp_managed_tenant.test", "id", testMspTenantResource.Id),
					resource.TestCheckResourceAttr("sccfm_msp_managed_tenant.test", "region", testMspTenantResource.Region),
				),
			},
		},
	})
}
