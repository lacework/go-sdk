package gcp

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/lacework/go-sdk/lwgenerate"
	"github.com/pkg/errors"
)

type GenerateGkeTfConfigurationArgs struct {
	ExistingServiceAccount    *ServiceAccount
	ExistingSinkName          string
	IntegrationName           string
	Labels                    map[string]string
	LaceworkProfile           string
	OrganizationId            string
	OrganizationIntegration   bool
	Prefix                    string
	ProjectId                 string
	PubSubSubscriptionLabels  map[string]string
	PubSubTopicLabels         map[string]string
	ServiceAccountCredentials string
	WaitTime                  string
}

type Modifier func(c *GenerateGkeTfConfigurationArgs)

func (args *GenerateGkeTfConfigurationArgs) Generate() (string, error) {
	if err := args.validate(); err != nil {
		return "", errors.Wrap(err, "invalid inputs")
	}

	requiredProviders, err := createRequiredProviders()
	if err != nil {
		return "", errors.Wrap(err, "failed to generate required providers")
	}

	gcpProvider, err := createGcpProvider(args.ServiceAccountCredentials, args.ProjectId)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate gcp provider")
	}

	laceworkProvider, err := createLaceworkProvider(args.LaceworkProfile)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate lacework provider")
	}

	gkeAuditLogModule, err := createGKEAuditLog(args)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate GKE Audit Log module")
	}

	hclBlocks := lwgenerate.CreateHclStringOutput(
		lwgenerate.CombineHclBlocks(
			requiredProviders,
			gcpProvider,
			laceworkProvider,
			gkeAuditLogModule),
	)
	return hclBlocks, nil
}

func (args *GenerateGkeTfConfigurationArgs) validate() error {
	if args.OrganizationIntegration && args.OrganizationId == "" {
		return errors.New("an Organization ID must be provided for an Organization Integration")
	}

	if !args.OrganizationIntegration && args.OrganizationId != "" {
		return errors.New("to provide an Organization ID, Organization Integration must be true")
	}

	if args.ExistingServiceAccount != nil {
		if args.ExistingServiceAccount.Name == "" ||
			args.ExistingServiceAccount.PrivateKey == "" {
			return errors.New(
				"when using an existing Service Account, existing name, and base64 encoded " +
					"JSON Private Key fields all must be set",
			)
		}
	}

	return nil
}

func NewGkeTerraform(mods ...Modifier) *GenerateGkeTfConfigurationArgs {
	config := &GenerateGkeTfConfigurationArgs{}

	for _, m := range mods {
		m(config)
	}

	return config
}

func WithGkeExistingServiceAccount(serviceAccount *ServiceAccount) Modifier {
	return func(c *GenerateGkeTfConfigurationArgs) {
		c.ExistingServiceAccount = serviceAccount
	}
}

func WithGkeExistingSinkName(name string) Modifier {
	return func(c *GenerateGkeTfConfigurationArgs) {
		c.ExistingSinkName = name
	}
}

func WithGkeIntegrationName(name string) Modifier {
	return func(c *GenerateGkeTfConfigurationArgs) {
		c.IntegrationName = name
	}
}

func WithGkeLabels(labels map[string]string) Modifier {
	return func(c *GenerateGkeTfConfigurationArgs) {
		c.Labels = labels
	}
}

func WithGkeLaceworkProfile(name string) Modifier {
	return func(c *GenerateGkeTfConfigurationArgs) {
		c.LaceworkProfile = name
	}
}

func WithGkeOrganizationId(id string) Modifier {
	return func(c *GenerateGkeTfConfigurationArgs) {
		c.OrganizationId = id
	}
}

func WithGkeOrganizationIntegration(enabled bool) Modifier {
	return func(c *GenerateGkeTfConfigurationArgs) {
		c.OrganizationIntegration = enabled
	}
}

func WithGkePrefix(prefix string) Modifier {
	return func(c *GenerateGkeTfConfigurationArgs) {
		c.Prefix = prefix
	}
}

func WithGkeProjectId(id string) Modifier {
	return func(c *GenerateGkeTfConfigurationArgs) {
		c.ProjectId = id
	}
}

func WithGkePubSubSubscriptionLabels(labels map[string]string) Modifier {
	return func(c *GenerateGkeTfConfigurationArgs) {
		c.PubSubSubscriptionLabels = labels
	}
}

func WithGkePubSubTopicLabels(labels map[string]string) Modifier {
	return func(c *GenerateGkeTfConfigurationArgs) {
		c.PubSubTopicLabels = labels
	}
}

func WithGkeServiceAccountCredentials(path string) Modifier {
	return func(c *GenerateGkeTfConfigurationArgs) {
		c.ServiceAccountCredentials = path
	}
}

func WithGkeWaitTime(waitTime string) Modifier {
	return func(c *GenerateGkeTfConfigurationArgs) {
		c.WaitTime = waitTime
	}
}

func createGKEAuditLog(args *GenerateGkeTfConfigurationArgs) (*hclwrite.Block, error) {
	var level string
	attributes := map[string]interface{}{}

	if args.OrganizationIntegration {
		level = "organization"
		attributes["integration_type"] = "ORGANIZATION"
		attributes["organization_id"] = args.OrganizationId

	} else {
		level = "project"
		attributes["integration_type"] = "PROJECT"
	}

	if args.ExistingSinkName != "" {
		attributes["existing_sink_name"] = args.ExistingSinkName
	}

	if args.ExistingServiceAccount != nil {
		attributes["use_existing_service_account"] = true
		attributes["service_account_name"] = args.ExistingServiceAccount.Name
		attributes["service_account_private_key"] = args.ExistingServiceAccount.PrivateKey
	}

	if args.IntegrationName != "" {
		attributes["lacework_integration_name"] = args.IntegrationName
	}

	if args.Prefix != "" {
		attributes["prefix"] = args.Prefix
	}

	if args.WaitTime != "" {
		attributes["wait_time"] = args.WaitTime
	}

	return lwgenerate.NewModule(
		fmt.Sprintf("gcp_%s_level_gke_audit_log", level),
		lwgenerate.GcpGKEAuditLogSource,
		lwgenerate.HclModuleWithAttributes(attributes),
		lwgenerate.HclModuleWithVersion(lwgenerate.GcpGKEAuditLogVersion),
	).ToBlock()
}
