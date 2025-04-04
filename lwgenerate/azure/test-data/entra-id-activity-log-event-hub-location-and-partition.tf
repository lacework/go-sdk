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

module "microsoft-entra-id-activity-log" {
  source                      = "lacework/microsoft-entra-id-activity-log/azure"
  version                     = "~> 0.3"
  application_id              = module.az_ad_application.application_id
  application_password        = module.az_ad_application.application_password
  location                    = "West US 2"
  num_partitions              = 2
  service_principal_id        = module.az_ad_application.service_principal_id
  use_existing_ad_application = false
}
