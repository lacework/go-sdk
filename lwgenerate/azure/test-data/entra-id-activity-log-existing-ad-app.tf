terraform {
  required_providers {
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

module "microsoft-entra-id-activity-log" {
  source                      = "lacework/microsoft-entra-id-activity-log/azure"
  version                     = "~> 0.2"
  application_id              = "testID"
  application_password        = "pass"
  service_principal_id        = "principal"
  use_existing_ad_application = true
}
