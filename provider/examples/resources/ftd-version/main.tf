data "sccfm_ftd_device" "ftd" {
  name = var.ftd_name
}

resource "sccfm_ftd_device_version" "ftd" {
  ftd_uid = data.sccfm_ftd_device.ftd.id
  software_version   = "7.3.1-19"
}