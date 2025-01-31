---
title: "lacework generate cloud-account aws controltower"
slug: lacework_generate_cloud-account_aws_controltower
hide_title: true
---

## lacework generate cloud-account aws controltower

Generate and/or execute Terraform code for ControlTower integration

### Synopsis

Use this command to generate Terraform code for deploying Lacework with Aws Cloudtrail and
ControlTower.

By default, this command interactively prompts for the required information to set up the new cloud account.
In interactive mode, this command will:

* Prompt for the required information to set up the integration
* Generate new Terraform code using the inputs
* Optionally, run the generated Terraform code:
  * If Terraform is already installed, the version is verified as compatible for use
  * If Terraform is not installed, or the version installed is not compatible, a new
    version will be installed into a temporary location
  * Once Terraform is detected or installed, the Terraform plan is executed
  * The command prompts you with the outcome of the plan and allows you to view more
    details or continue with Terraform apply
  * If confirmed, Terraform apply runs, completing the setup of the cloud account

This command can also be run in noninteractive mode.
See help output for more details on the parameter values required for Terraform code generation.


```
lacework generate cloud-account aws controltower [flags]
```

### Options

```
      --apply                            run terraform apply without executing plan or prompting
      --audit_account string             The audit account flag input in the format profile:region
  -h, --help                             help for controltower
      --iam_role_arn string              specify the arn of the existing iam role
      --iam_role_external_id string      specify the external id of the existing iam role
      --iam_role_name string             specify the name of the existing iam role
      --lacework_aws_account_id string   the Lacework AWS root account id
      --log_archive_account string       The log archive account flag input in the format profile:region
      --org_account_mapping string       Org account mapping json string. Example: '{"default_lacework_account":"main", "mapping": [{ "aws_accounts": ["123456789011"], "lacework_account": "sub-account-1"}]}'
      --output string                    location to write generated content
      --prefix string                    specify the prefix that will be used at the beginning of every generated resource
      --s3_bucket_arn string             the S3 Bucket for consolidated CloudTrail
      --sns_topic_arn string             the SNS Topic
      --sqs_queue_name string            specify the name of the sqs queue
```

### Options inherited from parent commands

```
  -a, --account string                            account subdomain of URL (i.e. <ACCOUNT>.lacework.net)
      --agentless                                 enable agentless integration
      --agentless_management_account_id string    AWS management account ID for Agentless integration
      --agentless_monitored_account_ids strings   AWS monitored account IDs for Agentless integrations; may contain account IDs, OUs, or the organization root (e.g. 123456789000,ou-abcd-12345678,r-abcd)
      --agentless_monitored_accounts strings      AWS monitored accounts for Agentless integrations; value format must be <aws profile>:<region>
      --agentless_scanning_accounts strings       AWS scanning accounts for Agentless integrations; value format must be <aws profile>:<region>
  -k, --api_key string                            access key id
  -s, --api_secret string                         secret access key
      --api_token string                          access token (replaces the use of api_key and api_secret)
      --aws_assume_role string                    specify aws assume role
      --aws_organization                          enable organization integration
      --aws_profile string                        specify aws profile
      --aws_region string                         specify aws region
      --aws_subaccount strings                    configure an additional aws account; value format must be <aws profile>:<region>
      --bucket_encryption_enabled                 enable S3 bucket encryption when creating bucket (default true)
      --bucket_name string                        specify bucket name when creating bucket
      --bucket_sse_key_arn string                 specify existing KMS encryption key arn for bucket
      --cloudtrail                                enable cloudtrail integration
      --cloudtrail_name string                    specify name of cloudtrail integration
      --cloudtrail_org_account_mapping string     Org account mapping json string. Example: '{"default_lacework_account":"main", "mapping": [{ "aws_accounts": ["123456789011"], "lacework_account": "sub-account-1"}]}'
      --config                                    enable config integration
      --config_cf_resource_prefix string          specify Cloudformation resource prefix for Config organization integration
      --config_lacework_access_key_id string      specify AWS access key ID for Config organization integration
      --config_lacework_account string            specify lacework account for Config organization integration
      --config_lacework_secret_key string         specify AWS secret key for Config organization integration
      --config_lacework_sub_account string        specify lacework sub-account for Config organization integration
      --config_organization_id string             specify AWS organization ID for Config organization integration
      --config_organization_units strings         specify AWS organization units for Config organization integration
      --consolidated_cloudtrail                   use consolidated trail
      --controltower                              enable Control Tower integration
      --controltower_audit_account string         specify AWS Control Tower Audit account; value format must be <aws profile>:<region>
      --controltower_kms_key_arn string           specify AWS Control Tower custom kMS key ARN
      --controltower_log_archive_account string   specify AWS Control Tower Log Archive account; value format must be <aws profile>:<region>
      --debug                                     turn on debug logging
      --existing_bucket_arn string                specify existing cloudtrail S3 bucket ARN
      --existing_iam_role_arn string              specify existing iam role arn to use
      --existing_iam_role_externalid string       specify existing iam role external_id to use
      --existing_iam_role_name string             specify existing iam role name to use
      --existing_sns_topic_arn string             specify existing SNS topic arn
      --json                                      switch commands output from human-readable to json format
      --nocache                                   turn off caching
      --nocolor                                   turn off colors
      --noninteractive                            turn off interactive mode (disable spinners, prompts, etc.)
      --organization                              access organization level data sets (org admins only)
  -p, --profile string                            switch between profiles configured at ~/.lacework.toml
      --sns_topic_encryption_enabled              enable encryption on SNS topic when creating one (default true)
      --sns_topic_encryption_key_arn string       specify existing KMS encryption key arn for SNS topic
      --sns_topic_name string                     specify SNS topic name if creating new one
      --sqs_encryption_enabled                    enable encryption on SQS queue when creating (default true)
      --sqs_encryption_key_arn string             specify existing KMS encryption key arn for SQS queue
      --subaccount string                         sub-account name inside your organization (org admins only)
      --use_s3_bucket_notification                enable S3 bucket notifications
```

### SEE ALSO

* [lacework generate cloud-account aws](lacework_generate_cloud-account_aws.md)	 - Generate and/or execute Terraform code for AWS integration

