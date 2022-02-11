// Author:: Darren Murray (<dmurray-lacework@lacework.net>)
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

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
)

func TestBuildTeamMembersTable(t *testing.T) {
	teamMemberDetails := buildTeamMemberDetailsTable(mockTeamMember)
	assert.Equal(t, teamMemberDetails, teamMemberOutput)
}

var (
	mockTeamMember = api.TeamMember{
		CustGuid: "TECHALLY_6EACE871E168F286275CBC83010FA6E2E080AA4BAE53489",
		Props: api.TeamMemberProps{
			AccountAdmin: false,
			Company:      "Lacework",
			CreatedTime:  "2022-01-25T12:48:15.394Z",
			FirstName:    "Test",
			JitCreated:   false,
			LastName:     "USER",
			OrgAdmin:     false,
			OrgUser:      false,
			UpdatedBy:    "",
			UpdatedTime:  nil,
		},
		UserEnabled: 0,
		UserGuid:    "",
		UserName:    "",
	}

	teamMemberOutput = `              TEAM MEMBER DETAILS               
------------------------------------------------
    FIRST NAME      Test                        
    LAST NAME       USER                        
    COMPANY         Lacework                    
    ACCOUNT ADMIN   false                       
    ORG ADMIN       false                       
    CREATED AT      2022-01-25T12:48:15.394Z    
    JIT CREATED     false                       
    UPDATED BY                                  
    UPDATED AT                                  
                                                

`
)
