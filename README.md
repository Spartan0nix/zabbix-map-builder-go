# zabbix-map-builder-go

[![Go Package](https://pkg.go.dev/badge/github.com/Spartan0nix/zabbix-map-builder-go?status.svg)](https://pkg.go.dev/github.com/Spartan0nix/zabbix-map-builder-go)
[![Go report](https://goreportcard.com/badge/github.com/Spartan0nix/zabbix-map-builder-go)](https://goreportcard.com/report/github.com/Spartan0nix/zabbix-map-builder-go)
![Test workflow](https://github.com/Spartan0nix/zabbix-map-builder-go/actions/workflows/test.yml/badge.svg)
![Build workflow](https://github.com/Spartan0nix/zabbix-map-builder-go/actions/workflows/build.yml/badge.svg)

## Table of contents

- [Description](#description)
- [Zabbix-nested-groups](#zabbix-nested-groups)
- [Usage](#usage)
  - [Fixtures (optional)](#fixtures-(optional))
  - [Required environment variables](#required-environment-variables)
  - [Install](#install)
  - [Run](#run)
  - [Completion](#completion)

## Description

This CLI tool is used to help administrator build a zabbix map using the given host mappings (network devices, etc.).

## Mapping format

Mapping required to be in json format and to respect the following format.
```json
[
    {
        "local_host": "router-1",
        "local_interface": "eth0",
        "local_trigger_pattern": "Interface eth0(): Link down",
        "local_image": "Firewall_(64)",
        "remote_host": "router-2",
        "remote_interface": "eth0",
        "remote_trigger_pattern": "Interface eth0(): Link down",
        "remote_image": "Switch_(64)"
    }
]
```

***\*_host :***

Name of the host to used for the mapping (first host : local // second host : remote).
The value needs to be the name of the host on Zabbix for the search query to match.

***\*_interface :***

Name of the host interface attached to the other host.
This field is currently not utilize.

***\*_image :***

Name of the image used for the host.
The value needs to be the name of the image on Zabbix for the search query to match.

## Usage

### Examples (optional)

Export files are available in the *examples* folder.

- Hosts
The *'zbx_export_hosts.json'* export will create hosts that can be utilize with the *'docker-compose.yml'* file.

- Mappings
The *'mapping.json'* file can be used to create an example map.

> Using the *docker-compose.test.yml* stack combine with the export files can give you a good preview of the possibilities available with this CLI tool.

### Required environment variables

To use this tool, you will need to set up the following variables :
- ZABBIX_URL
- ZABBIX_USER
- ZABBIX_PWD

You can simply export the variable in your current shell :

<u>Linux :</u>
```bash
export ZABBIX_URL="http://<zabbix-server-IP-or-DNS>:<port>/zabbix/api_jsonrpc.php"
export ZABBIX_USER="some-zabbix-user"
export ZABBIX_PWD="some-zabbix-user-password"
```
Adding this configuration to your ~/.bashrc or ~/.zshrc will make the configuration persistent between shell.

<u>Windows (example for powershell) :</u>
```powershell
$env:ZABBIX_URL="http://<zabbix-server-IP-or-DNS>:<port>/zabbix/api_jsonrpc.php"
$env:ZABBIX_USER="some-zabbix-user"
$env:ZABBIX_PWD="some-zabbix-user-password"
```

### Install

1. With a script (available in the *scripts* folder):

    ```bash
    bash scripts/install.sh
    ```

2. Manually :

    Each time a new release is created, the cli is compiled and the resultant binaries are pushed as assets (https://github.com/Spartan0nix/zabbix-map-builder-go/releases).

    ```bash
    # Create a temp installation folder
    mkdir /tmp/zabbix-map-builder

    # Retrieve the archive for release $RELEASE
    wget -O /tmp/zabbix-map-builder/zabbix-map-builder.tar.gz https://github.com/Spartan0nix/zabbix-map-builder-go/releases/download/$RELEASE/zabbix-map-builder-go_$RELEASE_linux_amd64.tar.gz

    # Remove previous install
    sudo rm /usr/local/bin/zabbix-map-builder

    # Extract the archive
    tar -C /tmp/zabbix-map-builder -xzf /tmp/zabbix-map-builder/zabbix-map-builder.tar.gz

    # Move the binairy
    sudo mv /tmp/zabbix-map-builder/zabbix-map-builder /usr/local/bin

    # Update permissions
    sudo chown $(id -un):$(id -gn) /usr/local/bin/zabbix-map-builder

    # Remove temp installation folder
    rm -r /tmp/zabbix-map-builder
    ```

### Uninstall

1. With a script (available in the *scripts* folder):

    ```bash
    bash scripts/uninstall.sh
    ```

2. Manually :

    ```bash
    sudo rm /usr/local/bin/zabbix-map-builder
    ```

### Run
```
This CLI tool is used to help administrator build a zabbix map using the given host mappings (network devices, etc.).

Usage:
   [flags]

Flags:
  -c, --color string           color in hexadecimal used for the links between each hosts (default "000000")
  -v, --debug                  enable debug logging verbosity
      --dry-run                output to the shell the map definition without created it on the server
  -f, --file string            input file
  -h, --help                   help for this command
      --name string            name of map
  -o, --output string          output the parameters used to create the map
      --stack-hosts bools      connect multiple links to a single host. If set to false, if mapping will have is own hosts. This can be useful for infrastructure with redundant connexion (default [true])
      --trigger-color string   color in hexadecimal used for the links between each hosts when a trigger is in problem state (default "DD0000")
```

### Completion

1. Zsh completion

    If shell completion is not enabled in your current shell (oh-my-zsh not running for example), add the following config to your .zshrc :

    ```bash
    echo "autoload -U compinit; compinit" >> ~/.zshrc
    ```

    - To load completions only in the current shell :
    ```bash
    source <(zabbix-map-builder completion zsh); compdef _zabbix-map-builder zabbix-map-builder
    ```

    - To make the configuration persistent between shells :
    ```bash
    zabbix-map-builder completion zsh > "${fpath[1]}/_zabbix-map-builder"
    ```

2. Bash completion

    To use completion scripts with bash, you will need to install the "bash-completion" package following your package manager recommendations.


    - To load completions only in the current shell
    ```bash
    source <(zabbix-map-builder completion bash)
    ```

    - To make the configuration persistent between shells :
    ```bash
    zabbix-map-builder completion bash > /etc/bash_completion.d/zabbix-map-builder
    ```

3. Other completions

    - Completion for fish and powershell are available but haven't been tested.