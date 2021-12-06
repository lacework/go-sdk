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
		api.WithOrgAccess(),
	)
	if err != nil {
		log.Fatal(err)
	}

	res, err := lacework.V2.TeamMembers.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, tm := range res.Data {
		fmt.Printf("%s:%s\n\n", tm.UserGuid, tm.UserName)
	}

	tmProps := api.TeamMemberProps{
		Company:   "Pokemon International Company",
		FirstName: "Vatasha",
		LastName:  "White",
	}

	tm := api.NewTeamMemberOrg("vatasha.white+ExampleOrgDemo@lacework.net", tmProps)

	response, err := lacework.V2.TeamMembers.CreateOrg(tm)
	if err != nil {
		log.Fatal(err)
	}

	userGuid := response.Data.Accounts[0].UserGuid

	// Update the user
	tm.Props.FirstName = "Vatasha Updated"

	tms, err := lacework.V2.TeamMembers.SearchUsername("vatasha.white+ExampleOrgDemo@lacework.net")
	if err != nil {
		log.Fatal(err)
	}

	for _, tm := range tms.Data {
		fmt.Printf("Team member is: %+v\n", tm)
	}

	response, err = lacework.V2.TeamMembers.UpdateOrg(tm)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("UpdateOrg response: %+v\n", response)

	tm.UserGuid = userGuid

	response, err = lacework.V2.TeamMembers.UpdateOrgById(tm)

	fmt.Printf("UpdateOrgById response: %+v\n", response)

	err = lacework.V2.TeamMembers.DeleteOrg(tms.Data[0].UserGuid)
	if err != nil {
		log.Fatal(err)
	}

}
