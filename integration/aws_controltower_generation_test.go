//go:build !windows && generation

package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Netflix/go-expect"
	"github.com/lacework/go-sdk/cli/cmd"
	"github.com/lacework/go-sdk/lwgenerate/aws_controltower"
	"github.com/stretchr/testify/assert"
)

func assertControlTowerTerraformSaved(t *testing.T, message string) {
	assert.Contains(t, message, "Terraform code saved in")
}

const (
	controltowerPath = "/lacework/aws_controltower/"
)

func runControlTowerGenerateTest(t *testing.T, conditions func(*expect.Console), args ...string) string {
	os.Setenv("HOME", tfPath)

	hcl_path := filepath.Join(tfPath, controltowerPath, "main.tf")

	runFakeTerminalTestFromDir(t, tfPath, conditions, args...)
	out, err := os.ReadFile(hcl_path)
	if err != nil {
		return fmt.Sprintf("main.tf not found: %s", err)
	}

	t.Cleanup(func() {
		os.Remove(hcl_path)
	})

	result := terraformValidate(filepath.Join(tfPath, controltowerPath))

	assert.True(t, result.Valid)

	return string(out)
}

func TestGenerationControlTowerBasic(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runControlTowerGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionAwsControlTowerCoreS3Bucket, "arn:aws:s3:::aws-controltower-logs-0123456789-us-east-1"},
				MsgRsp{cmd.QuestionAwsControlTowerCoreSnsTopic, "arn:aws:sns:us-east-1:0123456789:aws-controltower-AllConfigNotifications"},
				MsgRsp{cmd.QuestionAwsControlTowerCoreLogProfile, "AWSAdministratorAccess"},
				MsgRsp{cmd.QuestionAwsControlTowerCoreLogRegion, "us-east-1"},
				MsgRsp{cmd.QuestionAwsControlTowerCoreAuditProfile, "AWSAdministratorAccess"},
				MsgRsp{cmd.QuestionAwsControlTowerCoreAuditRegion, "us-east-1"},
				MsgRsp{cmd.QuestionAwsControlTowerConfigureAdvanced, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"ca",
		"aws",
		"controltower",
	)

	assertControlTowerTerraformSaved(t, final)

	buildTf, _ := aws_controltower.NewTerraform(
		"arn:aws:s3:::aws-controltower-logs-0123456789-us-east-1",
		"arn:aws:sns:us-east-1:0123456789:aws-controltower-AllConfigNotifications",
		aws_controltower.WithSubaccounts(
			aws_controltower.NewAwsSubAccount("AWSAdministratorAccess", "us-east-1", "log_archive"),
			aws_controltower.NewAwsSubAccount("AWSAdministratorAccess", "us-east-1", "audit"))).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationControlTowerPromptOrgAccountMappings(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runControlTowerGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionAwsControlTowerCoreS3Bucket, "arn:aws:s3:::aws-controltower-logs-0123456789-us-east-1"},
				MsgRsp{cmd.QuestionAwsControlTowerCoreSnsTopic, "arn:aws:sns:us-east-1:0123456789:aws-controltower-AllConfigNotifications"},
				MsgRsp{cmd.QuestionAwsControlTowerCoreLogProfile, "AWSAdministratorAccess"},
				MsgRsp{cmd.QuestionAwsControlTowerCoreLogRegion, "us-east-1"},
				MsgRsp{cmd.QuestionAwsControlTowerCoreAuditProfile, "AWSAdministratorAccess"},
				MsgRsp{cmd.QuestionAwsControlTowerCoreAuditRegion, "us-east-1"},
				MsgRsp{cmd.QuestionAwsControlTowerConfigureAdvanced, "y"},
				MsgMenu{cmd.ControlTowerAdvancedOptMappings, 5},
				MsgRsp{cmd.QuestionControlTowerOrgAccountMappingsLWDefaultAccount, "main"},
				MsgRsp{cmd.QuestionControlTowerOrgAccountMappingsLWAccount, "sub-account-1"},
				MsgMulitlineRsp{cmd.QuestionControlTowerOrgAccountMappingsAwsAccounts, []string{"123456789011"}},
				MsgRsp{cmd.QuestionControlTowerOrgAccountMappingAnotherAdvancedOpt, "n"},
				MsgRsp{cmd.QuestionControlTowerAnotherAdvancedOpt, "n"},

				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"ca",
		"aws",
		"controltower",
	)

	assertControlTowerTerraformSaved(t, final)

	orgAccountMappings := aws_controltower.OrgAccountMapping{
		DefaultLaceworkAccount: "main",
		Mapping: []aws_controltower.OrgAccountMap{
			{
				LaceworkAccount: "sub-account-1",
				AwsAccounts:     []string{"123456789011"},
			},
		},
	}

	buildTf, _ := aws_controltower.NewTerraform(
		"arn:aws:s3:::aws-controltower-logs-0123456789-us-east-1",
		"arn:aws:sns:us-east-1:0123456789:aws-controltower-AllConfigNotifications",
		aws_controltower.WithSubaccounts(
			aws_controltower.NewAwsSubAccount("AWSAdministratorAccess", "us-east-1", "log_archive"),
			aws_controltower.NewAwsSubAccount("AWSAdministratorAccess", "us-east-1", "audit")),
		aws_controltower.WithOrgAccountMappings(orgAccountMappings)).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationControlTowerNonInteractive(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runControlTowerGenerateTest(t,
		func(c *expect.Console) {
			final, _ = c.ExpectEOF()
		},
		"generate",
		"ca",
		"aws",
		"controltower",
		"--s3_bucket_arn",
		"arn:aws:s3:::aws-controltower-logs-0123456789-us-east-1",
		"--sns_topic_arn",
		"arn:aws:sns:us-east-1:0123456789:aws-controltower-AllConfigNotifications",
		"--audit_account",
		"AWSAdministratorAccess:us-east-1",
		"--log_archive_account",
		"AWSAdministratorAccess:us-east-2",
		"--org_account_mapping",
		"{\"default_lacework_account\":\"main\", \"mapping\": [{ \"aws_accounts\": [\"123456789011\"], \"lacework_account\": \"sub-account-1\"}]}",
		"--noninteractive",
	)

	assertControlTowerTerraformSaved(t, final)

	orgAccountMappings := aws_controltower.OrgAccountMapping{
		DefaultLaceworkAccount: "main",
		Mapping: []aws_controltower.OrgAccountMap{
			{
				LaceworkAccount: "sub-account-1",
				AwsAccounts:     []string{"123456789011"},
			},
		},
	}

	buildTf, _ := aws_controltower.NewTerraform(
		"arn:aws:s3:::aws-controltower-logs-0123456789-us-east-1",
		"arn:aws:sns:us-east-1:0123456789:aws-controltower-AllConfigNotifications",
		aws_controltower.WithSubaccounts(
			aws_controltower.NewAwsSubAccount("AWSAdministratorAccess", "us-east-2", "log_archive"),
			aws_controltower.NewAwsSubAccount("AWSAdministratorAccess", "us-east-1", "audit")),
		aws_controltower.WithOrgAccountMappings(orgAccountMappings),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationControlTowerPromptOptionalAttributes(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runControlTowerGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionAwsControlTowerCoreS3Bucket, "arn:aws:s3:::aws-controltower-logs-0123456789-us-east-1"},
				MsgRsp{cmd.QuestionAwsControlTowerCoreSnsTopic, "arn:aws:sns:us-east-1:0123456789:aws-controltower-AllConfigNotifications"},
				MsgRsp{cmd.QuestionAwsControlTowerCoreLogProfile, "AWSAdministratorAccess"},
				MsgRsp{cmd.QuestionAwsControlTowerCoreLogRegion, "us-east-1"},
				MsgRsp{cmd.QuestionAwsControlTowerCoreAuditProfile, "AWSAdministratorAccess"},
				MsgRsp{cmd.QuestionAwsControlTowerCoreAuditRegion, "us-east-1"},
				MsgRsp{cmd.QuestionAwsControlTowerConfigureAdvanced, "y"},
				MsgMenu{cmd.ControlTowerAdvancedOptMappings, 5},
				MsgRsp{cmd.QuestionControlTowerOrgAccountMappingsLWDefaultAccount, "main"},
				MsgRsp{cmd.QuestionControlTowerOrgAccountMappingsLWAccount, "sub-account-1"},
				MsgMulitlineRsp{cmd.QuestionControlTowerOrgAccountMappingsAwsAccounts, []string{"123456789011"}},
				MsgRsp{cmd.QuestionControlTowerOrgAccountMappingAnotherAdvancedOpt, "n"},
				MsgRsp{cmd.QuestionControlTowerAnotherAdvancedOpt, "y"},
				MsgMenu{cmd.ControlTowerConfigureExistingIamRoleOpt, 0},
				MsgRsp{cmd.QuestionAwsControlTowerCoreIamRoleName, "lw-role-name"},
				MsgRsp{cmd.QuestionAwsControlTowerCoreIamRoleArn, "arn:aws:iam::12345678:role/test"},
				MsgRsp{cmd.QuestionAwsControlTowerCoreIamRoleExternalID, "01234567"},
				MsgRsp{cmd.QuestionControlTowerAnotherAdvancedOpt, "y"},
				MsgMenu{cmd.ControlTowerIntegrationPrefixOpt, 3},
				MsgRsp{cmd.QuestionControlTowerPrefix, "prefix-"},
				MsgRsp{cmd.QuestionControlTowerAnotherAdvancedOpt, "y"},
				MsgMenu{cmd.ControlTowerIntegrationSqsOpt, 4},
				MsgRsp{cmd.QuestionControlTowerSqsQueueName, "lw-queue-name"},
				MsgRsp{cmd.QuestionControlTowerAnotherAdvancedOpt, "n"},

				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"ca",
		"aws",
		"controltower",
	)

	assertControlTowerTerraformSaved(t, final)

	orgAccountMappings := aws_controltower.OrgAccountMapping{
		DefaultLaceworkAccount: "main",
		Mapping: []aws_controltower.OrgAccountMap{
			{
				LaceworkAccount: "sub-account-1",
				AwsAccounts:     []string{"123456789011"},
			},
		},
	}

	buildTf, _ := aws_controltower.NewTerraform(
		"arn:aws:s3:::aws-controltower-logs-0123456789-us-east-1",
		"arn:aws:sns:us-east-1:0123456789:aws-controltower-AllConfigNotifications",
		aws_controltower.WithSubaccounts(
			aws_controltower.NewAwsSubAccount("AWSAdministratorAccess", "us-east-1", "log_archive"),
			aws_controltower.NewAwsSubAccount("AWSAdministratorAccess", "us-east-1", "audit")),
		aws_controltower.WithOrgAccountMappings(orgAccountMappings),
		aws_controltower.WithPrefix("prefix-"),
		aws_controltower.WithSqsQueueName("lw-queue-name"),
		aws_controltower.WithExisitingIamRole("arn:aws:iam::12345678:role/test",
			"lw-role-name", "01234567")).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationControlTowerExistingIamRoleNonInteractive(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runControlTowerGenerateTest(t,
		func(c *expect.Console) {
			final, _ = c.ExpectEOF()
		},
		"generate",
		"ca",
		"aws",
		"controltower",
		"--s3_bucket_arn",
		"arn:aws:s3:::aws-controltower-logs-0123456789-us-east-1",
		"--sns_topic_arn",
		"arn:aws:sns:us-east-1:0123456789:aws-controltower-AllConfigNotifications",
		"--audit_account",
		"AWSAdministratorAccess:us-east-1",
		"--log_archive_account",
		"AWSAdministratorAccess:us-east-2",
		"--org_account_mapping",
		"{\"default_lacework_account\":\"main\", \"mapping\": [{ \"aws_accounts\": [\"123456789011\"], \"lacework_account\": \"sub-account-1\"}]}",
		"--iam_role_arn",
		"arn:aws:iam::12345678:role/test",
		"--iam_role_name",
		"lw-role-name",
		"--iam_role_external_id",
		"01234567",
		"--noninteractive",
	)

	assertControlTowerTerraformSaved(t, final)

	orgAccountMappings := aws_controltower.OrgAccountMapping{
		DefaultLaceworkAccount: "main",
		Mapping: []aws_controltower.OrgAccountMap{
			{
				LaceworkAccount: "sub-account-1",
				AwsAccounts:     []string{"123456789011"},
			},
		},
	}

	buildTf, _ := aws_controltower.NewTerraform(
		"arn:aws:s3:::aws-controltower-logs-0123456789-us-east-1",
		"arn:aws:sns:us-east-1:0123456789:aws-controltower-AllConfigNotifications",
		aws_controltower.WithSubaccounts(
			aws_controltower.NewAwsSubAccount("AWSAdministratorAccess", "us-east-2", "log_archive"),
			aws_controltower.NewAwsSubAccount("AWSAdministratorAccess", "us-east-1", "audit")),
		aws_controltower.WithOrgAccountMappings(orgAccountMappings),
		aws_controltower.WithExisitingIamRole("arn:aws:iam::12345678:role/test",
			"lw-role-name", "01234567")).Generate()
	assert.Equal(t, buildTf, tfResult)
}
