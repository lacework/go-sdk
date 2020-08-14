package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/lacework/go-sdk/api"
)

var (
	account   string
	apiKey    string
	apiSecret string
)

func main() {
	flag.StringVar(&account, "account", "", "Lacework Account")
	flag.StringVar(&apiKey, "api_key", "", "Lacework API Key")
	flag.StringVar(&apiSecret, "api_secret", "", "Lacework API Secret")
	flag.Parse()

	lacework, err := api.NewClient(account, api.WithApiKeys(apiKey, apiSecret))
	if err != nil {
		log.Fatal(err)
	}

	response, err := lacework.Vulnerabilities.Host.ListCves()
	if err != nil {
		log.Fatal(err)
	}

	// Output:
	fmt.Println(response)
}
