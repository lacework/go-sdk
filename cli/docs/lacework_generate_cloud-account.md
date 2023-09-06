---
title: "lacework generate cloud-account"
slug: lacework_generate_cloud-account
hide_title: true
---

## lacework generate cloud-account

Generate cloud integration IaC

### Synopsis

Generate cloud-account IaC to deploy Lacework into a cloud environment.

This command creates Infrastructure as Code (IaC) in the form of Terraform HCL, with the option of running
Terraform and deploying Lacework into AWS, Azure, or GCP.


### Options

```
  -h, --help   help for cloud-account
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

* [lacework generate](lacework_generate.md)	 - Generate code to onboard your account
* [lacework generate cloud-account aws](lacework_generate_cloud-account_aws.md)	 - Generate and/or execute Terraform code for AWS integration
* [lacework generate cloud-account azure](lacework_generate_cloud-account_azure.md)	 - Generate and/or execute Terraform code for Azure integration
* [lacework generate cloud-account gcp](lacework_generate_cloud-account_gcp.md)	 - Generate and/or execute Terraform code for GCP integration
* [lacework generate cloud-account oci](lacework_generate_cloud-account_oci.md)	 - Generate and/or execute Terraform code for OCI integration

