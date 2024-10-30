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

// A library component provides one or more files that other components use
type Library interface {

	// Install downloads the library and deploys the files and index
	Install() error

	// Index returns the index of files that the library contains
	Index() []string

	// GetFile returns the content of one file from the library
	GetFile(string) ([]byte, error)
}

// @hazekamp figure out LibraryComponent (if component is a library how do we interact with it)
