FROM  mcr.microsoft.com/windows/nanoserver:1903
LABEL maintainer="tech-ally@lacework.net" \
      description="The Lacework CLI helps you manage the Lacework cloud security platform"

COPY ./LICENSE /
RUN Set-ExecutionPolicy Bypass -Scope Process -Force
RUN iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/lacework/go-sdk/master/cli/install.ps1'))
