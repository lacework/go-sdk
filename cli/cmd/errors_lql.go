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

package cmd

import (
	"fmt"

	"github.com/lacework/go-sdk/internal/failon"
)

type queryPolicyError struct {
	ExitCode    int
	Message     string
	Err         error
	FailonCount string
	Count       int
}

func NewQueryPolicyError(failonCount string, count int) *queryPolicyError {
	return &queryPolicyError{
		FailonCount: failonCount,
		Count:       count,
		// we use a default exit code that might change
		// during NonCompliant() or Compliant()
		ExitCode: 9,
	}
}

// Example of an error message sent to the end-user:
//
// ERROR (FAIL-ON): query matched fail_on_count expression [count:5] [expr:!=0] (exit code: 9)
func (e *queryPolicyError) Error() string {
	if e.ExitCode == 0 {
		return ""
	}
	return fmt.Sprintf("(FAIL-ON): %s (exit code: %d)", e.Message, e.ExitCode)
}

func (e *queryPolicyError) Unwrap() error {
	return e.Err
}

func (e *queryPolicyError) NonCompliant() bool {
	return !e.validate()
}

func (e *queryPolicyError) Compliant() bool {
	return e.validate()
}

// validate returns true if the error query is compliant, that is,
// when the provided count doesn't match the provided fail on count
// expression. It returns false if the query count matches
func (e *queryPolicyError) validate() bool {
	cli.Log.Debugw("validating policy",
		"count", e.Count,
		"fail_on_count", e.FailonCount,
	)

	co := failon.CountOperation{}
	if err := co.Parse(e.FailonCount); err != nil {
		e.ExitCode = 123
		e.Message = err.Error()
		return false
	}

	isFail, err := co.IsFail(e.Count)
	if err != nil {
		e.ExitCode = 123
		e.Message = err.Error()
		return false
	}

	if isFail {
		e.Message = fmt.Sprintf(
			"query matched fail_on_count expression. [count:%d] [expr:%s]",
			e.Count, e.FailonCount,
		)
		return false
	}

	e.Message = "Compliant policy"
	e.ExitCode = 0
	return true
}
