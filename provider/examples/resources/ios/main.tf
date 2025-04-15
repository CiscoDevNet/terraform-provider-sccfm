resource "sccfm_ios_device" "my_ios" {
  name               = var.device_name
  connector_name     = var.connector_name
  socket_address     = var.socket_address
  username           = var.username
  password           = var.password
  ignore_certificate = var.ignore_certificate
}