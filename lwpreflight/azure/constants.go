package azure

type IntegrationType string

const (
	Config      IntegrationType = "azure_config"
	ActivityLog IntegrationType = "azure_activity_log"
	Agentless   IntegrationType = "azure_agentless"
)

// Entra ID directory role template IDs
// https://learn.microsoft.com/en-us/entra/identity/role-based-access-control/permissions-reference
const (
	GlobalAdministratorRoleID           = "62e90394-69f5-4237-9190-012177145e10"
	ApplicationAdministratorRoleID      = "9b895d92-2cd3-44c7-9d02-a6ac2d5ea5c3"
	CloudApplicationAdministratorRoleID = "158c047a-c907-4556-b7ef-446551a6b5f7"
	PrivilegedRoleAdministratorRoleID   = "e8611ab8-c189-46e8-94e1-60213ab1f814"
)

var RequiredPermissions = map[IntegrationType][]string{
	Config: {
		"Microsoft.Authorization/roleAssignments/read",
		"Microsoft.Authorization/roleAssignments/write",
	},

	ActivityLog: {
		"Microsoft.Authorization/roleAssignments/read",
		"Microsoft.Authorization/roleAssignments/write",
		"Microsoft.Authorization/roleAssignments/delete",
		"Microsoft.Authorization/roleDefinitions/write",
		"Microsoft.Authorization/roleDefinitions/delete",
		"Microsoft.Resources/subscriptions/resourcegroups/read",
		"Microsoft.Resources/subscriptions/resourcegroups/write",
		"Microsoft.Resources/subscriptions/resourcegroups/delete",
		"Microsoft.Resources/deployments/read",
		"Microsoft.Resources/deployments/write",
		"Microsoft.Network/virtualNetworks/read",
		"Microsoft.Network/virtualNetworks/write",
		"Microsoft.Network/virtualNetworks/delete",
		"Microsoft.Network/virtualNetworks/subnets/read",
		"Microsoft.Network/virtualNetworks/subnets/write",
		"Microsoft.Network/virtualNetworks/subnets/delete",
		"Microsoft.Network/virtualNetworks/subnets/join/action",
		"Microsoft.Network/privateEndpoints/read",
		"Microsoft.Network/privateEndpoints/write",
		"Microsoft.Network/privateEndpoints/delete",
		"Microsoft.Storage/storageAccounts/read",
		"Microsoft.Storage/storageAccounts/write",
		"Microsoft.Storage/storageAccounts/delete",
		"Microsoft.Storage/storageAccounts/listKeys/action",
		"Microsoft.Storage/storageAccounts/PrivateEndpointConnectionsApproval/action",
		"Microsoft.Storage/storageAccounts/blobServices/read",
		"Microsoft.Storage/storageAccounts/fileServices/read",
		"Microsoft.EventGrid/eventSubscriptions/read",
		"Microsoft.EventGrid/eventSubscriptions/write",
		"Microsoft.EventGrid/eventSubscriptions/delete",
		"Microsoft.EventGrid/eventSubscriptions/getFullUrl/action",
		"Microsoft.Insights/diagnosticSettings/read",
		"Microsoft.Insights/diagnosticSettings/write",
		"Microsoft.Insights/diagnosticSettings/delete",
	},

	Agentless: {
		"Microsoft.App/jobs/*",
		"Microsoft.App/managedEnvironments/*",
		"Microsoft.Authorization/roleAssignments/*",
		"Microsoft.Authorization/roleDefinitions/*",
		"Microsoft.Compute/virtualMachines/read",
		"Microsoft.Compute/virtualMachines/delete",
		"Microsoft.Compute/virtualMachineScaleSets/read",
		"Microsoft.Compute/virtualMachineScaleSets/virtualMachines/read",
		"Microsoft.KeyVault/vaults/*",
		"Microsoft.KeyVault/locations/deletedVaults/purge/*",
		"Microsoft.KeyVault/locations/operationResults/*",
		"Microsoft.ManagedIdentity/userAssignedIdentities/*",
		"Microsoft.Network/natGateways/*",
		"Microsoft.Network/networkSecurityGroups/*",
		"Microsoft.Network/publicIPAddresses/*",
		"Microsoft.Network/virtualNetworks/*",
		"Microsoft.OperationalInsights/workspaces/*",
		"Microsoft.OperationalInsights/workspaces/sharedKeys/*",
		"Microsoft.Resources/subscriptions/resourcegroups/*",
		"Microsoft.Storage/storageAccounts/*",
		"Microsoft.Storage/storageAccounts/blobServices/*",
		"Microsoft.Storage/storageAccounts/fileServices/*",
		"Microsoft.Storage/storageAccounts/listKeys/*",
	},
}
