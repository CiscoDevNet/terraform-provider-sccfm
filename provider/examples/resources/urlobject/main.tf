resource "sccfm_url_object" "example" {
  name        = "cisco-website-tf-test"
  description = "Cisco corporate website"
  url         = "https://www.cisco.com"
}

data "sccfm_ftd_device" "my_ftd" {
  name = "siwarrie-ftd-to-upgrade-2025-09-08-1"
}

resource "sccfm_url_object" "with_overrides" {
  name        = "app-portal-tf-test"
  description = "Application portal with per-target overrides"
  url         = "https://portal.example.com"

  overrides = [
    {
      target_id = data.sccfm_ftd_device.my_ftd.id
      url       = "https://portal-staging.example.com"
    }
  ]
}
