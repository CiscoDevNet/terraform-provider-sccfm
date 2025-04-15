resource "sccfm_ftd_device" "my_ftd" {
  name               = var.device_name
  licenses           = var.licenses
  virtual            = var.is_virtual
  access_policy_name = var.access_policy_name
  performance_tier   = var.performance_tier
}

resource "null_resource" "ssh_into_ftds" {
  provisioner "local-exec" {
    command = <<EOT
      sshpass -p ${var.ftd_ssh_password} ssh -o StrictHostKeyChecking=no -p ${var.ftd_ssh_port} admin@${var.ftd_ssh_host} "${sccfm_ftd_device.my_ftd.generated_command}"
    EOT
  }
}

resource "sccfm_ftd_device_onboarding" "my_ftd" {
  depends_on = [null_resource.ssh_into_ftds]
  ftd_uid = sccfm_ftd_device.my_ftd.id
}