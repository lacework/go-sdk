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

module "az_config" {
  source                      = "lacework/config/azure"
  version                     = "~> 3.0"
  application_id              = "AD-Test-Application-ID"
  application_password        = "AD-Test-Password"
  lacework_integration_name   = "Test Config Rename"
  service_principal_id        = "AD-Test-Principal-ID"
  use_existing_ad_application = true
}

module "az_activity_log" {
  source                            = "lacework/activity-log/azure"
  version                           = "~> 3.0"
  application_id                    = "AD-Test-Application-ID"
  application_password              = "AD-Test-Password"
  infrastructure_encryption_enabled = true
  lacework_integration_name         = "Test Activity Log Rename"
  service_principal_id              = "AD-Test-Principal-ID"
  use_existing_ad_application       = true
}
