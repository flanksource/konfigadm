[![CircleCI](https://circleci.com/gh/moshloop/konfigadm.svg?style=svg)](https://circleci.com/gh/moshloop/konfigadm)
[![codecov](https://codecov.io/gh/moshloop/konfigadm/branch/master/graph/badge.svg)](https://codecov.io/gh/moshloop/konfigadm)
[![Go Report Card](https://goreportcard.com/badge/github.com/moshloop/konfigadm)](https://goreportcard.com/report/github.com/moshloop/konfigadm)

# konfigadm

konfigadm is a node instance configuration tool focused on bootstrapping nodes for container based environments

## Usage

```
Usage:
  konfigadm [command]

Available Commands:
  apply       Apply the configuration to the local machine or to remove machines over SSH
  cloud-init  Exports the configuration in cloud-init format
  help        Help about any command
  minify      Resolve all lookups and dependencies and export a single config file
  verify      Verify that the configuration has been applied correctly and is in a healthy state
  version     Print the version of konfigadm

Flags:
  -c, --config strings   Config files in YAML or JSON format
  -d, --detect           Detect tags to use
  -h, --help             help for konfigadm
  -v, --loglevel count   Increase logging level
  -t, --tag strings      Runtime tags to use, valid tags:  debian,ubuntu,redhat,rhel,centos,aws,vmware
  -e, --var strings      Extra Variables to in key=value format
```

## Design

![](./docs/flow.png)
