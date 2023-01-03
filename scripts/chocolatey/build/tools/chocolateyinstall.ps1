<#
.SYNOPSIS
Installs the Lacework CLI.

Authors: Technology Alliances <tech-ally@lacework.net>

.DESCRIPTION
This script installs the Lacework Command Line Interface.

#>

$ErrorActionPreference = "stop"

$version = $lacework_cli_version
$packageName = 'lacework-cli'
$toolsDir = "$( Split-Path -parent $MyInvocation.MyCommand.Definition )"
$url = "https://github.com/lacework/go-sdk/releases/download/v$lacework_cli_version/lacework-cli-windows-386.exe.zip"
$url64 = "https://github.com/lacework/go-sdk/releases/download/v$lacework_cli_version/lacework-cli-windows-amd64.exe.zip"

$packageArgs = @{
    packageName = $packageName
    unzipLocation = $toolsDir
    fileType = 'zip'
    url = $url
    url64bit = $url64

    softwareName = 'Lacework*'

    checksum = $checksum
    checksumType = 'sha256'
    checksum64 = $check64
    checksumType64 = 'sha256'
}

Install-ChocolateyZipPackage @packageArgs


Function Set-Environment-Variables
{
    $isAdmin = $false
    try
    {
        $currentPrincipal = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())
        $isAdmin = $currentPrincipal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
    }
    finally
    {
        if ($isAdmin)
        {
            $machinePath = [System.Environment]::GetEnvironmentVariable("PATH", "Machine")
            $machinePath = New-PathString -StartingPath $machinePath -Path $laceworkPath
            [System.Environment]::SetEnvironmentVariable("PATH", $machinePath, "Machine")

            ## Set Chocolatey environment variable
            [System.Environment]::SetEnvironmentVariable("LW_CHOCOLATEY_INSTALL", 1, "Machine")
        }
        else
        {
            $userPath = [System.Environment]::GetEnvironmentVariable("PATH", "User")
            $userPath = New-PathString -StartingPath $userPath -Path $laceworkPath
            [System.Environment]::SetEnvironmentVariable("PATH", $userPath, "User")

            ## Set Chocolatey environment variable
            [System.Environment]::SetEnvironmentVariable("LW_CHOCOLATEY_INSTALL", 1, "Machine")
        }
    }
}