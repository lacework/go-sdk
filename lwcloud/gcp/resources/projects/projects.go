//
// Author:: Nicholas Schmeller (<nick.schmeller@lacework.net>)
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

package resources

import (
	"context"
	"fmt"

	"github.com/lacework/go-sdk/lwcloud/gcp/helpers"
	folders "github.com/lacework/go-sdk/lwcloud/gcp/resources/folders"

	resourcemanager "cloud.google.com/go/resourcemanager/apiv3"
	resourcemanagerpb "cloud.google.com/go/resourcemanager/apiv3/resourcemanagerpb"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type ProjectInfo struct {
	Name        string
	Parent      string
	DisplayName string
	ProjectId   string
}

func enumerateTopLevelProjects(
	ctx context.Context, clientOption option.ClientOption, ParentId string,
	Ancestory string, skipList, allowList map[string]bool,
) ([]ProjectInfo, error) {

	var (
		client *resourcemanager.ProjectsClient
		err    error
	)

	if clientOption != nil {
		client, err = resourcemanager.NewProjectsClient(ctx, clientOption)
	} else {
		client, err = resourcemanager.NewProjectsClient(ctx)
	}
	if err != nil {
		return nil, fmt.Errorf("cannot enumerate projects in (%s) due to %s", Ancestory, err.Error())
	}
	defer client.Close()

	projects := make([]ProjectInfo, 0)

	req := &resourcemanagerpb.ListProjectsRequest{
		Parent: ParentId,
	}

	for {
		it := client.ListProjects(ctx, req)

		for {
			resp, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return nil, fmt.Errorf("cannot iterate projects in (%s) due to %s", Ancestory, err.Error())
			}

			if helpers.SkipEntry("projects/"+resp.ProjectId, skipList, allowList) {
				continue
			}

			pi := ProjectInfo{
				Name:        resp.Name,
				DisplayName: resp.DisplayName,
				Parent:      resp.Parent,
				ProjectId:   resp.ProjectId,
			}
			projects = append(projects, pi)
		}

		if req.GetPageToken() == "" {
			break
		}

	}
	return projects, nil
}

func EnumerateProjects(
	ctx context.Context, clientOptions option.ClientOption, ParentId string,
	Ancestory string, skipList, allowList map[string]bool,
) ([]ProjectInfo, error) {

	// find top level projects under Parent first
	projects, err := enumerateTopLevelProjects(ctx, clientOptions, ParentId, Ancestory, skipList, allowList)
	if err != nil {
		return projects, err
	}

	// find all sub folders first
	subFolders, err := folders.EnumerateFolders(ctx, clientOptions, ParentId, Ancestory, skipList, allowList)
	if err != nil {
		return projects, err
	}

	// list all projects in the nested folders
	var retError error
	for _, folder := range subFolders {
		nested_projects, err := enumerateTopLevelProjects(ctx, clientOptions, folder.Name,
			folder.Ancestory+" -> "+folder.DisplayName+" ("+folder.Name+")", skipList, allowList)
		if err != nil {
			// combine errors
			retError = helpers.CombineErrors(retError, err)
			continue
		}

		projects = append(projects, nested_projects...)
	}

	return projects, retError
}
