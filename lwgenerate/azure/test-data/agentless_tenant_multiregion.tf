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

provider "azurerm" {
  alias           = "agentless"
  subscription_id = "11111111-2222-3333-4444-111111111111"

  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

module "lacework_azure_agentless_scanning_tenant_west_us" {
  source                         = "lacework/agentless-scanning/azure"
  version                        = "~> 2.0"
  create_log_analytics_workspace = false
  global                         = true
  integration_level              = "TENANT"
  region                         = "West US"
  scanning_subscription_id       = "11111111-2222-3333-4444-111111111111"

  providers = {
    azurerm = azurerm.agentless
  }
}

module "lacework_azure_agentless_scanning_tenant_east_us" {
  source                         = "lacework/agentless-scanning/azure"
  version                        = "~> 2.0"
  create_log_analytics_workspace = false
  global                         = false
  global_module_reference        = module.lacework_azure_agentless_scanning_tenant_west_us
  integration_level              = "TENANT"
  region                         = "East US"
  scanning_subscription_id       = "11111111-2222-3333-4444-111111111111"

  providers = {
    azurerm = azurerm.agentless
  }
}
