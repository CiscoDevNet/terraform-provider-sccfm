terraform {
  required_providers {
    sccfm = {
      source = "CiscoDevnet/sccfm"
    }
  }
}

provider "sccfm" {
  base_url  = "<https://us.manage.security.cisco.com|https://eu.manage.security.cisco.com|https://apj.manage.security.cisco.com|https://aus.manage.security.cisco.com|https://in.manage.security.cisco.com>"
  api_token = file("${path.module}/api_token.txt")
}

data "sccfm_ios_device" "my_ios" {
  name = "<name-of-device>"
}
output "ios_sdc_name" {
  value = data.sccfm_ios_device.my_ios.connector_name
}
output "ios_name" {
  value = data.sccfm_ios_device.my_ios.name
}
output "ios_socket_address" {
  value = data.sccfm_ios_device.my_ios.socket_address
}
output "ios_host" {
  value = data.sccfm_ios_device.my_ios.host
}
output "ios_port" {
  value = data.sccfm_ios_device.my_ios.port
}
output "ios_ignore_certificate" {
  value = data.sccfm_ios_device.my_ios.ignore_certificate
}