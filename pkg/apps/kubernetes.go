package apps

import . "github.com/moshloop/configadm/pkg/types"
import log "github.com/sirupsen/logrus"

var Kubernetes Phase = kubernetes{}

type kubernetes struct{}

func (k kubernetes) ApplyPhase(sys *Config, ctx *SystemContext) ([]Command, Filesystem, error) {
	if sys.Kubernetes == nil {
		return []Command{}, Filesystem{}, nil
	}

	sys.PackageRepos = append(sys.PackageRepos, PackageRepo{
		Flags: []Flag{DEBIAN},
		URL:   "deb https://apt.kubernetes.io/ kubernetes-xenial main",
	})
	sys.PackageRepos = append(sys.PackageRepos, PackageRepo{
		Flags: []Flag{REDHAT},
		URL:   "https://packages.cloud.google.com/yum/repos/kubernetes-el7-x86_64",
	})
	sys.AddPackage("kubelet=="+sys.Kubernetes.Version, nil).
		AddPackage("kubeadm=="+sys.Kubernetes.Version, nil).
		AddPackage("kubectl=="+sys.Kubernetes.Version, nil)

	sys.Environment["KUBECONFIG"] = "/etc/kubernetes/admin.conf"
	log.Infof("%+v\n", sys.Packages)
	sys.Sysctls["vm.swapiness"] = "0"
	return []Command{}, Filesystem{}, nil

}
