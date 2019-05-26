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

