package cmd

import (
	"strconv"

	"github.com/lacework/go-sdk/api"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
)

func alertTimelineTable(timelines []api.AlertTimeline) (out [][]string) {
	for _, t := range timelines {
		out = append(out, []string{
			strconv.Itoa(t.ID),
			t.EntryType,
			t.Message.Value,
			t.EntryAuthorType,
			t.User.Name,
		})
	}
	return
}

func renderAlertTimelineTable(timelines []api.AlertTimeline) {
	cli.OutputHuman(
		renderCustomTable(
			[]string{"Timeline ID", "Entry Type", "Message", "Author Type", "Author"},
			alertTimelineTable(timelines),
			tableFunc(func(t *tablewriter.Table) {
				t.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
				t.SetAutoWrapText(false)
				t.SetBorder(false)
			}),
		),
	)
}

func showAlertTimeline(id int) error {
	cli.StartProgress(" Fetching alert timeline...")
	timelineResponse, err := cli.LwApi.V2.Alerts.GetTimeline(id)
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "unable to show alert")
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(timelineResponse.Data)
	}

	if len(timelineResponse.Data) == 0 {
		cli.OutputHuman("There are no timeline entries for the specified alert.\n")
		return nil
	}

	renderAlertTimelineTable(timelineResponse.Data)
	return nil
}
