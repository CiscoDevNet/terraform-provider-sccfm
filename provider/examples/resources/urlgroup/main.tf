resource "sccfm_url_object" "cisco" {
  name = "cisco-tf-test"
  url  = "https://www.cisco.com"
}

data "sccfm_ftd_device" "my_ftd" {
  name = "siwarrie-ftd-to-upgrade-2025-09-08-1"
}

resource "sccfm_url_group" "allowed_sites" {
  name        = "allowed-sites-tf-test"
  description = "Allowed website URLs"
  referenced_object_uids = [
    sccfm_url_object.cisco.id,
  ]
  values = [
    "https://www.example.com",
  ]

  overrides = [
    {
      target_id = data.sccfm_ftd_device.my_ftd.id
      values    = ["https://staging.example.com"]
      referenced_object_uids = [
        sccfm_url_object.cisco.id,
      ]
    }
  ]
}
