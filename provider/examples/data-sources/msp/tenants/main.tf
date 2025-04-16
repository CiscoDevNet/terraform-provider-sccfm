data "sccfm_msp_managed_tenant" "tenant" {
  name             = "CDO_tf-managed-starks-2"
}

output "tenant_display_name" {
  value = data.sccfm_msp_managed_tenant.tenant.display_name
}

output "tenant_region" {
  value = data.sccfm_msp_managed_tenant.tenant.region
}