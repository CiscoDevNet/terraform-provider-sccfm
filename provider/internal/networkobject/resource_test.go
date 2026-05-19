package networkobject_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/CiscoDevnet/terraform-provider-sccfm/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Step 1: Create with single IP
const testNetworkObjectResource_Create = `
resource "sccfm_network_object" "test" {
	name        = "tf-test-network-object"
	description = "created by acceptance test"
	value       = "10.0.0.1"
}`

// Step 2: Update value to CIDR and change description
const testNetworkObjectResource_UpdateValue = `
resource "sccfm_network_object" "test" {
	name        = "tf-test-network-object"
	description = "updated value to CIDR"
	value       = "10.0.0.0/24"
}`

// Step 3: Update value to IP range
const testNetworkObjectResource_UpdateToRange = `
resource "sccfm_network_object" "test" {
	name        = "tf-test-network-object"
	description = "updated value to range"
	value       = "10.0.0.1-10.0.0.10"
}`

// Step 4: Update name
const testNetworkObjectResource_UpdateName = `
resource "sccfm_network_object" "test" {
	name        = "tf-test-network-object-renamed"
	description = "updated value to range"
	value       = "10.0.0.1-10.0.0.10"
}`

// Step 5: Remove description (set to empty)
const testNetworkObjectResource_RemoveDescription = `
resource "sccfm_network_object" "test" {
	name  = "tf-test-network-object-renamed"
	value = "10.0.0.1-10.0.0.10"
}`

func TestAccNetworkObjectResource_FullLifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with single IP
			{
				Config: acctest.ProviderConfig() + testNetworkObjectResource_Create,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_network_object.test", "name", "tf-test-network-object"),
					resource.TestCheckResourceAttr("sccfm_network_object.test", "description", "created by acceptance test"),
					resource.TestCheckResourceAttr("sccfm_network_object.test", "value", "10.0.0.1"),
					resource.TestCheckResourceAttrSet("sccfm_network_object.test", "id"),
				),
			},
			// Update value to CIDR
			{
				Config: acctest.ProviderConfig() + testNetworkObjectResource_UpdateValue,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_network_object.test", "description", "updated value to CIDR"),
					resource.TestCheckResourceAttr("sccfm_network_object.test", "value", "10.0.0.0/24"),
				),
			},
			// Update value to IP range
			{
				Config: acctest.ProviderConfig() + testNetworkObjectResource_UpdateToRange,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_network_object.test", "description", "updated value to range"),
					resource.TestCheckResourceAttr("sccfm_network_object.test", "value", "10.0.0.1-10.0.0.10"),
				),
			},
			// Update name
			{
				Config: acctest.ProviderConfig() + testNetworkObjectResource_UpdateName,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_network_object.test", "name", "tf-test-network-object-renamed"),
				),
			},
			// Remove description
			{
				Config: acctest.ProviderConfig() + testNetworkObjectResource_RemoveDescription,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_network_object.test", "name", "tf-test-network-object-renamed"),
					resource.TestCheckResourceAttr("sccfm_network_object.test", "value", "10.0.0.1-10.0.0.10"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testNetworkObjectOverrideConfig_Create(targetId string) string {
	return fmt.Sprintf(`
resource "sccfm_network_object" "test_overrides" {
	name        = "tf-test-network-object-overrides"
	description = "with one override"
	value       = "10.0.0.1"

	overrides {
		target_id = %q
		value     = "10.0.0.2"
	}
}`, targetId)
}

func testNetworkObjectOverrideConfig_ModifyOverride(targetId string) string {
	return fmt.Sprintf(`
resource "sccfm_network_object" "test_overrides" {
	name        = "tf-test-network-object-overrides"
	description = "override value changed"
	value       = "10.0.0.1"

	overrides {
		target_id = %q
		value     = "192.168.1.1"
	}
}`, targetId)
}

func testNetworkObjectOverrideConfig_AddSecondOverride(targetId1, targetId2 string) string {
	return fmt.Sprintf(`
resource "sccfm_network_object" "test_overrides" {
	name        = "tf-test-network-object-overrides"
	description = "two overrides"
	value       = "10.0.0.1"

	overrides {
		target_id = %q
		value     = "192.168.1.1"
	}
	overrides {
		target_id = %q
		value     = "172.16.0.1"
	}
}`, targetId1, targetId2)
}

const testNetworkObjectOverrideConfig_RemoveAllOverrides = `
resource "sccfm_network_object" "test_overrides" {
	name        = "tf-test-network-object-overrides"
	description = "overrides removed"
	value       = "10.0.0.1"
}`

func TestAccNetworkObjectResource_OverrideLifecycle(t *testing.T) {
	targetId1 := os.Getenv("ACC_TEST_TARGET_ID_1")
	targetId2 := os.Getenv("ACC_TEST_TARGET_ID_2")
	if targetId1 == "" || targetId2 == "" {
		t.Skip("ACC_TEST_TARGET_ID_1 and ACC_TEST_TARGET_ID_2 must be set for override tests")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with one override
			{
				Config: acctest.ProviderConfig() + testNetworkObjectOverrideConfig_Create(targetId1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_network_object.test_overrides", "name", "tf-test-network-object-overrides"),
					resource.TestCheckResourceAttr("sccfm_network_object.test_overrides", "value", "10.0.0.1"),
					resource.TestCheckResourceAttr("sccfm_network_object.test_overrides", "overrides.#", "1"),
					resource.TestCheckResourceAttr("sccfm_network_object.test_overrides", "overrides.0.target_id", targetId1),
					resource.TestCheckResourceAttr("sccfm_network_object.test_overrides", "overrides.0.value", "10.0.0.2"),
				),
			},
			// Modify the override value
			{
				Config: acctest.ProviderConfig() + testNetworkObjectOverrideConfig_ModifyOverride(targetId1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_network_object.test_overrides", "overrides.#", "1"),
					resource.TestCheckResourceAttr("sccfm_network_object.test_overrides", "overrides.0.value", "192.168.1.1"),
					resource.TestCheckResourceAttr("sccfm_network_object.test_overrides", "description", "override value changed"),
				),
			},
			// Add a second override
			{
				Config: acctest.ProviderConfig() + testNetworkObjectOverrideConfig_AddSecondOverride(targetId1, targetId2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_network_object.test_overrides", "overrides.#", "2"),
					resource.TestCheckResourceAttr("sccfm_network_object.test_overrides", "description", "two overrides"),
				),
			},
			// Remove all overrides
			{
				Config: acctest.ProviderConfig() + testNetworkObjectOverrideConfig_RemoveAllOverrides,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_network_object.test_overrides", "description", "overrides removed"),
					resource.TestCheckResourceAttr("sccfm_network_object.test_overrides", "overrides.#", "0"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
