package lwcomponent_test

import (
	"testing"

	"github.com/Masterminds/semver"
	"github.com/lacework/go-sdk/lwcomponent"
	"github.com/stretchr/testify/assert"
)

func TestApiInfoId(t *testing.T) {
	version, err := semver.NewVersion("1.1.1")
	if err != nil {
		panic(err)
	}

	var id int32 = 23

	info := lwcomponent.NewAPIInfo(id, "test", version, "", 0)

	result := info.Id()
	assert.Equal(t, id, result)
}

func TestApiInfoLatestVersion(t *testing.T) {
	var expectedVer string = "1.2.3"

	version, err := semver.NewVersion(expectedVer)
	if err != nil {
		panic(err)
	}

	info := lwcomponent.NewAPIInfo(1, "test", version, "", 0)

	result := info.LatestVersion()
	assert.Equal(t, expectedVer, result.String())
}
