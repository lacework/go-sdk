package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlertLinkBuilder(t *testing.T) {
	c := NewDefaultState()
	c.Account = "jackbauer"
	c.Subaccount = ""

	expected := "https://jackbauer.lacework.net/ui/investigation/monitor/AlertInbox/12345/details?accountName=jackbauer"
	assert.Equal(t, expected, alertLinkBuilderWithCLI(c, 12345))
}

func TestAlertLinkBuilderWithSubaccount(t *testing.T) {
	c := NewDefaultState()
	c.Account = "jackbauer"
	c.Subaccount = "24"

	expected := "https://jackbauer.lacework.net/ui/investigation/monitor/AlertInbox/678910/details?accountName=24"
	assert.Equal(t, expected, alertLinkBuilderWithCLI(c, 678910))
}
