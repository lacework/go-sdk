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

// downloadFile is an internal helper that downloads a file to the provided file path
func DownloadFile(filepath string, url string, timeout time.Duration) error {
	var (
		resp     *http.Response
		err      error
		_timeout time.Duration = timeout
	)

	if _timeout == 0 {
		_timeout = defaultTimeout
	}

	client := &http.Client{Timeout: _timeout}

	resp, err = client.Get(url)
	if err != nil {
		for retry := 0; retry < defaultMaxRetry && os.IsTimeout(err); retry++ {
			resp, err = client.Get(url)
		}
		if err != nil {
			return err
		}
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
