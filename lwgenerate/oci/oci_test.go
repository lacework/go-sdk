package oci

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerationConfigNoArgs(t *testing.T) {
	_, err := NewTerraform(true).Generate()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid inputs")
}

func TestGenerationConfig(t *testing.T) {
	args := []OciTerraformModifier{
		WithTenantOcid("ocid1.tenancy...a"),
		WithUserEmail("a@b.c"),
	}
	hcl, err := NewTerraform(true, args...).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, hcl, ConfigResultBasic)
}

var ConfigResultBasic = `terraform {
  required_providers {
    lacework = {
      source  = "lacework/lacework"
      version = ">= 1.9.0"
    }
  }
}

module "oci_config" {
  source     = "lacework/config/oci"
  version    = "~> 0.2"
  tenancy_id = "ocid1.tenancy...a"
  user_email = "a@b.c"
}
`

func TestGenerationConfigCustomIntegrationName(t *testing.T) {
	args := []OciTerraformModifier{
		WithTenantOcid("ocid1.tenancy...a"),
		WithUserEmail("a@b.c"),
		WithConfigName("oci_test_config"),
	}
	hcl, err := NewTerraform(true, args...).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Contains(t, hcl, `integration_name = "oci_test_config"`)
}

func TestGenerationConfigCustomLaceworkProfile(t *testing.T) {
	args := []OciTerraformModifier{
		WithTenantOcid("ocid1.tenancy...a"),
		WithUserEmail("a@b.c"),
		WithLaceworkProfile("my_profile"),
	}
	hcl, err := NewTerraform(true, args...).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Contains(t, hcl, "provider \"lacework\" {\n  profile = \"my_profile\"\n}")
}
