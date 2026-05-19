package urlobject_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/CiscoDevnet/terraform-provider-sccfm/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testUrlObjectResource_Create = `
resource "sccfm_url_object" "test" {
	name        = "tf-test-url-object"
	description = "created by acceptance test"
	url         = "https://www.example.com"
}`

const testUrlObjectResource_UpdateUrl = `
resource "sccfm_url_object" "test" {
	name        = "tf-test-url-object"
	description = "url updated"
	url         = "https://www.updated-example.com"
}`

const testUrlObjectResource_UpdateName = `
resource "sccfm_url_object" "test" {
	name        = "tf-test-url-object-renamed"
	description = "url updated"
	url         = "https://www.updated-example.com"
}`

const testUrlObjectResource_RemoveDescription = `
resource "sccfm_url_object" "test" {
	name = "tf-test-url-object-renamed"
	url  = "https://www.updated-example.com"
}`

func TestAccUrlObjectResource_FullLifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create
			{
				Config: acctest.ProviderConfig() + testUrlObjectResource_Create,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_url_object.test", "name", "tf-test-url-object"),
					resource.TestCheckResourceAttr("sccfm_url_object.test", "description", "created by acceptance test"),
					resource.TestCheckResourceAttr("sccfm_url_object.test", "url", "https://www.example.com"),
					resource.TestCheckResourceAttrSet("sccfm_url_object.test", "id"),
				),
			},
			// Update URL and description
			{
				Config: acctest.ProviderConfig() + testUrlObjectResource_UpdateUrl,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_url_object.test", "description", "url updated"),
					resource.TestCheckResourceAttr("sccfm_url_object.test", "url", "https://www.updated-example.com"),
				),
			},
			// Update name
			{
				Config: acctest.ProviderConfig() + testUrlObjectResource_UpdateName,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_url_object.test", "name", "tf-test-url-object-renamed"),
				),
			},
			// Remove description
			{
				Config: acctest.ProviderConfig() + testUrlObjectResource_RemoveDescription,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_url_object.test", "name", "tf-test-url-object-renamed"),
					resource.TestCheckResourceAttr("sccfm_url_object.test", "url", "https://www.updated-example.com"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testUrlObjectOverrideConfig_Create(targetId string) string {
	return fmt.Sprintf(`
resource "sccfm_url_object" "test_overrides" {
	name        = "tf-test-url-object-overrides"
	description = "with one override"
	url         = "https://www.example.com"

	overrides {
		target_id = %q
		url       = "https://override.example.com"
	}
}`, targetId)
}

func testUrlObjectOverrideConfig_ModifyOverride(targetId string) string {
	return fmt.Sprintf(`
resource "sccfm_url_object" "test_overrides" {
	name        = "tf-test-url-object-overrides"
	description = "override changed"
	url         = "https://www.example.com"

	overrides {
		target_id = %q
		url       = "https://changed-override.example.com"
	}
}`, targetId)
}

const testUrlObjectOverrideConfig_RemoveOverrides = `
resource "sccfm_url_object" "test_overrides" {
	name        = "tf-test-url-object-overrides"
	description = "overrides removed"
	url         = "https://www.example.com"
}`

func TestAccUrlObjectResource_OverrideLifecycle(t *testing.T) {
	targetId := os.Getenv("ACC_TEST_TARGET_ID_1")
	if targetId == "" {
		t.Skip("ACC_TEST_TARGET_ID_1 must be set for override tests")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with override
			{
				Config: acctest.ProviderConfig() + testUrlObjectOverrideConfig_Create(targetId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_url_object.test_overrides", "overrides.#", "1"),
					resource.TestCheckResourceAttr("sccfm_url_object.test_overrides", "overrides.0.url", "https://override.example.com"),
				),
			},
			// Modify override
			{
				Config: acctest.ProviderConfig() + testUrlObjectOverrideConfig_ModifyOverride(targetId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_url_object.test_overrides", "overrides.#", "1"),
					resource.TestCheckResourceAttr("sccfm_url_object.test_overrides", "overrides.0.url", "https://changed-override.example.com"),
				),
			},
			// Remove overrides
			{
				Config: acctest.ProviderConfig() + testUrlObjectOverrideConfig_RemoveOverrides,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_url_object.test_overrides", "overrides.#", "0"),
					resource.TestCheckResourceAttr("sccfm_url_object.test_overrides", "description", "overrides removed"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
