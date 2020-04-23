//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2020, Lacework Inc.
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
	"time"

	"github.com/pkg/errors"
)

// EventsService is a service that interacts with the Events endpoints
// from the Lacework Server
type EventsService struct {
	client *Client
}

// List leverages ListRange and returns a list of events from the last 7 days
func (svc *EventsService) List() (EventsResponse, error) {
	var (
		now  = time.Now().UTC()
		from = now.AddDate(0, 0, -7) // 7 days from now
	)

	return svc.ListRange(from, now)
}

// ListRange returns a list of Lacework events during the specified date range.
//
// Requirements and specifications:
// * The dates format should be: yyyy-MM-ddTHH:mm:ssZ (example 2019-07-11T21:11:00Z)
// * The START_TIME and END_TIME must be specified in UTC
// * The difference between the START_TIME and END_TIME must not be greater than 7 days
// * The START_TIME must be less than or equal to three months from current date
// * The number of records produced is limited to 5000
func (svc *EventsService) ListRange(start, end time.Time) (
	response EventsResponse,
	err error,
) {
	if start.After(end) {
		err = errors.New("data range should have a start time before the end time")
		return
	}

	apiPath := fmt.Sprintf(
		"%s?START_TIME=%s&END_TIME=%s",
		apiEventsDateRange,
		start.UTC().Format(time.RFC3339),
		end.UTC().Format(time.RFC3339),
	)
	err = svc.client.RequestDecoder("GET", apiPath, nil, &response)
	return
}

// Details returns details about the specified event_id
func (svc *EventsService) Details(eventID string) (response EventDetailsResponse, err error) {
	if eventID == "" {
		err = errors.New("event_id cannot be empty")
		return
	}

	apiPath := fmt.Sprintf("%s?EVENT_ID=%s", apiEventsDetails, eventID)
	err = svc.client.RequestDecoder("GET", apiPath, nil, &response)
	return
}

type EventDetailsResponse struct {
	Events []EventDetails `json:"data"`
}

type EventDetails struct {
	EventID    string         `json:"event_id"`
	EventActor string         `json:"event_actor"`
	EventModel string         `json:"event_model"`
	EventType  string         `json:"event_type"`
	StartTime  time.Time      `json:"start_time"`
	EndTime    time.Time      `json:"end_time"`
	EntityMap  EventEntityMap `json:"entity_map"`
}

type EventEntityMap struct {
	User            []eventUserEntity            `json:"user,omitempty"`
	Application     []eventApplicationEntity     `json:"application,omitempty"`
	Machine         []eventMachineEntity         `json:"machine,omitempty"`
	Container       []eventContainerEntity       `json:"container,omitempty"`
	DnsName         []eventDnsNameEntity         `json:"DnsName,omitempty"`   // @afiune not in standard
	IpAddress       []eventIpAddressEntity       `json:"IpAddress,omitempty"` // @afiune not in standard
	Process         []eventProcessEntity         `json:"process,omitempty"`
	FileDataHash    []eventFileDataHashEntity    `json:"FileDataHash,omitempty"`    // @afiune not in standard
	FileExePath     []eventFileExePathEntity     `json:"FileExePath,omitempty"`     // @afiune not in standard
	SourceIpAddress []eventSourceIpAddressEntity `json:"SourceIpAddress,omitempty"` // @afiune not in standard
	API             []eventAPIEntity             `json:"api,omitempty"`
	Region          []eventRegionEntity          `json:"region,omitempty"`
	CTUser          []eventCTUserEntity          `json:"ct_user,omitempty"`
	Resource        []eventResourceEntity        `json:"resource,omitempty"`
	RecID           []eventRecIDEntity           `json:"RecId,omitempty"`           // @afiune not in standard
	CustomRule      []eventCustomRuleEntity      `json:"CustomRule,omitempty"`      // @afiune not in standard
	NewViolation    []eventNewViolationEntity    `json:"NewViolation,omitempty"`    // @afiune not in standard
	ViolationReason []eventViolationReasonEntity `json:"ViolationReason,omitempty"` // @afiune not in standard
}
type eventUserEntity struct {
	MachineHostname string `json:"machine_hostname"`
	Username        string `json:"username"`
}

type eventApplicationEntity struct {
	Application       string    `json:"application"`
	HasExternalConns  int32     `json:"has_external_conns"`
	IsClient          int32     `json:"is_client"`
	IsServer          int32     `json:"is_server"`
	EarliestKnownTime time.Time `json:"earliest_known_time"`
}

type eventMachineEntity struct {
	Hostname          string  `json:"hostname"`
	ExternalIp        string  `json:"external_ip"`
	InstanceID        string  `json:"instance_id"`
	InstanceName      string  `json:"instance_name"`
	CpuPercentage     float32 `json:"cpu_percentage"`
	InternalIpAddress string  `json:"internal_ip_address"`
}

