package main

import (
	"fmt"
	"os"

	"github.com/lacework/go-sdk/lwgenerate/gcp"
)

func basic() {
	hcl, err := gcp.NewTerraform(true, true).Generate()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Print(hcl)
}

func existingGcpServiceAccount() {
	hcl, err := gcp.NewTerraform(
		true,
		true,
		gcp.WithProjectId("example_project"),
		gcp.WithGcpServiceAccountCredentials("path/to/service/account/creds.json"),
		gcp.WithExistingServiceAccount(
			gcp.NewExistingServiceAccountDetails("foo", "123456789"),
		),
	).Generate()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Printf("\n-----\n%s", hcl)
}

func existingGcpBucketAndSink() {
	hcl, err := gcp.NewTerraform(
		true,
		true,
		gcp.WithProjectId("example_project"),
		gcp.WithGcpServiceAccountCredentials("path/to/service/account/creds.json"),
		gcp.WithExistingLogBucketName("existing_bucket_name"),
		gcp.WithExistingLogSinkName("existing_sink_name"),
	).Generate()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Printf("\n-----\n%s", hcl)
}

func skipCreateLaceworkIntegration() {
	hcl, err := gcp.NewTerraform(
		true,
		true,
		gcp.WithProjectId("example_project"),
		gcp.WithGcpServiceAccountCredentials("path/to/service/account/creds.json"),
		gcp.WithSkipCreateLaceworkIntegration(true),
	).Generate()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Printf("\n-----\n%s", hcl)
}

func gcpWithLaceworkProfile() {
	hcl, err := gcp.NewTerraform(
		true,
		true,
		gcp.WithProjectId("example_project"),
		gcp.WithGcpServiceAccountCredentials("path/to/service/account/creds.json"),
		gcp.WithLaceworkProfile("test-profile"),
	).Generate()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Printf("\n-----\n%s", hcl)
}

func configOnly() {
	hcl, err := gcp.NewTerraform(
		true,
		false,
		gcp.WithProjectId("example_project"),
		gcp.WithGcpServiceAccountCredentials("path/to/service/account/creds.json"),
	).Generate()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Printf("\n-----\n%s", hcl)
}

func auditLogOnly() {
	hcl, err := gcp.NewTerraform(
		false,
		true,
		gcp.WithProjectId("example_project"),
		gcp.WithGcpServiceAccountCredentials("path/to/service/account/creds.json"),
	).Generate()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Printf("\n-----\n%s", hcl)
}

func organizationGcpIntegration() {
	hcl, err := gcp.NewTerraform(
		true,
		false,
		gcp.WithProjectId("example_project"),
		gcp.WithGcpServiceAccountCredentials("path/to/service/account/creds.json"),
		gcp.WithOrganizationIntegration(true),
		gcp.WithOrganizationId("123456789"),
	).Generate()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Printf("\n-----\n%s", hcl)
}

func main() {
	basic()
	existingGcpServiceAccount()
	existingGcpBucketAndSink()
	gcpWithLaceworkProfile()
	configOnly()
	auditLogOnly()
	organizationGcpIntegration()
}
