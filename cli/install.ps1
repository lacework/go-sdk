<#
.SYNOPSIS
Installs the Lacework CLI.

Authors: Technology Alliances <tech-ally@lacework.net>

.DESCRIPTION
This script installs the Lacework Command Line Interface.

.Parameter Version
Specifies a version (ex: v0.1.0)
#>

param (
    [Alias("v")]
    [string]$Version
)

$ErrorActionPreference="stop"

Set-Variable GithubReleasesRootUrl -Option ReadOnly -value "https://github.com/lacework/go-sdk/releases" -Force
Set-Variable PackageName -Option ReadOnly -value "lacework-cli-windows-amd64.exe.zip" -Force

Function Get-File($url, $dst) {
    Write-Host "Downloading $url"
    try {
        [System.Net.ServicePointManager]::SecurityProtocol = [Enum]::ToObject([System.Net.SecurityProtocolType], 3072)
    } catch {
        Write-Error "TLS 1.2 is not supported on this operating system. Upgrade or patch your Windows installation."
    }
    $wc = New-Object System.Net.WebClient
    $wc.DownloadFile($url, $dst)
}

Function Get-WorkDir {
    $parent = [System.IO.Path]::GetTempPath()
    [string] $name = [System.Guid]::NewGuid()
    New-Item -ItemType Directory -Path (Join-Path $parent $name)
}

Function Get-Archive($version) {
    $url = $GithubReleasesRootUrl
    $package_name = $PackageName
    if(!$version -Or $version -eq "latest") {
        $lacework_cli_url="$url/latest/download/${package_name}"
    } else {
        $lacework_cli_url="$url/download/${version}/${package_name}"
    }
    $sha_url="$lacework_cli_url.sha256sum"
    $cli_dest = (Join-Path ($workdir) "lacework-cli.zip")
    $sha_dest = (Join-Path ($workdir) "lacework-cli.zip.shasum256")

    Get-File $lacework_cli_url $cli_dest
    $result = @{ "zip" = $cli_dest }

    try {
        Get-File $sha_url $sha_dest
        $result["shasum"] = (Get-Content $sha_dest).Split()[0]
    } catch {
        Write-Warning "No shasum exists for $version. Skipping validation."
    }
    $result
}

function Get-SHA256Converter {
    if($PSVersionTable.PSEdition -eq 'Core') {
        [System.Security.Cryptography.SHA256]::Create()
    } else {
        New-Object -TypeName Security.Cryptography.SHA256Managed
    }
}

Function Get-Sha256($src) {
    $converter = Get-SHA256Converter
    try {
        $bytes = $converter.ComputeHash(($in = (Get-Item $src).OpenRead()))
        return ([System.BitConverter]::ToString($bytes)).Replace("-", "").ToLower()
    } finally {
        # Older .Net versions do not expose Dispose()
        if($PSVersionTable.PSEdition -eq 'Core' -Or ($PSVersionTable.CLRVersion.Major -ge 4)) {
            $converter.Dispose()
        }
        if ($null -ne $in) { $in.Dispose() }
    }
}

Function Assert-Shasum($archive) {
    Write-Host "Verifying the shasum digest matches the downloaded archive"
    $actualShasum = Get-Sha256 $archive.zip
    if ($actualShasum -ne $archive.shasum) {
        Write-Error "Checksum '$($archive.shasum)' invalid."
    }
}

Function Install-Lacework-CLI {
    $laceworkPath = Join-Path $env:ProgramData Lacework
    if (-not (Test-Path $laceworkPath)) { New-Item $laceworkPath -ItemType Directory | Out-Null }
    $exe = (Get-ChildItem (Join-Path ($workdir) "bin"))
    $env:PATH = New-PathString -StartingPath $env:PATH -Path $laceworkPath

    try {
        Copy-Item "$($exe.FullName)" $laceworkPath -Force
    }
    catch {
        $exeOwner = Get-Acl (Join-Path $laceworkPath "lacework.exe") | Select-Object Owner
        Write-Error "Unable to install the Lacework CLI. The executable is owned by $exeOwner"
    }

    $isAdmin = $false
    try {
        $currentPrincipal = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())
        $isAdmin = $currentPrincipal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
    } finally {
        if ($isAdmin) {
            $machinePath = [System.Environment]::GetEnvironmentVariable("PATH", "Machine")
            $machinePath = New-PathString -StartingPath $machinePath -Path $laceworkPath
            [System.Environment]::SetEnvironmentVariable("PATH", $machinePath, "Machine")
        } else {
            $userPath = [System.Environment]::GetEnvironmentVariable("PATH", "User")
            $userPath = New-PathString -StartingPath $userPath -Path $laceworkPath
            [System.Environment]::SetEnvironmentVariable("PATH", $userPath, "User")
        }
    }
}

Function New-PathString([string]$StartingPath, [string]$Path) {
    if (-not [string]::IsNullOrEmpty($path)) {
        if (-not [string]::IsNullOrEmpty($StartingPath)) {
            [string[]]$PathCollection = "$path;$StartingPath" -split ';'
            $Path = ($PathCollection |
                    Select-Object -Unique |
                    Where-Object {-not [string]::IsNullOrEmpty($_.trim())}
            ) -join ';'
        }
        $path
    } else {
        $StartingPath
    }
}

Function Expand-Zip($zipPath) {
    $dest = $workdir
    try {
        [System.Reflection.Assembly]::LoadWithPartialName("System.IO.Compression.FileSystem") | Out-Null
        [System.IO.Compression.ZipFile]::ExtractToDirectory($zipPath, $dest)
    } catch {
        try {
            $shellApplication = New-Object -com shell.application
            $zipPackage = $shellApplication.NameSpace($zipPath)
            $destinationFolder = $shellApplication.NameSpace($dest)
            $destinationFolder.CopyHere($zipPackage.Items())
        } catch{
            Write-Error "Unable to unzip files on this OS"
        }
    }
}

Function Assert-Lacework-CLI {
    Write-Host "Verifying installed Lacework CLI version"
    try { lacework version } catch {
        Write-Error "Unable to verify that the Lacework CLI was succesfully installed"
    }
}

Write-Host "Installing the Lacework CLI"

$workdir = Get-WorkDir
New-Item $workdir -ItemType Directory -Force | Out-Null
try {
    $archive = Get-Archive $version
    if($archive.shasum) {
        Assert-Shasum $archive
    }
    Expand-zip $archive.zip
    Install-Lacework-CLI
    Assert-Lacework-CLI

    Write-Host "The Lacework CLI has been successfully installed."
} finally {
    try { Remove-Item $workdir -Recurse -Force } catch {
        Write-Warning "Unable to delete $workdir"
    }
}
