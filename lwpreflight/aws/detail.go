package aws

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	cloudtrailTypes "github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/lacework/go-sdk/v2/lwpreflight/logger"
)

type Details struct {
	Regions       []string // Regions that are enabled for the caller account
	ExistingTrail Trail    // Existing eligible trail for CloudTrail integration
	EKSClusters   []EKSCluster

	// Fields for org-level
	OrgAccess           bool
	OrgID               string
	IsManagementAccount bool
	ManagementAccountID string
	OrgAccountIDs       []string
	OrgUnitIDs          []string
	RootOrgUnitID       string
	ControlTowerAccess  bool
}

type Trail struct {
	Name        string
	S3BucketARN string
	SNSTopicARN string
	KMSKeyARN   string
}

type EKSCluster struct {
	Name   string
	Region string
}

func FetchDetails(p *Preflight) error {
	p.details = Details{
		OrgAccountIDs: []string{},
		OrgUnitIDs:    []string{},
		Regions:       []string{},
		EKSClusters:   []EKSCluster{},
	}

	err := fetchOrg(p)
	if err != nil {
		return err
	}

	err = fetchRegions(p)
	if err != nil {
		return err
	}

	err = fetchExistingTrail(p)
	if err != nil {
		logger.Log.Warn(err.Error())
	}

	err = fetchEKSClusters(p)
	if err != nil {
		return err
	}

	return nil
}

func fetchOrg(p *Preflight) error {
	p.verboseWriter.Write("Discovering organization information")

	ctx := context.Background()
	orgSvc := organizations.NewFromConfig(p.awsConfig)

	// Check if the caller can access org
	_, err := orgSvc.DescribeAccount(ctx, &organizations.DescribeAccountInput{
		AccountId: &p.caller.AccountID,
	})
	if err != nil {
		// Only respect errors if org level is enabled. Same for code below
		if p.isOrg {
			return err
		}
		return nil
	}
	p.details.OrgAccess = true

	// Get management account ID and org ID
	orgOutput, err := orgSvc.DescribeOrganization(ctx, nil)
	if err != nil {
		if p.isOrg {
			return err
		}
		return nil
	}
	p.details.ManagementAccountID = *orgOutput.Organization.MasterAccountId
	p.details.IsManagementAccount = *orgOutput.Organization.MasterAccountId == p.caller.AccountID
	p.details.OrgID = *orgOutput.Organization.Id

	if p.isOrg && !p.details.IsManagementAccount {
		return fmt.Errorf("The account %s is not a management account."+
			"Please use a management account to continue with organization level integration.",
			p.caller.AccountID,
		)
	}

	p.verboseWriter.Write("Discovering all accounts in the organization")

	// Get account IDs in the org
	accountsOutput, err := orgSvc.ListAccounts(ctx, nil)
	if err != nil {
		if p.isOrg {
			return err
		}
		return nil
	}
	for _, a := range accountsOutput.Accounts {
		p.details.OrgAccountIDs = append(p.details.OrgAccountIDs, *a.Id)
	}

	p.verboseWriter.Write("Discovering root organization unit")

	// Get root org unit ID and all org unit IDs
	rootsOutput, err := orgSvc.ListRoots(ctx, nil)
	if err != nil {
		if p.isOrg {
			return err
		}
		return nil
	}
	if len(rootsOutput.Roots) > 0 {
		p.verboseWriter.Write("Discovering all organization units")

		p.details.RootOrgUnitID = *rootsOutput.Roots[0].Id
		orgUnitsOutput, err := orgSvc.ListOrganizationalUnitsForParent(
			ctx,
			&organizations.ListOrganizationalUnitsForParentInput{
				ParentId: &p.details.RootOrgUnitID,
			},
		)
		if err != nil {
			if p.isOrg {
				return err
			}
			return nil
		}
		for _, ou := range orgUnitsOutput.OrganizationalUnits {
			p.details.OrgUnitIDs = append(p.details.OrgUnitIDs, *ou.Id)
		}
	}

	p.verboseWriter.Write("Discovering enabled services in the organization")

	// Check enabled services
	servicesOutput, err := orgSvc.ListAWSServiceAccessForOrganization(ctx, nil)
	if err != nil {
		if p.isOrg {
			return err
		}
		return nil
	}
	for _, service := range servicesOutput.EnabledServicePrincipals {
		if *service.ServicePrincipal == "controltower.amazonaws.com" {
			p.details.ControlTowerAccess = true
		}
	}

	return nil
}

func fetchRegions(p *Preflight) error {
	p.verboseWriter.Write("Discovering enabled regions")

	ec2Svc := ec2.NewFromConfig(p.awsConfig)
	output, err := ec2Svc.DescribeRegions(context.Background(), nil)
	if err != nil {
		return err
	}

	for _, r := range output.Regions {
		if *r.OptInStatus == "opt-in-not-required" || *r.OptInStatus == "opted-in" {
			p.details.Regions = append(p.details.Regions, *r.RegionName)
		}
	}

	return nil
}

