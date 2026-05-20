package networkobject_test

import (
	"testing"

	"github.com/CiscoDevnet/terraform-provider-sccfm/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testNetworkObjectDataSourceTemplate = `
resource "sccfm_network_object" "prereq" {
	name        = "tf-test-network-object-ds"
	description = "prereq for data source test"
	value       = "192.168.1.0/24"
}

data "sccfm_network_object" "test" {
	name = sccfm_network_object.prereq.name
}`

func TestAccNetworkObjectDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig() + testNetworkObjectDataSourceTemplate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.sccfm_network_object.test", "name", "tf-test-network-object-ds"),
					resource.TestCheckResourceAttr("data.sccfm_network_object.test", "description", "prereq for data source test"),
					resource.TestCheckResourceAttr("data.sccfm_network_object.test", "value", "192.168.1.0/24"),
					resource.TestCheckResourceAttrSet("data.sccfm_network_object.test", "id"),
				),
			},
		},
	})
}
