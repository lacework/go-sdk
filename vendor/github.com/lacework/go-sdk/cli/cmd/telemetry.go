package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/spf13/cobra"

	cdk "github.com/lacework/go-sdk/cli/cdk/go/proto/v1"
)

var (
	// telemetryUploadName is the name of the feature that telemetry is being uploaded for
	telemetryUploadName string

	// telemetryUploadFile is a path to a JSON file containing key value pairs to upload
	telemetryUploadFile string

	// telemetryCmd represents the telemetry command
	telemetryCmd = &cobra.Command{
		Hidden: true,
		Use:    "telemetry",
		Short:  "Manage telemetry sent by the Lacework CLI",
		Args:   cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runConfigureSetup()
		},
	}

	telemetryUploadCmd = &cobra.Command{
		Use:   "upload",
		Short: "Manually send some telemetry back to Lacework",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if telemetryUploadName == "" {
				return errors.New("--name flag is required for this command")
			}
			if telemetryUploadFile == "" {
				return errors.New("--data flag is required for this command")
			}
			return runUpload(telemetryUploadName, telemetryUploadFile)
		},
	}
)

func runUpload(name string, file string) error {
	err := cli.NewClient()
	if err != nil {
		return err
	}
	request, err := prepareUpload(name, file)
	if err != nil {
		return err
	}
	_, err = cli.Honeyvent(context.Background(), request)
	return err
}

func prepareUpload(name string, file string) (request *cdk.HoneyventRequest, err error) {
	jsonFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer func(jsonFile *os.File) {
		err = jsonFile.Close()
	}(jsonFile)
	telemetryBytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	var telemetryData map[string]string
	err = json.Unmarshal(telemetryBytes, &telemetryData)
	if err != nil {
		return nil, err
	}
	request = &cdk.HoneyventRequest{
		Feature:     name,
		FeatureData: telemetryData,
	}
	return request, nil
}

func init() {
	rootCmd.AddCommand(telemetryCmd)
	telemetryCmd.AddCommand(telemetryUploadCmd)

	telemetryUploadCmd.Flags().StringVar(&telemetryUploadName,
		"name", "", "Name of the feature the telemetry upload is for",
	)
	telemetryUploadCmd.Flags().StringVar(&telemetryUploadFile,
		"data", "", "Path to JSON file containing key-value pairs to upload",
	)
}
