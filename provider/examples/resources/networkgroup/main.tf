resource "sccfm_network_object" "web_server" {
  name  = "web-server-tf-test"
  value = "10.0.0.1"
}

resource "sccfm_network_object" "db_server" {
  name  = "db-server-tf-test"
  value = "10.0.0.2"
}

resource "sccfm_network_group" "servers" {
  name        = "all-servers-tf-test"
  description = "All server addresses"
  referenced_object_uids = [
    sccfm_network_object.web_server.id,
    sccfm_network_object.db_server.id,
  ]
  values = [
    "10.0.0.100",
    "10.0.0.101",
  ]

  overrides = [
    {
      target_id = "6dae078f-65a8-4cfc-8af8-9554976b5aae"
      values    = ["10.1.0.100", "10.1.0.101"]
      referenced_object_uids = [
        sccfm_network_object.web_server.id,
      ]
    }
  ]
}
