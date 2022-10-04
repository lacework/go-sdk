<#
.SYNOPSIS
Codefresh tests for the Lacework CLI.

.DESCRIPTION
This script runs on our Codefresh pipeline to test the Lacework Command Line Interface.
#>

$ErrorActionPreference="stop"

Write-Host "Installing Go tools."
$env:GOFLAGS = '-mod=readonly'
go install github.com/mitchellh/gox@v1.0.1
go install gotest.tools/gotestsum@v1.7.0

Write-Host "Building Lacework CLI binaries."
$env:GOFLAGS = '-mod=vendor'
gox -output="bin/lacework-cli-{{.OS}}-{{.Arch}}" -os="windows" -arch="amd64 386" github.com/lacework/go-sdk/cli

Write-Host "Running integrations tests."
$env:Path += ";$pwd\bin"

# Disable vulnerability tests until https://lacework.atlassian.net/browse/RAIN-37563 is resolved
gotestsum -- -v github.com/lacework/go-sdk/integration -timeout 30m -tags="account agent_token alert_rules compliance configure event help integration migration policy query version generation"
