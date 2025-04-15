terraform {
  required_providers {
    sccfm = {
      source = "CiscoDevnet/scc-firewall-manager"
    }
  }
}

provider "sccfm" {
  base_url  = "<https://us.manage.security.cisco.com|https://eu.manage.security.cisco.com|https://apj.manage.security.cisco.com|https://aus.manage.security.cisco.com|https://in.manage.security.cisco.com>"
  api_token = file("${path.module}/api_token.txt")
}

data "sccfm_asa_device" "my_asa" {
  name = "<name-of-device>"
}

output "asa_connector_type" {
  value = data.sccfm_asa_device.my_asa.connector_type
}
output "asa_sdc_name" {
  value = data.sccfm_asa_device.my_asa.sdc_name
}
output "asa_name" {
  value = data.sccfm_asa_device.my_asa.name
}
output "asa_socket_address" {
  value = data.sccfm_asa_device.my_asa.socket_address
}
output "asa_host" {
  value = data.sccfm_asa_device.my_asa.host
}
output "asa_port" {
  value = data.sccfm_asa_device.my_asa.port
}
output "asa_ignore_certificate" {
  value = data.sccfm_asa_device.my_asa.ignore_certificate
}
