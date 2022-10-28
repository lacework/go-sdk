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
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

// NOTE these flags will be automatically generated at runtime by cobra
var mockedGlobalFlags = []*pflag.Flag{
	&pflag.Flag{Name: "profile", Shorthand: "p", Value: &mockStringFlagValue{}},
	&pflag.Flag{Name: "debug", Shorthand: "", Value: &mockBoolFlagValue{}},
	&pflag.Flag{Name: "nocolor", Shorthand: "", Value: &mockBoolFlagValue{}},
	&pflag.Flag{Name: "nocache", Shorthand: "", Value: &mockBoolFlagValue{}},
	&pflag.Flag{Name: "noninteractive", Shorthand: "", Value: &mockBoolFlagValue{}},
}

func TestComponentFilterCLIFlagsFromComponentArgs(t *testing.T) {
	cases := []struct {
		Text string

		Args        []string
		GlobalFlags []*pflag.Flag

		expectedArgs  []string
		expectedFlags []string
	}{
		{"empty args and flags",
			[]string(nil), []*pflag.Flag(nil),
			[]string(nil), []string(nil)},

		{"only args without flags returns all args",
			[]string{"iac", "terraform-scan", "--verbose"}, []*pflag.Flag(nil),
			[]string{"iac", "terraform-scan", "--verbose"}, []string(nil)},

		{"only args with global flags returns args",
			[]string{"iac", "terraform-scan", "--verbose"}, mockedGlobalFlags,
			[]string{"iac", "terraform-scan", "--verbose"}, []string(nil)},

		{"args that only have global flags returns all flags",
			[]string{"--profile", "p2", "--debug"}, mockedGlobalFlags,
			[]string(nil), []string{"--profile", "p2", "--debug"}},

		{"args that have both arguments and global flags should split them correctly",
			[]string{"--profile", "p2", "iac", "terraform-scan", "--verbose", "--debug"}, mockedGlobalFlags,
			[]string{"iac", "terraform-scan", "--verbose"}, []string{"--profile", "p2", "--debug"}},

		{"complex args and flags with component commands and flags should split them correctly",
			[]string{"comp-cmd", "--nocolor", "--comp-flag", "comp-subcommand", "--comp-flag2", "-c", "-p", "foo"},
			mockedGlobalFlags,
			[]string{"comp-cmd", "--comp-flag", "comp-subcommand", "--comp-flag2", "-c"},
			[]string{"--nocolor", "-p", "foo"}},
	}

	for _, kase := range cases {
		t.Run(kase.Text, func(t *testing.T) {
			subjectArgs, subjectFlags := filterCLIFlagsFromComponentArgs(kase.Args, kase.GlobalFlags)
			if assert.Equal(t, len(subjectArgs), len(kase.expectedArgs)) {
				assert.Equal(t, kase.expectedArgs, subjectArgs)
			}
			if assert.Equal(t, len(subjectFlags), len(kase.expectedFlags)) {
				assert.Equal(t, kase.expectedFlags, subjectFlags)
			}
		})
	}
}

type mockStringFlagValue struct {
	value string
}

func (m *mockStringFlagValue) String() string {
	return m.value
}
func (m *mockStringFlagValue) Type() string {
	return "string"
}
func (m *mockStringFlagValue) Set(v string) error {
	m.value = v
	return nil
}

type mockBoolFlagValue struct {
	value bool
}

func (m *mockBoolFlagValue) String() string {
	if m.value {
		return "true"
	}
	return "false"
}
func (m *mockBoolFlagValue) Type() string {
	return "bool"
}
func (m *mockBoolFlagValue) Set(v string) error {
	if v == "true" {
		m.value = true
	}
	if v == "false" {
		m.value = false
	}
	return nil
}
