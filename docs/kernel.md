# Kernel

The kernel key can be used to specify the kernel version to be installed and
configured for boot.  The matching kernel headers package will also be
installed.  The value is expected to be a list of version-release strings.  To
standardardise input format, OS specific naming conventions such as including
target architecture and distribution tag are handled internally:

```yaml
kernel:

  - 5.3.11-100 #fedora
  - 3.10.0-1127.8.2 #centos7 centos8
  - 4.18.0-193 #centos8 rhel8
  - 4.19.0-8 #debian
  - 5.4.0-42 #ubuntu
  - 4.19.115-1 #photon
  - 4.14.165-131.185 #amazonLinux
```

Note that no additional repository setup is performed.  To include kernels
outside of default repositories, eg vaulted CentOS kernels or mainline kernels,
additional repository setup steps would be required.