type eventContainerEntity struct {
	ImageRepo        string    `json:"image_repo"`
	ImageTag         string    `json:"image_tag"`
	HasExternalConns int32     `json:"has_external_conns"`
	IsClient         int32     `json:"is_client"`
	IsServer         int32     `json:"is_server"`
	FirstSeenTime    time.Time `json:"first_seen_time"`
	PodNamespace     string    `json:"pod_namespace"`
	PodIpAddr        string    `json:"pod_ip_addr"`
}

type eventDnsNameEntity struct {
	Hostname      string  `json:"hostname"`
	PortList      []int32 `json:"port_list"`
	TotalInBytes  int32   `json:"total_in_bytes"`
	TotalOutBytes int32   `json:"total_out_bytes"`
}

type eventIpAddressEntity struct {
	IpAddress     string                 `json:"ip_address"`
	TotalInBytes  int32                  `json:"total_in_bytes"`
	TotalOutBytes int32                  `json:"total_out_bytes"`
	ThreatTags    []string               `json:"threat_tags"`
	ThreatSource  map[string]interface{} `json:"threat_source"`
	Country       string                 `json:"country"`
	Region        string                 `json:"region"`
	PortList      []int32                `json:"port_list"`
	FirstSeenTime time.Time              `json:"first_seen_time"`
}

type eventProcessEntity struct {
	Hostname         string    `json:"hostname"`
	ProcessId        int32     `json:"process_id"`
	ProcessStartTime time.Time `json:"process_start_time"`
	Cmdline          string    `json:"cmdline"`
	CpuPercentage    int32     `json:"cpu_percentage"`
}

type eventFileDataHashEntity struct {
	FiledataHash  string    `json:"filedata_hash"`
	MachineCount  int32     `json:"machine_count"`
	ExePathList   []string  `json:"exe_path_list"`
	FirstSeenTime time.Time `json:"first_seen_time"`
	IsKnownBad    int32     `json:"is_known_bad"`
}

type eventFileExePathEntity struct {
	ExePath          string    `json:"exe_path"`
	FirstSeenTime    time.Time `json:"first_seen_time"`
	LastFiledataHash string    `json:"last_filedata_hash"`
	LastPackageName  string    `json:"last_package_name"`
	LastVersion      string    `json:"last_version"`
	LastFileOwner    string    `json:"last_file_owner"`
}

type eventSourceIpAddressEntity struct {
	IpAddress string `json:"ip_address"`
	Region    string `json:"region"`
	Country   string `json:"country"`
}

type eventAPIEntity struct {
	Service string `json:"service"`
	Api     string `json:"api"`
}

type eventRegionEntity struct {
	Region      string   `json:"region"`
	AccountList []string `json:"account_list"`
}

type eventCTUserEntity struct {
	Username    string   `json:"username"`
	AccoutId    string   `json:"accout_id"`
	Mfa         int32    `json:"mfa"`
	ApiList     []string `json:"api_list"`
	RegionList  []string `json:"region_list"`
	PrincipalId string   `json:"principal_id"`
}

type eventResourceEntity struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type eventRecIDEntity struct {
	RecId        string `json:"rec_id"`
	AccountId    string `json:"account_id"`
	AccountAlias string `json:"account_alias"`
	Title        string `json:"title"`
	Status       string `json:"status"`
	EvalType     string `json:"eval_type"`
	EvalGuid     string `json:"eval_guid"`
}

type eventCustomRuleEntity struct {
	LastUpdatedTime time.Time `json:"last_updated_time"`
	LastUpdatedUser string    `json:"last_updated_user"`
	DisplayFilter   string    `json:"display_filter"`
	RuleGuid        string    `json:"rule_guid"`
}

type eventNewViolationEntity struct {
	RecId    string `json:"rec_id"`
	Reason   string `json:"reason"`
	Resource string `json:"resource"`
}

type eventViolationReasonEntity struct {
	RecId  string `json:"rec_id"`
	Reason string `json:"reason"`
}

type EventsResponse struct {
	Events []Event `json:"data"`
}

type Event struct {
	EventID   string    `json:"event_id"`
	EventType string    `json:"event_type"`
	Severity  string    `json:"severity"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

func (e *Event) SeverityString() string {
	switch e.Severity {
	case "1":
		return "Critical"
	case "2":
		return "High"
	case "3":
		return "Medium"
	case "4":
		return "Low"
	case "5":
		return "Info"
	default:
		return "Unknown"
	}
}

type EventsCount struct {
	Critical int
	High     int
	Medium   int
	Low      int
	Info     int
	Total    int
}

func (er *EventsResponse) GetEventsCount() EventsCount {
	counts := EventsCount{}
	for _, e := range er.Events {
		switch e.Severity {
		case "1":
			counts.Critical++
		case "2":
			counts.High++
		case "3":
			counts.Medium++
		case "4":
			counts.Low++
		case "5":
			counts.Info++
		}
		counts.Total++
	}
	return counts
}
