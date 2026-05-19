package networkgroup_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/CiscoDevnet/terraform-provider-sccfm/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testNetworkGroupResource_Create = `
resource "sccfm_network_group" "test" {
	name        = "tf-test-network-group"
	description = "created by acceptance test"
	values      = ["10.0.0.0/24", "192.168.1.0/24"]
}`

const testNetworkGroupResource_AddMember = `
resource "sccfm_network_group" "test" {
	name        = "tf-test-network-group"
	description = "member added"
	values      = ["10.0.0.0/24", "192.168.1.0/24", "172.16.0.0/16"]
}`

const testNetworkGroupResource_RemoveMember = `
resource "sccfm_network_group" "test" {
	name        = "tf-test-network-group"
	description = "member removed"
	values      = ["10.0.0.0/24"]
}`

const testNetworkGroupResource_UpdateName = `
resource "sccfm_network_group" "test" {
	name        = "tf-test-network-group-renamed"
	description = "member removed"
	values      = ["10.0.0.0/24"]
}`

const testNetworkGroupResource_RemoveDescription = `
resource "sccfm_network_group" "test" {
	name   = "tf-test-network-group-renamed"
	values = ["10.0.0.0/24"]
}`

func TestAccNetworkGroupResource_FullLifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with 2 members
			{
				Config: acctest.ProviderConfig() + testNetworkGroupResource_Create,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_network_group.test", "name", "tf-test-network-group"),
					resource.TestCheckResourceAttr("sccfm_network_group.test", "description", "created by acceptance test"),
					resource.TestCheckResourceAttr("sccfm_network_group.test", "values.#", "2"),
					resource.TestCheckResourceAttrSet("sccfm_network_group.test", "id"),
				),
			},
			// Add a member
			{
				Config: acctest.ProviderConfig() + testNetworkGroupResource_AddMember,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_network_group.test", "description", "member added"),
					resource.TestCheckResourceAttr("sccfm_network_group.test", "values.#", "3"),
				),
			},
			// Remove members (down to 1)
			{
				Config: acctest.ProviderConfig() + testNetworkGroupResource_RemoveMember,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_network_group.test", "description", "member removed"),
					resource.TestCheckResourceAttr("sccfm_network_group.test", "values.#", "1"),
				),
			},
			// Update name
			{
				Config: acctest.ProviderConfig() + testNetworkGroupResource_UpdateName,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_network_group.test", "name", "tf-test-network-group-renamed"),
				),
			},
			// Remove description
			{
				Config: acctest.ProviderConfig() + testNetworkGroupResource_RemoveDescription,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_network_group.test", "name", "tf-test-network-group-renamed"),
					resource.TestCheckResourceAttr("sccfm_network_group.test", "values.#", "1"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

const testNetworkGroupResource_Refs_Create = `
resource "sccfm_network_object" "member1" {
	name  = "tf-test-ng-member-1"
	value = "10.1.1.1"
}

resource "sccfm_network_group" "test_refs" {
	name                   = "tf-test-network-group-refs"
	description            = "one referenced object"
	referenced_object_uids = [sccfm_network_object.member1.id]
}`

const testNetworkGroupResource_Refs_AddSecond = `
resource "sccfm_network_object" "member1" {
	name  = "tf-test-ng-member-1"
	value = "10.1.1.1"
}

resource "sccfm_network_object" "member2" {
	name  = "tf-test-ng-member-2"
	value = "10.2.2.2"
}

resource "sccfm_network_group" "test_refs" {
	name                   = "tf-test-network-group-refs"
	description            = "two referenced objects"
	referenced_object_uids = [sccfm_network_object.member1.id, sccfm_network_object.member2.id]
}`

const testNetworkGroupResource_Refs_RemoveRef = `
resource "sccfm_network_object" "member1" {
	name  = "tf-test-ng-member-1"
	value = "10.1.1.1"
}

resource "sccfm_network_object" "member2" {
	name  = "tf-test-ng-member-2"
	value = "10.2.2.2"
}

resource "sccfm_network_group" "test_refs" {
	name                   = "tf-test-network-group-refs"
	description            = "ref removed, literal added"
	values                 = ["10.0.0.0/8"]
	referenced_object_uids = [sccfm_network_object.member1.id]
}`

func TestAccNetworkGroupResource_ReferencedObjects_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with one ref
			{
				Config: acctest.ProviderConfig() + testNetworkGroupResource_Refs_Create,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_network_group.test_refs", "name", "tf-test-network-group-refs"),
					resource.TestCheckResourceAttr("sccfm_network_group.test_refs", "referenced_object_uids.#", "1"),
				),
			},
			// Add second ref
			{
				Config: acctest.ProviderConfig() + testNetworkGroupResource_Refs_AddSecond,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_network_group.test_refs", "referenced_object_uids.#", "2"),
					resource.TestCheckResourceAttr("sccfm_network_group.test_refs", "description", "two referenced objects"),
				),
			},
			// Remove a ref, add a literal
			{
				Config: acctest.ProviderConfig() + testNetworkGroupResource_Refs_RemoveRef,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_network_group.test_refs", "referenced_object_uids.#", "1"),
					resource.TestCheckResourceAttr("sccfm_network_group.test_refs", "values.#", "1"),
					resource.TestCheckResourceAttr("sccfm_network_group.test_refs", "description", "ref removed, literal added"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testNetworkGroupOverrideConfig_Create(targetId string) string {
	return fmt.Sprintf(`
resource "sccfm_network_group" "test_overrides" {
	name        = "tf-test-network-group-overrides"
	description = "with override"
	values      = ["10.0.0.0/24"]

	overrides {
		target_id = %q
		values    = ["192.168.1.0/24"]
	}
}`, targetId)
}

func testNetworkGroupOverrideConfig_ModifyOverride(targetId string) string {
	return fmt.Sprintf(`
resource "sccfm_network_group" "test_overrides" {
	name        = "tf-test-network-group-overrides"
	description = "override modified"
	values      = ["10.0.0.0/24"]

	overrides {
		target_id = %q
		values    = ["172.16.0.0/16", "192.168.0.0/16"]
	}
}`, targetId)
}

const testNetworkGroupOverrideConfig_RemoveOverrides = `
resource "sccfm_network_group" "test_overrides" {
	name        = "tf-test-network-group-overrides"
	description = "overrides removed"
	values      = ["10.0.0.0/24"]
}`

func TestAccNetworkGroupResource_OverrideLifecycle(t *testing.T) {
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
				Config: acctest.ProviderConfig() + testNetworkGroupOverrideConfig_Create(targetId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_network_group.test_overrides", "overrides.#", "1"),
				),
			},
			// Modify override
			{
				Config: acctest.ProviderConfig() + testNetworkGroupOverrideConfig_ModifyOverride(targetId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_network_group.test_overrides", "overrides.#", "1"),
					resource.TestCheckResourceAttr("sccfm_network_group.test_overrides", "description", "override modified"),
				),
			},
			// Remove overrides
			{
				Config: acctest.ProviderConfig() + testNetworkGroupOverrideConfig_RemoveOverrides,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_network_group.test_overrides", "overrides.#", "0"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
