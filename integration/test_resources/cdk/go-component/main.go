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
	"fmt"
	"os"
	"time"

	"github.com/lacework/go-sdk/v2/api"
	componentCDKClient "github.com/lacework/go-sdk/v2/cli/cdk/client/go"
	"github.com/lacework/go-sdk/v2/lwlogger"
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
		"This is a dummy go-component for testing. Try the arguments 'run', 'fail', 'writecache', or 'readcache'",
	)
}

func app() error {
	// connect to the Lacework CDK server
	cdkClient, err := componentCDKClient.NewCDKClient("0.0.1")
	if err != nil {
		return err
	}
	defer cdkClient.Close()

	// create a client to access Lacework APIs
	lacework, err := api.NewClient(os.Getenv("LW_ACCOUNT"),
		api.WithSubaccount(os.Getenv("LW_SUBACCOUNT")),
		api.WithApiKeys(os.Getenv("LW_API_KEY"), os.Getenv("LW_API_SECRET")),
		api.WithToken(os.Getenv("LW_API_TOKEN")),
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
		err := cdkClient.Metric(
			"test_from_go_component",
			map[string]string{"arg": os.Args[1], "highest_severity": sev},
		).WithDuration(time.Since(now).Milliseconds()).Send()

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

	case "writecache":
		if err := cdkClient.WriteCacheAsset("test_this_cache",
			time.Now().Add(time.Hour*1),
			[]string{"data", "data", "data"},
		); err != nil {
			return err
		}
		return nil

	case "readcache":
		b, err := cdkClient.ReadCacheAsset("test_this_cache")
		if err != nil {
			return err
		}
		fmt.Println(string(b))
		return nil

	default:
		return help()
	}
}
