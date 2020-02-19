### Example chain for kubernetes

![](./kubernetes_app.png)


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

Specifying a kubernetes app is equivalent to:

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
