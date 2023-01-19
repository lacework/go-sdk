//
// Author:: Nicholas Schmeller (<nick.schmeller@lacework.net>)
// Copyright:: Copyright 2023, Lacework Inc.
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
	"context"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/stretchr/testify/assert"
)

type mockGetRoleAPI func(ctx context.Context, params *iam.GetRoleInput, optFns ...func(*iam.Options)) (*iam.GetRoleOutput, error)

func (m mockGetRoleAPI) GetRole(ctx context.Context, params *iam.GetRoleInput, optFns ...func(*iam.Options)) (*iam.GetRoleOutput, error) {
	return m(ctx, params, optFns...)
}

func TestGetRoleFromName(t *testing.T) {
	testRoleName := "test_role_name"
	cases := []struct {
		client   func(t *testing.T) IAMGetRoleAPI
		roleName string
	}{
		{
			client: func(t *testing.T) IAMGetRoleAPI {
				return mockGetRoleAPI(func(ctx context.Context, params *iam.GetRoleInput, optFns ...func(*iam.Options)) (*iam.GetRoleOutput, error) {
					return &iam.GetRoleOutput{
						Role: &types.Role{
							RoleName: aws.String(testRoleName),
						},
					}, nil

				})
			},
			roleName: testRoleName,
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			role, err := getRoleFromName(tt.client(t), tt.roleName)
			assert.NoError(t, err)
			assert.Equal(t, *role.RoleName, testRoleName)
		})
	}
}
