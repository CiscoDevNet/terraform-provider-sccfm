resource "sccfm_service_object" "http" {
  name     = "HTTP-tf-test"
  protocol = "TCP"
  value    = "80"
}

resource "sccfm_service_object" "https" {
  name     = "HTTPS-tf-test"
  protocol = "TCP"
  value    = "443"
}

resource "sccfm_service_group" "web_services" {
  name        = "web-services-tf-test"
  description = "HTTP and HTTPS services"
  referenced_object_uids = [
    sccfm_service_object.http.id,
    sccfm_service_object.https.id,
  ]

  values = [
    {
      protocol = "TCP"
      value    = "8080"
    }
  ]

  overrides = [
    {
      target_id = "6dae078f-65a8-4cfc-8af8-9554976b5aae"
      referenced_object_uids = [
        sccfm_service_object.http.id,
      ]
      values = [
        {
          protocol = "TCP"
          value    = "9090"
        }
      ]
    }
  ]
}
