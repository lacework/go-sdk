package aws

import (
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/smithy-go"
	"github.com/stretchr/testify/assert"
)

func TestSimulationDenialReason(t *testing.T) {
	t.Run("SCP denial wins over identity decision", func(t *testing.T) {
		got := simulationDenialReason(types.EvaluationResult{
			EvalDecision: types.PolicyEvaluationDecisionTypeImplicitDeny,
			OrganizationsDecisionDetail: &types.OrganizationsDecisionDetail{
				AllowedByOrganizations: false,
			},
		})
		assert.Equal(t, "blocked by Organizations SCP", got)
	})

	t.Run("permissions boundary denial wins over identity decision", func(t *testing.T) {
		got := simulationDenialReason(types.EvaluationResult{
			EvalDecision: types.PolicyEvaluationDecisionTypeImplicitDeny,
			PermissionsBoundaryDecisionDetail: &types.PermissionsBoundaryDecisionDetail{
				AllowedByPermissionsBoundary: false,
			},
		})
		assert.Equal(t, "blocked by permissions boundary", got)
	})

	t.Run("explicit deny in identity policy", func(t *testing.T) {
		got := simulationDenialReason(types.EvaluationResult{
			EvalDecision: types.PolicyEvaluationDecisionTypeExplicitDeny,
		})
		assert.Equal(t, "explicit deny in identity policy", got)
	})

	t.Run("implicit deny in identity policy", func(t *testing.T) {
		got := simulationDenialReason(types.EvaluationResult{
			EvalDecision: types.PolicyEvaluationDecisionTypeImplicitDeny,
		})
		assert.Equal(t, "no identity policy grants this action", got)
	})

	t.Run("SCP allowed but identity denies still reports identity reason", func(t *testing.T) {
		got := simulationDenialReason(types.EvaluationResult{
			EvalDecision: types.PolicyEvaluationDecisionTypeImplicitDeny,
			OrganizationsDecisionDetail: &types.OrganizationsDecisionDetail{
				AllowedByOrganizations: true,
			},
		})
		assert.Equal(t, "no identity policy grants this action", got)
	})
}

func TestIsAllowed(t *testing.T) {
	assert.True(t, isAllowed(types.PolicyEvaluationDecisionTypeAllowed))
	assert.False(t, isAllowed(types.PolicyEvaluationDecisionTypeExplicitDeny))
	assert.False(t, isAllowed(types.PolicyEvaluationDecisionTypeImplicitDeny))
}

func TestIsAWSAccessDenied(t *testing.T) {
	t.Run("nil is not access denied", func(t *testing.T) {
		assert.False(t, isAWSAccessDenied(nil))
	})

	t.Run("plain error is not access denied", func(t *testing.T) {
		assert.False(t, isAWSAccessDenied(errors.New("network blip")))
	})

	t.Run("IAM AccessDenied is detected", func(t *testing.T) {
		err := &smithy.GenericAPIError{Code: "AccessDenied", Message: "no"}
		assert.True(t, isAWSAccessDenied(err))
	})

	t.Run("AccessDeniedException is detected", func(t *testing.T) {
		err := &smithy.GenericAPIError{Code: "AccessDeniedException", Message: "no"}
		assert.True(t, isAWSAccessDenied(err))
	})

	t.Run("EC2-style UnauthorizedOperation is detected", func(t *testing.T) {
		err := &smithy.GenericAPIError{Code: "UnauthorizedOperation", Message: "no"}
		assert.True(t, isAWSAccessDenied(err))
	})

	t.Run("unrelated API error is not access denied", func(t *testing.T) {
		err := &smithy.GenericAPIError{Code: "InvalidInput", Message: "bad"}
		assert.False(t, isAWSAccessDenied(err))
	})

	t.Run("wrapped AccessDenied is still detected", func(t *testing.T) {
		inner := &smithy.GenericAPIError{Code: "AccessDenied", Message: "no"}
		wrapped := fmt.Errorf("calling iam:GetRole: %w", inner)
		assert.True(t, isAWSAccessDenied(wrapped))
	})
}

// Compile-time check that the EvaluationResult fields we depend on still
// have the types we expect from the IAM SDK. Guards against a vendored SDK
// upgrade silently changing the shape of the API we read.
var _ = func() *string {
	r := types.EvaluationResult{EvalActionName: aws.String("iam:GetRole")}
	return r.EvalActionName
}
