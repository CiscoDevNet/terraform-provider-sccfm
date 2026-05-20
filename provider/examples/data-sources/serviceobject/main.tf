data "sccfm_service_object" "existing" {
  name = "HTTP-tf-test"
}

output "service_object_id" {
  value = data.sccfm_service_object.existing.id
}

output "service_object_protocol" {
  value = data.sccfm_service_object.existing.protocol
}

output "service_object_value" {
  value = data.sccfm_service_object.existing.value
}

# Use the looked-up object in a group
resource "sccfm_service_group" "web_services" {
  name        = "web-services-from-ds-tf-test"
  description = "Group referencing an existing service object"
  referenced_object_uids = [
    data.sccfm_service_object.existing.id,
  ]
}
