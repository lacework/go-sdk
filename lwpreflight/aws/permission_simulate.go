package aws

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/smithy-go"
)

// simulateBatchSize caps how many actions are evaluated per
// SimulatePrincipalPolicy call. AWS allows more, but smaller batches keep
// individual responses well below the 1MB payload limit and avoid Marker
// pagination in the common case.
const simulateBatchSize = 50

// CheckPermissionsViaSimulation evaluates the required permissions against
// the caller's effective policies using IAM SimulatePrincipalPolicy. Unlike
// the static policy walk in CheckPermissions, this includes SCPs and
// permission boundaries because AWS performs the same evaluation it would for
// a real API call.
func CheckPermissionsViaSimulation(p *Preflight) error {
	if p.caller.IsRoot {
		// Root is unconstrained by IAM policies, SCPs, or boundaries.
		return nil
	}

	ctx := context.Background()
	iamSvc := iam.NewFromConfig(p.awsConfig)

	sourceArn, err := simulationPolicySourceArn(ctx, iamSvc, p.caller)
	if err != nil {
		return err
	}

	requiredPermissions := RequiredPermissions
	if p.isOrg {
		requiredPermissions = RequiredPermissionsForOrg
	}

	// Action -> integrations that require it. One simulation call can serve
	// multiple integrations and we attribute denials back to all of them.
	actionOwners := map[string][]IntegrationType{}
	for _, integrationType := range p.integrationTypes {
		for _, action := range requiredPermissions[integrationType] {
			actionOwners[action] = append(actionOwners[action], integrationType)
		}
	}
	if len(actionOwners) == 0 {
		return nil
	}

	actions := make([]string, 0, len(actionOwners))
	for action := range actionOwners {
		actions = append(actions, action)
	}
	sort.Strings(actions)

	p.verboseWriter.Write(fmt.Sprintf(
		"Simulating %d permission(s) against %s", len(actions), sourceArn,
	))

	for start := 0; start < len(actions); start += simulateBatchSize {
		end := start + simulateBatchSize
		if end > len(actions) {
			end = len(actions)
		}
		results, err := simulatePrincipalPolicyAll(ctx, iamSvc, sourceArn, actions[start:end])
		if err != nil {
			if isAWSAccessDenied(err) {
				return fmt.Errorf(
					"--simulate requires iam:SimulatePrincipalPolicy on the caller; "+
						"add it to the role's policy and re-run: %w", err)
			}
			return err
		}
		for _, result := range results {
			if result.EvalActionName == nil {
				continue
			}
			if isAllowed(result.EvalDecision) {
				continue
			}
			reason := simulationDenialReason(result)
			action := *result.EvalActionName
			for _, integrationType := range actionOwners[action] {
				p.errors[integrationType] = append(
					p.errors[integrationType],
					fmt.Sprintf("Required permission denied: %s (%s)", action, reason),
				)
			}
		}
	}

	return nil
}

// simulatePrincipalPolicyAll runs SimulatePrincipalPolicy for one batch of
// actions and follows the Marker through every page.
func simulatePrincipalPolicyAll(
	ctx context.Context, iamSvc *iam.Client, sourceArn string, actions []string,
) ([]types.EvaluationResult, error) {
	input := &iam.SimulatePrincipalPolicyInput{
		PolicySourceArn: aws.String(sourceArn),
		ActionNames:     actions,
	}
	out := []types.EvaluationResult{}
	for {
		page, err := iamSvc.SimulatePrincipalPolicy(ctx, input)
		if err != nil {
			return nil, err
		}
		out = append(out, page.EvaluationResults...)
		if page.Marker == nil || *page.Marker == "" {
			return out, nil
		}
		input.Marker = page.Marker
	}
}

// simulationPolicySourceArn picks the ARN to feed to SimulatePrincipalPolicy.
// SimulatePrincipalPolicy accepts user, group, or role ARNs but rejects
// assumed-role session ARNs, so we resolve the underlying role ARN via
// iam:GetRole when the caller is an assumed role.
func simulationPolicySourceArn(
	ctx context.Context, iamSvc *iam.Client, caller Caller,
) (string, error) {
	if !caller.IsAssumedRole() {
		return caller.ARN, nil
	}
	role, err := iamSvc.GetRole(ctx, &iam.GetRoleInput{RoleName: aws.String(caller.Name)})
	if err != nil {
		if isAWSAccessDenied(err) {
			return "", fmt.Errorf(
				"--simulate requires iam:GetRole on %s to translate the assumed-role "+
					"session into a role ARN; add iam:GetRole to the role's policy and "+
					"re-run: %w", caller.Name, err)
		}
		return "", fmt.Errorf("resolving role ARN for %s: %w", caller.Name, err)
	}
	if role.Role == nil || role.Role.Arn == nil {
		return "", fmt.Errorf("iam:GetRole returned no ARN for %s", caller.Name)
	}
	return *role.Role.Arn, nil
}

// isAWSAccessDenied reports whether err is an AWS API error indicating the
// caller lacks permission. IAM and STS use AccessDenied; some services use
// AccessDeniedException; EC2-style services use UnauthorizedOperation.
func isAWSAccessDenied(err error) bool {
	var apiErr smithy.APIError
	if !errors.As(err, &apiErr) {
		return false
	}
	switch apiErr.ErrorCode() {
	case "AccessDenied", "AccessDeniedException", "UnauthorizedOperation":
		return true
	}
	return false
}

func isAllowed(decision types.PolicyEvaluationDecisionType) bool {
	return decision == types.PolicyEvaluationDecisionTypeAllowed
}

// simulationDenialReason summarises why a simulated action was denied. SCPs
// and permission boundaries are checked first because their denials override
// any identity-policy allow.
func simulationDenialReason(r types.EvaluationResult) string {
	if r.OrganizationsDecisionDetail != nil && !r.OrganizationsDecisionDetail.AllowedByOrganizations {
		return "blocked by Organizations SCP"
	}
	if r.PermissionsBoundaryDecisionDetail != nil && !r.PermissionsBoundaryDecisionDetail.AllowedByPermissionsBoundary {
		return "blocked by permissions boundary"
	}
	switch r.EvalDecision {
	case types.PolicyEvaluationDecisionTypeExplicitDeny:
		return "explicit deny in identity policy"
	case types.PolicyEvaluationDecisionTypeImplicitDeny:
		return "no identity policy grants this action"
	}
	return string(r.EvalDecision)
}
