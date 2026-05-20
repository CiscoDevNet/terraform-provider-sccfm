package urlobject_test

import (
	"testing"

	"github.com/CiscoDevnet/terraform-provider-sccfm/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testUrlObjectDataSourceTemplate = `
resource "sccfm_url_object" "prereq" {
	name        = "tf-test-url-object-ds"
	description = "prereq for data source test"
	url         = "https://www.example.com"
}

data "sccfm_url_object" "test" {
	name = sccfm_url_object.prereq.name
}`

func TestAccUrlObjectDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig() + testUrlObjectDataSourceTemplate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.sccfm_url_object.test", "name", "tf-test-url-object-ds"),
					resource.TestCheckResourceAttr("data.sccfm_url_object.test", "url", "https://www.example.com"),
					resource.TestCheckResourceAttrSet("data.sccfm_url_object.test", "id"),
				),
			},
		},
	})
}
