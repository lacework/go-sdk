<#
.SYNOPSIS
Codefresh tests for the Lacework CLI.

.DESCRIPTION
This script runs on our Codefresh pipeline to test the Lacework Command Line Interface.
#>

Set-ExecutionPolicy Bypass -Scope Process -Force
iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/lacework/go-sdk/main/cli/install.ps1'))

lacework version

lacework int list -a $env:CI_V2_ACCOUNT -k $env:CI_API_KEY -s $env:CI_API_SECRET
