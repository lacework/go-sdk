<#
.SYNOPSIS
Uninstalls the Lacework CLI.

Authors: Technology Alliances <tech-ally@lacework.net>

.DESCRIPTION
This script removes any env vars created by Lacework CLI installation.

#>
$ErrorActionPreference = 'stop';

Uninstall-ChocolateyPackage -PackageName lacework-cli

[Environment]::SetEnvironmentVariable("LW_CHOCOLATEY_INSTALL", $null, "Machine")
