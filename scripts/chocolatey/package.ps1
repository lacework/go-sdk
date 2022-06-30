<#
.Synopsis
    Builds the Lacework CLI Chocolatey package
.Description
    This script builds a Chocolatey package for the Lacework CLI,
    that can be published to push.chocolatey.org
#>

$buildDir = Join-Path -Path $pwd.Path -ChildPath '\scripts\chocolatey\build'

# Get latest version
$version = (Invoke-WebRequest -UseBasicParsing 'https://api.github.com/repos/lacework/go-sdk/releases/latest' |`
  ConvertFrom-Json | Select-Object -ExpandProperty tag_name).Replace('v', '')

Write-Output "Latest lacework-cli version: $version"

# Replace version in nuspec template
$content = (Get-Content -Encoding utf8 -Raw "$buildDir\lacework-cli.nuspec").
        Replace('${version}', $version)
[System.IO.File]::WriteAllText("$buildDir\lacework-cli.nuspec", $content)

# Replace version in installer script
$installerContent = (Get-Content -Encoding utf8 -Raw "$buildDir\tools\chocolateyInstall.ps1").
        Replace('${lacework_cli_version}', $version)
[System.IO.File]::WriteAllText("$buildDir\tools\chocolateyInstall.ps1", $installerContent)

# Package and Push Chocolatey Package
cd $buildDir
$api_key = [System.Environment]::GetEnvironmentVariable("API_KEY")

C:\ProgramData\chocolatey\choco.exe apikey --key $api_key --source https://push.chocolatey.org/
C:\ProgramData\chocolatey\choco.exe pack
C:\ProgramData\chocolatey\choco.exe push lacework-cli.${version}.nupkg --source https://push.chocolatey.org/ -k $api_key