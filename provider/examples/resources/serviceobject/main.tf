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

resource "sccfm_service_object" "with_overrides" {
  name        = "app-service-tf-test"
  description = "Application service with per-target overrides"
  protocol    = "TCP"
  value       = "8080"

  overrides = [
    {
      target_id = "6dae078f-65a8-4cfc-8af8-9554976b5aae"
      protocol  = "TCP"
      value     = "9090"
    }
  ]
}
