//
// Author:: Lokesh Vadlamudi (<lvadlamudi@fortinet.com>)
// Copyright:: Copyright 2025, Fortinet Inc.
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

package api

import (
	"fmt"
)

type AlertObservationTimeline map[string]interface{}

type AlertObservationTimelineResponse struct {
	Data []AlertObservationTimeline `json:"data"`
}

func (svc *AlertsService) GetObservationTimeline(id int) (
	response AlertObservationTimelineResponse,
	err error,
) {
	err = svc.client.RequestDecoder(
		"GET",
		fmt.Sprintf(apiV2AlertsDetails, id, AlertObservationTimelineScope),
		nil,
		&response,
	)
	return
}
