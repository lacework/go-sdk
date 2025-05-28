package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type Caller struct {
	AccountID string
	ARN       string
	UserID    string
	Name      string // user name or role name
	IsRoot    bool
	IsAdmin   bool // true if the caller is root user or policies contain the action '*'
}

func (c *Caller) IsAssumedRole() bool {
	return strings.Contains(c.ARN, "assumed-role")
}

func FetchCaller(p *Preflight) error {
	stsSvc := sts.NewFromConfig(p.awsConfig)

	caller, err := stsSvc.GetCallerIdentity(context.Background(), nil)
	if err != nil {
		return err
	}

	resourceName, err := ParseResourceName(*caller.Arn)
	if err != nil {
		return err
	}

	isRoot := resourceName == "root"
	p.caller = Caller{
		AccountID: *caller.Account,
		ARN:       *caller.Arn,
		UserID:    *caller.UserId,
		Name:      resourceName,
		IsRoot:    isRoot,
		IsAdmin:   isRoot,
	}

	return nil
}
