---
title: "lacework generate k8s eks"
slug: lacework_generate_k8s_eks
hide_title: true
---

## lacework generate k8s eks

Generate and/or execute Terraform code for EKS integration

### Synopsis

Use this command to generate Terraform code for deploying Lacework into an EKS
environment.

By default, this command interactively prompts for the required information to set up the new cloud account.
In interactive mode, this command will:

* Prompt for the required information to set up the integration
* Generate new Terraform code using the inputs
* Optionally, run the generated Terraform code:
  * If Terraform is already installed, the version is verified as compatible for use
  * If Terraform is not installed, or the version installed is not compatible, a new version will be installed into a temporary location
  * Once Terraform is detected or installed, the Terraform plan is executed
  * The command prompts you with the outcome of the plan and allows you to view more details or continue with Terraform apply
  * If confirmed, Terraform apply runs, completing the setup of the cloud account

This command can also be run in noninteractive mode.
See help output for more details on the parameter values required for Terraform code generation.


```
lacework generate k8s eks [flags]
```

### Options

```
      --apply                                     run terraform apply without executing plan or prompting
      --aws_profile string                        specify aws profile
      --bucket_lifecycle_exp_days int             specify the s3 bucket lifecycle expiration days
      --bucket_sse_algorithm string               specify the encryption algorithm to use for S3 bucket server-side encryption
      --bucket_sse_key_arn string                 specify the kms key arn to be used for s3. (required when bucket_sse_algorithm is aws:kms & using an existing kms key)
      --custom_filter_pattern string              specify a custom cloudwatch log filter pattern
      --enable_bucket_versioning                  enable s3 bucket versioning (default true)
      --enable_encryption_s3                      enable encryption on s3 bucket (default true)
      --enable_firehose_encryption                enable firehose encryption (default true)
      --enable_force_destroy                      enable force destroy s3 bucket
      --enable_kms_key_rotation                   enable automatic kms key rotation (default true)
      --enable_mfa_delete_s3                      enable mfa delete on s3 bucket. Requires bucket versioning.
      --enable_sns_topic_encryption               enable encryption on the sns topic (default true)
      --existing_bucket_arn string                specify existing s3 bucket arn for the audit log
      --existing_ca_iam_role_arn string           specify existing cross account iam role arn to use
      --existing_ca_iam_role_external_id string   specify existing cross account iam role external_id to use
      --existing_cw_iam_role_arn string           specify existing cloudwatch iam role arn to use
      --existing_firehose_iam_role_arn string     specify existing firehose iam role arn to use
      --firehose_encryption_key_arn string        specify the kms key arn to be used with the Firehose
  -h, --help                                      help for eks
      --integration_name string                   specify the name of the eks audit integration
      --kms_key_deletion_days int                 specify the kms waiting period before deletion, in number of days
      --lacework_aws_account_id string            the Lacework AWS root account id
      --output string                             location to write generated content
      --prefix string                             specify the prefix that will be used at the beginning of every generated resource
      --region_clusters stringToString            configure eks clusters per aws region. To configure multiple regions pass the flag multiple times. Example format:  --region_clusters <region>="cluster,list" (default [])
      --sns_topic_encryption_key_arn string       specify the kms key arn to be used with the sns topic
      --use_existing_bucket                       use existing supplied s3 bucket (default false)
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

