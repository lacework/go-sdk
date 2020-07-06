//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2020, Lacework Inc.
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

package api_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestVulnerabilitiesScan(t *testing.T) {
	expectedStatus := "Scanning"
	expectedRequestID := "efd151c8-abcd-1234-5678-13e8cca93584"
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI("external/vulnerabilities/container/repository/images/scan",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method, "Scan should be a POST method")

			if assert.NotNil(t, r.Body) {
				body := httpBodySniffer(r)
				assert.Contains(t, body, "gcr.io", "container registry missing")
				assert.Contains(t, body, "example/repo", "wrong repository")
				assert.Contains(t, body, "v0.1.0-dev", "missing tag")
			}

			fmt.Fprintf(w, vulScanJsonResponse(expectedRequestID, expectedStatus))
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Vulnerabilities.Scan(
		"gcr.io",
		"example/repo",
		"v0.1.0-dev",
	)
	assert.Nil(t, err)
	if assert.NotNil(t, response) {
		assert.Equal(t, expectedStatus, response.Data.Status)
		assert.Equal(t, expectedRequestID, response.Data.RequestID)
	}
}

func TestVulnerabilitiesScanLaceworkError(t *testing.T) {
	expectedError := "Container registry not found"
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI("external/vulnerabilities/container/repository/images/scan",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method, "Scan should be a POST method")

			if assert.NotNil(t, r.Body) {
				body := httpBodySniffer(r)
				assert.Contains(t, body, "example.com", "container registry missing")
				assert.Contains(t, body, "example/repo", "wrong repository")
				assert.Contains(t, body, "v0.1.0-dev", "missing tag")
			}

			fmt.Fprintf(w, vulScanErrorJsonResponse(expectedError))
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Vulnerabilities.Scan(
		"example.com",
		"example/repo",
		"v0.1.0-dev",
	)
	assert.Nil(t, err)
	if assert.NotNil(t, response) {
		assert.False(t, response.Ok)
		assert.Equal(t, expectedError, response.Message)
	}
}

func TestVulnerabilitiesScan404Error(t *testing.T) {
	fakeServer := lacework.MockServer()
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Vulnerabilities.Scan(
		"example.com",
		"example/repo",
		"v0.1.0-dev",
	)
	assert.Empty(t, response)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "404 page not found")
	}
}

func TestVulnerabilitiesScanStatus(t *testing.T) {
	expectedStatus := "Scanning"
	expectedRequestID := "efd151c8-abcd-1234-5678-13e8cca93584"
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI("external/vulnerabilities/container/reqId/"+expectedRequestID,
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "ScanStatus should be a GET method")
			fmt.Fprintf(w, vulScanStatusJsonResponse(expectedStatus))
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Vulnerabilities.ScanStatus(expectedRequestID)
	assert.Nil(t, err)
	if assert.NotNil(t, response) {
		assert.Equal(t, expectedStatus, response.CheckStatus())
	}
}

func TestVulnerabilitiesScanStatusError(t *testing.T) {
	expectedRequestID := "efd151c8-abcd-1234-5678-13e8cca93584"
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI("external/vulnerabilities/container/reqId/"+expectedRequestID,
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "ScanStatus should be a GET method")
			fmt.Fprintf(w, vulScanErrorJsonResponse("something happened"))
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Vulnerabilities.ScanStatus(expectedRequestID)
	assert.Nil(t, err)
	if assert.NotNil(t, response) {
		assert.Equal(t, "there is a problem with the vulnerability scan: something happened", response.CheckStatus())
	}
}

func TestVulnerabilitiesReportFromID(t *testing.T) {
	var (
		imageID    = "sha256:01f5882aae5ea55e0dc1b49330b0a83e6be386acd502e6c3ff4b031a227c0dac"
		digestID   = "sha256:167ec3ad6d0368acc9d64b1d857c0208ac641da83209498d1f3da9f38c9ae9ec"
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI("external/vulnerabilities/container/imageId/"+imageID,
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "ReportFromID should be a GET method")
			fmt.Fprintf(w, vulReportJsonResponse(imageID, digestID))
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Vulnerabilities.ReportFromID(imageID)
	assert.Nil(t, err)
	var zero int32 = 0
	var uno int32 = 1
	if assert.NotNil(t, response) {
		assert.Equal(t, "Success", response.CheckStatus())
		assert.Equal(t, digestID, response.Data.Image.ImageInfo.ImageDigest)
		assert.Equal(t, imageID, response.Data.Image.ImageInfo.ImageID)
		assert.Equal(t, "test/repo", response.Data.Image.ImageInfo.Repository)
		assert.Equal(t, []string{"latest"}, response.Data.Image.ImageInfo.Tags)
		assert.Equal(t, uno, response.Data.TotalVulnerabilities)
		assert.Equal(t, zero, response.Data.CriticalVulnerabilities)
		assert.Equal(t, uno, response.Data.HighVulnerabilities)
		assert.Equal(t, zero, response.Data.MediumVulnerabilities)
		assert.Equal(t, zero, response.Data.LowVulnerabilities)
		assert.Equal(t, zero, response.Data.InfoVulnerabilities)
		assert.Equal(t, uno, response.Data.FixableVulnerabilities)

		assert.Equal(t, uno, response.Data.VulFixableCount("High"))
		assert.Equal(t, zero, response.Data.VulFixableCount("Info"))
	}
}

func TestVulnerabilitiesReportFromIDNotFound(t *testing.T) {
	fakeServer := lacework.MockServer()
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Vulnerabilities.ReportFromID("sha256:01f5882aae5ea55e0dc1b49330b0a83e6be386acd502e6c3ff4b031a227c0dac")
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "404 page not found")
	}
	if assert.NotNil(t, response) {
		assert.False(t, response.Ok)
		assert.Empty(t, response.Data)
	}
}

