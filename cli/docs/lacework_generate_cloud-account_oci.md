---
title: "lacework generate cloud-account oci"
slug: lacework_generate_cloud-account_oci
hide_title: true
---

## lacework generate cloud-account oci

Generate and/or execute Terraform code for OCI integration

### Synopsis

Use this command to generate Terraform code for deploying Lacework into an OCI tenant.

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
lacework generate cloud-account oci [flags]
```

### Options

```
      --apply                   run terraform apply without executing plan or prompting
      --config                  enable configuration integration
      --config_name string      specify name of configuration integration
  -h, --help                    help for oci
      --oci_user_email string   specify the email address to associate with the integration OCI user
      --output string           location to write generated content (default is ~/lacework/oci)
      --tenant_ocid string      specify the OCID of the tenant to integrate
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

