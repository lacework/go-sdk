package main

import (
	"fmt"
	"os"

	"github.com/lacework/go-sdk/lwconfig"
)

func main() {
	cPath, err := lwconfig.DefaultConfigPath()
	if err != nil {
		fmt.Printf("Unable to detect default config path location: %s\n", err)
		os.Exit(1)
	}

	profiles, err := lwconfig.LoadProfilesFrom(cPath)
	if err != nil {
		fmt.Printf("Error trying to load profiles: %s\n", err)
		os.Exit(1)
	}

	_, ok := profiles["default"]
	if !ok {
		fmt.Println("You have a default profile configured!")
	} else {
		fmt.Println("'default' profile not found")
	}
}
