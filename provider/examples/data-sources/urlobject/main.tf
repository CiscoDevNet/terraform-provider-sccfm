data "sccfm_url_object" "existing" {
  name = "cisco-website-tf-test"
}

output "url_object_id" {
  value = data.sccfm_url_object.existing.id
}

output "url_object_url" {
  value = data.sccfm_url_object.existing.url
}

# Use the looked-up object in a group
resource "sccfm_url_group" "sites" {
  name        = "allowed-sites-from-ds-tf-test"
  description = "Group referencing an existing URL object"
  referenced_object_uids = [
    data.sccfm_url_object.existing.id,
  ]
}
