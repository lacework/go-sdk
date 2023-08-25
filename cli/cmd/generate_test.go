package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateEmailAddress(t *testing.T) {
	tests := []struct {
		Name     string
		Data     string
		Expected bool
	}{
		{Name: "simple_email", Data: "alice@example.com", Expected: true},
		{Name: "dots_in_local_part", Data: "alice.alison@example.com", Expected: true},
		{Name: "dots_in_domain_part", Data: "alice@internal.example.com", Expected: true},
		{Name: "dots_in_domain_both", Data: "alice.alison@internal.example.com", Expected: true},
		{Name: "dashes", Data: "alice-alison@internal-example.com", Expected: true},
		{Name: "plus_sign", Data: "alice+newsletter@example.com", Expected: true},
		{Name: "numbers_sign", Data: "alice1234567890@example.com", Expected: true},
		{Name: "empty_string", Data: "", Expected: false},
		{Name: "local_only", Data: "alice", Expected: false},
		{Name: "domain_only", Data: "example.com", Expected: false},
		{Name: "at_sign_only", Data: "@", Expected: false},
		{Name: "two_at_signs", Data: "alice@something@example.com", Expected: false},
		{Name: "no_domain", Data: "alice@", Expected: false},
		{Name: "no_local", Data: "@example.com", Expected: false},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			err := validateEmailAddress(test.Data)
			assert.Equal(t, err == nil, test.Expected)
		})
	}
}
