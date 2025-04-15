data "sccfm_msp_managed_tenant" "tenant" {
  name             = "<your-tenant-name>"
}

resource "sccfm_msp_managed_tenant_users" "example" {
  tenant_uid = data.sccfm_msp_managed_tenant.tenant.id
  users = [
    {
      username = "<username1@example.com>",
      roles = ["ROLE_SUPER_ADMIN"]
      api_only_user = false
    },
    {
      username = "<username2@example.com>",
      roles = ["ROLE_ADMIN"]
      api_only_user = false
    },
    {
      username = "<api-only-user-name>",
      roles = ["ROLE_SUPER_ADMIN"]
      api_only_user = true
    }
  ]
}

resource "sccfm_msp_managed_tenant_user_api_token" "user_token" {
  tenant_uid = data.sccfm_msp_managed_tenant.tenant.id
  user_uid = sccfm_msp_managed_tenant_users.example.users[2].id
}

output "api_token" {
  value = sccfm_msp_managed_tenant_user_api_token.user_token.api_token
  sensitive = true
}