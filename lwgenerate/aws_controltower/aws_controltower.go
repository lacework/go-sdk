package aws_controltower

import (
	"encoding/json"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/lacework/go-sdk/lwgenerate"
	"github.com/pkg/errors"
)

type GenerateAwsControlTowerTfConfigurationArgs struct {
	// For AWS Subaccounts in consolidated CT setups
	SubAccounts []AwsSubAccount

	// ARN for the S3 bucket for consolidated CloudTrail logging
	S3BucketArn string

	// The SNS topic ARN
	SNSTopicArn string

	// The Aws profile of the log archive account
	LogArchiveProfile string

	// The Aws region of the log archive account
	LogArchiveRegion string

	// The Aws profile of the audit account
	AuditProfile string

	// The Aws region of the audit account
	AuditRegion string

	// The audit account flag input in the format profile:region
	AuditAccount string

	// The log archive account flag input in the format profile:region
	LogArchiveAccount string

	// A name for the cross account policy
	CrossAccountPolicyName string

	// Whether cloudtrail log file integrity validation is enabled
	EnableLogFileValidation bool

	// The length of the external ID to generate. Max length is 1224. Ignored when use_existing_iam_role is set to true
	ExternalIdLength int

	// The IAM role ARN is required when setting use_existing_iam_role to true
	IamRoleArn string

	// The external ID configured inside the IAM role is required when setting use_existing_iam_role to true
	IamRoleExternalID string

	// The IAM role name. Required to match with iam_role_arn if use_existing_iam_role is set to true
	IamRoleName string

	//The Lacework AWS account that the IAM role will grant access
	LaceworkAwsAccountID string

	// The name of the integration in Lacework.
	LaceworkIntegrationName string

	// The prefix that will be used at the beginning of every generated resource
	Prefix string

	// The SQS queue name
	SqsQueueName string

	// A map/dictionary of Tags to be assigned to created resources
	Tags map[string]string

	// Set this to true to use an existing IAM role from the log_archive AWS Account
	UseExistingIamRole bool

	// Amount of time to wait before the next resource is provisioned
	WaitTime int

	// The KMS key arn, if Control Tower was deployed with custom KMS key
	KmsKeyArn string

	// Mapping of AWS accounts to Lacework accounts within a Lacework organization
	OrgAccountMappings OrgAccountMapping

	// OrgAccountMapping json used for flag input
	OrgAccountMappingsJson string

	// Lacework Profile to use
	LaceworkProfile string

	// Lacework Organization
	LaceworkOrganizationLevel bool

	// The Lacework AWS Root Account ID
	LaceworkAccountID string
}

type OrgAccountMapping struct {
	DefaultLaceworkAccount string          `json:"default_lacework_account"`
	Mapping                []OrgAccountMap `json:"mapping"`
}

type OrgAccountMap struct {
	LaceworkAccount string   `json:"lacework_account"`
	AwsAccounts     []string `json:"aws_accounts"`
}

func (args GenerateAwsControlTowerTfConfigurationArgs) GetSubAccounts() []AwsSubAccount {
	return args.SubAccounts
}

func (args *GenerateAwsControlTowerTfConfigurationArgs) GetLaceworkProfile() string {
	return args.LaceworkProfile
}

func (args *GenerateAwsControlTowerTfConfigurationArgs) Generate() (string, error) {
	// Validate inputs
	if err := args.validate(); err != nil {
		return "", errors.Wrap(err, "invalid inputs")
	}

	// Create blocks
	requiredProviders, err := createRequiredProviders()
	if err != nil {
		return "", errors.Wrap(err, "failed to generate required providers")
	}

	awsProvider, err := createAwsProvider(args)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate aws provider")
	}

	laceworkProvider, err := createLaceworkProvider(args)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate lacework provider")
	}

	// controlTower
	controlTowerModule, err := createCloudTrailControlTower(args)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate aws controlTower module")
	}

	// Render
	hclBlocks := lwgenerate.CreateHclStringOutput(
		lwgenerate.CombineHclBlocks(
			requiredProviders,
			awsProvider,
			laceworkProvider,
			controlTowerModule,
		),
	)
	return hclBlocks, nil
}

type AwsSubAccount struct {
	// The name of the AwsProfile to use (in AWS configuration)
	AwsProfile string

	// The AwsRegion this profile should use if any resources are created
	AwsRegion string

	// The Alias of the provider block
	Alias string
}

func NewAwsSubAccount(profile string, region string, alias ...string) AwsSubAccount {
	subaccount := AwsSubAccount{AwsProfile: profile, AwsRegion: region}
	if len(alias) > 0 {
		subaccount.Alias = alias[0]
	}
	return subaccount
}

func (args GenerateAwsControlTowerTfConfigurationArgs) validate() error {
	// Validate s3 bucket arn has been set
	if args.S3BucketArn == "" {
		return errors.New("s3 bucket arn must be set")
	}
	// Validate sns topic arn has been set
	if args.SNSTopicArn == "" {
		return errors.New("sns topic arn must be set")
	}
	// Validate log and audit accounts archive
	if len(args.SubAccounts) == 0 {
		return errors.New("log archive and audit accounts must be set")
	}

	// Validate existing role IAM values, if set
	if args.UseExistingIamRole {
		if args.IamRoleArn == "" ||
			args.IamRoleName == "" ||
			args.IamRoleExternalID == "" {
			return errors.New("when using an existing IAM role, existing role ARN, name, and external ID all must be set")
		}
	}

	return nil
}

