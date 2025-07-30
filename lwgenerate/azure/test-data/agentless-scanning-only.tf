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
  subscription_id = "11111111-2222-3333-4444-111111111111"
  features {
  }
}

module "lacework_azure_agentless_scanning_subscription_west_us" {
  source                         = "lacework/agentless-scanning/azure"
  version                        = "~> 1.6"
  create_log_analytics_workspace = false
  global                         = true
  included_subscriptions         = ["/subscriptions/11111111-2222-3333-4444-111111111111"]
  integration_level              = "SUBSCRIPTION"
  region                         = "West US"
  scanning_subscription_id       = "11111111-2222-3333-4444-111111111111"
}
