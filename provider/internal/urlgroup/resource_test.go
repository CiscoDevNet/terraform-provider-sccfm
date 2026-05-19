package urlgroup_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/CiscoDevnet/terraform-provider-sccfm/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testUrlGroupResource_Create = `
resource "sccfm_url_group" "test" {
	name        = "tf-test-url-group"
	description = "created by acceptance test"
	values      = ["https://www.example.com", "https://www.test.com"]
}`

const testUrlGroupResource_AddMember = `
resource "sccfm_url_group" "test" {
	name        = "tf-test-url-group"
	description = "member added"
	values      = ["https://www.example.com", "https://www.test.com", "https://www.new.com"]
}`

const testUrlGroupResource_RemoveMember = `
resource "sccfm_url_group" "test" {
	name        = "tf-test-url-group"
	description = "member removed"
	values      = ["https://www.example.com"]
}`

const testUrlGroupResource_UpdateName = `
resource "sccfm_url_group" "test" {
	name        = "tf-test-url-group-renamed"
	description = "member removed"
	values      = ["https://www.example.com"]
}`

const testUrlGroupResource_RemoveDescription = `
resource "sccfm_url_group" "test" {
	name   = "tf-test-url-group-renamed"
	values = ["https://www.example.com"]
}`

func TestAccUrlGroupResource_FullLifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with 2 URLs
			{
				Config: acctest.ProviderConfig() + testUrlGroupResource_Create,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_url_group.test", "name", "tf-test-url-group"),
					resource.TestCheckResourceAttr("sccfm_url_group.test", "description", "created by acceptance test"),
					resource.TestCheckResourceAttr("sccfm_url_group.test", "values.#", "2"),
					resource.TestCheckResourceAttrSet("sccfm_url_group.test", "id"),
				),
			},
			// Add a URL
			{
				Config: acctest.ProviderConfig() + testUrlGroupResource_AddMember,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_url_group.test", "description", "member added"),
					resource.TestCheckResourceAttr("sccfm_url_group.test", "values.#", "3"),
				),
			},
			// Remove URLs (down to 1)
			{
				Config: acctest.ProviderConfig() + testUrlGroupResource_RemoveMember,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_url_group.test", "description", "member removed"),
					resource.TestCheckResourceAttr("sccfm_url_group.test", "values.#", "1"),
				),
			},
			// Update name
			{
				Config: acctest.ProviderConfig() + testUrlGroupResource_UpdateName,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_url_group.test", "name", "tf-test-url-group-renamed"),
				),
			},
			// Remove description
			{
				Config: acctest.ProviderConfig() + testUrlGroupResource_RemoveDescription,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_url_group.test", "name", "tf-test-url-group-renamed"),
					resource.TestCheckResourceAttr("sccfm_url_group.test", "values.#", "1"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testUrlGroupOverrideConfig_Create(targetId string) string {
	return fmt.Sprintf(`
resource "sccfm_url_group" "test_overrides" {
	name        = "tf-test-url-group-overrides"
	description = "with override"
	values      = ["https://www.example.com"]

	overrides {
		target_id = %q
		values    = ["https://override.example.com"]
	}
}`, targetId)
}

func testUrlGroupOverrideConfig_ModifyOverride(targetId string) string {
	return fmt.Sprintf(`
resource "sccfm_url_group" "test_overrides" {
	name        = "tf-test-url-group-overrides"
	description = "override modified"
	values      = ["https://www.example.com"]

	overrides {
		target_id = %q
		values    = ["https://changed-override.example.com", "https://second-override.example.com"]
	}
}`, targetId)
}

const testUrlGroupOverrideConfig_RemoveOverrides = `
resource "sccfm_url_group" "test_overrides" {
	name        = "tf-test-url-group-overrides"
	description = "overrides removed"
	values      = ["https://www.example.com"]
}`

func TestAccUrlGroupResource_OverrideLifecycle(t *testing.T) {
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
				Config: acctest.ProviderConfig() + testUrlGroupOverrideConfig_Create(targetId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_url_group.test_overrides", "overrides.#", "1"),
				),
			},
			// Modify override
			{
				Config: acctest.ProviderConfig() + testUrlGroupOverrideConfig_ModifyOverride(targetId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_url_group.test_overrides", "overrides.#", "1"),
					resource.TestCheckResourceAttr("sccfm_url_group.test_overrides", "description", "override modified"),
				),
			},
			// Remove overrides
			{
				Config: acctest.ProviderConfig() + testUrlGroupOverrideConfig_RemoveOverrides,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_url_group.test_overrides", "overrides.#", "0"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
