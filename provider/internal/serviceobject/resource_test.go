package serviceobject_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/CiscoDevnet/terraform-provider-sccfm/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testServiceObjectResource_Create = `
resource "sccfm_service_object" "test" {
	name        = "tf-test-service-object"
	description = "created by acceptance test"
	protocol    = "TCP"
	value       = "443"
}`

const testServiceObjectResource_UpdatePort = `
resource "sccfm_service_object" "test" {
	name        = "tf-test-service-object"
	description = "port updated"
	protocol    = "TCP"
	value       = "8443"
}`

const testServiceObjectResource_UpdateProtocol = `
resource "sccfm_service_object" "test" {
	name        = "tf-test-service-object"
	description = "protocol changed to UDP"
	protocol    = "UDP"
	value       = "53"
}`

const testServiceObjectResource_UpdateName = `
resource "sccfm_service_object" "test" {
	name        = "tf-test-service-object-renamed"
	description = "protocol changed to UDP"
	protocol    = "UDP"
	value       = "53"
}`

const testServiceObjectResource_PortRange = `
resource "sccfm_service_object" "test" {
	name        = "tf-test-service-object-renamed"
	description = "port range"
	protocol    = "TCP"
	value       = "8000-9000"
}`

func TestAccServiceObjectResource_FullLifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create TCP 443
			{
				Config: acctest.ProviderConfig() + testServiceObjectResource_Create,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_service_object.test", "name", "tf-test-service-object"),
					resource.TestCheckResourceAttr("sccfm_service_object.test", "description", "created by acceptance test"),
					resource.TestCheckResourceAttr("sccfm_service_object.test", "protocol", "TCP"),
					resource.TestCheckResourceAttr("sccfm_service_object.test", "value", "443"),
					resource.TestCheckResourceAttrSet("sccfm_service_object.test", "id"),
				),
			},
			// Update port
			{
				Config: acctest.ProviderConfig() + testServiceObjectResource_UpdatePort,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_service_object.test", "description", "port updated"),
					resource.TestCheckResourceAttr("sccfm_service_object.test", "value", "8443"),
					resource.TestCheckResourceAttr("sccfm_service_object.test", "protocol", "TCP"),
				),
			},
			// Change protocol to UDP
			{
				Config: acctest.ProviderConfig() + testServiceObjectResource_UpdateProtocol,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_service_object.test", "protocol", "UDP"),
					resource.TestCheckResourceAttr("sccfm_service_object.test", "value", "53"),
				),
			},
			// Update name
			{
				Config: acctest.ProviderConfig() + testServiceObjectResource_UpdateName,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_service_object.test", "name", "tf-test-service-object-renamed"),
				),
			},
			// Update to port range
			{
				Config: acctest.ProviderConfig() + testServiceObjectResource_PortRange,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_service_object.test", "description", "port range"),
					resource.TestCheckResourceAttr("sccfm_service_object.test", "protocol", "TCP"),
					resource.TestCheckResourceAttr("sccfm_service_object.test", "value", "8000-9000"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

const testServiceObjectResource_ICMP_Create = `
resource "sccfm_service_object" "test_icmp" {
	name        = "tf-test-service-object-icmp"
	description = "ICMP service without port"
	protocol    = "ICMP"
}`

const testServiceObjectResource_ICMP_UpdateDescription = `
resource "sccfm_service_object" "test_icmp" {
	name        = "tf-test-service-object-icmp"
	description = "ICMP description updated"
	protocol    = "ICMP"
}`

func TestAccServiceObjectResource_ICMP(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create ICMP (no port)
			{
				Config: acctest.ProviderConfig() + testServiceObjectResource_ICMP_Create,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_service_object.test_icmp", "protocol", "ICMP"),
					resource.TestCheckResourceAttrSet("sccfm_service_object.test_icmp", "id"),
				),
			},
			// Update description only
			{
				Config: acctest.ProviderConfig() + testServiceObjectResource_ICMP_UpdateDescription,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_service_object.test_icmp", "description", "ICMP description updated"),
					resource.TestCheckResourceAttr("sccfm_service_object.test_icmp", "protocol", "ICMP"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testServiceObjectOverrideConfig_Create(targetId string) string {
	return fmt.Sprintf(`
resource "sccfm_service_object" "test_overrides" {
	name        = "tf-test-service-object-overrides"
	description = "with one override"
	protocol    = "TCP"
	value       = "443"

	overrides {
		target_id = %q
		protocol  = "TCP"
		value     = "8443"
	}
}`, targetId)
}

func testServiceObjectOverrideConfig_ModifyOverride(targetId string) string {
	return fmt.Sprintf(`
resource "sccfm_service_object" "test_overrides" {
	name        = "tf-test-service-object-overrides"
	description = "override port changed"
	protocol    = "TCP"
	value       = "443"

	overrides {
		target_id = %q
		protocol  = "TCP"
		value     = "9443"
	}
}`, targetId)
}

const testServiceObjectOverrideConfig_RemoveOverrides = `
resource "sccfm_service_object" "test_overrides" {
	name        = "tf-test-service-object-overrides"
	description = "overrides removed"
	protocol    = "TCP"
	value       = "443"
}`

func TestAccServiceObjectResource_OverrideLifecycle(t *testing.T) {
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
				Config: acctest.ProviderConfig() + testServiceObjectOverrideConfig_Create(targetId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_service_object.test_overrides", "overrides.#", "1"),
					resource.TestCheckResourceAttr("sccfm_service_object.test_overrides", "overrides.0.value", "8443"),
				),
			},
			// Modify override
			{
				Config: acctest.ProviderConfig() + testServiceObjectOverrideConfig_ModifyOverride(targetId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_service_object.test_overrides", "overrides.#", "1"),
					resource.TestCheckResourceAttr("sccfm_service_object.test_overrides", "overrides.0.value", "9443"),
				),
			},
			// Remove overrides
			{
				Config: acctest.ProviderConfig() + testServiceObjectOverrideConfig_RemoveOverrides,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_service_object.test_overrides", "overrides.#", "0"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
