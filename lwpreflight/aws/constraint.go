package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
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
