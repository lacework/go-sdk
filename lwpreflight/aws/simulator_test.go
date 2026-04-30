// Author:: Fortinet
// Copyright:: Copyright 2026, Fortinet
// License:: Apache License, Version 2.0

package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/stretchr/testify/assert"
)

// fakeSimulator records each SimulatePrincipalPolicy call and returns
// canned decisions keyed by action name.
type fakeSimulator struct {
	decisions map[string]types.PolicyEvaluationDecisionType
	calls     []iam.SimulatePrincipalPolicyInput
}

func (f *fakeSimulator) SimulatePrincipalPolicy(
	_ context.Context,
	input *iam.SimulatePrincipalPolicyInput,
	_ ...func(*iam.Options),
) (*iam.SimulatePrincipalPolicyOutput, error) {
	f.calls = append(f.calls, *input)

	results := make([]types.EvaluationResult, 0, len(input.ActionNames))
	for _, action := range input.ActionNames {
		decision := types.PolicyEvaluationDecisionTypeAllowed
		if d, ok := f.decisions[action]; ok {
			decision = d
		}
		a := action
		results = append(results, types.EvaluationResult{
			EvalActionName: &a,
			EvalDecision:   decision,
		})
	}
	return &iam.SimulatePrincipalPolicyOutput{EvaluationResults: results}, nil
}

func TestSimulatePrincipalPermissions_Classification(t *testing.T) {
	stub := &fakeSimulator{
		decisions: map[string]types.PolicyEvaluationDecisionType{
			"s3:GetObject":        types.PolicyEvaluationDecisionTypeAllowed,
			"s3:DeleteBucket":     types.PolicyEvaluationDecisionTypeExplicitDeny,
			"kms:CreateKey":       types.PolicyEvaluationDecisionTypeImplicitDeny,
			"ec2:DescribeRegions": types.PolicyEvaluationDecisionTypeAllowed,
		},
	}

	denied, err := simulatePrincipalPermissions(
		context.Background(),
		stub,
		"arn:aws:iam::123456789012:role/test",
		[]string{"s3:GetObject", "s3:DeleteBucket", "kms:CreateKey", "ec2:DescribeRegions"},
	)
	assert.NoError(t, err)

	// allowed actions should not appear in denied map
	_, ok := denied["s3:GetObject"]
	assert.False(t, ok, "allowed action should not be denied")
	_, ok = denied["ec2:DescribeRegions"]
	assert.False(t, ok)

	assert.Contains(t, denied, "s3:DeleteBucket")
	assert.Contains(t, denied["s3:DeleteBucket"], "explicit deny")

	assert.Contains(t, denied, "kms:CreateKey")
	assert.Contains(t, denied["kms:CreateKey"], "no allow")
}

func TestSimulatePrincipalPermissions_Batches(t *testing.T) {
	// Build 130 actions; with simulatorBatchSize=50 we expect 3 calls
	// (50 + 50 + 30).
	actions := make([]string, 130)
	for i := range actions {
		actions[i] = "svc:Action"
	}

	stub := &fakeSimulator{decisions: map[string]types.PolicyEvaluationDecisionType{}}
	_, err := simulatePrincipalPermissions(
		context.Background(),
		stub,
		"arn:aws:iam::123456789012:role/test",
		actions,
	)
	assert.NoError(t, err)
	assert.Len(t, stub.calls, 3, "expected 3 batches for 130 actions at batch size 50")
	assert.Len(t, stub.calls[0].ActionNames, 50)
	assert.Len(t, stub.calls[1].ActionNames, 50)
	assert.Len(t, stub.calls[2].ActionNames, 30)
}

func TestSimulatePrincipalPermissions_PassesResourceArn(t *testing.T) {
	stub := &fakeSimulator{}
	_, err := simulatePrincipalPermissions(
		context.Background(),
		stub,
		"arn:aws:iam::123456789012:role/test",
		[]string{"s3:GetObject"},
	)
	assert.NoError(t, err)
	assert.Equal(t, []string{"*"}, stub.calls[0].ResourceArns,
		"simulator should evaluate against ResourceArns=[\"*\"]")
	assert.Equal(t,
		"arn:aws:iam::123456789012:role/test",
		aws.ToString(stub.calls[0].PolicySourceArn),
	)
}

type paginatedSimulator struct {
	page  int
	calls int
}

func (p *paginatedSimulator) SimulatePrincipalPolicy(
	_ context.Context,
	input *iam.SimulatePrincipalPolicyInput,
	_ ...func(*iam.Options),
) (*iam.SimulatePrincipalPolicyOutput, error) {
	p.calls++

	// Return all actions allowed on first page, then one denied on second.
	switch p.page {
	case 0:
		p.page++
		marker := "next"
		return &iam.SimulatePrincipalPolicyOutput{
			EvaluationResults: []types.EvaluationResult{
				{EvalActionName: aws.String("s3:GetObject"), EvalDecision: types.PolicyEvaluationDecisionTypeAllowed},
			},
			IsTruncated: true,
			Marker:      &marker,
		}, nil
	default:
		return &iam.SimulatePrincipalPolicyOutput{
			EvaluationResults: []types.EvaluationResult{
				{EvalActionName: aws.String("s3:DeleteBucket"), EvalDecision: types.PolicyEvaluationDecisionTypeExplicitDeny},
			},
		}, nil
	}
}

func TestSimulatePrincipalPermissions_FollowsPagination(t *testing.T) {
	stub := &paginatedSimulator{}
	denied, err := simulatePrincipalPermissions(
		context.Background(),
		stub,
		"arn:aws:iam::123456789012:role/test",
		[]string{"s3:GetObject", "s3:DeleteBucket"},
	)
	assert.NoError(t, err)
	assert.Equal(t, 2, stub.calls, "expected pagination follow-up call")
	assert.Contains(t, denied, "s3:DeleteBucket")
	assert.NotContains(t, denied, "s3:GetObject")
}

func TestPolicySourceArnFromCaller(t *testing.T) {
	cases := []struct {
		name string
		in   Caller
		want string
	}{
		{
			name: "user ARN passes through",
			in: Caller{
				ARN:       "arn:aws:iam::123456789012:user/alice",
				AccountID: "123456789012",
				Name:      "alice",
			},
			want: "arn:aws:iam::123456789012:user/alice",
		},
		{
			name: "role ARN passes through",
			in: Caller{
				ARN:       "arn:aws:iam::123456789012:role/my-role",
				AccountID: "123456789012",
				Name:      "my-role",
			},
			want: "arn:aws:iam::123456789012:role/my-role",
		},
		{
			name: "assumed-role STS ARN collapses to underlying role",
			in: Caller{
				ARN:       "arn:aws:sts::123456789012:assumed-role/preflight_ro/aws-go-sdk-1234",
				AccountID: "123456789012",
				Name:      "preflight_ro",
			},
			want: "arn:aws:iam::123456789012:role/preflight_ro",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := policySourceArnFromCaller(&tc.in)
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestPolicySourceArnFromCaller_Invalid(t *testing.T) {
	_, err := policySourceArnFromCaller(&Caller{ARN: "not-an-arn"})
	assert.Error(t, err)
}
