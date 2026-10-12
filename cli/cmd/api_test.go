//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2021, Lacework Inc.
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

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPICleanupEndpoint(t *testing.T) {

	cases := []struct {
		endpoint         string
		expectedEndpoint string
	}{
		{"", ""},
		{"/foo", "foo"},
		{"/external/bar", "external/bar"},
		{"/v1/bubulubu", "v1/bubulubu"},
		{"/v2/bubulubu", "v2/bubulubu"},

		{"/api/v2/schemas", "v2/schemas"},
		{"/api/v1/external/foo", "v1/external/foo"},
		{"api/v1/endpoint", "v1/endpoint"},
		{"api/v2/coolendpoint", "v2/coolendpoint"},
	}

	for i, kase := range cases {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			assert.Equal(t, kase.expectedEndpoint, cleanupEndpoint(kase.endpoint))
		})
	}
}

func TestAPIReadDataFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "policies.json")
	body := `[{ "policyId" : "lacework-global-118", "enabled" : false }]`
	assert.NoError(t, os.WriteFile(path, []byte(body), 0o600))

	data, err := readApiDataFile(path)
	assert.NoError(t, err)
	assert.Equal(t, body, string(data))

	assert.NoError(t, os.WriteFile(path, append([]byte{0xEF, 0xBB, 0xBF}, body...), 0o600))
	data, err = readApiDataFile(path)
	assert.NoError(t, err)
	assert.Equal(t, body, string(data))

	assert.NoError(t, os.WriteFile(path, []byte{0xFF, 0xFE, '[', 0, ']', 0}, 0o600))
	_, err = readApiDataFile(path)
	assert.ErrorContains(t, err, "UTF-16")

	assert.NoError(t, os.WriteFile(path, nil, 0o600))
	_, err = readApiDataFile(path)
	assert.ErrorContains(t, err, "empty")

	_, err = readApiDataFile(filepath.Join(t.TempDir(), "missing.json"))
	assert.ErrorContains(t, err, "unable to read data file")
}

func TestAPIRunCommandDataValidation(t *testing.T) {
	dir := t.TempDir()
	badFile := filepath.Join(dir, "bad.json")
	assert.NoError(t, os.WriteFile(badFile, []byte("{not json"), 0o600))
	emptyFile := filepath.Join(dir, "empty.json")
	assert.NoError(t, os.WriteFile(emptyFile, nil, 0o600))
	missingFile := filepath.Join(dir, "none.json")

	const inlineHint = "put the JSON in a file and use '--file'"

	// every case fails before the request is sent, so no API client is needed
	cases := []struct {
		name     string
		args     []string
		data     string
		file     string
		expected string
		hint     bool // the PowerShell hint is only for inline --data
	}{
		{"invalid json in --data", []string{"post", "/api/v2/Foo"}, "{not json", "", "not valid JSON", true},
		{"invalid json from --file", []string{"post", "/api/v2/Foo"}, "", badFile, "the JSON in file '" + badFile + "' is not valid", false},
		{"invalid json from @file", []string{"patch", "/api/v2/Foo"}, "@" + badFile, "", "the JSON in file '" + badFile + "' is not valid", false},
		{"missing file after @", []string{"post", "/api/v2/Foo"}, "@", "", "missing file name", false},
		{"missing data file", []string{"post", "/api/v2/Foo"}, "", missingFile, "unable to read data file", false},
		{"patch without data", []string{"patch", "/api/v2/Foo"}, "", "", "'--data' or '--file' parameter for patch", false},
		{"get with data", []string{"get", "/api/v2/Foo"}, "{}", "", "'--data' or '--file' only", false},
		{"get with missing file", []string{"get", "/api/v2/Foo"}, "", missingFile, "'--data' or '--file' only", false},
		{"get with empty @file", []string{"get", "/api/v2/Foo"}, "@" + emptyFile, "", "'--data' or '--file' only", false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			apiData, apiDataFile = c.data, c.file
			t.Cleanup(func() { apiData, apiDataFile = "", "" })

			err := runApiCommand(nil, c.args)
			assert.ErrorContains(t, err, c.expected)
			if c.hint {
				assert.ErrorContains(t, err, inlineHint)
			} else {
				assert.NotContains(t, err.Error(), inlineHint)
			}
		})
	}
}

func TestAPIRunCommandStdinValidation(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{"invalid json", "{not json", "the JSON read from stdin is not valid"},
		{"empty", "", "the data file is empty"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r, w, err := os.Pipe()
			assert.NoError(t, err)
			_, err = w.WriteString(c.input)
			assert.NoError(t, err)
			assert.NoError(t, w.Close())

			stdin := os.Stdin
			os.Stdin = r
			apiDataFile = "-"
			t.Cleanup(func() { os.Stdin, apiDataFile = stdin, ""; r.Close() })

			err = runApiCommand(nil, []string{"post", "/api/v2/Foo"})
			assert.ErrorContains(t, err, c.expected)
			assert.NotContains(t, err.Error(), "put the JSON in a file")
		})
	}
}
