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

module "azure-microsoft-entra-id-activity-log" {
  source                      = "lacework/entra-id-activity-log/azure"
  version                     = "~> 1.0"
  application_id              = "testID"
  application_password        = "pass"
  service_principal_id        = "principal"
  use_existing_ad_application = true
}
