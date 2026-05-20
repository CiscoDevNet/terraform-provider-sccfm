resource "sccfm_service_object" "http" {
  name        = "HTTP-tf-test"
  description = "HTTP traffic"
  protocol    = "TCP"
  value       = "80"
}

resource "sccfm_service_object" "https" {
  name        = "HTTPS-tf-test"
  description = "HTTPS traffic"
  protocol    = "TCP"
  value       = "443"
}

resource "sccfm_service_object" "high_ports" {
  name        = "high-ports-tf-test"
  description = "High port range"
  protocol    = "TCP"
  value       = "8000-9000"
}

data "sccfm_ftd_device" "my_ftd" {
  name = "siwarrie-ftd-to-upgrade-2025-09-08-1"
}

resource "sccfm_service_object" "with_overrides" {
  name        = "app-service-tf-test"
  description = "Application service with per-target overrides"
  protocol    = "TCP"
  value       = "8080"

  overrides = [
    {
      target_id = data.sccfm_ftd_device.my_ftd.id
      protocol  = "TCP"
      value     = "9090"
    }
  ]
}
