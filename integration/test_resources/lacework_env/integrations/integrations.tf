terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
  backend "gcs" {
    bucket = var.storage_bucket_name
    prefix = var.state_folder
  }
}

provider "lacework" {}
provider "aws" {}
provider "google" {}
provider "azuread" {}
provider "azurerm" {
  subscription_id = var.az_subscription
  features {}
}

resource "lacework_agent_access_token" "token" {
  name        = "codefresh-int-test-token"
  description = "this token is used for our ci/cd tests (do-not-update)"
}

# Tech Ally Docker Hub, required for go-sdk/integration/container_vulnerability_test.go
resource "lacework_integration_docker_hub" "techally_dockerhub" {
  name     = "TF tech-ally docker"
  username = var.docker_hub_user
  password = var.docker_hub_pass
}

# Lacework AWS config integration, required for go-sdk/integration/compliance_aws_test.go
module "aws_config" {
  source = "lacework/config/aws"

  lacework_aws_account_id = var.lacework_aws_account_id
}

# Lacework GCP config integration, required for go-sdk/integration/compliance_gcp_test.go
module "gcp_organization_level_config" {
  source = "lacework/config/gcp"

  org_integration = var.org_integration
  organization_id = var.organization_id
  project_id      = var.project_id
}

module "az_config" {
  source = "lacework/config/azure"
}
