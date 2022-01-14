package cmd

import (
	"strings"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/stretchr/testify/assert"
)

func buildPlanOutput(resourceChanges ...*tfjson.ResourceChange) *tfjson.Plan {
	return &tfjson.Plan{
		FormatVersion:    "0.2",
		TerraformVersion: "1.0.9",
		ResourceChanges:  resourceChanges,
	}
}

func resourceChange(address string, actionType ...tfjson.Action) *tfjson.ResourceChange {
	addressDetails := strings.Split(address, ".")
	return &tfjson.ResourceChange{
		Address:       address,
		ModuleAddress: addressDetails[0],
		Mode:          "managed",
		Type:          addressDetails[len(addressDetails)-2],
		Name:          addressDetails[len(addressDetails)-1],
		ProviderName:  "registry.terraform.io/foo/foo",
		Change: &tfjson.Change{
			Actions: actionType,
			Before:  nil,
			After:   nil,
		},
	}
}

// Test parsing Terraform Plan output
func TestGenerationExecutionPlanOutputParsing(t *testing.T) {
	t.Run("single state should report correctly", func(t *testing.T) {
		result := parseTfPlanOutput(buildPlanOutput(resourceChange("module.aws.config.aws_iam_policy.lacework_audit_policy", tfjson.ActionCreate)))
		assert.Equal(t, 1, result.create)
	})
	t.Run("multiple state of same type should report correctly", func(t *testing.T) {
		result := parseTfPlanOutput(
			buildPlanOutput(
				resourceChange("module.aws.config.aws_iam_policy.lacework_audit_policy", tfjson.ActionCreate),
				resourceChange("module.aws.config.aws_iam_policy.lacework_audit_policy1", tfjson.ActionCreate),
				resourceChange("module.aws.config.aws_iam_policy.lacework_audit_policy2", tfjson.ActionCreate),
			))
		assert.Equal(t, 3, result.create)
	})
	t.Run("replace state should report correctly", func(t *testing.T) {
		result := parseTfPlanOutput(
			buildPlanOutput(
				resourceChange("module.aws.config.aws_iam_policy.lacework_audit_policy3", tfjson.ActionDelete, tfjson.ActionCreate), // Replace
			))
		assert.Equal(t, 1, result.replace)
	})
	t.Run("multiple states should report correctly", func(t *testing.T) {
		result := parseTfPlanOutput(
			buildPlanOutput(
				resourceChange("module.aws.config.aws_iam_policy.lacework_audit_policy", tfjson.ActionCreate),
				resourceChange("module.aws.config.aws_iam_policy.lacework_audit_policy1", tfjson.ActionDelete),
				resourceChange("module.aws.config.aws_iam_policy.lacework_audit_policy2", tfjson.ActionUpdate),
				resourceChange("module.aws.config.aws_iam_policy.lacework_audit_policy3", tfjson.ActionDelete, tfjson.ActionCreate), // Replace
			))
		assert.Equal(t, 1, result.create)
		assert.Equal(t, 1, result.deleted)
		assert.Equal(t, 1, result.update)
		assert.Equal(t, 1, result.replace)
	})
	t.Run("correctly parse resource change states for human output", func(t *testing.T) {
		plan := buildPlanOutput(
			resourceChange("module.test.sub.resource.name_goes_here", tfjson.ActionCreate))
		assert.Equal(t, `Resource details: 
  resource.name_goes_here will be created

More details can be viewed by running:

  cd /place
  /var/tmp/tf show tfplan.json

`,
			buildHumanReadablePlannedActions("/place", "/var/tmp/tf", plan.ResourceChanges))
	})
	t.Run("correctly parse resource changes states for human output", func(t *testing.T) {
		plan := buildPlanOutput(
			resourceChange("module.aws.config.aws_iam_policy.lacework_audit_policy", tfjson.ActionCreate),
			resourceChange("module.aws.config.aws_iam_policy.lacework_audit_policy1", tfjson.ActionDelete),
			resourceChange("module.aws.config.aws_iam_policy.lacework_audit_policy2", tfjson.ActionUpdate),
			resourceChange("module.aws.config.aws_iam_policy.lacework_audit_policy3", tfjson.ActionDelete, tfjson.ActionCreate))
		assert.Equal(t, `Resource details: 
  aws_iam_policy.lacework_audit_policy will be created
  aws_iam_policy.lacework_audit_policy1 will be destroyed
  aws_iam_policy.lacework_audit_policy2 will be update
  aws_iam_policy.lacework_audit_policy3 will be replaced

More details can be viewed by running:

  cd /tmp
  terraform show tfplan.json

`,
			buildHumanReadablePlannedActions("/tmp", "terraform", plan.ResourceChanges))
	})
}