func TestVulnerabilitiesReportFromDigest(t *testing.T) {
	var (
		imageID    = "sha256:01f5882aae5ea55e0dc1b49330b0a83e6be386acd502e6c3ff4b031a227c0dac"
		digestID   = "sha256:167ec3ad6d0368acc9d64b1d857c0208ac641da83209498d1f3da9f38c9ae9ec"
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI("external/vulnerabilities/container/imageDigest/"+digestID,
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "ReportFromDigest should be a GET method")
			fmt.Fprintf(w, vulReportJsonResponse(imageID, digestID))
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Vulnerabilities.ReportFromDigest(digestID)
	assert.Nil(t, err)
	var zero int32 = 0
	var uno int32 = 1
	if assert.NotNil(t, response) {
		assert.Equal(t, "Success", response.CheckStatus())
		assert.Equal(t, digestID, response.Data.Image.ImageInfo.ImageDigest)
		assert.Equal(t, imageID, response.Data.Image.ImageInfo.ImageID)
		assert.Equal(t, "test/repo", response.Data.Image.ImageInfo.Repository)
		assert.Equal(t, []string{"latest"}, response.Data.Image.ImageInfo.Tags)
		assert.Equal(t, uno, response.Data.TotalVulnerabilities)
		assert.Equal(t, zero, response.Data.CriticalVulnerabilities)
		assert.Equal(t, uno, response.Data.HighVulnerabilities)
		assert.Equal(t, zero, response.Data.MediumVulnerabilities)
		assert.Equal(t, zero, response.Data.LowVulnerabilities)
		assert.Equal(t, zero, response.Data.InfoVulnerabilities)
		assert.Equal(t, uno, response.Data.FixableVulnerabilities)

		assert.Equal(t, uno, response.Data.VulFixableCount("High"))
		assert.Equal(t, zero, response.Data.VulFixableCount("Info"))
	}
}

func TestVulnerabilitiesListEvaluations(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		"external/vulnerabilities/container/GetEvaluationsForDateRange",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "ListEvaluations or ListEvaluationsDateRange should be a GET method")

			start, ok := r.URL.Query()["START_TIME"]
			if assert.True(t, ok,
				"START_TIME parameter missing") {

				end, ok := r.URL.Query()["END_TIME"]
				if assert.True(t, ok,
					"END_TIME parameter missing") {

					// verify that start and end times are 7 days apart
					// and that the start time is before the end time
					startTime, err := time.Parse(time.RFC3339, start[0])
					assert.Nil(t, err)
					endTime, err := time.Parse(time.RFC3339, end[0])
					assert.Nil(t, err)

					assert.True(t,
						startTime.Before(endTime),
						"the start time should not be after the end time",
					)
					assert.True(t,
						startTime.AddDate(0, 0, 7).Add(-(time.Minute * time.Duration(2))).Equal(endTime),
						"the data range is not 7 days apart",
					)
					fmt.Fprintf(w, vulContainerEvaluationsResponse(startTime.UnixNano()))
				}
			}
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Vulnerabilities.ListEvaluations()
	assert.Nil(t, err)
	if assert.NotNil(t, response) {
		if assert.Equal(t, 2, len(response.Evaluations)) {
			eval := response.Evaluations[0]
			assert.Equal(t, "PASSED", eval.EvalStatus)
			assert.Equal(t, "EvalBySQL", eval.EvalType)
			assert.Equal(t, "492c2f55cf3073e3978138e599bd2074", eval.EvalGuid)
			assert.Equal(t, "sha256:4393dcffe989b8785a15c696ae5201b8f82744f7f63743ab470e4ca956dad660", eval.ImageDigest)
			assert.Equal(t, "sha256:14a3076d0885a4ab36d52a6834583bc07b6530f7940bff378a67d33c2ee0002b", eval.ImageID)
			assert.Equal(t, "2020-06-25T21:01:18Z", eval.ImageCreatedTime.UTC().Format(time.RFC3339))
			assert.Equal(t, "2020-07-01T16:00:30Z", eval.ImageScanTime.ToTime().UTC().Format(time.RFC3339))
			assert.Equal(t, "2020-07-01T16:00:51Z", eval.StartTime.UTC().Format(time.RFC3339))
			assert.Equal(t, "index.docker.io", eval.ImageRegistry)
			assert.Equal(t, "techallylw/lacework-cli", eval.ImageRepo)
			assert.Equal(t, "Success", eval.ImageScanStatus)
			assert.Equal(t, "", eval.ImageNamespace)
			assert.Equal(t, "", eval.ImageScanErrorMsg)
			assert.Equal(t, "102087642", eval.ImageSize)
			assert.Equal(t, []string{"centos-8"}, eval.ImageTags)
			assert.Equal(t, "0", eval.NdvContainers)
			assert.Equal(t, "2", eval.NumFixes)
			assert.Equal(t, "0", eval.NumVulnerabilitiesSeverity1)
			assert.Equal(t, "2", eval.NumVulnerabilitiesSeverity2)
			assert.Equal(t, "0", eval.NumVulnerabilitiesSeverity3)
			assert.Equal(t, "0", eval.NumVulnerabilitiesSeverity4)
			assert.Equal(t, "0", eval.NumVulnerabilitiesSeverity5)

			eval = response.Evaluations[1]
			assert.Equal(t, "2020-04-17T23:13:43Z", eval.ImageCreatedTime.UTC().Format(time.RFC3339))
			assert.Equal(t, "2020-04-17T23:14:07Z", eval.ImageScanTime.ToTime().UTC().Format(time.RFC3339))
			assert.Equal(t, "2020-07-02T02:19:06Z", eval.StartTime.UTC().Format(time.RFC3339))
			assert.Equal(t, "index.docker.io", eval.ImageRegistry)
			assert.Equal(t, "techallylw/lacework-cli", eval.ImageRepo)
			assert.Equal(t, "Unsupported", eval.ImageScanStatus)
			assert.Equal(t, "Image Distro is not supported.", eval.ImageScanErrorMsg)
			assert.Equal(t, []string{"latest"}, eval.ImageTags)
			assert.Equal(t, "0", eval.NdvContainers)
			assert.Equal(t, "", eval.NumFixes)
			assert.Equal(t, "", eval.NumVulnerabilitiesSeverity1)
			assert.Equal(t, "", eval.NumVulnerabilitiesSeverity2)
			assert.Equal(t, "", eval.NumVulnerabilitiesSeverity3)
			assert.Equal(t, "", eval.NumVulnerabilitiesSeverity4)
			assert.Equal(t, "", eval.NumVulnerabilitiesSeverity5)
		}
	}
}

