package gcp

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"os"

	"github.com/lacework/go-sdk/internal/file"
	"github.com/pkg/errors"
)

type ServiceAccount struct {
	Name       string
	PrivateKey string
}

func NewServiceAccount(name string, privateKey string) *ServiceAccount {
	return &ServiceAccount{
		Name:       name,
		PrivateKey: privateKey,
	}
}

func (s *ServiceAccount) IsPartial() bool {
	if s == nil {
		return false
	}

	if s.Name == "" && s.PrivateKey == "" {
		return false
	}

	if s.Name != "" && s.PrivateKey != "" {
		return false
	}

	return true
}

func ValidateServiceAccountCredentialsFile(credFile string) error {
	if file.FileExists(credFile) {
		jsonFile, err := os.Open(credFile) // guardrails-disable-line
		if err != nil {
			return errors.Wrap(err, "issue opening credentials file")
		}
		defer jsonFile.Close()

		byteValue, err := io.ReadAll(jsonFile)
		if err != nil {
			return errors.Wrap(err, "unable to parse credentials file")
		}

		var credFileContent map[string]interface{}
		err = json.Unmarshal(byteValue, &credFileContent)
		if err != nil {
			return errors.Wrap(err, "unable to parse credentials file")
		}
		credFileContent, valid := ValidateSaCredFileContent(credFileContent)
		if !valid {
			return errors.New("invalid Service Account credentials file. " +
				"The private_key and client_email fields MUST be present.")
		}
	} else {
		return errors.New("provided credentials file does not exist")
	}
	return nil
}

func ValidateSaCredFileContent(credFileContent map[string]interface{}) (map[string]interface{}, bool) {
	if credFileContent["private_key"] != nil && credFileContent["client_email"] != nil {
		privateKey, ok := credFileContent["private_key"].(string)
		if !ok {
			return credFileContent, false
		}
		err := ValidateStringIsBase64(privateKey)
		if err != nil {
			privateKey := base64.StdEncoding.EncodeToString([]byte(privateKey))
			credFileContent["private_key"] = privateKey
			return credFileContent, true
		}
	}
	return credFileContent, false
}

func ValidateStringIsBase64(val interface{}) error {
	switch value := val.(type) {
	case string:
		_, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			return errors.New("provided private key is not base64 encoded")
		}
	default:
		return errors.New("value must be a string")
	}

	return nil
}

func ValidateServiceAccountCredentials(val interface{}) error {
	if value, ok := val.(string); ok {
		if value == "" {
			return nil
		} else {
			return ValidateServiceAccountCredentialsFile(value)
		}
	}

	return errors.New("value must be a string")
}
