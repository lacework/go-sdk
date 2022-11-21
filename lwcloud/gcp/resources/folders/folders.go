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

	resourcemanager "cloud.google.com/go/resourcemanager/apiv3"
	"github.com/lacework/go-sdk/lwcloud/gcp/helpers"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	resourcemanagerpb "google.golang.org/genproto/googleapis/cloud/resourcemanager/v3"
)

type FolderInfo struct {
	Name        string
	Parent      string
	DisplayName string
	Ancestory   string
}

func EnumerateFolders(ctx context.Context, clientOptions option.ClientOption, ParentId string, Ancestory string, skipList, allowList map[string]bool) ([]FolderInfo, error) {

	var (
		client *resourcemanager.FoldersClient
		err    error
	)

	if clientOptions != nil {
		client, err = resourcemanager.NewFoldersClient(ctx, clientOptions)
	} else {
		client, err = resourcemanager.NewFoldersClient(ctx)
	}

	if err != nil {
		return nil, fmt.Errorf("cannot enumerate folders in (%s) due to %s", Ancestory, err.Error())
	}
	defer client.Close()

	req := &resourcemanagerpb.ListFoldersRequest{
		Parent: ParentId,
	}

	folders := make([]FolderInfo, 0)

	for {

		it := client.ListFolders(ctx, req)

		for {
			resp, err := it.Next()
			if err == iterator.Done {
				break
			}

			if err != nil {
				return nil, fmt.Errorf("cannot iterate folders in ancestory (%s) due to %s", Ancestory, err.Error())
			}

			if helpers.SkipEntry(resp.Name, skipList, allowList) {
				// log.WithFields(ctx, log.Fields{"Folder": resp.Name, "Ancestory": Ancestory}).Error("Skipping Folder")
				continue
			}

			fi := FolderInfo{
				Name:        resp.Name,
				DisplayName: resp.DisplayName,
				Parent:      resp.Parent,
				Ancestory:   Ancestory,
			}
			folders = append(folders, fi)

			// search for folders recursively; ignore errors
			subFolders, _ := EnumerateFolders(ctx, clientOptions, resp.Name, Ancestory+" -> "+resp.DisplayName+" ("+resp.Name+")", skipList, allowList)
			if len(subFolders) != 0 {
				folders = append(folders, subFolders...)
			}
		}

		if req.GetPageToken() == "" {
			break
		}
	}

	return folders, nil
}
