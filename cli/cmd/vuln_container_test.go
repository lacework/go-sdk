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
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestUserFriendlyErrorFromOnDemandCtrVulnScanRepositoryNotFound(t *testing.T) {
	err := userFriendlyErrorForOnDemandCtrVulnScan(
		errors.New("Could not successfully send scan request to available integrations for given repo and label"),
		"my-registry.example.com", "image", "tag",
	)
	if assert.NotNil(t, err) {
		assert.Contains(t,
			err.Error(),
			"container image 'image:tag' not found in registry 'my-registry.example.com'")
		assert.Contains(t,
			err.Error(),
			"To view all container registries configured in your account use the command:")
		assert.Contains(t,
			err.Error(),
			"lacework vulnerability container list-registries")
	}
}
