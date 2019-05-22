package apps

import . "github.com/moshloop/konfigadm/pkg/types"

var Kubernetes Phase = kubernetes{}

type kubernetes struct{}

func (k kubernetes) ApplyPhase(sys *Config, ctx *SystemContext) ([]Command, Filesystem, error) {
	if sys.Kubernetes == nil {
		return []Command{}, Filesystem{}, nil
	}

	sys.
		AppendPackageRepo(PackageRepo{
			URL:             "https://apt.kubernetes.io/",
			VersionCodeName: "kubernetes-xenial",
			GPGKey:          "https://packages.cloud.google.com/apt/doc/apt-key.gpg",
		}, &DEBIAN)
		// .
		// AppendPackageRepo(PackageRepo{
		// 	URL:    "https://packages.cloud.google.com/yum/repos/kubernetes-el7-x86_64",
		// 	GPGKey: "https://packages.cloud.google.com/yum/doc/yum-key.gpg https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg",
		// }, &REDHAT)

	sys.AddPackage("kubelet=="+sys.Kubernetes.Version+"-00", nil).
		AddPackage("kubeadm=="+sys.Kubernetes.Version+"-00", nil).
		AddPackage("kubectl=="+sys.Kubernetes.Version+"-00", nil)

	sys.Environment["KUBECONFIG"] = "/etc/kubernetes/admin.conf"
	// sys.Sysctls["vm.swappiness"] = "0"
	return []Command{}, Filesystem{}, nil

}
