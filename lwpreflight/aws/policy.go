package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/lacework/go-sdk/v2/lwpreflight/logger"
	"github.com/mitchellh/mapstructure"
)

type Policy struct {
	Version   string
	Statement []StatementEntry
}

type StatementEntry struct {
	Effect   string
	Action   []string
	Resource []string
}

func FetchPolicies(p *Preflight) error {
	var ctx = context.Background()

	if p.caller.IsRoot {
		logger.Log.Info("Skip fetching IAM policies for root user")
		p.verboseWriter.Write("Skip fetching IAM policies for root user")
		return nil
	}

	iamSvc := iam.NewFromConfig(p.awsConfig)
	documents := []string{}

	if p.caller.IsAssumedRole() {
		p.verboseWriter.Write(fmt.Sprintf("Discovering managed IAM policies for %s", p.caller.Name))
		docs, err := fetchManangedRolePolicies(ctx, iamSvc, p.caller.Name)
		if err != nil {
			return err
		}
		documents = append(documents, docs...)

		p.verboseWriter.Write(fmt.Sprintf("Discovering inline IAM policies for %s", p.caller.Name))
		docs, err = fetchInlineRolePolicies(ctx, iamSvc, p.caller.Name)
		if err != nil {
			return err
		}
		documents = append(documents, docs...)
	} else {
		p.verboseWriter.Write(fmt.Sprintf("Discovering managed IAM policies for %s", p.caller.Name))

		docs, err := fetchManagedUserPolicies(ctx, iamSvc, p.caller.Name)
		if err != nil {
			return err
		}
		documents = append(documents, docs...)

		p.verboseWriter.Write(fmt.Sprintf("Discovering inline IAM policies for %s", p.caller.Name))
		docs, err = fetchInlineUserPolicies(ctx, iamSvc, p.caller.Name)
		if err != nil {
			return err
		}
		documents = append(documents, docs...)

		p.verboseWriter.Write(fmt.Sprintf("Discovering IAM groups for %s", p.caller.Name))
		docs, err = fetchUserGroupPolicies(ctx, iamSvc, p.caller.Name)
		if err != nil {
			return err
		}
		documents = append(documents, docs...)
	}

	// Look through policy statements and set permissions
	for _, document := range documents {
		policy, err := decodePolicyDocument(document)
		if err != nil {
			return err
		}
		for _, statement := range policy.Statement {
			for _, a := range statement.Action {
				if a == "*" {
					p.caller.IsAdmin = true
					return nil
				}
				if strings.Contains(a, "*") {
					p.permissionsWithWildcard = append(p.permissionsWithWildcard, a)
				} else {
					p.permissions[a] = true
				}
			}
		}
	}

	return nil
}

func fetchManangedRolePolicies(ctx context.Context, svc *iam.Client, roleName string) ([]string, error) {
	output, err := svc.ListAttachedRolePolicies(ctx, &iam.ListAttachedRolePoliciesInput{
		RoleName: aws.String(roleName),
	})
	if err != nil {
		return nil, err
	}

	documents := []string{}
	for _, policy := range output.AttachedPolicies {
		doc, err := fetchPolicyDocument(ctx, svc, policy.PolicyArn)
		if err != nil {
			return nil, err
		}
		documents = append(documents, doc)
	}

	return documents, nil
}

func fetchInlineRolePolicies(ctx context.Context, svc *iam.Client, roleName string) ([]string, error) {
	output, err := svc.ListRolePolicies(ctx, &iam.ListRolePoliciesInput{RoleName: aws.String(roleName)})
	if err != nil {
		return nil, err
	}

	documents := []string{}
	for _, policyName := range output.PolicyNames {
		rolePolicy, err := svc.GetRolePolicy(ctx, &iam.GetRolePolicyInput{
			PolicyName: aws.String(policyName),
			RoleName:   aws.String(roleName),
		})
		if err != nil {
			return nil, err
		}
		doc, err := url.QueryUnescape(*rolePolicy.PolicyDocument)
		if err != nil {
			return nil, err
		}
		documents = append(documents, doc)
	}

	return documents, nil
}

