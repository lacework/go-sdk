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
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	DefaultMaxRetry = 3
	defaultTimeout  = 5 * time.Minute
)

// Default timeout is 5 minutes
// Retry 3 times (4 requests total)
// Resty default RetryWaitTime is 100ms
// Exponential backoff to a maximum of RetryWaitTime of 2s
func DownloadFile(filepath string, url string, size int64, timeout time.Duration) (err error) {
	client := resty.New()

	if timeout == 0 {
		timeout = defaultTimeout
	}

	client.SetTimeout(timeout)
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

	var stop bool = false

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
				fmt.Printf("Downloaded: %.0fkb\n", float64(size)/(1<<10))
			} else {
				var percent float64 = float64(size) / float64(totalSize) * 100
				fmt.Printf("Downloaded: %.0f", percent)
				fmt.Println("%")
			}
		}

		time.Sleep(time.Second)
	}
}
