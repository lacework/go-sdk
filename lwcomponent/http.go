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

package lwcomponent

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/lacework/go-sdk/lwlogger"
)

const (
	DefaultMaxRetry = 3
)

var log = lwlogger.New("INFO")

// Retry 3 times (4 requests total)
// Resty default RetryWaitTime is 100ms
// Exponential backoff to a maximum of RetryWaitTime of 2s
func DownloadFile(path string, url string) (err error) {
	client := resty.New()

	download_timeout := os.Getenv("CDK_DOWNLOAD_TIMEOUT_MINUTES")
	if download_timeout != "" {
		val, err := strconv.Atoi(download_timeout)

		if err == nil {
			client.SetTimeout(time.Duration(val) * time.Minute)
		}
	}

	client.SetRetryCount(DefaultMaxRetry)

	client.OnError(func(req *resty.Request, err error) {
		if v, ok := err.(*resty.ResponseError); ok {
			log.Warn(fmt.Sprintf("Failed to download component: %s: %s", v.Response.Body(), v.Err))
		}

		log.Warn(fmt.Sprintf("Failed to download component: %s", err.Error()))
	})

	_, err = client.R().SetOutput(path).Get(url)

	return
}
