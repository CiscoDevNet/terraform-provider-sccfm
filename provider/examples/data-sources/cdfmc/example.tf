terraform {
  required_providers {
    sccfm = {
      source = "CiscoDevnet/scc-firewall-manager"
    }
  }
}

provider "sccfm" {
  base_url  = "<https://us.manage.security.cisco.com|https://eu.manage.security.cisco.com|https://apj.manage.security.cisco.com|https://aus.manage.security.cisco.com|https://in.manage.security.cisco.com>"
  api_token = file("${path.module}/api_token.txt")
}

data "sccfm_cdfmc" "current" {
}

output "cdfmc_hostname" {
  value = data.sccfm_cdfmc.current.hostname
}

output "cdfmc_software_version" {
  value = data.sccfm_cdfmc.current.software_version
}

output "cdfmc_uid" {
  value = data.sccfm_cdfmc.current.id
}

output "cdfmc_domain_uuid" {
  value = data.sccfm_cdfmc.current.domain_uuid
}