resource "sccfm_url_object" "example" {
  name        = "cisco-website-tf-test"
  description = "Cisco corporate website"
  url         = "https://www.cisco.com"
}

resource "sccfm_url_object" "with_overrides" {
  name        = "app-portal-tf-test"
  description = "Application portal with per-target overrides"
  url         = "https://portal.example.com"

  overrides = [
    {
      target_id = "6dae078f-65a8-4cfc-8af8-9554976b5aae"
      url       = "https://portal-staging.example.com"
    }
  ]
}
