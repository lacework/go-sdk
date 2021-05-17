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

package cmd

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type countOperationParseTest struct {
	Name     string
	Input    string
	Expected countOperation
	Return   error
}

var countOperationParseTests = []countOperationParseTest{
	countOperationParseTest{
		Name:     "garbage input",
		Input:    "random garbage",
		Expected: countOperation{},
		Return:   errors.New("count operation (random garbage) is invalid"),
	},
	countOperationParseTest{
		Name:     "bad input",
		Input:    ">x",
		Expected: countOperation{},
		Return:   errors.New("count operation (>x) is invalid"),
	},
	countOperationParseTest{
		Name:     "greater",
		Input:    ">10",
		Expected: countOperation{">", 10},
		Return:   nil,
	},
	countOperationParseTest{
		Name:     "greater-equal",
		Input:    ">=10",
		Expected: countOperation{">=", 10},
		Return:   nil,
	},
	countOperationParseTest{
		Name:     "less", // whitespace test
		Input:    " < 890 ",
		Expected: countOperation{"<", 890},
		Return:   nil,
	},
	countOperationParseTest{
		Name:     "less-equal",
		Input:    "<=890",
		Expected: countOperation{"<=", 890},
		Return:   nil,
	},
	countOperationParseTest{
		Name:     "equal",
		Input:    "=1",
		Expected: countOperation{"=", 1},
		Return:   nil,
	},
	countOperationParseTest{
		Name:     "equal-equal",
		Input:    "==1",
		Expected: countOperation{"==", 1},
		Return:   nil,
	},
	countOperationParseTest{
		Name:     "not-equal",
		Input:    "!=7",
		Expected: countOperation{"!=", 7},
		Return:   nil,
	},
}

func TestCountOperationParse(t *testing.T) {
	for _, copt := range countOperationParseTests {
		t.Run(copt.Name, func(t *testing.T) {
			var coActual countOperation
			err := coActual.parse(copt.Input)
			assert.Equal(t, copt.Expected, coActual)
			if copt.Return == nil {
				assert.Nil(t, err)
				return
			}
			assert.Equal(t, copt.Return.Error(), err.Error())
		})
	}
}

type countOperationIsFailTest struct {
	Name        string
	Input       countOperation
	Count       int
	ReturnBool  bool
	ReturnError error
}

var countOperationIsFailTests = []countOperationIsFailTest{
	countOperationIsFailTest{
		Name:        "greater-false",
		Input:       countOperation{">", 0},
		Count:       0,
		ReturnBool:  false,
		ReturnError: nil,
	},
	countOperationIsFailTest{
		Name:        "greater-true",
		Input:       countOperation{">", 0},
		Count:       1,
		ReturnBool:  true,
		ReturnError: nil,
	},
	countOperationIsFailTest{
		Name:        "greater-equal-false",
		Input:       countOperation{">=", 0},
		Count:       -1,
		ReturnBool:  false,
		ReturnError: nil,
	},
	countOperationIsFailTest{
		Name:        "greater-equal-true",
		Input:       countOperation{">=", 0},
		Count:       0,
		ReturnBool:  true,
		ReturnError: nil,
	},
	countOperationIsFailTest{
		Name:        "less-false",
		Input:       countOperation{"<", 1},
		Count:       1,
		ReturnBool:  false,
		ReturnError: nil,
	},
	countOperationIsFailTest{
		Name:        "less-true",
		Input:       countOperation{"<", 1},
		Count:       0,
		ReturnBool:  true,
		ReturnError: nil,
	},
	countOperationIsFailTest{
		Name:        "less-equal-false",
		Input:       countOperation{"<=", 0},
		Count:       1,
		ReturnBool:  false,
		ReturnError: nil,
	},
	countOperationIsFailTest{
		Name:        "less-equal-true",
		Input:       countOperation{"<=", 1},
		Count:       1,
		ReturnBool:  true,
		ReturnError: nil,
	},
	countOperationIsFailTest{
		Name:        "equal-false",
		Input:       countOperation{"=", 1},
		Count:       0,
		ReturnBool:  false,
		ReturnError: nil,
	},
	countOperationIsFailTest{
		Name:        "equal-true",
		Input:       countOperation{"=", 1},
		Count:       1,
		ReturnBool:  true,
		ReturnError: nil,
	},
	countOperationIsFailTest{
		Name:        "not-equal-false",
		Input:       countOperation{"!=", 1},
		Count:       1,
		ReturnBool:  false,
		ReturnError: nil,
	},
	countOperationIsFailTest{
		Name:        "not-equal-true",
		Input:       countOperation{"!=", 1},
		Count:       0,
		ReturnBool:  true,
		ReturnError: nil,
	},
	countOperationIsFailTest{
		Name:        "error",
		Input:       countOperation{},
		Count:       0,
		ReturnBool:  true,
		ReturnError: errors.New("count operation () is invalid"),
	},
}

func TestCountOperationIsFail(t *testing.T) {
	for _, coift := range countOperationIsFailTests {
		t.Run(coift.Name, func(t *testing.T) {
			actualIsFail, err := coift.Input.isFail(coift.Count)
			assert.Equal(t, coift.ReturnBool, actualIsFail)
			if coift.ReturnError == nil {
				assert.Nil(t, err)
				return
			}
			assert.Equal(t, coift.ReturnError.Error(), err.Error())
		})
	}
}
