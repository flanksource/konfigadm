[![CircleCI](https://circleci.com/gh/moshloop/konfigadm.svg?style=svg)](https://circleci.com/gh/moshloop/konfigadm)
[![codecov](https://codecov.io/gh/moshloop/konfigadm/branch/master/graph/badge.svg)](https://codecov.io/gh/moshloop/konfigadm)
[![Go Report Card](https://goreportcard.com/badge/github.com/moshloop/konfigadm)](https://goreportcard.com/report/github.com/moshloop/konfigadm)

[docs](www.moshloop.com/konfigadm)
# konfigadm

konfigadm is a node instance configuration tool focused on bootstrapping nodes for container based environments

## Usage

```
Usage:
  konfigadm [command]

Available Commands:
  apply       Apply the configuration to the local machine
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


## Testing

* jrei/systemd-ubuntu:16.04
* jrei/systemd-ubuntu:18.04
* jrei/systemd-debian:9
* jrei/systemd-debian:latest
* jrei/systemd-centos:7
* jrei/systemd-fedora:latest


## Features


* Dependency free and easily embeddable into an image builder.
* Has built-in higher-order abstractions for kubernetes, containers, cri, cni, etc.
* Declarative, The order of operations cannot be changed, there are no implicit or explicit dependencies between items, no conditionals (besides for os/cloud tags) or control flows
* Typed, can validate the configuration(e.g. docker image name is valid, systemd.unit file only includes valid keys, and the values are typed correctly)
* Supports multiple operating systems and package managers.
* Abstractions and many of the built-in elements are easily unit-testable due to the use of virtual filesystem and command execution list.
* Automatic testing based on the declarations (If I have declared a service, I can automatically test if that service is running and health (or crash-looping), If I declare a container runtime, ensure that I can connect to it via it’s client )
* Generate cloud-init or shell scripts to be used by other systems


## Contributing

Make sure both unit and integration tests pass:

```bash
make
```

You can run unit tests only via:

```bash
make test
```

And only integration tests via:

```bash
make integration
```
