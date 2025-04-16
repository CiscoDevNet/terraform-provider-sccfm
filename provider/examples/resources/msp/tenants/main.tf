resource "sccfm_msp_managed_tenant" "tenant" {
  name         = var.tenant_name
  display_name = var.tenant_display_name
}

resource "sccfm_msp_managed_tenant" "existing_tenant" {
  api_token = var.existing_tenant_api_token
}