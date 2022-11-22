//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2022, Lacework Inc.
// License:: Apache License, Version 2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/lacework/go-sdk/api"
	cdk "github.com/lacework/go-sdk/cli/cdk/go/proto/v1"
	"github.com/lacework/go-sdk/lwlogger"
	"github.com/pkg/errors"
)

var log = lwlogger.New("").Sugar()

func main() {
	if err := app(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR %s\n", err)
		os.Exit(1)
	}
}

func help() error {
	return errors.New(
		"This is a dummy go-component for testing. Try the arguments 'run' and 'fail'",
	)
}

func app() error {
	// connect to the Lacework CDK server
	cdkClient, grpcConn, err := Connect()
	if err != nil {
		return err
	}
	defer grpcConn.Close()

	// create a client to access Lacework APIs
	lacework, err := api.NewClient(os.Getenv("LW_ACCOUNT"),
		api.WithSubaccount(os.Getenv("LW_SUBACCOUNT")),
		api.WithApiKeys(os.Getenv("LW_API_KEY"), os.Getenv("LW_API_SECRET")),
		api.WithToken(os.Getenv("LW_API_TOKEN")),
		api.WithApiV2(),
	)
	if err != nil {
		return errors.Wrap(err, "One or more missing configuration")
	}

	if len(os.Args) <= 1 {
		return help()
	}

	var (
		now        = time.Now().UTC()
		highestSev = "unknown"
	)

	defer func(sev string) {
		_, err := cdkClient.Honeyvent(context.Background(), &cdk.HoneyventRequest{
			DurationMs: time.Since(now).Milliseconds(),
			Feature:    "test_from_go_component",
			FeatureData: map[string]string{
				"arg":              os.Args[1],
				"highest_severity": sev,
			},
		})
		if err != nil {
			log.Error("unable to send honeyvent", "error", err)
		}
	}(highestSev)

	switch os.Args[1] {

	case "run":
		fmt.Println("Running...")
		response, err := lacework.V2.Alerts.List()
		if err != nil {
			return errors.Wrap(err, "unable to access alerts")
		}
		if len(response.Data) != 0 {
			response.Data.SortBySeverity()
			highestSev = response.Data[0].Severity
		}
		fmt.Printf("Highest Severity: %s\n", highestSev)

		return nil

	case "fail":
		return errors.New("Purposely failing...")

	default:
		return help()
	}
}
