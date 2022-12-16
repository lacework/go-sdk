package cmd

import (
	"testing"

	"github.com/lacework/go-sdk/api"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func testGetRemediationTemplateIDsOK() ([]string, error) {
	return []string{"LW_foo", "lacework-global-40", "lwcustom-11"}, nil
}

func testGetRemediationTemplateIDsFail() ([]string, error) {
	return []string{}, errors.New("ah-ah-ah you didn't say the magic workd")
}

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
	alertsActual, err := filterFixableAlerts(alertsIn, testGetRemediationTemplateIDsOK)
	assert.Nil(t, err)
	assert.Equal(t, alertsExpected, alertsActual)

	alertsActual, err = filterFixableAlerts(alertsIn, testGetRemediationTemplateIDsFail)
	assert.NotNil(t, err)
	assert.Equal(t, alertsIn, alertsActual)
}
