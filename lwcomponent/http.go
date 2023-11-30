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
	"log"
	"os"
	"strconv"
	"time"

	"github.com/abiosoft/colima/util/terminal"
	"github.com/go-resty/resty/v2"
)

const (
	DefaultMaxRetry = 3
)

// Retry 3 times (4 requests total)
// Resty default RetryWaitTime is 100ms
// Exponential backoff to a maximum of RetryWaitTime of 2s
func DownloadFile(filepath string, url string, size int64) (err error) {
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
			fmt.Println(v.Response.Body())
		}
	})

	_, err = os.Create(filepath)
	if err != nil {
		return
	}

	done := make(chan int64)

	go downloadProgress(done, filepath, size)

	response, err := client.R().SetOutput(filepath).Get(url)

	done <- response.Size()

	return
}

func downloadProgress(done chan int64, filepath string, totalSize int64) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	vterm := terminal.NewVerboseWriter(10)
	defer vterm.Close()

	var (
		previous float64 = 0
		stop     bool    = false
	)

	for !stop {
		select {
		case <-done:
			stop = true
		default:
			info, err := file.Stat()
			if err != nil {
				log.Fatal(err)
			}

			size := info.Size()
			if size == 0 {
				size = 1
			}

			if totalSize == 0 {
				mb := float64(size) / (1 << 20)

				if mb > previous {
					vterm.Write([]byte(fmt.Sprintf("Downloaded: %.0fmb\n", mb)))

					previous = mb
				}
			} else {
				percent := float64(size) / float64(totalSize) * 100

				if percent > previous {
					vterm.Write([]byte(fmt.Sprintf("Downloaded: %.0f", percent)))
					vterm.Write([]byte("%"))

					previous = percent
				}
			}
		}

		time.Sleep(time.Second)
	}
}
