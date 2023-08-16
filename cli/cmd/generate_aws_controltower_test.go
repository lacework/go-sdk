package cmd

import (
	"testing"

	"github.com/lacework/go-sdk/lwgenerate/aws_controltower"
	"github.com/stretchr/testify/assert"
)

func toggleControlTowerNonInteractive() {
	cli.noCache = !cli.noCache
	cli.nonInteractive = !cli.nonInteractive
}

func TestGenerateBasicControlTowerArgs(t *testing.T) {
	toggleControlTowerNonInteractive()
	defer toggleControlTowerNonInteractive()

	data := aws_controltower.GenerateAwsControlTowerTfConfigurationArgs{}
	err := promptAwsControlTowerGenerate(
		&data,
		&AwsControlTowerGenerateCommandExtraState{Output: "/tmp"},
	)

	assert.Nil(t, err)
}
