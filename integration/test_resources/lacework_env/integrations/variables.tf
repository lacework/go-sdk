variable "az_subscription" {
  type        = string
  default     = ""
  description = "The Azure subscription ID"
}

variable "storage_bucket_name" {
  type        = string
  default     = ""
  description = "The name of the gcs where tf state is stored"
}

variable "state_folder" {
  type        = string
  default     = ""
  description = "The folder inside the gcs where tf state is stored"
}

variable "org_integration" {
  type        = bool
  default     = false
  description = "If set to true, configure an organization level integration"
}

variable "organization_id" {
  type        = string
  default     = ""
  description = "The GCP organization ID, required if org_integration is set to true"
}

variable "project_id" {
  type        = string
  default     = ""
  description = "The GCP project ID"
}

variable "docker_hub_user" {
  type        = string
  default     = ""
  description = "The username for dockerhub"
}

variable "docker_hub_pass" {
  type        = string
  default     = ""
  description = "The password for dockerhub"
}

variable "lacework_aws_account_id" {
  type        = string
  default     = ""
  description = "The Lacework AWS account that the IAM role will grant access"
}
