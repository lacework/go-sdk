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

	res, err := lacework.V2.TeamMembers.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, tm := range res.Data {
		fmt.Printf("%s:%s\n\n", tm.UserGuid, tm.UserName)
	}

	tms, err := lacework.V2.TeamMembers.SearchUsername("vatasha.white+ExampleDemo@lacework.net")
	if err != nil {
		log.Fatal(err)
	}

	for _, tm := range tms.Data {
		fmt.Println("team member is as follows")
		fmt.Printf("%+v", tm)

	}

	tmProps := api.TeamMemberProps{
		Company:   "Pokemon International Company",
		FirstName: "Vatasha",
		LastName:  "White",
	}

	tm := api.NewTeamMember("vatasha.white+exampledemo@lacework.net", tmProps)

	response, err := lacework.V2.TeamMembers.Create(tm)
	if err != nil {
		log.Fatal(err)
	}

	// Update the user
	tm.Props.FirstName = "Vatasha Updated"
	tm.UserGuid = response.Data.UserGuid

	response, err = lacework.V2.TeamMembers.Update(tm)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Update response: %+v", response)
}
