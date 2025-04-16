variable "access_policy_name" {
  default = "Default Access Control Policy"
}

variable "licenses" {
  default = ["BASE", "THREAT"]
  type = list(string)
}

variable "device_name" {
  type = string
}

variable "is_virtual" {
  type = bool
}

variable "performance_tier" {
  type = string
  description = "Required only for virtual FTDs. The performance tier of the virtual FTD."
}

variable "labels" {
  type = list(string)
}

variable "ftd_ssh_host" {
  type        = string
  description = "The host to use to SSH to the FTD to paste the CLI"
}

variable "ftd_ssh_port" {
  type = number
  default = 22
}

variable "ftd_ssh_password" {
  description = "The SSH password for the FTD device"
  type        = string
}