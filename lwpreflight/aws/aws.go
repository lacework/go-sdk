package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/lacework/go-sdk/v2/lwpreflight/verbosewriter"
)

type Preflight struct {
	awsConfig               aws.Config
	isOrg                   bool
	integrationTypes        []IntegrationType
	permissions             map[string]bool
	permissionsWithWildcard []string
	tasks                   []func(p *Preflight) error

	caller  Caller
	details Details
	errors  map[IntegrationType][]string

	verboseWriter verbosewriter.WriteCloser
}

type Result struct {
	Caller  Caller
	Details Details
	Errors  map[IntegrationType][]string
}

type Params struct {
	Agentless       bool
	Config          bool
	CloudTrail      bool
	IsOrg           bool // If it's org-level integration
	Region          string
	Profile         string
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string // Optional for temporary credentials
}

func New(params Params) (*Preflight, error) {
	opts := []func(*config.LoadOptions) error{}

	if params.Region != "" {
		opts = append(opts, config.WithRegion(params.Region))
	}
	if params.Profile != "" {
		opts = append(opts, config.WithSharedConfigProfile(params.Profile))
	}
	if params.AccessKeyID != "" && params.SecretAccessKey != "" {
		opts = append(opts, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				params.AccessKeyID,
				params.SecretAccessKey,
				params.SessionToken,
			),
		))
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	integrationTypes := []IntegrationType{}
	tasks := []func(p *Preflight) error{
		FetchCaller,
		FetchPolicies,
		CheckPermissions,
		FetchDetails,
	}

	if params.Agentless {
		integrationTypes = append(integrationTypes, Agentless)
		tasks = append(tasks, CheckVPCQuota)
	}
	if params.Config {
		integrationTypes = append(integrationTypes, Config)
	}
	if params.CloudTrail {
		integrationTypes = append(integrationTypes, CloudTrail)
	}

	preflight := &Preflight{
		awsConfig:               cfg,
		isOrg:                   params.IsOrg,
		integrationTypes:        integrationTypes,
		permissions:             map[string]bool{},
		permissionsWithWildcard: []string{},
		tasks:                   tasks,
		details:                 Details{},
		errors:                  map[IntegrationType][]string{},
		verboseWriter:           verbosewriter.New(),
	}

	return preflight, nil
}

// Overwrite the default verbose writer
func (p *Preflight) SetVerboseWriter(vw verbosewriter.WriteCloser) {
	p.verboseWriter = vw
}

func (p *Preflight) Run() (*Result, error) {
	defer p.verboseWriter.Close()

	for _, task := range p.tasks {
		err := task(p)
		if err != nil {
			p.verboseWriter.Write(fmt.Sprintf("Error running preflight task: %s", err.Error()))
			return nil, err
		}
	}
	result := &Result{
		Caller:  p.caller,
		Details: p.details,
		Errors:  p.errors,
	}
	return result, nil
}
