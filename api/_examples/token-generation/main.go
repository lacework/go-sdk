package main

import (
	"fmt"
	"log"

	"github.com/lacework/go-sdk/api"
)

func main() {
	lacework, err := api.NewClient("account")
	if err != nil {
		log.Fatal(err)
	}

	tokenRes, err := lacework.GenerateTokenWithKeys("KEY", "SECRET")
	if err != nil {
		log.Fatal(err)
	}

	// Output: YOUR-ACCESS-TOKEN
	fmt.Println(tokenRes.Token())
}
