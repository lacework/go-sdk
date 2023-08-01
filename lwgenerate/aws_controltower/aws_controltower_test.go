package aws_controltower

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerationControlTowerCloudTrail(t *testing.T) {
	hcl, err := NewTerraform(
		"arn:aws:s3:::aws-controltower-logs-0123456789-us-east-1",
		"arn:aws:sns:us-east-1:0123456789:aws-controltower-AllConfigNotifications",
		WithSubaccounts(
			NewAwsSubAccount("AWSAdministratorAccess", "us-east-1", "log_archive"),
			NewAwsSubAccount("AWSAdministratorAccess", "us-east-1", "audit")),
		WithLaceworkOrgLevel()).Generate()

	assert.NoError(t, err)
	assert.Equal(t, hcl, moduleCloudtrailControlTowerBasic)
}

func TestGenerationControlTowerCloudTrailWithOrgAccountMappings(t *testing.T) {
	orgAccountMappings := OrgAccountMapping{
		DefaultLaceworkAccount: "main",
		Mapping: []OrgAccountMap{
			{
				LaceworkAccount: "sub-account-1",
				AwsAccounts:     []string{"123456789011"},
			},
			{
				LaceworkAccount: "sub-account-2",
				AwsAccounts:     []string{"123456789012"},
			},
		},
	}

	hcl, err := NewTerraform(
		"arn:aws:s3:::aws-controltower-logs-0123456789-us-east-1",
		"arn:aws:sns:us-east-1:0123456789:aws-controltower-AllConfigNotifications",
		WithSubaccounts(
			NewAwsSubAccount("AWSAdministratorAccess", "us-east-1", "log_archive"),
			NewAwsSubAccount("AWSAdministratorAccess", "us-east-1", "audit")),
		WithLaceworkOrgLevel(),
		WithOrgAccountMappings(orgAccountMappings)).Generate()

	assert.NoError(t, err)
	assert.Equal(t, hcl, moduleCloudtrailControlTowerWithAccountMappings)
}

func TestCreateHclAttributeFromStruct(t *testing.T) {
	orgAccountMappings := OrgAccountMapping{
		DefaultLaceworkAccount: "main",
		Mapping: []OrgAccountMap{
			{
				LaceworkAccount: "sub-account-1",
				AwsAccounts:     []string{"123456789011"},
			},
			{
				LaceworkAccount: "sub-account-2",
				AwsAccounts:     []string{"123456789012"},
			},
		},
	}

	expected := map[string]interface{}{ // support deeply nested data types
		"default_lacework_account": "main",
		"mapping": []map[string]interface{}{
			{
				"lacework_account": "sub-account-1",
				"aws_accounts":     []any{"123456789011"},
			},
			{
				"lacework_account": "sub-account-2",
				"aws_accounts":     []any{"123456789012"},
			},
		},
	}

	attribute, err := orgAccountMappings.ToMap()
	assert.NoError(t, err)
	assert.Equal(t, attribute, expected)
}

var moduleCloudtrailControlTowerBasic = `terraform {
  required_providers {
    lacework = {
      source  = "lacework/lacework"
      version = "~> 1.0"
    }
  }
}

provider "aws" {
  alias   = "log_archive"
  profile = "AWSAdministratorAccess"
  region  = "us-east-1"
}

provider "aws" {
  alias   = "audit"
  profile = "AWSAdministratorAccess"
  region  = "us-east-1"
}

provider "lacework" {
  organization = true
}

module "lacework_aws_controltower" {
  source        = "lacework/cloudtrail-controltower/aws"
  version       = "~> 0.3"
  s3_bucket_arn = "arn:aws:s3:::aws-controltower-logs-0123456789-us-east-1"
  sns_topic_arn = "arn:aws:sns:us-east-1:0123456789:aws-controltower-AllConfigNotifications"

  providers = {
    aws.audit       = aws.audit
    aws.log_archive = aws.log_archive
  }
}
`

var moduleCloudtrailControlTowerWithAccountMappings = `terraform {
  required_providers {
    lacework = {
      source  = "lacework/lacework"
      version = "~> 1.0"
    }
  }
}

provider "aws" {
  alias   = "log_archive"
  profile = "AWSAdministratorAccess"
  region  = "us-east-1"
}

provider "aws" {
  alias   = "audit"
  profile = "AWSAdministratorAccess"
  region  = "us-east-1"
}

provider "lacework" {
  organization = true
}

module "lacework_aws_controltower" {
  source  = "lacework/cloudtrail-controltower/aws"
  version = "~> 0.3"
  org_account_mappings = [{
    default_lacework_account = "main"
    mapping = [{
      aws_accounts     = ["123456789011"]
      lacework_account = "sub-account-1"
      }, {
      aws_accounts     = ["123456789012"]
      lacework_account = "sub-account-2"
    }]
  }]
  s3_bucket_arn = "arn:aws:s3:::aws-controltower-logs-0123456789-us-east-1"
  sns_topic_arn = "arn:aws:sns:us-east-1:0123456789:aws-controltower-AllConfigNotifications"

  providers = {
    aws.audit       = aws.audit
    aws.log_archive = aws.log_archive
  }
}
`
