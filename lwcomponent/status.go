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

import "github.com/fatih/color"

type Status int

const (
	UnknownStatus Status = iota
	Development
	NotInstalled
	Installed
	UpdateAvailable
	Deprecated
)

func (s Status) Color() *color.Color {
	switch s {
	case Development:
		return color.New(color.FgBlue, color.Bold)
	case NotInstalled:
		return color.New(color.FgWhite, color.Bold)
	case Installed:
		return color.New(color.FgGreen, color.Bold)
	case UpdateAvailable:
		return color.New(color.FgYellow, color.Bold)
	case Deprecated:
		return color.New(color.FgRed, color.Bold)
	default:
		return color.New(color.FgRed, color.Bold)
	}
}

func (s Status) String() string {
	switch s {
	case Development:
		return "Development"
	case NotInstalled:
		return "Not Installed"
	case Installed:
		return "Installed"
	case UpdateAvailable:
		return "Update Available"
	case Deprecated:
		return "Deprecated"
	default:
		return "Unknown"
	}
}

// IsInstalled returns true if the component is installed on disk
//
// TODO: @jon-stewart: remove - is in wrong place
func (c Component) IsInstalled() bool {
	switch c.Status() {
	case Installed, UpdateAvailable:
		return true
	default:
		return false
	}
}
