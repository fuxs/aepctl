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

* [Access Control API](https://www.adobe.io/apis/experienceplatform/home/api-reference.html#!acpdr/swagger-specs/access-control.yaml) (External API documentation)
* [Offer Decisioning](https://experienceleague.adobe.com/docs/offer-decisioning/using/api-reference/getting-started.html?lang=en#api-reference) (External API documentation)
* [Schema Registry](doc/sr.md) commands
* [Identity Service](doc/is.md) commands
* [Query Service](doc/qs.md) commands

# Quick Start

1. Install `aepctl`
   * macOS
        ```terminal
        brew install fuxs/formulae/aepctl
        ```
    * Windows (requires PowerShell)
        ```terminal
        Invoke-WebRequest https://www.bungenstock.de/aepctl/releases/latest/windows/amd64/aepctl.exe -OutFile aepctl.exe
        ```
        Add `aepctl.exe` to your `PATH`
2. Create an [Adobe I/O Project](https://console.adobe.io/projects) ([detailed documentation](doc/new_project.md))
3. Provide a `config.yaml` file with the following command
    ```terminal
    aepctl configure
    ```
    Paste authentication credentials from the Adobe I/O project (click on *Service
      Account(JWT)* of your Adobe I/O project) and select a private key file. ([detailed documentation](doc/configuration.md))
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
aepctl in order to ease the input with complex IDs or names.

The zsh requires some code for the completion function which must be stored in a
file with the name `_aepctl`. This file must be located in a subdirectory of the
`$fpath` environment variable. Sounds too complicated? Just follow the next
steps:

Execute the following helper command:
```terminal
aepctl zsh
```
This creates the `_aepctl` file in your home directory `~/.aepctl/zsh_completion`

Now you have to add two lines to the `.zshrc` file. The first line adds the
directory of the created `_aepctl` file to the `$fpath` environment variable.
The second line with the `compinit` function activates the extended completion
system of zsh.

```terminal
cat <<EOT >> ~/.zshrc
fpath=(~/.aepctl/zsh_completion "${fpath[@]}")
autoload -U compinit; compinit
EOT
```

Some zsh frameworks like oh-my-zsh are calling `compinit` on their own. If you
use oh-my-zsh then you must update the `fpath` before the source command in
`.zshrc`.

A valid configuration could look like this:

```bash
fpath=(~/.aepctl/zsh_completion "${fpath[@]}")
source $ZSH/oh-my-zsh.s
```

### Completion for other shells

Please call the following command for other shells and follow the instructions:

```terminal
aepctl help completion
```


## Windows

### PowerShell
Open the PowerShell and follow the instructions:

1. Download the latest pre-release [aepctl
latest](https://www.bungenstock.de/aepctl/releases/latest/windows/amd64/aepctl.exe) or use `Invoke-WebRequest`:

```terminal
Invoke-WebRequest https://www.bungenstock.de/aepctl/releases/latest/windows/amd64/aepctl.exe -OutFile aepctl.exe
```

2. The validation of the binary is optional (go to step 4 if you want to skip
   it). Download the [SHA256
   file](https://www.bungenstock.de/aepctl/releases/latest/windows/amd64/aepctl.exe.sha256)
   or use `Invoke-WebRequest`:

```terminal
Invoke-WebRequest https://www.bungenstock.de/aepctl/releases/latest/windows/amd64/aepctl.exe.sha256 -OutFile aepctl.exe.sha256
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