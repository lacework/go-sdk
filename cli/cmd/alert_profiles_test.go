// Author:: Darren Murray (<dmurray-lacework@lacework.net>)
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

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
)

func TestBuildAlertProfilesTable(t *testing.T) {
	alertProfileDetails := buildAlertProfileDetailsTable(mockAlertProfile)
	assert.Equal(t, alertProfileDetails, alertProfileOutput)
}

func TestFilterAlertProfiles(t *testing.T) {
	filtered := filterAlertProfiles(mockAlertProfileResponse)
	assert.Equal(t, len(filtered), 1)
	assert.Equal(t, filtered[0], "LW_1")
}

var (
	mockAlertProfile = api.AlertProfile{
		Guid:    "LW_HE_FILES_DEFAULT_PROFILE",
		Extends: "LW_LPP_BaseProfile",
		Fields: []api.AlertProfileField{{Name: "FILE_ACCESSED_TIME"},
			{Name: "INODE"},
			{Name: "HARD_LINK_COUNT"},
			{Name: "MID"},
			{Name: "SIZE"},
			{Name: "RECORD_CREATED_TIME"},
			{Name: "BLOCK_COUNT"},
			{Name: "_EVENT_COUNT"},
			{Name: "_POLICY_ID"},
			{Name: "_SEVERITY"},
			{Name: "FILE_NAME"},
			{Name: "_OCCURRENCE"},
			{Name: "LINK_DEST_PATH"},
			{Name: "OWNER_GID"},
			{Name: "_SEVERITY"},
			{Name: "OWNER_USERNAME"},
			{Name: "FILE_TYPE"},
			{Name: "METADATA_HASH"},
			{Name: "OWNER_UID"},
			{Name: "_PRIMARY_TAG"},
			{Name: "PATH"},
			{Name: "_POLICY_DESCRIPTION"},
			{Name: "IS_LINK"},
			{Name: "BLOCK_SIZE"},
			{Name: "_RISK"},
			{Name: "FILE_PERMISSIONS"},
			{Name: "LINK_ABS_DEST_PATH"},
			{Name: "_POLICY_TITLE"},
			{Name: "_POLICY_SEVERITY"},
			{Name: "FILE_CREATED_TIME"},
			{Name: "FILE_MODIFIED_TIME"},
			{Name: "FILEDATA_HASH"},
		},

		DescriptionKeys: []api.AlertProfileDescriptionKeys{
			{Name: "MID", Spec: "{{MID}}"},
			{Name: "_OCCURRENCE", Spec: "{{_OCCURRENCE}}"},
			{Name: "_POLICY_ID", Spec: "{{_POLICY_ID}}"},
			{Name: "OWNER_USERNAME", Spec: "{{OWNER_USERNAME}}"},
			{Name: "_POLICY_TITLE", Spec: "{{_POLICY_TITLE}}"},
			{Name: "FILE_NAME", Spec: "{{FILE_NAME}}"},
			{Name: "FILE_MODIFIED_TIME", Spec: "{{FILE_MODIFIED_TIME}}"},
			{Name: "HARD_LINK_COUNT", Spec: "{{HARD_LINK_COUNT}}"},
			{Name: "LINK_ABS_DEST_PATH", Spec: "{{LINK_ABS_DEST_PATH}}"},
			{Name: "FILE_TYPE", Spec: "{{FILE_TYPE}}"},
			{Name: "IS_LINK", Spec: "{{IS_LINK}}"},
			{Name: "OWNER_GID", Spec: "{{OWNER_GID}}"},
			{Name: "INODE", Spec: "{{INODE}}"},
			{Name: "METADATA_HASH", Spec: "{{METADATA_HASH}}"},
			{Name: "BLOCK_COUNT", Spec: "{{BLOCK_COUNT}}"},
			{Name: "_PRIMARY_TAG", Spec: "{{_PRIMARY_TAG}}"},
			{Name: "FILE_ACCESSED_TIME", Spec: "{{FILE_ACCESSED_TIME}}"},
			{Name: "OWNER_UID", Spec: "{{OWNER_UID}}"},
			{Name: "FILEDATA_HASH", Spec: "{{FILEDATA_HASH}}"},
			{Name: "PATH", Spec: "{{PATH}}"},
			{Name: "SIZE", Spec: "{{SIZE}}"},
			{Name: "FILE_CREATED_TIME", Spec: "{{FILE_CREATED_TIME}}"},
			{Name: "BLOCK_SIZE", Spec: "{{BLOCK_SIZE}}"},
			{Name: "RECORD_CREATED_TIME", Spec: "{{RECORD_CREATED_TIME}}"},
			{Name: "FILE_PERMISSIONS", Spec: "{{FILE_PERMISSIONS}}"},
			{Name: "_POLICY_DESCRIPTION", Spec: "{{_POLICY_DESCRIPTION}}"},
			{Name: "LINK_DEST_PATH", Spec: "{{LINK_DEST_PATH}}"},
		},
		Alerts: []api.AlertTemplate{
			{
				Name:        "HE_File_Violation",
				EventName:   "LW Host Entity File Violation Alert",
				Description: "{{_OCCURRENCE}} Violation for file {{PATH}} on machine {{MID}}",
				Subject:     "{{_OCCURRENCE}} violation detected for file {{PATH}} on machine {{MID}}",
			}, {
				Name:        "HE_File_PolicyChanged",
				EventName:   "LW Host Entity File Violation Alert",
				Description: "{{_OCCURRENCE}} Violation for file {{PATH}} on machine {{MID}}",
				Subject:     "{{_OCCURRENCE}} violation detected for file {{PATH}} on machine {{MID}}",
			},
			{
				Name:        "HE_File_NewViolation",
				EventName:   "LW Host Entity File Violation Alert",
				Description: "{{_OCCURRENCE}} Violation for file {{PATH}} on machine {{MID}}",
				Subject:     "{{_OCCURRENCE}} violation detected for file {{PATH}} on machine {{MID}}",
			},
		},
	}

	alertProfileOutput = `                                                        ALERT TEMPLATES                                                         
--------------------------------------------------------------------------------------------------------------------------------
            NAME                      EVENT NAME                      DESCRIPTION                        SUBJECT                
  ------------------------+--------------------------------+--------------------------------+---------------------------------  
    HE_File_Violation       LW Host Entity File Violation    {{_OCCURRENCE}} Violation        {{_OCCURRENCE}} violation         
                            Alert                            for file {{PATH}} on machine     detected for file {{PATH}} on     
                                                             {{MID}}                          machine {{MID}}                   
    HE_File_PolicyChanged   LW Host Entity File Violation    {{_OCCURRENCE}} Violation        {{_OCCURRENCE}} violation         
                            Alert                            for file {{PATH}} on machine     detected for file {{PATH}} on     
                                                             {{MID}}                          machine {{MID}}                   
    HE_File_NewViolation    LW Host Entity File Violation    {{_OCCURRENCE}} Violation        {{_OCCURRENCE}} violation         
                            Alert                            for file {{PATH}} on machine     detected for file {{PATH}} on     
                                                             {{MID}}                          machine {{MID}}                   
                                                                                                                                
                                                                      FIELDS                                                                       
---------------------------------------------------------------------------------------------------------------------------------------------------
  FILE_ACCESSED_TIME, INODE, HARD_LINK_COUNT, MID, SIZE, RECORD_CREATED_TIME, BLOCK_COUNT, _EVENT_COUNT, _POLICY_ID, _SEVERITY                     
  FILE_NAME, _OCCURRENCE, LINK_DEST_PATH, OWNER_GID, _SEVERITY, OWNER_USERNAME, FILE_TYPE, METADATA_HASH, OWNER_UID, _PRIMARY_TAG                  
  PATH, _POLICY_DESCRIPTION, IS_LINK, BLOCK_SIZE, _RISK, FILE_PERMISSIONS, LINK_ABS_DEST_PATH, _POLICY_TITLE, _POLICY_SEVERITY, FILE_CREATED_TIME  
  FILE_MODIFIED_TIME, FILEDATA_HASH                                                                                                                
                                                                                                                                                   

Fields can be used inside an alert template subject or description by enclosing in double brackets. For example: '{{FIELD_NAME}}'
`
)

var mockAlertProfileResponse = api.AlertProfilesResponse{
	Data: []api.AlertProfile{{
		Guid:            "CUSTOM_1",
		Extends:         "LW_BASE",
		Fields:          nil,
		DescriptionKeys: nil,
		Alerts: []api.AlertTemplate{{
			Name:        "TestAlert1",
			EventName:   "TestEvent",
			Description: "",
			Subject:     "",
		}},
	}, {
		Guid:            "CUSTOM_2",
		Extends:         "LW_BASE",
		Fields:          nil,
		DescriptionKeys: nil,
		Alerts:          nil},
		{
			Guid:            "LW_1",
			Extends:         "LW_BASE",
			Fields:          nil,
			DescriptionKeys: nil,
			Alerts: []api.AlertTemplate{{
				Name:        "TestAlert1",
				EventName:   "TestEvent",
				Description: "",
				Subject:     "",
			}},
		}},
}
