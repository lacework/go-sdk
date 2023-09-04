package validate

import (
	"net/mail"

	"github.com/pkg/errors"
)

// validate email address against RFC 5322
func EmailAddress(val interface{}) error {
	switch value := val.(type) {
	case string:
		// This validates the email format against RFC 5322
		_, err := mail.ParseAddress(value)
		if err != nil {
			return errors.Wrap(err, "supplied email address is not a valid email address format")
		}
	default:
		return errors.New("value must be a string")
	}
	return nil
}
