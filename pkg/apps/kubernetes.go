package apps

import (
	"strings"

	"github.com/flanksource/konfigadm/pkg/types"
)

var Kubernetes types.Phase = kubernetes{}

type kubernetes struct{}

func (k kubernetes) ApplyPhase(sys *types.Config, ctx *types.SystemContext) ([]types.Command, types.Filesystem, error) {
	if sys.Kubernetes == nil {
		return []types.Command{}, types.Filesystem{}, nil
	}
	sys.
		AppendPackageRepo(types.PackageRepo{
			Name:            "kubernetes",
			URL:             "https://apt.kubernetes.io/",
			VersionCodeName: "kubernetes-xenial",
			GPGKey:          "https://packages.cloud.google.com/apt/doc/apt-key.gpg",
		}, types.DEBIAN_LIKE).
		AppendPackageRepo(types.PackageRepo{
			Name:   "kubernetes",
			URL:    "https://packages.cloud.google.com/yum/repos/kubernetes-el7-x86_64",
			GPGKey: "https://packages.cloud.google.com/yum/doc/yum-key.gpg https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg",
		}, types.REDHAT_LIKE).
		AppendPackageRepo(types.PackageRepo{
			Name:   "kubernetes",
			URL:    "https://packages.cloud.google.com/yum/repos/kubernetes-el7-x86_64",
			GPGKey: "https://packages.cloud.google.com/yum/doc/yum-key.gpg https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg",
		}, types.FEDORA)

	sys.AppendPackages(&types.REDHAT_LIKE,
		types.Package{Name: "kubelet", Version: withDefaultPatch(sys.Kubernetes.Version, "0"), Mark: true},
		types.Package{Name: "kubeadm", Version: withDefaultPatch(sys.Kubernetes.Version, "0"), Mark: true},
		types.Package{Name: "kubectl", Version: withDefaultPatch(sys.Kubernetes.Version, "0"), Mark: true})

	sys.AppendPackages(&types.FEDORA,
		types.Package{Name: "kubelet", Version: withDefaultPatch(sys.Kubernetes.Version, "0"), Mark: true},
		types.Package{Name: "kubeadm", Version: withDefaultPatch(sys.Kubernetes.Version, "0"), Mark: true},
		types.Package{Name: "kubectl", Version: withDefaultPatch(sys.Kubernetes.Version, "0"), Mark: true})

	sys.AppendPackages(&types.DEBIAN_LIKE,
		types.Package{Name: "kubelet", Version: withDefaultPatch(sys.Kubernetes.Version, "00"), Mark: true},
		types.Package{Name: "kubeadm", Version: withDefaultPatch(sys.Kubernetes.Version, "00"), Mark: true},
		types.Package{Name: "kubectl", Version: withDefaultPatch(sys.Kubernetes.Version, "00"), Mark: true})

	sys.AddPackage("socat ebtables ntp libseccomp nfs-utils", &types.REDHAT_LIKE)
	sys.AddPackage("socat ebtables ntp libseccomp2 nfs-common", &types.DEBIAN_LIKE)

	sys.Environment["KUBECONFIG"] = "/etc/kubernetes/admin.conf"
	sys.Sysctls["vm.swappiness"] = "0"
	sys.Sysctls["net.bridge.bridge-nf-call-iptables"] = "1"
	sys.Sysctls["net.bridge.bridge-nf-call-ip6tables"] = "1"
	sys.Sysctls["net.ipv4.ip_forward"] = "1"

	fs := types.Filesystem{}
	fs["/etc/modules-load.d/kubernetes.conf"] = types.File{Content: "overlay\nbr_netfilter"}
	return []types.Command{}, fs, nil
}

func withDefaultPatch(version, patch string) string {
	if strings.Contains(version, "-") {
		return version
	}
	return version + "-" + patch

}
