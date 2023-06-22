package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lacework/go-sdk/api"
)

func main() {
	lacework, err := api.NewClient(os.Getenv("LW_ACCOUNT"),
		api.WithSubaccount(os.Getenv("LW_SUBACCOUNT")),
		api.WithApiKeys(os.Getenv("LW_API_KEY"), os.Getenv("LW_API_SECRET")),
		api.WithApiV2(),
	)
	if err != nil {
		log.Fatal(err)
	}

	ociCfgData := api.OciCfgData{
		Credentials: api.OciCfgCredentials{
			Fingerprint: "01:02:03:04:05:06:07:08:09:0A::0B:0C:0D:0E:0F",
			PrivateKey:  "",
		},
		HomeRegion: "us-sanjose-1",
		TenantID:   "ocid1.tenancy.oc1..abcdefghijklmnopqrstuvwxyz",
		TenantName: "oci_tenant",
		UserOCID:   "ocid1.user.oc1..abcdefghijklmnopqrstuvwxyz",
	}

	ociCloudAccount := api.NewCloudAccount(
		"oci-cfg-from-golang",
		api.OciCfgCloudAccount,
		ociCfgData,
	)

	ociCfgResponse, err := lacework.V2.CloudAccounts.Create(ociCloudAccount)
	if err != nil {
		log.Fatal(err)
	}

	// Output: OciCfg Cloud Account created: THE-INTEGRATION-GUID
	fmt.Printf("Cloud Account created: %s", ociCfgResponse.Data.IntgGuid)
}
