terraform {
  required_providers {
    sccfm = {
      source = "CiscoDevnet/sccfm"
    }
  }
}

provider "sccfm" {
  base_url = "<https://us.manage.security.cisco.com|https://eu.manage.security.cisco.com|https://apj.manage.security.cisco.com|https://aus.manage.security.cisco.com|https://in.manage.security.cisco.com>"
  api_token = file("${path.module}/api_token.txt")
}


resource "sccfm_tenant_settings" "tenant_settings" {
  deny_cisco_support_access_to_tenant_enabled = false
}

output "example_tenant_id" {
  value = sccfm_tenant_settings.tenant_settings.id
}

output "example_change_request_support_enabled" {
  value = sccfm_tenant_settings.tenant_settings.change_request_support_enabled
}

output "example_auto_accept_device_changes_enabled" {
  value = sccfm_tenant_settings.tenant_settings.auto_accept_device_changes_enabled
}

output "example_web_analytics_enabled" {
  value = sccfm_tenant_settings.tenant_settings.web_analytics_enabled
}

output "example_scheduled_deployments_enabled" {
  value = sccfm_tenant_settings.tenant_settings.scheduled_deployments_enabled
}

output "example_deny_cisco_support_access_to_tenant_enabled" {
  value = sccfm_tenant_settings.tenant_settings.deny_cisco_support_access_to_tenant_enabled
}

output "example_multi_cloud_defense_enabled" {
  value = sccfm_tenant_settings.tenant_settings.multi_cloud_defense_enabled
}

output "example_auto_discover_on_prem_fmcs" {
  value = sccfm_tenant_settings.tenant_settings.auto_discover_on_prem_fmcs_enabled
}

output "example_conflict_detection_interval" {
  value = sccfm_tenant_settings.tenant_settings.conflict_detection_interval
}
