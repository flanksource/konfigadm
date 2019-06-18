<p align="center">
<a href="https://circleci.com/gh/moshloop/konfigadm"><img src="https://circleci.com/gh/moshloop/konfigadm.svg?style=svg"></a>
<a href="https://codecov.io/gh/moshloop/konfigadm"><img src="https://codecov.io/gh/moshloop/konfigadm/branch/master/graph/badge.svg"></a>
<a href="https://goreportcard.com/report/github.com/moshloop/konfigadm"><img src="https://goreportcard.com/badge/github.com/moshloop/konfigadm"></a>
<img src="https://img.shields.io/badge/OS-ubuntu%20%7C%20debian%20%7C%20centos%20%7C%20redhat%20%7C%20fedora-lightgrey.svg"/></a>
</p>

<p align="center">
  <a href="#installation">Installation</a> •
  <a href="#installing-kubernetes">Kubernetes</a> •
  <a href="#features">Key Features</a> •
  <a href="#compatibility">Compatibility</a> •
  <a href="https://www.moshloop.com/konfigadm"> Full Documentation </a>
</p>

konfigadm is a declarative configuration management tool focused on bootstrapping nodes for container based environments.

`konfigadm` differs from tools like ansible, puppet and chef in that it uses a Virtual Filesystem and Shell (VFS), configuration items don't actually make changes, they only update the VFS, changes to the system are only made from the VFS in the final step.

The declarative model (VFS) has a number of advantages over imperative tools where changes are parsed and applied in the same step:

* Unit testing across multiple different configurations becomes trivial
* The VFS can be "compliled" into cloud-init, or inspected to verify state after application
* "Operators" can easily be built that emit other configuration items, these are easily unit and integration testable

`konfigadm` differs to CFEngine in that it does not have any control flows (explicit or implicit) and the "model" only includes 2 low-level items, shell commands and files.

## Goals

* Avoid entire classes of common problems in configuration management by reducing the feature-set
* Simple mental model, that is easy to learn and test against by leveraging declarative configuration following the kubernetes model
* Batteries included for kubernetes environments
* Usable within tools as a go module or standalone as a CLI
* Support the notion of different environments (OS, Cloud) natively without the need for control flow

## Non-Goals

* Be a generic replacement for ansible / cloud-init
* Remote execution (The output of konfigadm can be fed into a remote execution engine)
* Control flow (If/when, loops, etc)

### Common Problems with ansible/puppet/chef

* Trying to install ansible, but requiring a specific version of python
* Installing python, but then requiring a specific version of a pip module in order to execute a  specific task
* Having a playbook work today, but break tomorrow because ansible or python modules have been upgraded / downgraded
* Running an ansible task and it reports a failure or success incorrectly, but you can't work out why, or how to reproduce it outside of ansible
* Trying to support different operating systems, but then having to maintain separate playbooks
* Calling the wrong playbook, because the when expression is wrong / error prone
* Adding a conditional flow to a task, but specifying an incorrect / invalid condition
* Attempting to call sub tasks / roles, but then having different variables reported from the sub-task

Many of these issues could be mitigated by unit testing, however ansible is fundamentally not unit-testable, modules almost exclusively update  global state (by installing packages, updating files or running commands) and there is no way of checking what a module has done or wants todo. The dry-run capability is not universally supported and is self-reported.

Integration testing with tools like molecule is possible, but cannot verify that no other changes were made.


### Common Problems with cloud-init

* Trying to template out a cloud-init YAML file but then getting indentation wrong
* Getting YAML indentation right, but then specifying a field in the wrong section
* Attempting to use a field, but not having the requisite packages installed for it to work
* Trying to support multiple environments within a single or group of YAML files


## Installation

### Ubuntu / Debian

```bash
wget https://github.com/moshloop/konfigadm/releases/download/0.2.0/konfigadm.deb
dpkg -i konfigadm.deb
```

### Centos / Fedora / Redhat

```bash
rpm -i https://github.com/moshloop/konfigadm/releases/download/0.2.0/konfigadm.rpm
```

### Binary

```bash
wget -O /usr/bin/konfigadm https://github.com/moshloop/konfigadm/releases/download/0.2.0/konfigadm && chmod +x /usr/bin/konfigadm
```

