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

package intgguid

// This package generates Globally Unique Identifiers (GUID)
// matching Lacework INTD_GUID that is used for integrations

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	charset               = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomSeed *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func New() string {
	return NewPrefix("MOCK")
}

func NewPrefix(prefix string) string {
	return fmt.Sprintf("%s_%s", prefix, randomString(47))
}

func stringFromCharset(length int, charset string) string {
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = charset[randomSeed.Intn(len(charset))]
	}
	return string(bytes)
}

func randomString(length int) string {
	return stringFromCharset(length, charset)
}
