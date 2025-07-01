package azure

import (
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/lacework/go-sdk/v2/lwpreflight/verbosewriter"
)

type azureConfig struct {
	cred           azcore.TokenCredential
	subscriptionID string
	tenantID       string
	region         string
}

type Preflight struct {
	azureConfig             azureConfig
	integrationTypes        []IntegrationType
	tasks                   []func(p *Preflight) error
	permissions             map[string]bool
	permissionsWithWildcard []string

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
	Agentless      bool
	Config         bool
	ActivityLog    bool
	SubscriptionID string
	TenantID       string
	ClientID       string
	ClientSecret   string
	Region         string
}

func New(params Params) (*Preflight, error) {
	integrationTypes := []IntegrationType{}
	tasks := []func(p *Preflight) error{
		FetchCaller,
		FetchPolicies,
		CheckPermissions,
		FetchDetails,
	}

	if params.Config {
		integrationTypes = append(integrationTypes, Config)
	}
	if params.ActivityLog {
		integrationTypes = append(integrationTypes, ActivityLog)
	}
	if params.Agentless {
		integrationTypes = append(integrationTypes, Agentless)
		tasks = append(tasks, CheckVNetQuota)
	}

	if params.SubscriptionID == "" {
		return nil, errors.New("SubscriptionID must be provided")
	}

	// Initialize credentials
	var cred azcore.TokenCredential
	var err error
	if params.ClientID != "" && params.ClientSecret != "" {
		cred, err = azidentity.NewClientSecretCredential(
			params.TenantID,
			params.ClientID,
			params.ClientSecret,
			nil,
		)
	} else {
		cred, err = azidentity.NewDefaultAzureCredential(nil)
	}
	if err != nil {
		return nil, err
	}

	cfg := azureConfig{
		cred:           cred,
		subscriptionID: params.SubscriptionID,
		tenantID:       params.TenantID,
		region:         params.Region,
	}

	preflight := &Preflight{
		azureConfig:             cfg,
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
	for _, task := range p.tasks {
		err := task(p)
		if err != nil {
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
