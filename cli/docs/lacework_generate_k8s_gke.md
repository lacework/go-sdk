---
title: "lacework generate k8s gke"
slug: lacework_generate_k8s_gke
hide_title: true
---

## lacework generate k8s gke

Generate and/or execute Terraform code for GKE integration

### Synopsis

Use this command to generate Terraform code for deploying Lacework into a GKE environment.

By default, this command interactively prompts for the required information to setup the new cloud account.
In interactive mode, this command will:

* Prompt for the required information to setup the integration
* Generate new Terraform code using the inputs
* Optionally, run the generated Terraform code:
  * If Terraform is already installed, the version is verified as compatible for use
	* If Terraform is not installed, or the version installed is not compatible, a new
    version will be installed into a temporary location
	* Once Terraform is detected or installed, Terraform plan will be executed
	* The command will prompt with the outcome of the plan and allow to view more details
    or continue with Terraform apply
	* If confirmed, Terraform apply will be run, completing the setup of the cloud account

This command can also be run in noninteractive mode.
See help output for more details on the parameter value(s) required for Terraform code generation.


```
lacework generate k8s gke [flags]
```

### Options

```
      --apply                                         run terraform apply without executing plan or prompting
      --existing_service_account_name string          specify existing service account name
      --existing_service_account_private_key string   specify existing service account private key (base64 encoded)
      --existing_sink_name string                     specify existing sink name
  -h, --help                                          help for gke
      --integration_name string                       specify a custom integration name
      --organization_id string                        specify the organization id (only set if organization_integration is set)
      --organization_integration                      enable organization integration
      --output string                                 location to write generated content (default is ~/lacework/gcp)
      --prefix string                                 prefix that will be used at the beginning of every generated resource
      --project_id string                             specify the project id to be used to provision lacework resources (required)
      --service_account_credentials string            specify service account credentials JSON file path (leave blank to make use of google credential ENV vars)
      --wait_time string                              amount of time to wait before the next resource is provisioned
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

* [lacework generate k8s](lacework_generate_k8s.md)	 - Generate Kubernetes integration IaC

