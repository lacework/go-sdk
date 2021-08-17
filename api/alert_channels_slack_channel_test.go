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

package api_test

func singleSlackChannelAlertChannel(id string) string {
	return `
{
  "data": {
    "createdOrUpdatedBy": "salim.afiunemaya@lacework.net",
    "createdOrUpdatedTime": "2021-06-01T18:10:40.745Z",
    "data": {
      "slackUrl": "https://hooks.slack.com/services/abc/foo/secret123"
    },
    "enabled": 1,
    "intgGuid": "` + id + `",
    "isOrg": 0,
    "name": "#tech-ally-notify",
    "state": {
      "details": {},
      "lastSuccessfulTime": 1627895573122,
      "lastUpdatedTime": 1627895573122,
      "ok": true
    },
    "type": "SlackChannel"
  }
}
  `
}
