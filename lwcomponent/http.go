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
	"io"
	"net/http"
	"os"
	"time"
)

const (
	defaultTimeout  = 30 * time.Second
	defaultMaxRetry = 2
)

func DownloadFile(filepath string, url string, timeout time.Duration) (err error) {
	if timeout == 0 {
		timeout = defaultTimeout
	}

	client := &http.Client{Timeout: timeout}

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = downloadFile(client, url, file)

	for retry := 0; retry < defaultMaxRetry && os.IsTimeout(err); retry++ {
		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			return err
		}

		err = downloadFile(client, url, file)
	}

	return
}

// client.Get returns on receiving HTTP headers and we stream the HTTP data to the output file.
// A timeout will interrupt io.Copy.
func downloadFile(client *http.Client, url string, file *os.File) (err error) {
	var (
		resp *http.Response
	)

	resp, err = client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)

	return
}
