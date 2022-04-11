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

type Type string

const (
	EmptyType Type = ""

	// the component is a binary
	BinaryType = "BINARY"

	// will this component be accessible via the CLI
	CommandType = "CLI_COMMAND"

	// the component is a library, only provides content for the CLI or other components
	LibraryType = "LIBRARY"

	// the component is standalone, should be available in $PATH
	StandaloneType = "STANDALONE"
)

func (c Component) IsExecutable() bool {
	switch c.Type {
	case BinaryType, CommandType:
		return true
	default:
		return false
	}
}

func (c Component) IsCommandType() bool {
	return c.Type == CommandType
}
