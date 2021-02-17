# aepctl

aepctl is a command line tool for the [Adobe Experience
Platform](https://experienceleague.adobe.com/docs/experience-platform/landing/home.html)
implementing a part of the [REST
API](https://www.adobe.io/apis/experienceplatform/home/api-reference.html).

This is the initial state of this project and no release is available.

# Overview

aepctl is a complement to the existing web interface and has been developed for
advanced users as well as developers. In combination with activated syntax
completion, aepctl accelerates the execution of repeating tasks, prototyping and
learning the APIs.

# Status of Implementation

At the moment the following APIs are implemented:

* [Access Control API](https://www.adobe.io/apis/experienceplatform/home/api-reference.html#!acpdr/swagger-specs/access-control.yaml)
* [Offer Decisioning](https://experienceleague.adobe.com/docs/offer-decisioning/using/api-reference/getting-started.html?lang=en#api-reference)

Planed APIs:

* Catalog Service API

# Quick Start

1. Install `aepctl`
   * macOS
        ```terminal
        brew install fuxs/formulae/aepctl
        ```
    * Windows (requires PowerShell)
        ```terminal
        Invoke-WebRequest https://www.bungenstock.de/aepctl/releases/v0.1.1/windows/amd64/aepctl.exe -OutFile aepctl.exe
        ```
        Add `aepctl.exe` to your `PATH`
2. Create an [Adobe I/O Project](https://console.adobe.io/projects)
3. Provide a  `config.yaml` file with the following command
    ```terminal
    aepctl configure
    ```
    Paste authentication credentials from the Adobe I/O project (click on *Service
      Account(JWT)* of your Adobe I/O project)
4. Test the configuration by getting an access token.
    ```terminal
    aepctl get token
    ```
   
# Installation
## macOS

The recommended installation method is [homebrew](https://brew.sh/). Visit the
website and install it if you haven't already.

Run the following command to install aepctl:

```terminal
brew install fuxs/formulae/aepctl
```
### zsh completion
The zsh is the default shell since macOS 10.15 Catalina and provides strong
completion capabilities. It is recommended to activate completions for
aepctl in order to ease the input with complex ids or names.

## Windows

### PowerShell
Open the PowerShell and follow the instructions:

1. Download the latest pre-release [aepctl
v0.1.1](https://www.bungenstock.de/aepctl/releases/v0.1.1/windows/amd64/aepctl.exe) or use `Invoke-WebRequest`:

```terminal
Invoke-WebRequest https://www.bungenstock.de/aepctl/releases/v0.1.1/windows/amd64/aepctl.exe -OutFile aepctl.exe
```

2. The validation of the binary is optional (go to step 4 if you want to skip
   it). Download the [SHA256
   file](https://www.bungenstock.de/aepctl/releases/v0.1.1/windows/amd64/aepctl.exe.sha256)
   or use `Invoke-WebRequest`:

```terminal
Invoke-WebRequest https://www.bungenstock.de/aepctl/releases/v0.1.1/windows/amd64/aepctl.exe.sha256 -OutFile aepctl.exe.sha256
```
3. Check the integrity with the following command. You should get the value
   `True` as result.

```terminal
(Get-FileHash aepctl.exe).Hash -eq (Get-Content aepctl.exe.sha256)
```

4. Add the aepctl.exe to your `PATH`

5. Test that everything is working by opening the help.
```terminal
aepctl --help
```
# License
`aepctl` is released under the Apache 2.0 license.