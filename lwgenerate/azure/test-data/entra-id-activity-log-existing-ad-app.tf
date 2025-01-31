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

module "microsoft-entra-id-activity-log" {
  source                      = "lacework/microsoft-entra-id-activity-log/azure"
  version                     = "~> 0.3"
  application_id              = "testID"
  application_password        = "pass"
  service_principal_id        = "principal"
  use_existing_ad_application = true
}
