package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateOciTenantOcid(t *testing.T) {
	tests := []struct {
		Name     string
		Data     string
		Expected bool
	}{
		{Name: "id_field_only", Data: "ocid1.tenancy...a", Expected: true},
		{Name: "id_and_region_fields", Data: "ocid1.tenancy..b.a", Expected: true},
		{Name: "all_fields", Data: "ocid1.tenancy.c.b.a", Expected: true},
		{Name: "all_fields_long_id", Data: "ocid1.tenancy.c.b.aaaaaabbbbbbbccccccc1111112222222xxxxxx333333aaaaabbbbbccccc", Expected: true},
		{Name: "all_fields_plus_future_use_field", Data: "ocid1.tenancy.c.b.x.a", Expected: true},
		{Name: "wrong_version_field", Data: "www.tenancy.c.b.a", Expected: false},
		{Name: "wrong_type_field", Data: "ocid1.instance.c.b.a", Expected: false},
		{Name: "too_many_fields", Data: "ocid1.instance.c.b.a.x.y", Expected: false},
		{Name: "too_few_fields", Data: "ocid1.instance.c.b", Expected: false},
		{Name: "empty_string", Data: "", Expected: false},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			err := validateOciTenantOcid(test.Data)
			assert.Equal(t, err == nil, test.Expected)
		})
	}
}

func TestValidateOciUserEmail(t *testing.T) {
	tests := []struct {
		Name     string
		Data     string
		Expected bool
	}{
		{Name: "valid_email", Data: "alice@example.com", Expected: true},
		{Name: "invalid_email", Data: "www.example.com", Expected: false},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			err := validateOciUserEmail(test.Data)
			assert.Equal(t, err == nil, test.Expected)
		})
	}
}
