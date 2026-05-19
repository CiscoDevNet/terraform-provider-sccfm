package urlgroup_test

import (
	"testing"

	"github.com/CiscoDevnet/terraform-provider-sccfm/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testUrlGroupDataSourceTemplate = `
resource "sccfm_url_group" "prereq" {
	name        = "tf-test-url-group-ds"
	description = "prereq for data source test"
	values      = ["https://www.example.com"]
}

data "sccfm_url_group" "test" {
	name = sccfm_url_group.prereq.name
}`

func TestAccUrlGroupDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig() + testUrlGroupDataSourceTemplate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.sccfm_url_group.test", "name", "tf-test-url-group-ds"),
					resource.TestCheckResourceAttrSet("data.sccfm_url_group.test", "id"),
				),
			},
		},
	})
}
