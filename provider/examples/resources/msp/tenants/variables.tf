variable "tenant_name" {
  description = "The name of the tenant"
  type        = string
}

variable "tenant_display_name" {
  description = "The display name for the tenant"
  type        = string
}

variable "existing_tenant_api_token" {
  description = "The API token for the existing tenant"
  type        = string
  sensitive   = true
}