package cmd

import (
	"github.com/lacework/go-sdk/api"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
)

func alertIntegrationsTable(integrations []api.AlertIntegration) (out [][]string) {
	for _, i := range integrations {
		out = append(out, []string{
			i.ID,
			i.IntgGUID,
			i.Type,
			i.Channel.Status(),
			i.Channel.StateString(),
		})
	}
	return
}

func renderAlertIntegrationsTable(integrations []api.AlertIntegration) {
	cli.OutputHuman(
		renderCustomTable(
			[]string{"Alert Integration ID", "Alert Channel GUID", "Type", "Status", "State"},
			alertIntegrationsTable(integrations),
			tableFunc(func(t *tablewriter.Table) {
				t.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
				t.SetAutoWrapText(false)
				t.SetBorder(false)
			}),
		),
	)
}

func showAlertIntegrations(id int) error {
	cli.StartProgress(" Fetching alert integrations...")
	integrationsResponse, err := cli.LwApi.V2.Alerts.GetIntegrations(id)
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "unable to show alert")
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(integrationsResponse.Data)
	}

	if len(integrationsResponse.Data) == 0 {
		cli.OutputHuman("There are no integration details available for the specified alert.\n")
		return nil
	}

	renderAlertIntegrationsTable(integrationsResponse.Data)
	return nil
}
