package serviceobject_test

import (
	"testing"

	"github.com/CiscoDevnet/terraform-provider-sccfm/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testServiceObjectDataSourceTemplate = `
resource "sccfm_service_object" "prereq" {
	name        = "tf-test-service-object-ds"
	description = "prereq for data source test"
	protocol    = "TCP"
	value       = "80"
}

data "sccfm_service_object" "test" {
	name = sccfm_service_object.prereq.name
}`

func TestAccServiceObjectDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig() + testServiceObjectDataSourceTemplate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.sccfm_service_object.test", "name", "tf-test-service-object-ds"),
					resource.TestCheckResourceAttr("data.sccfm_service_object.test", "protocol", "TCP"),
					resource.TestCheckResourceAttr("data.sccfm_service_object.test", "value", "80"),
					resource.TestCheckResourceAttrSet("data.sccfm_service_object.test", "id"),
				),
			},
		},
	})
}
