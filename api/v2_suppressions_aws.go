//
// Author:: Ross Moles (<ross.moles@lacework.net>)
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

package api

// AwsSuppressionsV2 is a service that interacts with the V2 Suppressions
// endpoints from the Lacework Server
type AwsSuppressionsV2 struct {
	client *Client
}

func (svc *AwsSuppressionsV2) List() (map[string]SuppressionV2, error) {
	return svc.client.V2.Suppressions.list(AwsSuppression)
}
