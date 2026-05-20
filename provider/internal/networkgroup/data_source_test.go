package networkgroup_test

import (
	"testing"

	"github.com/CiscoDevnet/terraform-provider-sccfm/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testNetworkGroupDataSourceTemplate = `
resource "sccfm_network_group" "prereq" {
	name        = "tf-test-network-group-ds"
	description = "prereq for data source test"
	values      = ["10.0.0.0/24"]
}

data "sccfm_network_group" "test" {
	name = sccfm_network_group.prereq.name
}`

func TestAccNetworkGroupDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig() + testNetworkGroupDataSourceTemplate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.sccfm_network_group.test", "name", "tf-test-network-group-ds"),
					resource.TestCheckResourceAttrSet("data.sccfm_network_group.test", "id"),
				),
			},
		},
	})
}
