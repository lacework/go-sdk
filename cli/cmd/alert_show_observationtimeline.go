package cmd

import (
	"github.com/pkg/errors"
)

func showAlertObservationTimeline(id int) error {
	cli.StartProgress(" Fetching alert observation timeline...")
	observationtimelineResponse, err := cli.LwApi.V2.Alerts.GetObservationTimeline(id)
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "unable to show alert")
	}

	return cli.OutputJSON(observationtimelineResponse.Data)
}
