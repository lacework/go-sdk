package cmd

import (
	"github.com/pkg/errors"
)

func showAlertEvents(id int) error {
	cli.StartProgress(" Fetching alert events...")
	eventsResponse, err := cli.LwApi.V2.Alerts.GetEvents(id)
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "unable to show alert")
	}

	return cli.OutputJSON(eventsResponse.Data)
}
