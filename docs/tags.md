# Runtime Tags

```bash
konfigadm minify -c config.yml --tags ubuntu
# tags are detected by default when using the apply command
konfigadm apply -c config
```

Similar to go build tags, runtime tags provide a way of deciding what gets run, the following tags are provided by default:


* `centos`
* `ubuntu`
* `fedora`
* `debian` matched for all debian based distros (ubuntu)
* `rhel`
* `redhat` matched for all redhat based distros (centos, fedora, rhel, amazon linux)
* `amazonLinux`
* `aws` matched when running inside Amazon Web Services
* `azure` matched when running inside Azure
* `vmware` matched when running on a vSphere Hypervisor
* `kvm` matched when running on a KVM Hypervisor

Tags can be applied to the following elements:
* packages
* packageRepos
* packageKeys
* commands, pre_commands, post_commands

Multiple tags can be specified in which case all tags must match.
```yaml
packages:
  # only install aws-cli on debian based system running in AWS
  - aws-cli #+debian +aws
```

Tags can be negated using `!`

```yaml
pre_commands:
  # attach a rhel subscription, but only if we are not running in AWS
  - subscription manager attach #+rhel !aws
```