func fetchExistingTrail(p *Preflight) error {
	var trail *cloudtrailTypes.Trail
	var err error

	if p.details.ControlTowerAccess {
		trail, err = fetchControlTowerTrail(p)
	} else {
		trail, err = fetchEligibleTrail(p)
	}

	if err != nil {
		return err
	}

	p.details.ExistingTrail = Trail{
		Name:        *trail.Name,
		S3BucketARN: fmt.Sprintf("arn:aws:s3:::%s", *trail.S3BucketName),
		SNSTopicARN: *trail.SnsTopicARN,
		KMSKeyARN:   *trail.KmsKeyId,
	}

	return nil
}

/*
To determine if an existing trail is eligible CloudTrail integration:
 1. The trail should be multi-region
 2. Is org-level integration?
    - If yes, the trail should be org trail
    - If no, the trail should NOT be org trail
 3. Is SNS enabled?
    - If yes, select the trail and done
    - If no, select the trail as a candidate(Will use s3 bucket notification instead)
 4. No need to check KMS
*/
func fetchEligibleTrail(p *Preflight) (*cloudtrailTypes.Trail, error) {
	p.verboseWriter.Write("Discovering existing eligible CloudTrail")

	ctx := context.Background()

	trailSvc := cloudtrail.NewFromConfig(p.awsConfig)
	trails, err := trailSvc.ListTrails(ctx, &cloudtrail.ListTrailsInput{})
	if err != nil {
		return nil, err
	}

	var eligibleTrail *cloudtrailTypes.Trail

	for _, trailInfo := range trails.Trails {
		trailSvc = cloudtrail.NewFromConfig(p.awsConfig, func(o *cloudtrail.Options) {
			o.Region = *trailInfo.HomeRegion
		})
		trailOutput, err := trailSvc.GetTrail(ctx, &cloudtrail.GetTrailInput{Name: trailInfo.Name})
		if err != nil {
			return nil, err
		}
		trail := trailOutput.Trail
		if !*trail.IsMultiRegionTrail || p.isOrg != *trail.IsOrganizationTrail {
			continue
		}
		if trail.SnsTopicARN != nil {
			// Found the most eligible trail
			return trail, nil
		}
		if eligibleTrail == nil {
			eligibleTrail = trail
		}
	}

	if eligibleTrail == nil {
		return nil, errors.New("can not find any existing eligible trail")
	}

	return eligibleTrail, nil
}

func fetchControlTowerTrail(p *Preflight) (*cloudtrailTypes.Trail, error) {
	p.verboseWriter.Write("Discovering existing eligible CloudTrail for Control Tower")

	ctx := context.Background()

	trailSvc := cloudtrail.NewFromConfig(p.awsConfig)
	trailsOutput, err := trailSvc.ListTrails(ctx, &cloudtrail.ListTrailsInput{})
	if err != nil {
		return nil, err
	}

	trailRegion := ""
	trailName := "aws-controltower-BaselineCloudTrail"
	for _, trail := range trailsOutput.Trails {
		if *trail.Name == trailName {
			trailRegion = *trail.HomeRegion
			break
		}
	}
	if trailRegion == "" {
		return nil, errors.New("can not find the trail with name \"aws-controltower-BaselineCloudTrail\"")
	}

	trailSvc = cloudtrail.NewFromConfig(p.awsConfig, func(o *cloudtrail.Options) {
		o.Region = trailRegion
	})
	trailOutput, err := trailSvc.GetTrail(ctx, &cloudtrail.GetTrailInput{Name: &trailName})
	if err != nil {
		return nil, err
	}
	trail := trailOutput.Trail

	if trail.S3BucketName == nil || *trail.S3BucketName == "" {
		return nil, errors.New("CloudTrail S3 bucket must be set when using Control Tower")
	}
	if trail.SnsTopicARN == nil || *trail.SnsTopicARN == "" {
		return nil, errors.New("CloudTrail SNS topic must be set when using Control Tower")
	}

	return trail, nil
}

func fetchEKSClusters(p *Preflight) error {
	p.verboseWriter.Write("Discovering EKS clusters")

	var numRegions = len(p.details.Regions)
	var wg sync.WaitGroup
	var ch = make(chan EKSCluster, numRegions)

	wg.Add(numRegions)

	// Collect EKS cluster information for each region
	for _, region := range p.details.Regions {
		cfg := p.awsConfig
		cfg.Region = region
		go func(cfg aws.Config, ch chan<- EKSCluster) {
			eksSvc := eks.NewFromConfig(cfg)
			output, err := eksSvc.ListClusters(context.Background(), nil)
			if err != nil {
				logger.Log.Warnf(
					"Discovering EKS Clusters: unable to check region %s. ERROR: %s",
					region, err.Error(),
				)
			} else {
				for _, name := range output.Clusters {
					ch <- EKSCluster{Name: name, Region: region}
				}
			}
			wg.Done()
		}(cfg, ch)
	}

	// Wait until we discover all clusters from all regions
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Receive EKS cluster information
	for cluster := range ch {
		p.details.EKSClusters = append(p.details.EKSClusters, cluster)
	}

	return nil
}
