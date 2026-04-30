// Author:: Fortinet
// Copyright:: Copyright 2026, Fortinet
// License:: Apache License, Version 2.0

package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

// simulatorBatchSize is the number of action names we send per
// SimulatePrincipalPolicy call. The API has no documented hard limit on
// ActionNames length, but it caps the response page size — keeping batches
// small avoids pagination and stays well under any service-side caps.
const simulatorBatchSize = 50

// SimulatePrincipalPolicyAPI is the subset of the IAM client surface used by
// the simulator-based permission check. Defining it as an interface lets us
// mock the AWS call in unit tests.
type SimulatePrincipalPolicyAPI interface {
	SimulatePrincipalPolicy(
		ctx context.Context,
		params *iam.SimulatePrincipalPolicyInput,
		optFns ...func(*iam.Options),
	) (*iam.SimulatePrincipalPolicyOutput, error)
}

// CheckPermissionsViaSimulator is the simulator-backed alternative to
// CheckPermissions. Instead of parsing the caller's IAM policies locally and
// matching them against the required list, it asks AWS to evaluate each
// required action against the caller's effective authorization — which
// includes Service Control Policies, permissions boundaries, and other deny
// gates that local parsing cannot see. Results are written into p.errors in
// the same shape CheckPermissions produces, so downstream rendering does not
// change.
func CheckPermissionsViaSimulator(p *Preflight) error {
	if p.caller.IsAdmin {
		return nil
	}

	sourceArn, err := policySourceArnFromCaller(&p.caller)
	if err != nil {
		return fmt.Errorf("unable to derive simulator source ARN: %w", err)
	}

	iamSvc := iam.NewFromConfig(p.awsConfig)

	for _, integrationType := range p.integrationTypes {
		p.verboseWriter.Write(fmt.Sprintf(
			"Simulating permissions for %s against %s", integrationType, sourceArn,
		))

		requiredPermissions := RequiredPermissions
		if p.isOrg {
			requiredPermissions = RequiredPermissionsForOrg
		}

		actions := requiredPermissions[integrationType]
		if len(actions) == 0 {
			continue
		}

		denied, err := simulatePrincipalPermissions(
			context.Background(), iamSvc, sourceArn, actions,
		)
		if err != nil {
			return err
		}

		for _, action := range actions {
			if reason, blocked := denied[action]; blocked {
				p.errors[integrationType] = append(
					p.errors[integrationType],
					fmt.Sprintf("Required permission missing: %s (%s)", action, reason),
				)
			}
		}
	}

	return nil
}

// simulatePrincipalPermissions evaluates each action against sourceArn using
// iam:SimulatePrincipalPolicy. It returns a map keyed by action name for any
// action that did NOT evaluate to "allowed" — the value is a short
// human-readable reason ("explicit deny" vs. "no allow"). Actions that
// evaluated to "allowed" are absent from the returned map.
func simulatePrincipalPermissions(
	ctx context.Context,
	client SimulatePrincipalPolicyAPI,
	sourceArn string,
	actions []string,
) (map[string]string, error) {
	denied := map[string]string{}

	for start := 0; start < len(actions); start += simulatorBatchSize {
		end := start + simulatorBatchSize
		if end > len(actions) {
			end = len(actions)
		}
		batch := actions[start:end]

		// ResourceArns=["*"] simulates against any resource. SCPs/policies
		// that deny only specific resource patterns may produce false
		// negatives here — that is a documented limitation we accept for v1.
		input := &iam.SimulatePrincipalPolicyInput{
			PolicySourceArn: aws.String(sourceArn),
			ActionNames:     batch,
			ResourceArns:    []string{"*"},
		}

		marker := (*string)(nil)
		for {
			input.Marker = marker
			out, err := client.SimulatePrincipalPolicy(ctx, input)
			if err != nil {
				return nil, fmt.Errorf("SimulatePrincipalPolicy failed: %w", err)
			}

			for _, result := range out.EvaluationResults {
				if result.EvalActionName == nil {
					continue
				}
				action := *result.EvalActionName
				switch result.EvalDecision {
				case types.PolicyEvaluationDecisionTypeAllowed:
					// granted — skip
				case types.PolicyEvaluationDecisionTypeExplicitDeny:
					denied[action] = "explicit deny — likely SCP, permissions boundary, or policy"
				default:
					// implicitDeny or any future value AWS adds — treat as missing.
					denied[action] = "no allow grants this action"
				}
			}

			if out.IsTruncated && out.Marker != nil {
				marker = out.Marker
				continue
			}
			break
		}
	}

	return denied, nil
}

// policySourceArnFromCaller returns an IAM ARN suitable for the
// PolicySourceArn parameter of SimulatePrincipalPolicy. STS assumed-role
// session ARNs (arn:aws:sts::123:assumed-role/Role/Session) are not accepted
// by the simulator — we collapse them to the underlying IAM role ARN
// (arn:aws:iam::123:role/Role). User and role IAM ARNs are returned as-is.
func policySourceArnFromCaller(c *Caller) (string, error) {
	parsed, err := arn.Parse(c.ARN)
	if err != nil {
		return "", err
	}
	if c.IsAssumedRole() {
		return fmt.Sprintf("arn:aws:iam::%s:role/%s", parsed.AccountID, c.Name), nil
	}
	return c.ARN, nil
}
