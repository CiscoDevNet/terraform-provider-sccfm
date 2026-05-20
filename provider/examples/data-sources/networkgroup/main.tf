data "sccfm_network_group" "existing" {
  name = "all-servers-tf-test"
}

output "network_group_id" {
  value = data.sccfm_network_group.existing.id
}

output "network_group_values" {
  value = data.sccfm_network_group.existing.values
}

output "network_group_referenced_uids" {
  value = data.sccfm_network_group.existing.referenced_object_uids
}
