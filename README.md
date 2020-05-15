# Sensu Go OpsGenie Handler
[![Sensu Bonsai Asset](https://img.shields.io/badge/Bonsai-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/nixwiz/sensu-opsgenie-handler)
![Go Test](https://github.com/nixwiz/sensu-opsgenie-handler/workflows/Go%20Test/badge.svg)
![goreleaser](https://github.com/nixwiz/sensu-opsgenie-handler/workflows/goreleaser/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/nixwiz/sensu-opsgenie-handler)](https://goreportcard.com/report/github.com/nixwiz/sensu-opsgenie-handler)

# Sensu Go OpsGenie Handler

## Table of Contents
- [Overview](#overview)
- [Files](#files)
- [Usage examples](#usage-examples)
- [Configuration](#configuration)
  - [Asset registration](#asset-registration)
  - [Handler definition](#handler-definition)
  - [Environment Variables](#environment-variables)
  - [Argument Annotations](#argument-annotations)
  - [Proxy support](#proxy-support)
- [Installation from source](#installation-from-source)
- [Contributing](#contributing)

## Overview

The Sensu Go OpsGenie Handler is a [Sensu Event Handler][3] which manages
[OpsGenie][2] incidents, for alerting operators. With this handler,
[Sensu][1] can trigger OpsGenie incidents.

This handler was inspired by [pagerduty plugin][6].

## Files

N/A

## Usage Examples

Help:
```
The Sensu Go OpsGenie handler for incident management

Usage:
  sensu-opsgenie-handler [flags]
  sensu-opsgenie-handler [command]

Available Commands:
  help        Help about any command
  version     Print the version number of this plugin

Flags:
  -a, --auth string                  The OpsGenie V2 API authentication token, use default from OPSGENIE_AUTHTOKEN env var
  -h, --help                         help for sensu-opsgenie-handler
  -i, --includeEventInNote           Include the event JSON in the payload sent to OpsGenie
  -m, --messageTemplate string       The template for the message to be sent (default "{{.Entity.Name}}/{{.Check.Name}}")
  -l, --messageLimit int             The maximum length of the message field (default 100)
  -d, --descriptionTemplate string   The template for the description to be sent (default "{{.Check.Output}}")
  -L, --descriptionLimit int         The maximum length of the description field (default 100)
  -s, --sensuDashboard string        The OpsGenie Handler will use it to create a source Sensu Dashboard URL. Use OPSGENIE_SENSU_DASHBOARD. Example: http://sensu-dashboard.example.local/c/~/n (default "disabled")
  -t, --team string                  The OpsGenie V2 API Team, use default from OPSGENIE_TEAM env var
  -u, --url string                   The OpsGenie V2 API URL, use default from OPSGENIE_APIURL env var (default "https://api.opsgenie.com")

Use "sensu-opsgenie-handler [command] --help" for more information about a command.

```

To configure OpsGenie Sensu Integration follow these first part in [OpsGenie Docs][5].

#### To use Opsgenie Priority

Please add this annotations inside sensu-agent:
```sh
# /etc/sensu/agent.yml example
annotations:
  opsgenie_priority: "P1"
```

Or inside check:
```yml
---
type: CheckConfig
api_version: core/v2
metadata:
  name: interval_check
  namespace: default
  annotations:
    opsgenie_priority: P2
    documentation": https://docs.sensu.io/sensu-go/latest
spec:
  command: check-cpu.sh -w 75 -c 90
  subscriptions:
  - system
  handlers:
  - opsgenie
  interval: 60
  publish: true
```

## Configuration
### Sensu Go
#### Asset registration

[Sensu Assets][7] are the best way to make use of this plugin. If you're not using an asset, please
consider doing so! If you're using sensuctl 5.13 with Sensu Backend 5.13 or later, you can use the
following command to add the asset:

```
sensuctl asset add nixwiz/sensu-opsgenie-handler
```

If you're using an earlier version of sensuctl, you can find the asset on the [Bonsai Asset Index][8].


#### Handler definition

```yml
type: Handler
api_version: core/v2
metadata:
  name: opsgenie
  namespace: default
spec:
  type: pipe
  command: sensu-opsgenie-handler
  env_vars:
  - OPSGENIE_TEAM=TEAM_NAME
  - OPSGENIE_APIURL=https://api.opsgenie.com
  timeout: 10
  runtime_assets:
  - nixwiz/sensu-opsgenie-handler
  filters:
  - is_incident
  secrets:
  - name: OPSGENIE_AUTHTOKEN
    secret: opgsgenie_authtoken
```

### Environment Variables

Most arguments for this handler are available to be set via environment variables.  However, any
arguments specified directly on the command line override the corresponding environment variable.

|Argument             |Environment Variable         |
|---------------------|-----------------------------|
|--url                |OPSGENIE_APIURL              |
|--auth               |OPSGENIE_AUTHTOKEN           |
|--team               |OPSGENIE_TEAM                |
|--withAnnotations    |OPSGENIE_ANNOTATIONS         |
|--sensuDashboard     |OPSGENIE_SENSU_DASHBOARD     |
|--messageTemplate    |OPSGENIE_MESSAGE_TEMPLATE    |
|--messageLimit       |OPSGENIE_MESSAGE_LIMIT       |
|--descriptionTemplate|OPSGENIE_DESCRIPTION_TEMPLATE|
|--descriptionLimit   |OPSGENIE_DESCRIPTION_LIMIT   |

**Security Note:** Care should be taken to not expose the auth token for this handler by specifying it
on the command line or by directly setting the environment variable in the handler definition.  It is
suggested to make use of [secrets management][10] to surface it as an environment variable.  The
handler definition above references it as a secret.  Below is an example secrets definition that make
use of the built-in [env secrets provider][11].

```yml
---
type: Secret
api_version: secrets/v1
metadata:
  name: opsgenie_authtoken
spec:
  provider: env
  id: OPSGENIE_AUTHTOKEN
```

### Argument Annotations

All arguments for this handler are tunable on a per entity or check basis based on annotations.  The
annotations keyspace for this handler is `sensu.io/plugins/sensu-opsgenie-handler/config`.

#### Examples

To change the team argument for a particular check, for that checks's metadata add the following:

```yml
type: CheckConfig
api_version: core/v2
metadata:
  annotations:
    sensu.io/plugins/sensu-opsgenie-handler/config/team: WebOps
[...]
```

### Proxy Support

This handler supports the use of the environment variables HTTP_PROXY,
HTTPS_PROXY, and NO_PROXY (or the lowercase versions thereof). HTTPS_PROXY takes
precedence over HTTP_PROXY for https requests.  The environment values may be
either a complete URL or a "host[:port]", in which case the "http" scheme is assumed.

### Sensu Core

See [this plugin][9]

## Installation from source

Download the latest version of the sensu-opsgenie-handler from [releases][4],
or create an executable script from this source.

From the local path of the sensu-opsgenie-handler repository:
```
go build -o /usr/local/bin/sensu-opsgenie-handler main.go
```

## Contributing

See https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md

[1]: https://github.com/sensu/sensu-go
[2]: https://www.opsgenie.com/ 
[3]: https://docs.sensu.io/sensu-go/5.0/reference/handlers/#how-do-sensu-handlers-work
[4]: https://github.com/nixwiz/sensu-opsgenie-handler/releases
[5]: https://docs.opsgenie.com/docs/sensu-integration#section-add-sensu-integration-in-opsgenie
[6]: https://github.com/sensu/sensu-pagerduty-handler
[7]: https://docs.sensu.io/sensu-go/latest/reference/assets/
[8]: https://bonsai.sensu.io/
[9]: https://github.com/sensu-plugins/sensu-plugins-opsgenie
[10]: https://docs.sensu.io/sensu-go/latest/guides/secrets-management/
[11]: https://docs.sensu.io/sensu-go/latest/guides/secrets-management/#use-env-for-secrets-management