## Getting Started

### Installing Kubernetes

```bash
konfigadm apply -c - <<-EOF
kubernetes:
  version: 1.14.2
container_runtime:
  type: docker
commands:
  - kubeadm init
EOF
```

[![asciicast](https://asciinema.org/a/250079.png)](https://asciinema.org/a/250079)

### Building a kubernetes image

```bash
konfigadm build-image --image ubuntu:1804 -c - <<-EOF
kubernetes:
  version: 1.14.2
container_runtime:
  type: docker
cleanup: true
EOF
```

Cloud Images are downloaded and then configured with `--build-driver` (Currently only libguestfs is supported)

[![asciicast](https://asciinema.org/a/252399.svg)](https://asciinema.org/a/252399)

## Compatibility

Compatibility is tested via the docker systemd images created by [jrei](https://github.com/j8r/dockerfiles/tree/master/systemd), All example fixtures are first verified as false, applied, and then verified as true.

To run integration tests:

```bash
make ubuntu
```

**Compatibility Matrix**

| OS                      | Status                                                       | Tags                       |
| ----------------------- | ------------------------------------------------------------ | -------------------------- |
| Ubuntu 16.04            | ![](https://img.shields.io/badge/-PASSED-brightgreen.svg?logo=circleci) | `#ubuntu debian-like`      |
| Ubuntu 18.04            | ![](https://img.shields.io/badge/-PASSED-brightgreen.svg?logo=circleci) | `#ubuntu debian-like`      |
| Centos 7                | ![](https://img.shields.io/badge/-PASSED-brightgreen.svg?logo=circleci) | `#centos redhat-like`      |
| Debian 9                | ![](https://img.shields.io/badge/-PASSED-brightgreen.svg?logo=circleci) | `#debian debian-like`      |
| Fedora 29               | ![](https://img.shields.io/badge/-FAILED-red.svg?logo=circleci) | `#fedora `                 |
| Amazon Linux            | ![](https://img.shields.io/badge/-UNTESTED-gray.svg) but should work         | `#amazonLinux redhat-like` |
| Redhat Enterprise Linux | ![](https://img.shields.io/badge/-UNTESTED-gray.svg) but should work             | `#rhel redhat-like`        |


## Features

![](./docs/konfigadm.png)

konfigadm is a node instance configuration tool focused on bootstrapping nodes for container based environments

* **Dependency Free** and easily embeddable into an image builder.
* **Declarative**, The order of operations cannot be changed, there are no implicit or explicit dependencies between items, no conditionals (besides for os/cloud tags) or control flows
* **Typed**, can validate the configuration (e.g. docker image name is valid, systemd.unit file only includes valid keys, and the values are typed correctly)
* Has built-in higher-order abstractions for kubernetes, containers, cri, cni, etc.
* Supports multiple operating systems and package managers
* Abstractions and many of the built-in elements are easily unit-testable due to the use of virtual filesystem and command execution list.
* Automatic testing / verification based on intent, not just command success code
* Generate cloud-init or shell scripts to be used by other systems

## TODO

* Incremental mode
* Merge duplicate command dependencies (e.g. installing curl)
* Support templating everywhere (currently only supported in files)
* Packer/QEMU/VirtualBox/Fusion drivers for building images
* AMI/OVA Image upload
* Multi-OS cleanup scripts for building images

## Usage

```
Usage:
  konfigadm [command]

Available Commands:
  apply       Apply the configuration to the local machine
  build-image Build a new image using the specified image and konfig
  cloud-init  Exports the configuration in cloud-init format
  help        Help about any command
  minify      Resolve all lookups and dependencies and export a single config file
  verify      Verify that the configuration has been applied and is in a healthy state
  version     Print the version of konfigadm

Flags:
  -c, --config strings   Config files in YAML or JSON format
  -d, --detect           Detect tags to use
  -h, --help             help for konfigadm
  -v, --loglevel count   Increase logging level
  -t, --tag strings      Runtime tags to use, valid tags: debian,ubuntu,redhat,rhel,fedora,redhat-like,debian-like,centos,aws,vmware
  -e, --var strings      Extra Variables to in key=value format
```


