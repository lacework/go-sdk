New-Item -ItemType SymbolicLink -Target "C:\devcli\lacework-cli-windows-amd64.exe" -Path "C:\Users\vagrant\lacework.exe"

# add LW_ environment variables
$empty = '@@empty@@'
if ( $args[0] -ne $empty )
{
    [System.Environment]::SetEnvironmentVariable('LW_ACCOUNT', $args[0], [System.EnvironmentVariableTarget]::Machine)
}
if ( $args[1] -ne $empty )
{
    [System.Environment]::SetEnvironmentVariable('LW_API_KEY', $args[1], [System.EnvironmentVariableTarget]::Machine)
}
if ( $args[2] -ne $empty )
{
    [System.Environment]::SetEnvironmentVariable('LW_API_SECRET', $args[2], [System.EnvironmentVariableTarget]::Machine)
}
if ( $args[3] -ne $empty )
{
    [System.Environment]::SetEnvironmentVariable('LW_API_TOKEN', $args[3], [System.EnvironmentVariableTarget]::Machine)
}
