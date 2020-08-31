package apps

import (
	"fmt"
	"strings"

	"github.com/flanksource/commons/logger"
	"github.com/flanksource/commons/net"
	"github.com/flanksource/konfigadm/pkg/phases"
	"github.com/flanksource/konfigadm/pkg/types"
	"github.com/flanksource/konfigadm/pkg/utils"
	"github.com/flanksource/konfigadm/resources"
)

var CRI types.Phase = cri{}

type cri struct{}

func (c cri) ApplyPhase(sys *types.Config, ctx *types.SystemContext) ([]types.Command, types.Filesystem, error) {
	if sys.ContainerRuntime.Type == "" {
		return []types.Command{}, types.Filesystem{}, nil
	}

	if sys.ContainerRuntime.Type == "docker" {
		return c.Docker(sys, ctx)
	} else if sys.ContainerRuntime.Type == "containerd" {
		return c.Containerd(sys, ctx)
	} else {
		return []types.Command{}, types.Filesystem{}, fmt.Errorf("unknown container runtime %s", sys.ContainerRuntime.Type)
	}
}

func (c cri) Verify(sys *types.Config, results *types.VerifyResults, flags ...types.Flag) bool {
	verify := true
	if sys.ContainerRuntime.Type == "" {
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

func addDockerRepos(sys *types.Config) {
	sys.AppendPackageRepo(types.PackageRepo{
		Name:    "docker-ce",
		URL:     "https://download.docker.com/linux/ubuntu/",
		GPGKey:  "https://download.docker.com/linux/ubuntu/gpg",
		Channel: "stable",
	}, types.UBUNTU)

	sys.AppendPackageRepo(types.PackageRepo{
		Name:    "docker-ce",
		URL:     "https://download.docker.com/linux/debian",
		GPGKey:  "https://download.docker.com/linux/debian/gpg",
		Channel: "stable",
	}, types.DEBIAN)

	sys.AppendPackageRepo(types.PackageRepo{
		Name:   "docker-ce",
		URL:    "https://download.docker.com/linux/centos/7/x86_64/stable",
		GPGKey: "https://download.docker.com/linux/centos/gpg",
	}, types.RHEL, types.CENTOS)

	sys.AppendPackageRepo(types.PackageRepo{
		Name:   "docker-ce",
		URL:    "https://download.docker.com/linux/fedora/\\$releasever/\\$basearch/nightly",
		GPGKey: "https://download.docker.com/linux/fedora/gpg",
	}, types.FEDORA)

	sys.AppendPackageRepo(types.PackageRepo{
		Name:   "amz2extra-docker",
		GPGKey: "file:///etc/pki/rpm-gpg/RPM-GPG-KEY-amazon-linux-2",
		ExtraArgs: map[string]string{
			"mirrorlist": "http://amazonlinux.\\$awsregion.\\$awsdomain/\\$releasever/extras/docker/latest/\\$basearch/mirror.list",
		},
	}, types.AMAZON_LINUX)

}

var containerdChecksumCache = map[string]string{
	"1.3.4": "61e65c9589e5abfded1daa353e6dfb4b8c2436199bbc5507fc45809a3bb80c1d  containerd-1.3.4.linux-amd64.tar.gz",
	"1.3.3": "b76d54ca86b69871266c29d0f1ad56f37892ab4879b82d34909ab94918b83d16  containerd-1.3.3.linux-amd64.tar.gz",
}

func (c cri) Containerd(sys *types.Config, ctx *types.SystemContext) ([]types.Command, types.Filesystem, error) {
	fs := types.Filesystem{}
	fs["/etc/systemd/system/containerd.service"] = types.File{Content: resources.ContainerdService}
	sys.AddPackage("device-mapper-persistent-data lvm2 libseccomp", &types.REDHAT_LIKE)
	sys.AddPackage("libseccomp2", &types.DEBIAN_LIKE)
	sys.AddPackage("runc", &types.DEBIAN_LIKE)
	if sys.ContainerRuntime.Version == "" {
		sys.ContainerRuntime.Version = "1.3.3"
	}
	sys.ContainerRuntime.Version = strings.TrimPrefix(sys.ContainerRuntime.Version, "v")
	checksum, checksumFound := containerdChecksumCache[sys.ContainerRuntime.Version]
	if !checksumFound {
		checksumB, err := net.GET(fmt.Sprintf("https://github.com/containerd/containerd/releases/download/v%s/containerd-%s.linux-amd64.tar.gz.sha256sum", sys.ContainerRuntime.Version, sys.ContainerRuntime.Version))
		if err != nil {
			logger.Warnf("Failed to get checksum for containerd: %s", err)
		} else {
			checksum = strings.TrimSpace(string(checksumB))
		}

	}
	sys.AddTarPackage(types.TarPackage{
		URL:          fmt.Sprintf("https://github.com/containerd/containerd/releases/download/v%s/containerd-%s.linux-amd64.tar.gz", sys.ContainerRuntime.Version, sys.ContainerRuntime.Version),
		Checksum:     checksum,
		ChecksumType: "sha256",
		Destination:  "/usr/local",
	})
	sys.AddCommand("mkdir -p /etc/containerd && containerd config default > /etc/containerd/config.toml")
	sys.AddCommand("systemctl enable containerd || true && systemctl restart containerd ")
	sys.Environment["CONTAINER_RUNTIME_ENDPOINT"] = "unix:///var/run/containerd/containerd.sock"
	sys.AddCommand("export CONTAINER_RUNTIME_ENDPOINT=unix:///var/run/containerd/containerd.sock")
	for _, image := range sys.ContainerRuntime.Images {
		sys.AddCommand(fmt.Sprintf("crictl pull %s", image))
	}
	return []types.Command{}, fs, nil
}

var versioner map[string]func(version string) string = make(map[string]func(version string) string)

func init() {
	fn := func(version string) string { return version }
	versioner[types.FEDORA.String()] = fn
	versioner[types.REDHAT_LIKE.String()] = fn
	versioner[types.UBUNTU.String()] = func(version string) string {
		if strings.Contains(version, "~") {
			return version
		}
		id := "$(. /etc/os-release && echo $ID)"
		codename := "$(lsb_release -cs)"
		if strings.Contains(version, "18.06") || strings.Contains(version, "18.03") {
			return fmt.Sprintf("%s~ce~3-0~%s", version, id)
		}
		// docker versions 18.09+ use a new version syntax
		return fmt.Sprintf("5:%s~3-0~%s-%s", version, id, codename)
	}
	versioner[types.DEBIAN.String()] = versioner[types.UBUNTU.String()]
}

func (c cri) Docker(sys *types.Config, ctx *types.SystemContext) ([]types.Command, types.Filesystem, error) {
	addDockerRepos(sys)
	version := "19.03.12"
	if sys.ContainerRuntime.Version != "" {
		version = sys.ContainerRuntime.Version
	}

	for tag, fn := range versioner {
		if !strings.Contains(version, "18.06") && !strings.Contains(version, "18.03") {
			//docker-ce-cli package is required from 18.09
			sys.AppendPackages(nil, types.Package{
				Name:    "docker-ce-cli",
				Version: fn(version),
				Mark:    true,
				Flags:   []types.Flag{*types.GetTag(tag), types.NOT_AMAZON_LINUX},
			})
		}

		sys.AppendPackages(nil, types.Package{
			Name:    "docker-ce",
			Version: fn(version),
			Mark:    true,
			Flags:   []types.Flag{*types.GetTag(tag), types.NOT_AMAZON_LINUX},
		})
	}
	fs := types.Filesystem{}
	// Amazon Linux has a non-standard mechanism for installing with limited support
	// for specific docker versions
	sys.AppendPackages(nil, types.Package{
		Name:    "docker",
		Version: "",
		Flags:   []types.Flag{types.AMAZON_LINUX},
	})
// just to trigger pipeline
	sys.AddPackage("device-mapper-persistent-data lvm2", &types.FEDORA)
	sys.AddPackage("device-mapper-persistent-data lvm2", &types.REDHAT_LIKE)
	sys.AddCommand("systemctl enable docker || true && systemctl start docker || true")

	for _, image := range sys.ContainerRuntime.Images {
		sys.AddCommand(fmt.Sprintf("docker pull %s", image))
	}

	return []types.Command{}, fs, nil
}
