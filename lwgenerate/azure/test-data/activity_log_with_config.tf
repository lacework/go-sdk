terraform {
  required_providers {
    lacework = {
      source  = "lacework/lacework"
      version = "~> 2.0"
    }
  }
}

provider "azuread" {
}

provider "azurerm" {
  subscription_id = "test-subscription"
  features {
  }
}

module "az_ad_application" {
  source  = "lacework/ad-application/azure"
  version = "~> 2.0"
}

module "az_config" {
  source                      = "lacework/config/azure"
  version                     = "~> 3.0"
  application_id              = module.az_ad_application.application_id
  application_password        = module.az_ad_application.application_password
  service_principal_id        = module.az_ad_application.service_principal_id
  use_existing_ad_application = true
}

module "az_activity_log" {
  source                            = "lacework/activity-log/azure"
  version                           = "~> 3.0"
  application_id                    = module.az_ad_application.application_id
  application_password              = module.az_ad_application.application_password
  infrastructure_encryption_enabled = true
  service_principal_id              = module.az_ad_application.service_principal_id
  use_existing_ad_application       = true
}
