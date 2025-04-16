resource "sccfm_user" "new_api_only_user" {
  name             = "api_user"
  is_api_only_user = true
  role             = "ROLE_ADMIN"
}

resource "sccfm_api_token" "new_api_only_user_api_token" {
  username = sccfm_user.new_api_only_user.generated_username
}

output "api_only_user_api_token_value" {
  value     = sccfm_api_token.new_api_only_user_api_token.api_token
  sensitive = true
}