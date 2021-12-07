package cmd

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateSubAccountFlag(t *testing.T) {
	ok, _ := regexp.MatchString(ValidateSubAccountFlagRegEx, "simplestring")
	assert.False(t, ok, "subaccounts must be supplied in the format of `<profile>:<region>[,]`")

	ok, _ = regexp.MatchString(ValidateSubAccountFlagRegEx, "profile:region")
	assert.True(t, ok, "subaccounts can be supplied as single <profile>:<region> string or csv of them")

	ok, _ = regexp.MatchString(ValidateSubAccountFlagRegEx, "profile:region,profile2:region2")
	assert.True(t, ok, "subaccounts can be supplied as single <profile>:<region> string or csv of them")

	ok, _ = regexp.MatchString(ValidateSubAccountFlagRegEx, "profile:$#%RFEWF")
	assert.False(t, ok, "subaccount profile or region cannot contain special characters")
}