func TestVulnerabilitiesListEvaluationsDateRangeError(t *testing.T) {
	var (
		now    = time.Now().UTC()
		from   = now.AddDate(0, 0, -7) // 7 days from now
		c, err = api.NewClient("test", api.WithToken("TOKEN"))
	)
	assert.Nil(t, err)

	// a tipical user input error could be that they provide the
	// date range the other way around, from should be the start
	// time, and now should be the end time
	response, err := c.Vulnerabilities.ListEvaluationsDateRange(now, from)
	assert.Empty(t, response)
	if assert.NotNil(t, err) {
		assert.Equal(t,
			"data range should have a start time before the end time",
			err.Error(), "error message mismatch",
		)
	}
}

func vulScanJsonResponse(reqID, status string) string {
	return `
		{
			"data": { "requestId": "` + reqID + `", "Status": "` + status + `" },
			"ok": true,
			"message": "SUCCESS"
		}
	`
}

func vulScanStatusJsonResponse(status string) string {
	return `
		{
			"data": { "scan_status": "` + status + `" },
			"ok": true,
			"message": "SUCCESS"
		}
	`
}

func vulReportJsonResponse(id, digest string) string {
	return `
  {
  "data": {
    "image": {
      "image_info": {
        "image_digest": "` + digest + `",
        "image_id": "` + id + `",
        "registry": "index.docker.io",
        "repository": "test/repo",
        "size": 102087642,
        "tags": ["latest"]
      },
      "image_layers": [
        {
          "hash": "sha256:6910e5a164f725142d77994b247ba20040477fbab49a721bdbe8d61cf855ac23",
          "created_by": "ADD file:84700c11fcc969ac08ef25f115513d76c7b72a4118c01fbc86ef0a6056fdebeb in / ",
          "packages": [
            {
              "name": "bind",
              "namespace": "centos:8",
              "version": "32:9.11.13-3.el8",
              "vulnerabilities": [
                {
                  "name": "CVE-2020-8617",
                  "description": "description of the vulnerability",
                  "severity": "high",
                  "metadata": {
                    "NVD": {
                      "CVSSv3": {
                        "Score": 7.5,
                        "ExploitabilityScore": 3.9,
                        "ImpactScore": 3.6,
                        "Vectors": "CVSS:3.0/AV:N/AC:L/PR:N/UI:N/S:U/C:N/I:N/A:H"
                      },
                      "CVSSv2": {
                        "Score": 5,
                        "PublishedDateTime": "2020-05-19T14:15Z",
                        "Vectors": "AV:N/AC:L/Au:N/C:N/I:N/A:P"
                      }
                    }
                  },
                  "fix_version": "32:9.11.13-5.el8_2"
                }
              ]
            }
          ]
        }
      ]
    },
    "scan_status": "Success",
    "total_vulnerabilities": 1,
    "high_vulnerabilities": 1,
    "fixable_vulnerabilities": 1
  },
  "ok": true,
  "message": "SUCCESS"
}
	`
}

