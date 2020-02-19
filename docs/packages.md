# Packages

Packages can include modifiers:<br><br>
**-** removed if installed<br>
**+** update to the latest<br>
**=** mark to prevent future automatic updates<br>


```yaml
packages:
  - socat
  - -docker-common
  - -docker
  - =docker-ce==18.06
```

Packages can also leverage runtime flags:

```yaml
packages:
  - netcat #+debian
  - nmap-ncat #+redhat
  - open-vm-tools #+vmware
  - aws-cli #+aws
  - azure-cli #+azure
```
