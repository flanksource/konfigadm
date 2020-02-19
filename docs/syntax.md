[![Build Status](https://travis-ci.org/moshloop/konfigadm.svg?branch=master)](https://travis-ci.org/moshloop/konfigadm)
[![codecov](https://codecov.io/gh/moshloop/konfigadm/branch/master/graph/badge.svg)](https://codecov.io/gh/moshloop/konfigadm)
[![Go Report Card](https://goreportcard.com/badge/github.com/moshloop/konfigadm)](https://goreportcard.com/report/github.com/moshloop/konfigadm)

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

## Design

![](./docs/flow.png)

### Mental Models

`konfigadm` intentionally reuses mental models and concepts from kubernetes, golang and ansible these include:

* Kubernetes declarative model for specifying intent
* Operators for providing higher-order abstractions
* Go build tags in comments for specifying behavior based on OS, Cloud etc..
* Ansible's way of defining variables and allowing for merging of multiple variable files.


### Apps

Apps provide an abstraction over low-level native and primitive elements, They describe high-level intent for using an application that may require multiple elements to configure.


#### Kubernetes

The kubernetes config element is the primary purpose of `konfigadm`, configuring machines so that they have all pre-requisites met for running `kubeadm`

* Install and mark the specific versions of `kubeadm`, `kubelet`, `kubectl`, `kubernetes-cni`
* Install a container runtime if not specified
* Prepull images required to run kubernetes
* Set any sysctl values that are required

```bash
konfigadm apply -c k8s.yml
```

`k8s.yml`
```yaml
kubernetes:
  version: 1.14.1
```

The config can also be specified via stdin: `echo "kubernetes: {version: 1.14.1}" | konfigadm minify -c -`

```yaml
kubernetes:
  version: 1.14.1
```

#### Container Runtimes (CRI)

```yaml
cri:
 version: 18.6.0
 type: docker
 config:
   log-driver: json-file
   log-opts:
     max-size: 1000m
     max-file": 3
```

### Native

Native elements, are not application specific they include packages, repositories, keys, containers, sysctls and environment variables.

e.g running `echo "kubernetes: {version: 1.14.1}" | konfigadm minify -c -` will result in the application being transformed into native elements.
```yaml
packageRepos:
 - deb https://apt.kubernetes.io/ kubernetes-xenial main #+debian
 - https://packages.cloud.google.com/yum/repos/kubernetes-el7-x86_64 #+redhat
gpg:
 - https://packages.cloud.google.com/apt/doc/apt-key.gpg #+debian
 - https://packages.cloud.google.com/yum/doc/yum-key.gpg #+redhat
 - https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg #+redhat
packages:
 - kubelet=1.14.1
 - kubeadm=1.14.1
 - kubectl=1.14.1
sysctls:
 vm.swapinness: 1
```

Native elements are verifiable, i.e. if you specify a container runtime then `konfigadm` will verify that the runtime has a service enabled and started and that `root` can connect to the daemon and list running containers.

### Primitives
Primitives are the low-level commands and files that are need to implement native items.

For example a `package: [curl]` native element would create a `apt-get install -y curl` primitive command on debian systems and `yum install -y curl` on redhat systems

The relationship between the 3 kinds is similar to Deployment, ReplicaSet and Pod. Apps insert and/or update native elements, native elements are then “compiled” down to primitives.







```

### Services



## Natives





