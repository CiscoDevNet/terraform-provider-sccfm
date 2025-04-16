
resource "sccfm_asa_device" "my_asa" {
  name               = var.asa_name
  socket_address     = var.asa_hostname
  username           = var.asa_username
  password           = var.asa_password
  connector_type     = "CDG"
  ignore_certificate = false
}

output "asa_connector_type" {
  value = sccfm_asa_device.my_asa.connector_type
}
output "asa_sdc_name" {
  value = sccfm_asa_device.my_asa.connector_name
}

output "asa_host" {
  value = sccfm_asa_device.my_asa.host
}
output "asa_port" {
  value = sccfm_asa_device.my_asa.port
}
