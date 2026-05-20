data "sccfm_network_object" "existing" {
  name = "web-server-tf-test"
}

output "network_object_id" {
  value = data.sccfm_network_object.existing.id
}

output "network_object_value" {
  value = data.sccfm_network_object.existing.value
}

# Use the looked-up object in a group
resource "sccfm_network_group" "servers" {
  name        = "all-servers-from-ds-tf-test"
  description = "Group referencing an existing object"
  referenced_object_uids = [
    data.sccfm_network_object.existing.id,
  ]
}