type AwsControlTowerTerraformModifier func(c *GenerateAwsControlTowerTfConfigurationArgs)

// NewTerraform returns an instance of the GenerateAwsControlTowerTfConfigurationArgs struct.
//
// Note: Additional configuration details may be set using modifiers of the AwsControlTowerTerraformModifier type
//
// Basic usage: Initialize a new AwsControlTowerTerraformModifier struct, with a non-default AWS profile set. Then
// use generate to create a string output of the required HCL.
//
//	hcl, err := aws_controltower.NewTerraform("us-east-1")
//	  .WithAwsProfile("mycorp-profile")).Generate()
func NewTerraform(s3BucketArn string, snsTopicArn string,
	mods ...AwsControlTowerTerraformModifier) *GenerateAwsControlTowerTfConfigurationArgs {
	config := &GenerateAwsControlTowerTfConfigurationArgs{
		S3BucketArn: s3BucketArn,
		SNSTopicArn: snsTopicArn,
	}
	for _, m := range mods {
		m(config)
	}
	return config
}

func WithCrossAccountPolicyName(name string) AwsControlTowerTerraformModifier {
	return func(c *GenerateAwsControlTowerTfConfigurationArgs) {
		c.CrossAccountPolicyName = name
	}
}

func WithLaceworkAccountID(account string) AwsControlTowerTerraformModifier {
	return func(c *GenerateAwsControlTowerTfConfigurationArgs) {
		c.LaceworkAccountID = account
	}
}

func WithLaceworkProfile(profile string) AwsControlTowerTerraformModifier {
	return func(c *GenerateAwsControlTowerTfConfigurationArgs) {
		c.LaceworkProfile = profile
	}
}

func WithEnableLogFileValidation() AwsControlTowerTerraformModifier {
	return func(c *GenerateAwsControlTowerTfConfigurationArgs) {
		c.EnableLogFileValidation = true
	}
}

func WithExternalIdLength(length int) AwsControlTowerTerraformModifier {
	return func(c *GenerateAwsControlTowerTfConfigurationArgs) {
		c.ExternalIdLength = length
	}
}

func WithWaitTime(waitTime int) AwsControlTowerTerraformModifier {
	return func(c *GenerateAwsControlTowerTfConfigurationArgs) {
		c.WaitTime = waitTime
	}
}

func WithTags(tags map[string]string) AwsControlTowerTerraformModifier {
	return func(c *GenerateAwsControlTowerTfConfigurationArgs) {
		c.Tags = tags
	}
}

func WithKmsKeyArn(arn string) AwsControlTowerTerraformModifier {
	return func(c *GenerateAwsControlTowerTfConfigurationArgs) {
		c.KmsKeyArn = arn
	}
}

func WithExisitingIamRole(arn string, name string, externalID string) AwsControlTowerTerraformModifier {
	return func(c *GenerateAwsControlTowerTfConfigurationArgs) {
		c.IamRoleArn = arn
		c.IamRoleExternalID = externalID
		c.IamRoleName = name
		c.UseExistingIamRole = true
	}
}

func WithPrefix(prefix string) AwsControlTowerTerraformModifier {
	return func(c *GenerateAwsControlTowerTfConfigurationArgs) {
		c.Prefix = prefix
	}
}

func WithSqsQueueName(name string) AwsControlTowerTerraformModifier {
	return func(c *GenerateAwsControlTowerTfConfigurationArgs) {
		c.SqsQueueName = name
	}
}

func WithOrgAccountMappings(mapping OrgAccountMapping) AwsControlTowerTerraformModifier {
	return func(c *GenerateAwsControlTowerTfConfigurationArgs) {
		c.OrgAccountMappings = mapping
	}
}

func WithLaceworkIntegrationName(name string) AwsControlTowerTerraformModifier {
	return func(c *GenerateAwsControlTowerTfConfigurationArgs) {
		c.LaceworkIntegrationName = name
	}
}

func WithLaceworkOrgLevel() AwsControlTowerTerraformModifier {
	return func(c *GenerateAwsControlTowerTfConfigurationArgs) {
		c.LaceworkOrganizationLevel = true
	}
}

func WithSubaccounts(subaccounts ...AwsSubAccount) AwsControlTowerTerraformModifier {
	return func(c *GenerateAwsControlTowerTfConfigurationArgs) {
		c.SubAccounts = subaccounts
	}
}

