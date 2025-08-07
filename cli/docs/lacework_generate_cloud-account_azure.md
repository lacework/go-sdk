---
title: "lacework generate cloud-account azure"
slug: lacework_generate_cloud-account_azure
hide_title: true
---

## lacework generate cloud-account azure

Generate and/or execute Terraform code for Azure integration

### Synopsis

Use this command to generate Terraform code for deploying Lacework into new Azure environment.

By default, this command will function interactively, prompting for the required information to setup
the new cloud account. In interactive mode, this command will:

* Prompt for the required information to setup the integration
* Generate new Terraform code using the inputs
* Optionally, run the generated Terraform code:
  * If Terraform is already installed, the version will be confirmed suitable for use
  * If Terraform is not installed, or the version installed is not suitable, a new version will be
    installed into a temporary location
  * Once Terraform is detected or installed, Terraform plan will be executed
  * The command will prompt with the outcome of the plan and allow to view more details or continue
    with Terraform apply
  * If confirmed, Terraform apply will be run, completing the setup of the cloud account


```
lacework generate cloud-account azure [flags]
```

### Options

```
      --activity_log                                    enable activity log integration
      --activity_log_integration_name string            specify a custom activity log integration name
      --ad_create                                       create new active directory integration (default true)
      --ad_id string                                    existing active directory application id
      --ad_pass string                                  existing active directory application password
      --ad_pid string                                   existing active directory application service principle id
      --agentless                                       enable agentless integration
      --agentless_subscription_ids strings              comma-separated list of subscription IDs for Agentless scanning (e.g., 'sub1,sub2,sub3')
      --all_subscriptions subscription ids              grant read access to ALL subscriptions within Tenant (overrides subscription ids)
      --apply                                           run terraform apply for the generated hcl
      --configuration                                   enable configuration integration
      --configuration_name string                       specify a custom configuration integration name
      --create_log_analytics_workspace                  enable creation of Log Analytics Workspace for agentless scanning
      --entra_id_activity_log                           enable Entra ID activity log integration
      --entra_id_activity_log_integration_name string   specify a custom Entra ID activity log integration name
      --event_hub_location string                       specify the location where the Event Hub for logging will reside
      --event_hub_partition_count int                   specify the number of partitions for the Event Hub (default 1)
      --existing_storage                                use existing storage account
      --global                                          enable global agentless scanning
  -h, --help                                            help for azure
      --integration_level string                        specify the agentless integration level (e.g., 'SUBSCRIPTION', 'TENANT')
      --location string                                 specify azure region where storage account logging resides
      --management_group                                management group level integration
      --management_group_id string                      specify management group id. Required if mgmt_group provided
      --output string                                   location to write generated content (default is ~/lacework/azure)
      --regions strings                                 comma-separated list of Azure regions for agentless scanning (e.g., 'East US,West US')
      --storage_account_name string                     specify storage account name
      --storage_resource_group string                   specify storage resource group
      --subscription_id string                          specify the Azure Subscription ID to be used to provision Lacework resources
      --subscription_ids strings                        list of subscriptions to grant read access; format is id1,id2,id3
```

### Options inherited from parent commands

```
  -a, --account string      account subdomain of URL (i.e. <ACCOUNT>.lacework.net)
  -k, --api_key string      access key id
  -s, --api_secret string   secret access key
      --api_token string    access token (replaces the use of api_key and api_secret)
      --debug               turn on debug logging
      --json                switch commands output from human-readable to json format
      --nocache             turn off caching
      --nocolor             turn off colors
      --noninteractive      turn off interactive mode (disable spinners, prompts, etc.)
      --organization        access organization level data sets (org admins only)
  -p, --profile string      switch between profiles configured at ~/.lacework.toml
      --subaccount string   sub-account name inside your organization (org admins only)
```

### SEE ALSO

* [lacework generate cloud-account](lacework_generate_cloud-account.md)	 - Generate cloud integration IaC

