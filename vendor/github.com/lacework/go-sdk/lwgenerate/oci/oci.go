package oci

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/lacework/go-sdk/lwgenerate"
	"github.com/pkg/errors"
)

type GenerateOciTfConfigurationArgs struct {
	// Should we configure CSPM integration in LW?
	Config bool

	// Optional name for config
	ConfigName string

	// Lacework profile to use
	LaceworkProfile string

	// Tenant OCID
	TenantOcid string

	// OCI user email
	OciUserEmail string
}

type OciTerraformModifier func(c *GenerateOciTfConfigurationArgs)

const (
	LaceworkProviderVersion = ">= 1.9.0"
)

// Set the Lacework profile to use for integration
func WithLaceworkProfile(name string) OciTerraformModifier {
	return func(c *GenerateOciTfConfigurationArgs) {
		c.LaceworkProfile = name
	}
}

// Set the name Lacework will use for the name
func WithConfigName(name string) OciTerraformModifier {
	return func(c *GenerateOciTfConfigurationArgs) {
		c.ConfigName = name
	}
}

// Set the OCID of the tenant to be integrated
func WithTenantOcid(ocid string) OciTerraformModifier {
	return func(c *GenerateOciTfConfigurationArgs) {
		c.TenantOcid = ocid
	}
}

// Set the email for the OCI user created for the integration
func WithUserEmail(email string) OciTerraformModifier {
	return func(c *GenerateOciTfConfigurationArgs) {
		c.OciUserEmail = email
	}
}

// NewTerraform returns an instance of the GenerateOciTfConfigurationArgs struct
//
// Note: Additional configuration details may be set using modifiers of the OciTerraformModifier type
//
// Basic usage:
// Initialize a new OciTerraformModifier struct then use generate to
// create a string output of the required HCL.
//
//	hcl, err := aws.NewTerraform(
//		true,
//	  	oci.WithTenancyOcid("ocid1.tenancy...abc"),
//		oci.WithUserEmail("a@b.c"),
//	).Generate()
func NewTerraform(enableConfig bool, mods ...OciTerraformModifier,
) *GenerateOciTfConfigurationArgs {
	config := &GenerateOciTfConfigurationArgs{Config: enableConfig}
	for _, m := range mods {
		m(config)
	}
	return config
}

// Generate new Terraform code based on the supplied args.
func (args *GenerateOciTfConfigurationArgs) Generate() (string, error) {
	// Validate inputs
	if err := args.validate(); err != nil {
		return "", errors.Wrap(err, "invalid inputs")
	}

	// Required providers block
	requiredProviders, err := createRequiredProviders()
	if err != nil {
		return "", errors.Wrap(err, "failed to generate required providers")
	}

	// provider lacework block
	laceworkProvider, err := createLaceworkProvider(args)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate lacework provider")
	}

	configModule, err := createConfig(args)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate oci config module")
	}

	// render HCL
	hclBlocks := lwgenerate.CreateHclStringOutput(
		lwgenerate.CombineHclBlocks(
			requiredProviders,
			laceworkProvider,
			configModule,
		),
	)
	return hclBlocks, nil
}

// Ensure all combinations of inputs our valid for supported spec
func (args *GenerateOciTfConfigurationArgs) validate() error {
	if !args.Config {
		return errors.New("config integration must be enabled to continue")
	}

	if args.TenantOcid == "" {
		return errors.New("tenant OCID must be set")
	}

	if args.OciUserEmail == "" {
		return errors.New("OCI user email must be set")
	}

	return nil
}

func createRequiredProviders() (*hclwrite.Block, error) {
	return lwgenerate.CreateRequiredProviders(
		lwgenerate.NewRequiredProvider("lacework",
			lwgenerate.HclRequiredProviderWithSource(lwgenerate.LaceworkProviderSource),
			lwgenerate.HclRequiredProviderWithVersion(LaceworkProviderVersion),
		),
	)
}

func createLaceworkProvider(args *GenerateOciTfConfigurationArgs) (*hclwrite.Block, error) {
	if args.LaceworkProfile != "" {
		return lwgenerate.NewProvider(
			"lacework",
			lwgenerate.HclProviderWithAttributes(map[string]interface{}{"profile": args.LaceworkProfile}),
		).ToBlock()
	}
	return nil, nil
}

func createConfig(args *GenerateOciTfConfigurationArgs) (*hclwrite.Block, error) {
	if !args.Config {
		return nil, nil
	}

	attributes := map[string]interface{}{}

	// Set the attributes
	attributes["tenancy_id"] = args.TenantOcid
	attributes["user_email"] = args.OciUserEmail
	if args.ConfigName != "" {
		attributes["integration_name"] = args.ConfigName
	}

	// Create and return the module
	modDetails := []lwgenerate.HclModuleModifier{
		lwgenerate.HclModuleWithVersion(lwgenerate.OciConfigVersion),
		lwgenerate.HclModuleWithAttributes(attributes),
	}
	return lwgenerate.NewModule(
		"oci_config",
		lwgenerate.OciConfigSource,
		modDetails...,
	).ToBlock()
}
