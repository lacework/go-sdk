package cmd

import (
	"strconv"

	"github.com/lacework/go-sdk/api"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
)

func relatedAlertsTable(relatedAlerts api.RelatedAlerts) (out [][]string) {
	for _, relatedAlert := range relatedAlerts.SortRankDescending() {
		out = append(out, []string{
			relatedAlert.ID,
			relatedAlert.Name,
			relatedAlert.Severity,
			relatedAlert.StartTime,
			relatedAlert.EndTime,
			strconv.Itoa(relatedAlert.Rank),
		})
	}

	return
}

func renderRelatedAlertsTable(relatedAlerts api.RelatedAlerts) {
	cli.OutputHuman(
		renderCustomTable(
			[]string{"Alert ID", "Name", "Severity", "Start Time", "End Time", "Rank"},
			relatedAlertsTable(relatedAlerts),
			tableFunc(func(t *tablewriter.Table) {
				t.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
				t.SetAutoWrapText(false)
				t.SetBorder(false)
			}),
		),
	)
}

func showRelatedAlerts(id int) error {
	cli.StartProgress(" Fetching related alerts...")
	relatedResponse, err := cli.LwApi.V2.Alerts.GetRelatedAlerts(id)
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "unable to show alert")
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(relatedResponse.Data)
	}

	if len(relatedResponse.Data) == 0 {
		cli.OutputHuman("There are no related alerts associated with the specified alert.\n")
		return nil
	}

	renderRelatedAlertsTable(relatedResponse.Data)
	return nil
}
