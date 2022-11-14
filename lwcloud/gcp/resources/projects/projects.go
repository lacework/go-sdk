package resources

import (
	"context"
	"fmt"

	"github.com/lacework/go-sdk/lwcloud/gcp/helpers"
	folders "github.com/lacework/go-sdk/lwcloud/gcp/resources/folders"

	resourcemanager "cloud.google.com/go/resourcemanager/apiv3"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	resourcemanagerpb "google.golang.org/genproto/googleapis/cloud/resourcemanager/v3"
)

type ProjectInfo struct {
	Name        string
	Parent      string
	DisplayName string
	ProjectId   string
}

func enumerateTopLevelProjects(ctx context.Context, clientOption option.ClientOption, ParentId string, Ancestory string, skipList, allowList map[string]bool) ([]ProjectInfo, error) {
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
				// log.WithFields(ctx, log.Fields{"Project": resp.ProjectId, "Ancestory": Ancestory}).Error("Skipping Project")
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

func EnumerateProjects(ctx context.Context, clientOptions option.ClientOption, ParentId string, Ancestory string, skipList, allowList map[string]bool) ([]ProjectInfo, error) {

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
		nested_projects, err := enumerateTopLevelProjects(ctx, clientOptions, folder.Name, folder.Ancestory+" -> "+folder.DisplayName+" ("+folder.Name+")", skipList, allowList)
		if err != nil {
			// combine errors
			retError = helpers.CombineErrors(retError, err)
			continue
		}

		projects = append(projects, nested_projects...)
	}

	return projects, retError
}
