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
