package cdfmc_test

import (
	"testing"

	"github.com/CiscoDevnet/terraform-provider-sccfm/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testCdFmc = struct {
	Hostname string
}{
	Hostname: acctest.Env.CdFmcDataSourceHostname(),
}

const testCdFmcTemplate = `
data "sccfm_cdfmc" "test" {}`

var testCdfmcConfig = acctest.MustParseTemplate(testCdFmcTemplate, testCdFmc)

func TestAccCdFmcDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: acctest.ProviderConfig() + testCdfmcConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.sccfm_cdfmc.test", "hostname", testCdFmc.Hostname),
				),
			},
		},
	})
}
