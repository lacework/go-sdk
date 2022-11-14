package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	helpers "github.com/lacework/go-sdk/lwcloud/gcp/helpers"
	projects "github.com/lacework/go-sdk/lwcloud/gcp/resources/projects"

	models "github.com/lacework/go-sdk/lwcloud/gcp/resources/models"

	compute "cloud.google.com/go/compute/apiv1"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
)

func EnumerateInstancesInProject(ctx context.Context, clientOption option.ClientOption, region string, ProjectId string) ([]models.InstanceDetails, error) {

	var (
		client *compute.InstancesClient
		err    error
	)

	if clientOption != nil {
		client, err = compute.NewInstancesRESTClient(ctx, clientOption)
	} else {
		client, err = compute.NewInstancesRESTClient(ctx)
	}

	if err != nil {
		return nil, err
	}
	defer client.Close()

	var filter string
	if region != "" {
		filter = fmt.Sprintf("zone eq .*%s.*", region)
	}

	req := &computepb.AggregatedListInstancesRequest{
		Project: ProjectId,
		Filter:  &filter,
	}

	instances := make([]models.InstanceDetails, 0)

	for {
		it := client.AggregatedList(ctx, req)

		for {
			resp, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return nil, err
			}

			if resp.Value == nil || len(resp.Value.Instances) == 0 {
				continue
			}

			for _, instance := range resp.Value.Instances {

				launchTime, _ := time.Parse(time.RFC3339, instance.GetCreationTimestamp())
				instanceIdStr := fmt.Sprintf("%d", instance.GetId())

				diskIds := make([]string, len(instance.GetDisks()))
				for i, disk := range instance.GetDisks() {
					diskIds[i] = disk.GetSource()
				}

				zone := filepath.Base(instance.GetZone())
				zoneStart := strings.LastIndex(zone, "-")
				region := zone[:zoneStart]

				tags := make(map[string]string)

				privateIp, externalIp, vpcId := getNetworkInfo(instance.GetNetworkInterfaces(), tags)

				instanceInfo := models.InstanceDetails{
					InstanceID: instanceIdStr,
					Type:       instance.GetMachineType(),
					State:      instance.GetStatus(),
					Name:       instance.GetName(),
					Zone:       zone,
					Region:     region,
					ImageID:    instance.GetSourceMachineImage(),
					AccountID:  filepath.Base(ProjectId),
					VpcID:      vpcId,
					PublicIP:   externalIp,
					PrivateIP:  privateIp,
					LaunchTime: launchTime,
					Tags:       tags,
					Props:      nil,
				}
				instances = append(instances, instanceInfo)
			}
		}

		if req.GetPageToken() == "" {
			break
		}
	}

	return instances, nil
}

func EnumerateInstancesInOrg(ctx context.Context, clientOption option.ClientOption, region string, OrgId string, skipList map[string]bool, allowList map[string]bool) (map[string][]models.InstanceDetails, error) {

	projects, err := projects.EnumerateProjects(ctx, clientOption, OrgId, OrgId, skipList, allowList)
	if err != nil {
		return nil, err
	}

	m := make(map[string][]models.InstanceDetails, 0)

	for _, project := range projects {

		if helpers.SkipEntry("projects/"+project.ProjectId, skipList, allowList) {
			continue
		}

		projectInstances, err := EnumerateInstancesInProject(ctx, clientOption, region, project.ProjectId)
		if err != nil {
			// TODO log error and continue
			continue
		}

		m[project.Name] = projectInstances
	}

	return m, nil
}

type NwIntfInfo struct {
	Ipaddr        string                    `json:"ipAddr,omitempty"`
	Ipv6addr      string                    `json:"ipv6Addr,omitempty"`
	Kind          string                    `json:"kind,omitempty"`
	Name          string                    `json:"name,omitempty"`
	NicType       string                    `json:"nicType,omitempty"`
	Network       string                    `json:"network,omitempty"`
	SubNetwork    string                    `json:"subNetwork,omitempty"`
	AccessConfigs []*computepb.AccessConfig `json:"accessConfigs,omitempty"`
}

func getNetworkInfo(nwIntfs []*computepb.NetworkInterface, tags map[string]string) (string, string, string) {
	privateIp := ""
	externalIp := ""
	vpcId := ""

	for _, intf := range nwIntfs {

		nwInfo := NwIntfInfo{
			Ipaddr:        intf.GetNetworkIP(),
			Ipv6addr:      intf.GetIpv6Address(),
			Kind:          intf.GetKind(),
			Name:          intf.GetName(),
			NicType:       intf.GetNicType(),
			Network:       intf.GetNetwork(),
			SubNetwork:    intf.GetSubnetwork(),
			AccessConfigs: intf.GetAccessConfigs(),
		}

		if nwInfo.Ipaddr != "" && privateIp == "" {
			privateIp = nwInfo.Ipaddr
		}

		if nwInfo.Ipv6addr != "" && privateIp == "" {
			privateIp = nwInfo.Ipv6addr
		}

		if nwInfo.Network != "" && vpcId == "" {
			vpcId = nwInfo.Network
		}

		accessConfigs := intf.GetAccessConfigs()
		if len(accessConfigs) != 0 {
			for _, accessConfig := range accessConfigs {
				natIp := accessConfig.GetNatIP()
				externalIpv6Length := accessConfig.GetExternalIpv6()

				if natIp != "" && externalIp == "" {
					externalIp = natIp
				}

				if externalIpv6Length != "" && externalIp == "" {
					externalIp = externalIpv6Length
				}
			}
		}

		nwIntfJson, err := json.Marshal(nwInfo)
		if err == nil {
			tags["InterfaceInfo:"+nwInfo.Name] = string(nwIntfJson)
		}
	}

	return privateIp, externalIp, vpcId
}
