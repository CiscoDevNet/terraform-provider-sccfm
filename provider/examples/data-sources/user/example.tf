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


data "sccfm_user" "example_user" {
  name = "<your-username@example.com>"
  is_api_only_user = false
}

data "sccfm_user" "api_only_user" {
  name = "<api-only-user>@<tenant-name>"
  is_api_only_user = true
}

output "example_user_uid" {
  value = data.sccfm_user.example_user.id
}

output "example_user_is_api_only_user" {
  value = data.sccfm_user.example_user.is_api_only_user
}

output "example_user_role" {
  value = data.sccfm_user.example_user.role
}

output "api_only_user_uid" {
  value = data.sccfm_user.api_only_user.id
}

output "api_only_user_is_api_only_user" {
  value = data.sccfm_user.api_only_user.is_api_only_user
}

output "api_only_user_role" {
  value = data.sccfm_user.api_only_user.role
}