variable "duo_admin_panel_name" {
  description = "The name of the Duo Admin Panel"
  type        = string
}

variable "duo_admin_panel_host" {
  description = "The host of the Duo Admin Panel"
  type        = string
}

variable "duo_admin_panel_integration_key" {
  description = "The integration key for the Duo Admin Panel"
  type        = string
}

variable "duo_admin_panel_secret_key" {
  description = "The secret key for the Duo Admin Panel"
  type        = string
}

variable "duo_admin_panel_labels" {
  description = "Labels for the Duo Admin Panel"
  type        = list(string)
  default     = []
}