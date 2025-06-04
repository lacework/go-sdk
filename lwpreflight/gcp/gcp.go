package gcp

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/oauth2"
	"google.golang.org/api/option"
)

type Preflight struct {
	gcpClientOption  option.ClientOption
	projectID        string
	orgID            string
	integrationTypes []IntegrationType
	permissions      map[string]bool
	tasks            []func(p *Preflight) error

	caller  Caller
	details Details
	errors  map[IntegrationType][]string
}

type Result struct {
	Caller  Caller
	Details Details
	Errors  map[IntegrationType][]string
}

type Params struct {
	Agentless       bool
	AuditLog        bool
	Config          bool
	Region          string
	OrgID           string // Org-level integration if non-empty
	ProjectID       string
	AccessToken     string
	CredentialsFile string // Path to the credential JSON file
	CredentialsJSON string // Content of the credential JSON file
}

func New(params Params) (*Preflight, error) {
	integrationTypes := []IntegrationType{}
	tasks := []func(p *Preflight) error{
		FetchCaller,
		FetchPolicies,
		CheckPermissions,
		FetchDetails,
	}

	if params.Agentless {
		integrationTypes = append(integrationTypes, Agentless)
	}
	if params.AuditLog {
		integrationTypes = append(integrationTypes, AuditLog)
	}
	if params.Config {
		integrationTypes = append(integrationTypes, Config)
	}

	if params.ProjectID == "" {
		return nil, errors.New("ProjectID must be provided")
	}

	preflight := &Preflight{
		projectID:        params.ProjectID,
		orgID:            params.OrgID,
		integrationTypes: integrationTypes,
		permissions:      map[string]bool{},
		tasks:            tasks,
		details:          Details{},
		errors:           map[IntegrationType][]string{},
	}

	if params.AccessToken != "" {
		token := &oauth2.Token{AccessToken: params.AccessToken}
		preflight.gcpClientOption = option.WithTokenSource(oauth2.StaticTokenSource(token))
	} else if params.CredentialsJSON != "" {
		preflight.gcpClientOption = option.WithCredentialsJSON([]byte(params.CredentialsJSON))
	} else {
		var credentialsFile string
		if params.CredentialsFile != "" {
			credentialsFile = params.CredentialsFile
		} else if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") != "" {
			credentialsFile = os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
		}
		if credentialsFile == "" {
			return nil, errors.New(fmt.Sprint(
				"AccessToken, CredentialsFile or CredentialsJSON must be provided. ",
				"Alternatively, set the GOOGLE_APPLICATION_CREDENTIALS environment variable.",
			))
		}
		preflight.gcpClientOption = option.WithCredentialsFile(credentialsFile)
	}

	return preflight, nil
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
