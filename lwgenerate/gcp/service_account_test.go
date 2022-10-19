package gcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateSaCredFile(t *testing.T) {
	t.Run("JSON credentials file with client_email and private_key is valid", func(t *testing.T) {
		err := ValidateServiceAccountCredentialsFile("test_resources/private_key_client_email_valid.json")
		assert.Equal(t, err, nil)
	})

	t.Run("JSON credentials file without client_email is not  valid", func(t *testing.T) {
		err := ValidateServiceAccountCredentialsFile("test_resources/creds_no_client_email.json")
		assert.EqualError(t, err, "invalid Service Account credentials file. The private_key and client_email fields MUST be present.")
	})

	t.Run("JSON credentials file without private_key is not valid", func(t *testing.T) {
		err := ValidateServiceAccountCredentialsFile("test_resources/creds_no_private_key.json")
		assert.EqualError(t, err, "invalid Service Account credentials file. The private_key and client_email fields MUST be present.")
	})

	t.Run("invalid JSON file", func(t *testing.T) {
		err := ValidateServiceAccountCredentialsFile("test_resources/invalid_json.json")
		assert.EqualError(t, err, "unable to parse credentials file: invalid character '}' looking for beginning of object key string")
	})

	t.Run("non existent JSON file", func(t *testing.T) {
		err := ValidateServiceAccountCredentialsFile("test_resources/generate_gcp_test_data/foo.json")
		assert.EqualError(t, err, "provided credentials file does not exist")
	})
}
