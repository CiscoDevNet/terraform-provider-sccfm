package servicegroup_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/CiscoDevnet/terraform-provider-sccfm/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testServiceGroupResource_Create = `
resource "sccfm_service_group" "test" {
	name        = "tf-test-service-group"
	description = "created by acceptance test"

	values = [
		{ protocol = "TCP", value = "80" },
		{ protocol = "TCP", value = "443" },
	]
}`

const testServiceGroupResource_AddMember = `
resource "sccfm_service_group" "test" {
	name        = "tf-test-service-group"
	description = "member added"

	values = [
		{ protocol = "TCP", value = "80" },
		{ protocol = "TCP", value = "443" },
		{ protocol = "UDP", value = "53" },
	]
}`

const testServiceGroupResource_RemoveMember = `
resource "sccfm_service_group" "test" {
	name        = "tf-test-service-group"
	description = "member removed"

	values = [
		{ protocol = "TCP", value = "443" },
	]
}`

const testServiceGroupResource_UpdateName = `
resource "sccfm_service_group" "test" {
	name        = "tf-test-service-group-renamed"
	description = "member removed"

	values = [
		{ protocol = "TCP", value = "443" },
	]
}`

const testServiceGroupResource_RemoveDescription = `
resource "sccfm_service_group" "test" {
	name = "tf-test-service-group-renamed"

	values = [
		{ protocol = "TCP", value = "443" },
	]
}`

func TestAccServiceGroupResource_FullLifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with 2 services
			{
				Config: acctest.ProviderConfig() + testServiceGroupResource_Create,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_service_group.test", "name", "tf-test-service-group"),
					resource.TestCheckResourceAttr("sccfm_service_group.test", "description", "created by acceptance test"),
					resource.TestCheckResourceAttr("sccfm_service_group.test", "values.#", "2"),
					resource.TestCheckResourceAttrSet("sccfm_service_group.test", "id"),
				),
			},
			// Add a service (UDP 53)
			{
				Config: acctest.ProviderConfig() + testServiceGroupResource_AddMember,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_service_group.test", "description", "member added"),
					resource.TestCheckResourceAttr("sccfm_service_group.test", "values.#", "3"),
				),
			},
			// Remove services (down to 1)
			{
				Config: acctest.ProviderConfig() + testServiceGroupResource_RemoveMember,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_service_group.test", "description", "member removed"),
					resource.TestCheckResourceAttr("sccfm_service_group.test", "values.#", "1"),
				),
			},
			// Update name
			{
				Config: acctest.ProviderConfig() + testServiceGroupResource_UpdateName,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_service_group.test", "name", "tf-test-service-group-renamed"),
				),
			},
			// Remove description
			{
				Config: acctest.ProviderConfig() + testServiceGroupResource_RemoveDescription,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_service_group.test", "name", "tf-test-service-group-renamed"),
					resource.TestCheckResourceAttr("sccfm_service_group.test", "values.#", "1"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testServiceGroupOverrideConfig_Create(targetId string) string {
	return fmt.Sprintf(`
resource "sccfm_service_group" "test_overrides" {
	name        = "tf-test-service-group-overrides"
	description = "with override"

	values = [
		{ protocol = "TCP", value = "80" },
	]

	overrides = [
		{
			target_id = %q
			values = [
				{ protocol = "TCP", value = "8080" },
			]
		},
	]
}`, targetId)
}

func testServiceGroupOverrideConfig_ModifyOverride(targetId string) string {
	return fmt.Sprintf(`
resource "sccfm_service_group" "test_overrides" {
	name        = "tf-test-service-group-overrides"
	description = "override modified"

	values = [
		{ protocol = "TCP", value = "80" },
	]

	overrides = [
		{
			target_id = %q
			values = [
				{ protocol = "TCP", value = "9090" },
				{ protocol = "UDP", value = "53" },
			]
		},
	]
}`, targetId)
}

const testServiceGroupOverrideConfig_RemoveOverrides = `
resource "sccfm_service_group" "test_overrides" {
	name        = "tf-test-service-group-overrides"
	description = "overrides removed"

	values = [
		{ protocol = "TCP", value = "80" },
	]
}`

func TestAccServiceGroupResource_OverrideLifecycle(t *testing.T) {
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
				Config: acctest.ProviderConfig() + testServiceGroupOverrideConfig_Create(targetId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_service_group.test_overrides", "overrides.#", "1"),
				),
			},
			// Modify override
			{
				Config: acctest.ProviderConfig() + testServiceGroupOverrideConfig_ModifyOverride(targetId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_service_group.test_overrides", "overrides.#", "1"),
					resource.TestCheckResourceAttr("sccfm_service_group.test_overrides", "description", "override modified"),
				),
			},
			// Remove overrides
			{
				Config: acctest.ProviderConfig() + testServiceGroupOverrideConfig_RemoveOverrides,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_service_group.test_overrides", "overrides.#", "0"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
