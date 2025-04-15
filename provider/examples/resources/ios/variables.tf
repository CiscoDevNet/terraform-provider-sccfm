variable "device_name" {
  description = "The name of the device"
  type        = string
}

variable "connector_name" {
  description = "The name of the SDC connector (not required if connector type is CDG)"
  type        = string
  default     = null
}

variable "socket_address" {
  description = "The host and port of the device in the format <host>:<port>"
  type        = string
}

variable "username" {
  description = "The username for the device"
  type        = string
}

variable "password" {
  description = "The password for the device"
  type        = string
  sensitive   = true
}

variable "ignore_certificate" {
  description = "Whether to ignore certificate validation (true or false)"
  type        = bool
}