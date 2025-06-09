package aws

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func a(path string) string {
	return fmt.Sprintf("arn:aws:iam::123456789012:%s", path)
}

func TestParseResourceName(t *testing.T) {
	t.Run("should correctly parse assumed role", func(t *testing.T) {
		name, err := ParseResourceName(a("assumed-role/preflight_ro/aws-go-sdk-00000000000"))
		assert.NoError(t, err)
		assert.Equal(t, "preflight_ro", name)
	})

	t.Run("should correctly extract simple path", func(t *testing.T) {
		name, err := ParseResourceName(a("user/MyUser"))
		assert.NoError(t, err)
		assert.Equal(t, "MyUser", name)
	})

	t.Run("should correctly extract deep path", func(t *testing.T) {
		name, err := ParseResourceName(a("role/application_abc/component_xyz/RDSAccess"))
		assert.NoError(t, err)
		assert.Equal(t, "RDSAccess", name)
	})

	t.Run("should correctly extract root", func(t *testing.T) {
		name, err := ParseResourceName(a("root"))
		assert.NoError(t, err)
		assert.Equal(t, "root", name)
	})

	t.Run("should correctly extract value with spaces", func(t *testing.T) {
		name, err := ParseResourceName(a("u2f/user/JohnDoe/default (U2F security key)"))
		assert.NoError(t, err)
		assert.Equal(t, "default (U2F security key)", name)
	})
}
