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
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

// Used to store the list of available filters from a CLI command
//
// E.g. get available filters for a cobra.Command.Long
//
// ```go
//
//	dummyCmdState = struct {
//	    // The available filters
//	    AvailableFilters CmdFilters
//
//	    // List of filters to apply
//	    Filters []string
//		}{}
//
//	dummyCmdState := &cobra.Command{
//	    Long: `The available keys for this command are:
//
// ` + stringSliceToMarkdownList(
//
//	dummyCmdState.AvailableFilters.GetFiltersFrom(
//	    api.MachineDetailEntity{},
//	 ),
//
// )}
// ```
type CmdFilters struct {
	Filters []string
}

func (f *CmdFilters) GetFiltersFrom(T interface{}) []string {
	if len(f.Filters) == 0 {
		f.Filters = getFiltersFrom(T, "")
	}

	return f.Filters
}

func getFiltersFrom(T interface{}, prefix string) []string {
	var (
		filters = []string{}
		rt      = reflect.TypeOf(T)
		rv      = reflect.Indirect(reflect.ValueOf(T))
	)

	for i := 0; i < rt.NumField(); i++ {
		v := rv.Field(i)

		// only use a field if it has a 'json' tag
		if fieldJSON, ok := rt.Field(i).Tag.Lookup("json"); ok {
			// split fieldJSON by comma to handle omitempty/omitzero modifiers
			fieldJSONSlice := strings.Split(fieldJSON, ",")
			if len(fieldJSONSlice) > 0 {
				fieldJSON = fieldJSONSlice[0]
			}

			// if there is any prefix, we need to append it to the JSON field
			if prefix != "" {
				fieldJSON = fmt.Sprintf("%s.%s", prefix, fieldJSON)
			}

			// if the field is a struct we recursively get the fields inside
			if v.Kind() == reflect.Struct {
				filters = append(filters, getFiltersFrom(v.Interface(), fieldJSON)...)
			} else {
				filters = append(filters, fieldJSON)
			}

		}
	}

	return filters
}

// validateKeyValuePairs returns and error if any filter is malformed
func validateKeyValuePairs(filters []string) error {
	for _, pair := range filters {
		kv := strings.Split(pair, ":")
		if len(kv) != 2 || kv[0] == "" || kv[1] == "" {
			return errors.Errorf("malformed filter '%s'. Expected format 'key:value'", pair)
		}
	}
	return nil
}

// stringSliceToMarkdownList display a list of filters in Markdown format.
//
// E.g. The list []string{"a","b","c"} will return
//   - a
//   - b
//   - c
func stringSliceToMarkdownList(filters []string) string {
	if len(filters) == 0 {
		return ""
	}
	return fmt.Sprintf("    * %s", strings.Join(filters, "\n    * "))
}
