package cmd

import (
	"fmt"

	"github.com/lacework/go-sdk/api"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
)

func alertDetailsTable(alert api.Alert) (
	category, description, subject, entity [][]string,
) {
	category = append(category, []string{
		alert.DerivedFields.Source,
		alert.DerivedFields.Category,
		alert.DerivedFields.SubCategory,
		alert.PolicyID,
	})

	description = append(description, []string{
		alert.Info.Description,
	})

	subject = append(subject, []string{
		alert.Info.Subject,
	})

	entity = append(entity, []string{
		fmt.Sprintf(
			"Entity map coming soon. Use 'lacework alert show %d --json' for full details.",
			alert.ID,
		),
	})

	return
}

func renderAlertDetailsTable(alert api.Alert) {
	cat, desc, subj, entity := alertDetailsTable(alert)
	cli.OutputHuman(
		renderCustomTable(
			[]string{"Source", "Category", "Sub Category", "Policy ID"},
			cat,
			tableFunc(func(t *tablewriter.Table) {
				t.SetAutoWrapText(false)
				t.SetBorder(false)
			}),
		),
	)
	cli.OutputHuman("\n")
	cli.OutputHuman(
		renderCustomTable(
			[]string{"Description"},
			desc,
			tableFunc(func(t *tablewriter.Table) {
				t.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
				t.SetAutoWrapText(false)
				t.SetBorder(false)
			}),
		),
	)
	cli.OutputHuman("\n")
	cli.OutputHuman(
		renderCustomTable(
			[]string{"Subject"},
			subj,
			tableFunc(func(t *tablewriter.Table) {
				t.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
				t.SetAutoWrapText(false)
				t.SetBorder(false)
			}),
		),
	)
	cli.OutputHuman("\n")
	cli.OutputHuman(
		renderCustomTable(
			[]string{"Entities"},
			entity,
			tableFunc(func(t *tablewriter.Table) {
				t.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
				t.SetAutoWrapText(false)
				t.SetBorder(false)
			}),
		),
	)
}

func showAlertDetails(id int) error {
	cli.StartProgress(" Fetching alert details...")
	details, err := cli.LwApi.V2.Alerts.GetDetails(id)
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "unable to show alert")
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(details.Data)
	}

	alerts := api.Alerts{details.Data.Alert}
	renderAlertListTable(alerts)
	cli.OutputHuman("\n")
	renderAlertDetailsTable(details.Data.Alert)

	return nil
}