func vulScanErrorJsonResponse(message string) string {
	return `
		{
			"data": {},
			"ok": false,
			"message": "` + message + `"
		}
	`
}

func vulContainerEvaluationsResponse(t int64) string {
	return `
{
  "data": [
      {
      "EVAL_GUID": "492c2f55cf3073e3978138e599bd2074",
      "EVAL_STATUS": "PASSED",
      "EVAL_TYPE": "EvalBySQL",
      "IMAGE_CREATED_TIME": 1593118878546,
      "IMAGE_DIGEST": "sha256:4393dcffe989b8785a15c696ae5201b8f82744f7f63743ab470e4ca956dad660",
      "IMAGE_ID": "sha256:14a3076d0885a4ab36d52a6834583bc07b6530f7940bff378a67d33c2ee0002b",
      "IMAGE_NAMESPACE": null,
      "IMAGE_REGISTRY": "index.docker.io",
      "IMAGE_REPO": "techallylw/lacework-cli",
      "IMAGE_SCAN_ERROR_MSG": "",
      "IMAGE_SCAN_STATUS": "Success",
      "IMAGE_SCAN_TIME": 1593619230613,
      "IMAGE_SIZE": "102087642",
      "IMAGE_TAGS": [
        "centos-8"
      ],
      "NDV_CONTAINERS": "0",
      "NUM_FIXES": "2",
      "NUM_VULNERABILITIES_SEVERITY_1": "0",
      "NUM_VULNERABILITIES_SEVERITY_2": "2",
      "NUM_VULNERABILITIES_SEVERITY_3": "0",
      "NUM_VULNERABILITIES_SEVERITY_4": "0",
      "NUM_VULNERABILITIES_SEVERITY_5": "0",
      "START_TIME": 1593619251414
    },
       {
      "EVAL_GUID": "6c590a95af27068ff7cec5f327044ce9",
      "EVAL_STATUS": "PASSED",
      "EVAL_TYPE": "EvalBySQL",
      "IMAGE_CREATED_TIME": 1587165223000,
      "IMAGE_DIGEST": "sha256:412e1e517c9a6bc2dbcd7bcaaaf5bc09d139ab6b1d0d841c23515a1add1a6eb5",
      "IMAGE_ID": "sha256:b24f6152db4a18a12a3ba60d298bb3ef5004e8e2c573b91327bf2c7ddc0f4b12",
      "IMAGE_NAMESPACE": null,
      "IMAGE_REGISTRY": "index.docker.io",
      "IMAGE_REPO": "techallylw/lacework-cli",
      "IMAGE_SCAN_ERROR_MSG": "Image Distro is not supported.",
      "IMAGE_SCAN_STATUS": "Unsupported",
      "IMAGE_SCAN_TIME": 1587165247448,
      "IMAGE_SIZE": "7341696",
      "IMAGE_TAGS": [
        "latest"
      ],
      "NDV_CONTAINERS": "0",
      "NUM_FIXES": null,
      "NUM_VULNERABILITIES_SEVERITY_1": null,
      "NUM_VULNERABILITIES_SEVERITY_2": null,
      "NUM_VULNERABILITIES_SEVERITY_3": null,
      "NUM_VULNERABILITIES_SEVERITY_4": null,
      "NUM_VULNERABILITIES_SEVERITY_5": null,
      "START_TIME": 1593656346316
    }
  ],
  "message": "SUCCESS",
  "ok": true
}
`
}