func createCloudTrailControlTower(args *GenerateAwsControlTowerTfConfigurationArgs) (*hclwrite.Block, error) {
	attributes := map[string]interface{}{}
	modDetails := []lwgenerate.HclModuleModifier{
		lwgenerate.HclModuleWithVersion(lwgenerate.AwsCloudTrailControlTowerVersion)}

	//required args
	if args.S3BucketArn != "" {
		attributes["s3_bucket_arn"] = args.S3BucketArn
	}
	if args.SNSTopicArn != "" {
		attributes["sns_topic_arn"] = args.SNSTopicArn
	}

	// optional args
	if args.LaceworkAccountID != "" {
		attributes["lacework_aws_account_id"] = args.LaceworkAccountID
	}

	if args.CrossAccountPolicyName != "" {
		attributes["cross_account_policy_name"] = args.CrossAccountPolicyName
	}
	if args.EnableLogFileValidation {
		attributes["enable_log_file_validation"] = args.EnableLogFileValidation
	}
	if args.ExternalIdLength != 0 {
		attributes["external_id_length"] = args.ExternalIdLength
	}
	if args.WaitTime != 0 {
		attributes["wait_time"] = args.WaitTime
	}
	if len(args.Tags) != 0 {
		attributes["tags"] = args.Tags
	}
	if args.SqsQueueName != "" {
		attributes["sqs_queue_name"] = args.SqsQueueName
	}
	if args.Prefix != "" {
		attributes["prefix"] = args.Prefix
	}
	if !args.OrgAccountMappings.IsEmpty() {
		orgAccountMappings, err := args.OrgAccountMappings.ToMap()
		if err != nil {
			return nil, errors.Wrap(err, "unable to parse 'org_account_mappings'")
		}
		attributes["org_account_mappings"] = []map[string]any{orgAccountMappings}
	}
	if args.LaceworkIntegrationName != "" {
		attributes["lacework_integration_name"] = args.LaceworkIntegrationName
	}
	if args.KmsKeyArn != "" {
		attributes["kms_key_arn"] = args.KmsKeyArn
	}

	// existing iam role
	if args.UseExistingIamRole {
		attributes["use_existing_iam_role"] = true
		attributes["iam_role_arn"] = args.IamRoleArn
		attributes["iam_role_name"] = args.IamRoleName
		attributes["iam_role_external_id"] = args.IamRoleExternalID
	}

	modDetails = append(modDetails, lwgenerate.HclModuleWithProviderDetails(
		map[string]string{"aws.audit": "aws.audit", "aws.log_archive": "aws.log_archive"}))

	modDetails = append(modDetails,
		lwgenerate.HclModuleWithAttributes(attributes),
	)

	return lwgenerate.NewModule(
		"lacework_aws_controltower",
		lwgenerate.AwsCloudTrailControlTowerSource,
		modDetails...,
	).ToBlock()

}

func createRequiredProviders() (*hclwrite.Block, error) {
	return lwgenerate.CreateRequiredProviders(
		lwgenerate.NewRequiredProvider("lacework",
			lwgenerate.HclRequiredProviderWithSource(lwgenerate.LaceworkProviderSource),
			lwgenerate.HclRequiredProviderWithVersion(lwgenerate.LaceworkProviderVersion)))
}

func createLaceworkProvider(args *GenerateAwsControlTowerTfConfigurationArgs) (*hclwrite.Block, error) {
	lwProviderAttributes := map[string]any{}

	if args.LaceworkProfile != "" {
		lwProviderAttributes["profile"] = args.LaceworkProfile
	}

	if args.LaceworkOrganizationLevel {
		lwProviderAttributes["organization"] = true
	}

	if len(lwProviderAttributes) > 0 {
		return lwgenerate.NewProvider(
			"lacework", lwgenerate.HclProviderWithAttributes(lwProviderAttributes)).ToBlock()
	}

	return nil, nil
}

func createAwsProvider(args *GenerateAwsControlTowerTfConfigurationArgs) ([]*hclwrite.Block, error) {
	var blocks []*hclwrite.Block
	if len(args.SubAccounts) > 0 {
		for _, subaccount := range args.SubAccounts {
			alias := subaccount.AwsProfile
			if subaccount.Alias != "" {
				alias = subaccount.Alias
			}
			attrs := map[string]interface{}{
				"alias":   alias,
				"profile": subaccount.AwsProfile,
				"region":  subaccount.AwsRegion,
			}
			providerBlock, err := lwgenerate.NewProvider(
				"aws",
				lwgenerate.HclProviderWithAttributes(attrs),
			).ToBlock()

			if err != nil {
				return nil, err
			}

			blocks = append(blocks, providerBlock)
		}
	}

	return blocks, nil
}

func (orgMap *OrgAccountMapping) IsEmpty() bool {
	return len(orgMap.Mapping) == 0 && orgMap.DefaultLaceworkAccount == ""
}

func (orgMap *OrgAccountMapping) ToMap() (map[string]any, error) {
	var mappings []map[string]any
	mappingsJsonString, err := json.Marshal(orgMap.Mapping)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(mappingsJsonString, &mappings)
	if err != nil {
		return nil, err
	}

	orgMap.Mapping = []OrgAccountMap{}

	result := make(map[string]any)
	jsonString, err := json.Marshal(orgMap)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonString, &result)
	if err != nil {
		return nil, err
	}

	result["mapping"] = mappings
	return result, nil
}
