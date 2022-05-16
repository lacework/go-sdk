package lwcomponent_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/lwcomponent"
)

func TestIsNotFound(t *testing.T) {
	err := errors.New("some other error")
	assert.Equal(t, false, lwcomponent.IsNotFound(err))

	err = lwcomponent.ErrComponentNotFound
	assert.Equal(t, true, lwcomponent.IsNotFound(err))
}
