data "sccfm_service_group" "existing" {
  name = "web-services-tf-test"
}

output "service_group_id" {
  value = data.sccfm_service_group.existing.id
}

output "service_group_values" {
  value = data.sccfm_service_group.existing.values
}

output "service_group_referenced_uids" {
  value = data.sccfm_service_group.existing.referenced_object_uids
}
