package apps

import (
	"fmt"
	"strings"

	"github.com/moshloop/konfigadm/pkg/phases"
	. "github.com/moshloop/konfigadm/pkg/types"
	"github.com/moshloop/konfigadm/pkg/utils"
)

var CRI Phase = cri{}

type cri struct{}

func (c cri) ApplyPhase(sys *Config, ctx *SystemContext) ([]Command, Filesystem, error) {
	if sys.ContainerRuntime == nil {
		return []Command{}, Filesystem{}, nil
	}

	if sys.ContainerRuntime.Type == "docker" {
		return c.Docker(sys, ctx)
	} else if sys.ContainerRuntime.Type == "containerd" {
		return c.Containerd(sys, ctx)
	} else {
		return []Command{}, Filesystem{}, fmt.Errorf("unknown container runtime %s", sys.ContainerRuntime.Type)
	}
	return []Command{}, Filesystem{}, nil
}

func (c cri) Verify(sys *Config, results *VerifyResults, flags ...Flag) bool {
	verify := true
	if sys.ContainerRuntime == nil {
		return true
	}

	if sys.ContainerRuntime.Type == "docker" {
		verify = verify && phases.VerifyService("docker", results)
		out, ok := utils.SafeExec("docker ps")
		if ok {
			results.Pass("docker ps returned %d containers", len(strings.Split(out, "\n"))-2)
		} else {
			verify = false
			results.Fail("docker ps failed with: %s", out)
		}

	} else if sys.ContainerRuntime.Type == "containerd" {
		verify = verify && phases.VerifyService("containerd", results)
		out, ok := utils.SafeExec("ctr c list")
		if ok {
			results.Pass("ctr c list returned %d containers", len(strings.Split(out, "\n"))-2)
		} else {
			verify = false
			results.Fail("ctr c list failed with: %s", out)
		}
	} else {
		results.Fail("Unknown runtime %s", sys.ContainerRuntime.Type)
		return false
	}
	return verify
}

func addDockerRepos(sys *Config) {
	sys.AppendPackageRepo(PackageRepo{
		Name:    "docker-ce",
		URL:     "https://download.docker.com/linux/ubuntu/",
		GPGKey:  "https://download.docker.com/linux/ubuntu/gpg",
		Channel: "stable",
	}, UBUNTU)

	sys.AppendPackageRepo(PackageRepo{
		Name:    "docker-ce",
		URL:     "https://download.docker.com/linux/debian",
		GPGKey:  "https://download.docker.com/linux/debian/gpg",
		Channel: "stable",
	}, DEBIAN)

	sys.AppendPackageRepo(PackageRepo{
		Name:   "docker-ce",
		URL:    "https://download.docker.com/linux/centos/7/x86_64/stable",
		GPGKey: "https://download.docker.com/linux/centos/gpg",
	}, REDHAT_LIKE)

	sys.AppendPackageRepo(PackageRepo{
		Name:   "docker-ce",
		URL:    "https://download.docker.com/linux/fedora/\\$releasever/\\$basearch/nightly",
		GPGKey: "https://download.docker.com/linux/fedora/gpg",
	}, FEDORA)

}

func (c cri) Containerd(sys *Config, ctx *SystemContext) ([]Command, Filesystem, error) {
	addDockerRepos(sys)
	sys.AddPackage("containerd.io device-mapper-persistent-data lvm2 libseccomp", &REDHAT_LIKE)
	sys.AddPackage("containerd.io device-mapper-persistent-data lvm2 libseccomp", &FEDORA)
	sys.AddPackage("containerd.io", &DEBIAN_LIKE)
	sys.AddCommand("mkdir -p /etc/containerd && containerd config default > /etc/containerd/config.toml")
	sys.AddCommand("systemctl enable containerd && systemctl restart containerd")
	sys.Environment["CONTAINER_RUNTIME_ENDPOINT"] = "unix:///var/run/containerd/containerd.sock"
	for _, image := range sys.ContainerRuntime.Images {
		sys.AddCommand(fmt.Sprintf("crictl pull %s", image))
	}
	return []Command{}, Filesystem{}, nil
}

func (c cri) Docker(sys *Config, ctx *SystemContext) ([]Command, Filesystem, error) {
	addDockerRepos(sys)
	version := "19.03.2"
	if sys.ContainerRuntime.Version != "" {
		version = sys.ContainerRuntime.Version
	}

	id := "$(. /etc/os-release && echo $ID)"
	codename := "$(. /etc/os-release && echo $VERSION_CODENAME)"

	if strings.Contains(version, "18.06") || strings.Contains(version, "18.03")  {
		if !strings.Contains(version, "~") {
			// Expand the logical docker version into a full version string
			version = fmt.Sprintf("%s~ce~3-0~%s", version, id)
		}
	} else {
			if !strings.Contains(version, "~") {
				// Expand the logical docker version into a full version string
				// docker versions 18.09+ use a new version syntax
				version = fmt.Sprintf("5:%s~3-0~%s-%s", version, id, codename)
			}
		//docker-ce-cli package is required from 18.09
		sys.AppendPackages(nil,		Package{Name: fmt.Sprintf("docker-ce-cli=%s", version), Mark: true})
	}
	sys.AppendPackages(nil,		Package{Name: fmt.Sprintf("docker-ce=%s", version), Mark: true})
	sys.AddPackage("device-mapper-persistent-data lvm2", &FEDORA)
	sys.AddPackage("device-mapper-persistent-data lvm2", &REDHAT_LIKE)
	sys.AddCommand("systemctl enable docker && systemctl start docker")

	for _, image := range sys.ContainerRuntime.Images {
		sys.AddCommand(fmt.Sprintf("docker pull %s", image))
	}

	return []Command{}, Filesystem{}, nil
}
