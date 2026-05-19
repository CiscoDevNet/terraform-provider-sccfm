resource "sccfm_network_object" "example_host" {
  name        = "web-server-tf-test"
  description = "Web server IP address"
  value       = "10.0.0.1"
}

resource "sccfm_network_object" "example_subnet" {
  name        = "internal-network-tf-test"
  description = "Internal subnet"
  value       = "10.10.0.0/24"
}

resource "sccfm_network_object" "example_range" {
  name        = "dhcp-range-tf-test"
  description = "DHCP address pool"
  value       = "192.168.1.100-192.168.1.200"
}

resource "sccfm_network_object" "example_ipv6" {
  name        = "ipv6-host-tf-test"
  description = "IPv6 host address"
  value       = "2001:db8::1"
}

# resource "sccfm_network_object" "with_overrides" {
#   name        = "app-server-tf-test"
#   description = "Application server with per-target overrides"
#   value       = "10.0.0.50"
#
#   overrides = [
#     {
#       target_id = "6dae078f-65a8-4cfc-8af8-9554976b5aae"
#       value     = "10.1.0.50"
#     },
#     {
#       target_id = "1ba692de-8166-4cee-90b5-6f817c3c2008"
#       value     = "10.2.0.50"
#     }
#   ]
# }
