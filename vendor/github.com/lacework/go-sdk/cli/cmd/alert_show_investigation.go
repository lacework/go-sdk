package cmd

import (
	"github.com/lacework/go-sdk/api"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
)

func alertInvestigationTable(investigations []api.AlertInvestigation) (out [][]string) {
	for _, i := range investigations {
		out = append(out, []string{
			i.Question,
			i.Answer,
		})
	}
	return
}

func renderAlertInvestigationTable(investigations []api.AlertInvestigation) {
	cli.OutputHuman(
		renderCustomTable(
			[]string{"Question", "Answer"},
			alertInvestigationTable(investigations),
			tableFunc(func(t *tablewriter.Table) {
				t.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
				t.SetAutoWrapText(false)
				t.SetBorder(false)
			}),
		),
	)
}

func showAlertInvestigation(id int) error {
	cli.StartProgress(" Fetching alert investigation...")
	investigationResponse, err := cli.LwApi.V2.Alerts.GetInvestigation(id)
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "unable to show alert")
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(investigationResponse.Data)
	}

	if len(investigationResponse.Data) == 0 {
		cli.OutputHuman("There are no investigation details available for the specified alert.\n")
		return nil
	}

	renderAlertInvestigationTable(investigationResponse.Data)
	return nil
}
