terraform {
  required_providers {
    sccfm = {
      source = "CiscoDevnet/sccfm"
    }
  }
}

provider "sccfm" {
  base_url  = "<https://us.manage.security.cisco.com|https://eu.manage.security.cisco.com|https://apj.manage.security.cisco.com|https://aus.manage.security.cisco.com|https://in.manage.security.cisco.com>"
  api_token = file("${path.module}/api_token.txt")
}

data "sccfm_tenant" "current" {
}

output "current_tenant_uid" {
  value = data.sccfm_tenant.current.id
}

output "current_tenant_name" {
  value = data.sccfm_tenant.current.name
}

output "current_tenant_human_readable_name" {
  value = data.sccfm_tenant.current.human_readable_name
}

output "current_tenant_subscription_type" {
  value = data.sccfm_tenant.current.subscription_type
}