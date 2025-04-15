resource "sccfm_user" "new_user" {
  name             = "<your-emaiL>@example.com"
  first_name       = "<firstname>"
  last_name        = "<lastname>"
  is_api_only_user = false
  role             = "<ROLE_READ_ONLY|ROLE_DEPLOY_ONLY|ROLE_EDIT_ONLY|ROLE_ADMIN|ROLE_SUPER_ADMIN>"
}

resource "sccfm_user" "new_api_only_user" {
  name             = "<name-of-api-only-user>"
  is_api_only_user = true
  role             = "<ROLE_READ_ONLY|ROLE_DEPLOY_ONLY|ROLE_EDIT_ONLY|ROLE_ADMIN|ROLE_SUPER_ADMIN>"
}
