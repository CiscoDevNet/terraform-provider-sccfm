resource "sccfm_url_object" "cisco" {
  name = "cisco-tf-test"
  url  = "https://www.cisco.com"
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
      target_id = "6dae078f-65a8-4cfc-8af8-9554976b5aae"
      values    = ["https://staging.example.com"]
      referenced_object_uids = [
        sccfm_url_object.cisco.id,
      ]
    }
  ]
}
