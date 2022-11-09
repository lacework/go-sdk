terraform {
  required_providers {
    azuread = {
      source  = "hashicorp/azuread"
      version = "~> 2.16"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 2.91.0"
    }
    lacework = {
      source  = "lacework/lacework"
      version = "~> 1.0"
    }
  }
}

provider "azuread" {
}

provider "azurerm" {
  features {
  }
}

module "az_ad_application" {
  source  = "lacework/ad-application/azure"
  version = "~> 1.0"
}

module "az_activity_log" {
  source                      = "lacework/activity-log/azure"
  version                     = "~> 1.0"
  application_id              = module.az_ad_application.application_id
  application_password        = module.az_ad_application.application_password
  lacework_integration_name   = "Test Activity Log Rename"
  service_principal_id        = module.az_ad_application.service_principal_id
  use_existing_ad_application = true
}
