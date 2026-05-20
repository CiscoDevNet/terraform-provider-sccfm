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

data "sccfm_ftd_device" "my_ftd" {
  name = "siwarrie-ftd-to-upgrade-2025-09-08-1"
}

resource "sccfm_network_object" "with_overrides" {
  name        = "app-server-tf-test"
  description = "Application server with per-target overrides"
  value       = "10.0.0.50"

  overrides = [
    {
      target_id = data.sccfm_ftd_device.my_ftd.id
      value     = "10.1.0.50"
    }
  ]
}
