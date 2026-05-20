package servicegroup_test

import (
	"testing"

	"github.com/CiscoDevnet/terraform-provider-sccfm/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testServiceGroupDataSourceTemplate = `
resource "sccfm_service_group" "prereq" {
	name        = "tf-test-service-group-ds"
	description = "prereq for data source test"

	values = [
		{ protocol = "TCP", value = "80" },
	]
}

data "sccfm_service_group" "test" {
	name = sccfm_service_group.prereq.name
}`

func TestAccServiceGroupDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig() + testServiceGroupDataSourceTemplate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.sccfm_service_group.test", "name", "tf-test-service-group-ds"),
					resource.TestCheckResourceAttrSet("data.sccfm_service_group.test", "id"),
				),
			},
		},
	})
}
