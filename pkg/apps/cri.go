package apps

import (
	"fmt"
	"strings"

	"github.com/moshloop/konfigadm/pkg/phases"
	"github.com/moshloop/konfigadm/pkg/utils"

	. "github.com/moshloop/konfigadm/pkg/types"
)

var CRI Phase = cri{}

type cri struct{}

func (c cri) ApplyPhase(sys *Config, ctx *SystemContext) ([]Command, Filesystem, error) {
	if sys.ContainerRuntime == nil {
		return []Command{}, Filesystem{}, nil
	}

	if sys.ContainerRuntime.Type == "docker" {
		return c.Docker(sys, ctx)
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

	} else {
		results.Fail("Unknown runtime %s", sys.ContainerRuntime.Type)
		return false
	}
	return verify
}

func (c cri) Docker(sys *Config, ctx *SystemContext) ([]Command, Filesystem, error) {
	sys.AppendPackageRepo(PackageRepo{
		URL:     "https://download.docker.com/linux/ubuntu",
		GPGKey:  "https://download.docker.com/linux/ubuntu/gpg",
		Channel: "stable",
	}, &DEBIAN)

	sys.AddPackage("docker-ce", &DEBIAN)
	sys.AddCommand("systemctl enable docker && systemctl start docker")
	return []Command{}, Filesystem{}, nil
}
