package cmd

import (
	"testing"

	"github.com/lacework/go-sdk/api"
	"github.com/stretchr/testify/assert"
)

func TestFilterFixableAlerts(t *testing.T) {
	alertsIn := api.Alerts{
		{
			PolicyID: "not-fixable",
		},
		{
			PolicyID: "lacework-global-40",
		},
		{
			PolicyID: "dev7-lwcustom-11",
		},
	}
	alertsExpected := api.Alerts{
		{
			PolicyID: "lacework-global-40",
		},
		{
			PolicyID: "dev7-lwcustom-11",
		},
	}
	alertsActual := filterFixableAlerts(
		alertsIn, []string{"LW_foo", "lacework-global-40", "lwcustom-11"})
	assert.Equal(t, alertsExpected, alertsActual)

	alertsActual = filterFixableAlerts(alertsIn, []string{})
	assert.Equal(t, api.Alerts{}, alertsActual)
}
