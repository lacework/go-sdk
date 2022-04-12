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

	res, err := lacework.V2.AlertProfiles.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range res.Data {
		fmt.Printf("%s\n", p.Guid)
	}

	var profileRes api.AlertProfileResponse
	err = lacework.V2.AlertProfiles.Get(res.Data[0].Guid, profileRes)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(profileRes)

	profile := api.NewAlertProfile("CUSTOM_PROFILE_EXAMPLE",
		"LW_HE_FILES_DEFAULT_PROFILE",
		[]api.AlertProfileAlert{{Name: "HE_File_Violation",
			EventName:   "LW Host Entity File Violation Alert",
			Description: "{{_OCCURRENCE}} Violation for file {{PATH}} on machine {{MID}}",
			Subject:     "{{_OCCURRENCE}} violation detected for file {{PATH}} on machine {{MID}}"},
		})

	response, err := lacework.V2.AlertProfiles.Create(profile)
	if err != nil {
		log.Fatal(err)
	}

	// Output: Alert Profile created: GUID
	fmt.Printf("Alert Profile created: %s", response.Data.Guid)

	err = lacework.V2.AlertProfiles.Delete(response.Data.Guid)
	if err != nil {
		log.Fatal(err)
	}

	// Output: Alert Profile deleted
	fmt.Println("Alert Profile deleted")

	// create an alert template
	alertTemplate := api.AlertProfileAlert{
		Name:        "Alert Template",
		EventName:   "My Example Alert",
		Description: "This is a test alert template",
		Subject:     "Violation for Testing",
	}

	templateResponse, err := lacework.V2.AlertProfiles.Templates.Create(response.Data.Guid, alertTemplate)
	if err != nil {
		log.Fatal(err)
	}

	// Output: Alert Template created: GUID
	fmt.Printf("Alert Template created: %s", templateResponse.Data.Guid)

	err = lacework.V2.AlertProfiles.Templates.Delete(response.Data.Guid, templateResponse.Data.Guid)
	if err != nil {
		log.Fatal(err)
	}

	// Output: Alert Template deleted
	fmt.Println("Alert Template deleted")
}
