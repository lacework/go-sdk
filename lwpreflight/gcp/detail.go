package gcp

import (
	"context"
	"fmt"

	scheduler "google.golang.org/api/cloudscheduler/v1"
)

type Details struct {
	SchedulerRegions []string // Supported regions for Cloud Scheduler. Used for Agentless.
}

func FetchDetails(p *Preflight) error {
	p.details = Details{
		SchedulerRegions: []string{},
	}

	err := fetchSchedulerRegions(p)
	if err != nil {
		return err
	}

	return nil
}

func fetchSchedulerRegions(p *Preflight) error {
	schedulerSvc, err := scheduler.NewService(context.Background(), p.gcpClientOption)
	if err != nil {
		return err
	}

	response, err := schedulerSvc.Projects.Locations.List(
		fmt.Sprintf("projects/%s", p.projectID),
	).Do()
	if err != nil {
		return err
	}

	for _, location := range response.Locations {
		p.details.SchedulerRegions = append(p.details.SchedulerRegions, location.LocationId)
	}

	return nil
}
