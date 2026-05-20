data "sccfm_url_group" "existing" {
  name = "allowed-sites-tf-test"
}

output "url_group_id" {
  value = data.sccfm_url_group.existing.id
}

output "url_group_values" {
  value = data.sccfm_url_group.existing.values
}

output "url_group_referenced_uids" {
  value = data.sccfm_url_group.existing.referenced_object_uids
}
