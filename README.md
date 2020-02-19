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
  <a href="DESIGN.md">Design</a> •
  <a href="https://github.com/moshloop/konfigadm-images/releases">Prebuilt Images</a> •
  <a href="https://www.moshloop.com/konfigadm"> Full Documentation </a>
</p>

`konfigadm` is a declarative configuration management tool and image builder focused on bootstrapping nodes for container based environments.

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

## Installation

### Ubuntu / Debian

```bash
wget https://github.com/flanksource/konfigadm/releases/download/v0.4.2/konfigadm.deb
dpkg -i konfigadm.deb
```

### Centos / Fedora / Redhat

```bash
rpm -i https://github.com/flanksource/konfigadm/releases/download/v0.4.2/konfigadm.rpm
```

### Binary

```bash
wget -O /usr/bin/konfigadm https://github.com/flanksource/konfigadm/releases/download/v0.4.2/konfigadm && sudo chmod +x /usr/bin/konfigadm
```

## Getting Started

### Installing Kubernetes

```bash
sudo konfigadm apply -c - <<-EOF
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
sudo konfigadm build-image --image ubuntu:1804 -c - <<-EOF
kubernetes:
  version: 1.14.2
container_runtime:
  type: docker
cleanup: true
EOF
```

Cloud Images are downloaded and then configured with `--build-driver` 2 drivers are supported:

1. `qemu` (default) - Launches the image with KVM and attaches a cloud-init ISO to configure on boot
2. `libguestfs` - Uses virt-customize to launch an appliance and chroot into the disk, does not require cloud-init in the image, but also cannot test/verify systemd based services due to the chroot.

[![asciicast](https://asciinema.org/a/252399.svg)](https://asciinema.org/a/252399)


## Features

* **Dependency Free** and easily embeddable into an image builder.
* **Declarative**, The order of operations cannot be changed, there are no implicit or explicit dependencies between items, no conditionals (besides for os/cloud tags) or control flows
* **Typed**, can validate the configuration (e.g. docker image name is valid, systemd.unit file only includes valid keys, and the values are typed correctly)
* Has built-in higher-order abstractions for kubernetes, containers, cri, cni, etc.
* Supports multiple operating systems and package managers
* Abstractions and many of the built-in elements are easily unit-testable due to the use of virtual filesystem and command execution list.
* Automatic testing / verification based on intent, not just command success code
* Generate cloud-init or shell scripts to be used by other systems

## Compatibility

Compatibility is tested via the docker systemd images created by [jrei](https://github.com/j8r/dockerfiles/tree/master/systemd), All example fixtures are first verified as false, applied, and then verified as true.

To run integration tests:

```bash
make ubuntu
```

**Compatibility Matrix**

| Target   | OS                      | Status                                                       | Tags                       |
| -------- | ----------------------- | ------------------------------------------------------------ | -------------------------- |
| ubuntu16 | Ubuntu 16.04            | ![](https://img.shields.io/badge/-PASSED-brightgreen.svg?logo=circleci) | `#ubuntu debian-like`      |
| ubuntu   | Ubuntu 18.04            | ![](https://img.shields.io/badge/-PASSED-brightgreen.svg?logo=circleci) | `#ubuntu debian-like`      |
| centos   | Centos 7                | ![](https://img.shields.io/badge/-PASSED-brightgreen.svg?logo=circleci) | `#centos redhat-like`      |
| debian9  | Debian 9                | ![](https://img.shields.io/badge/-PASSED-brightgreen.svg?logo=circleci) | `#debian debian-like`      |
| fedora   | Fedora Latest           | ![](https://img.shields.io/badge/-FAILED-red.svg?logo=circleci) | `#fedora `                 |
|          | Amazon Linux            | ![](https://img.shields.io/badge/-UNTESTED-gray.svg) but should work | `#amazonLinux redhat-like` |
|          | Redhat Enterprise Linux | ![](https://img.shields.io/badge/-UNTESTED-gray.svg) but should work | `#rhel redhat-like`        |

## TODO

* Incremental mode
* Merge duplicate command dependencies (e.g. installing curl)
* Support templating everywhere (currently only supported in files)
* Packer/QEMU/VirtualBox/Fusion drivers for building images
* AMI/OVA Image upload
* Multi-OS cleanup scripts for building images