func fetchManagedUserPolicies(ctx context.Context, svc *iam.Client, userName string) ([]string, error) {
	output, err := svc.ListAttachedUserPolicies(ctx, &iam.ListAttachedUserPoliciesInput{
		UserName: &userName,
	})
	if err != nil {
		return nil, err
	}

	documents := []string{}
	for _, policy := range output.AttachedPolicies {
		doc, err := fetchPolicyDocument(ctx, svc, policy.PolicyArn)
		if err != nil {
			return nil, err
		}
		documents = append(documents, doc)
	}

	return documents, nil
}

func fetchInlineUserPolicies(ctx context.Context, svc *iam.Client, userName string) ([]string, error) {
	output, err := svc.ListUserPolicies(ctx, &iam.ListUserPoliciesInput{UserName: aws.String(userName)})
	if err != nil {
		return nil, err
	}

	documents := []string{}
	for _, policyName := range output.PolicyNames {
		policy, err := svc.GetUserPolicy(ctx, &iam.GetUserPolicyInput{
			PolicyName: aws.String(policyName),
			UserName:   aws.String(userName),
		})
		if err != nil {
			return nil, err
		}
		doc, err := url.QueryUnescape(*policy.PolicyDocument)
		if err != nil {
			return nil, err
		}
		documents = append(documents, doc)
	}

	return documents, nil
}

func fetchUserGroupPolicies(ctx context.Context, svc *iam.Client, userName string) ([]string, error) {
	output, err := svc.ListGroupsForUser(ctx, &iam.ListGroupsForUserInput{
		UserName: aws.String(userName),
	})
	if err != nil {
		return nil, err
	}

	documents := []string{}
	for _, group := range output.Groups {
		docs, err := fetchManagedGroupPolicies(ctx, svc, group.GroupName)
		if err != nil {
			return nil, err
		}
		documents = append(documents, docs...)

		docs, err = fetchInlineGroupPolicies(ctx, svc, group.GroupName)
		if err != nil {
			return nil, err
		}
		documents = append(documents, docs...)
	}

	return documents, nil
}

func fetchManagedGroupPolicies(ctx context.Context, svc *iam.Client, groupName *string) ([]string, error) {
	output, err := svc.ListAttachedGroupPolicies(ctx, &iam.ListAttachedGroupPoliciesInput{
		GroupName: groupName,
	})
	if err != nil {
		return nil, err
	}

	documents := []string{}
	for _, policy := range output.AttachedPolicies {
		doc, err := fetchPolicyDocument(ctx, svc, policy.PolicyArn)
		if err != nil {
			return nil, err
		}
		documents = append(documents, doc)
	}

	return documents, nil
}

func fetchInlineGroupPolicies(ctx context.Context, svc *iam.Client, groupName *string) ([]string, error) {
	output, err := svc.ListGroupPolicies(ctx, &iam.ListGroupPoliciesInput{
		GroupName: groupName,
	})
	if err != nil {
		return nil, err
	}

	documents := []string{}
	for _, policyName := range output.PolicyNames {
		policy, err := svc.GetGroupPolicy(ctx, &iam.GetGroupPolicyInput{
			GroupName:  groupName,
			PolicyName: aws.String(policyName),
		})
		if err != nil {
			return nil, err
		}
		doc, err := url.QueryUnescape(*policy.PolicyDocument)
		if err != nil {
			return nil, err
		}
		documents = append(documents, doc)
	}

	return documents, nil
}

func fetchPolicyDocument(ctx context.Context, svc *iam.Client, policyArn *string) (string, error) {
	policy, err := svc.GetPolicy(ctx, &iam.GetPolicyInput{
		PolicyArn: policyArn,
	})
	if err != nil {
		return "", err
	}

	policyVersion, err := svc.GetPolicyVersion(ctx, &iam.GetPolicyVersionInput{
		PolicyArn: policyArn,
		VersionId: policy.Policy.DefaultVersionId,
	})
	if err != nil {
		return "", err
	}

	document, err := url.QueryUnescape(*policyVersion.PolicyVersion.Document)
	if err != nil {
		return "", err
	}

	return document, nil
}

// Converts policy document JSON string to Policy struct
func decodePolicyDocument(doc string) (Policy, error) {
	data := make(map[string]interface{})
	policy := Policy{}

	err := json.Unmarshal([]byte(doc), &data)
	if err != nil {
		return policy, err
	}

	decoder, err := mapstructure.NewDecoder(
		&mapstructure.DecoderConfig{WeaklyTypedInput: true, Result: &policy},
	)
	if err != nil {
		return policy, err
	}

	err = decoder.Decode(data)
	return policy, err
}
