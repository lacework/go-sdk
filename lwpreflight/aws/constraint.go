package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/servicequotas"
)

// For AWS Agentless only
func CheckVPCQuota(p *Preflight) error {
	p.verboseWriter.Write(fmt.Sprintf("Discovering VPC quota for region %s", p.awsConfig.Region))

	ctx := context.Background()

	quotaSvc := servicequotas.NewFromConfig(p.awsConfig)
	quotaOutput, err := quotaSvc.GetServiceQuota(ctx, &servicequotas.GetServiceQuotaInput{
		QuotaCode:   aws.String("L-F678F1CE"), // Quota code for VPCs per Region
		ServiceCode: aws.String("vpc"),
	})
	if err != nil {
		return err
	}

	ec2Svc := ec2.NewFromConfig(p.awsConfig)
	vpcsOutput, err := ec2Svc.DescribeVpcs(ctx, nil)
	if err != nil {
		return err
	}

	if len(vpcsOutput.Vpcs) >= int(*quotaOutput.Quota.Value) {
		p.errors[Agentless] = append(
			p.errors[Agentless],
			fmt.Sprintf("VPC Quota limit exceeded in region %s", p.awsConfig.Region),
		)
	}

	return nil
}

func CheckCloudTrailOrgServiceEnabled(p *Preflight) error {
	p.verboseWriter.Write("Discovering enabled services in the organization")

	ctx := context.Background()
	orgSvc := organizations.NewFromConfig(p.awsConfig)

	servicesOutput, err := orgSvc.ListAWSServiceAccessForOrganization(ctx, nil)
	if err != nil {
		return err
	}

	for _, service := range servicesOutput.EnabledServicePrincipals {
		switch *service.ServicePrincipal {
		case "controltower.amazonaws.com":
			p.details.ControlTowerAccess = true
		case "cloudtrail.amazonaws.com":
			p.details.CloudTrailOrgServiceEnabled = true
		}
	}

	if p.details.ControlTowerAccess || p.details.CloudTrailOrgServiceEnabled {
		return nil
	}

	p.verboseWriter.Write("Verifying CloudTrail is enabled as a trusted service in the organization")
	p.errors[CloudTrail] = append(
		p.errors[CloudTrail],
		"CloudTrail is not enabled as a trusted service in the AWS Organization. "+
			"Enable it from the management account: "+
			"aws organizations enable-aws-service-access --service-principal cloudtrail.amazonaws.com",
	)

	return nil
}
