package phases

import (
	"testing"

	"github.com/onsi/gomega"
)

func TestPackageRuntimeFlag(t *testing.T) {

}

func TestPackageDebian(t *testing.T) {
	cfg, g := NewFixture("packages.yml", t).WithFlags(DEBIAN).Build()
	g.Expect(cfg).To(ContainPackage("netcat"))
	g.Expect(cfg).NotTo(ContainPackage("nmap-netcat"))

}

func TestPackageRedhat(t *testing.T) {
	cfg, g := NewFixture("packages.yml", t).WithFlags(REDHAT).Build()
	g.Expect(cfg).NotTo(ContainPackage("netcat"))
	g.Expect(cfg).To(ContainPackage("nmap-netcat"))
}

func TestPackageUninstall(t *testing.T) {

}

func TestPackageInstall(t *testing.T) {
	cfg, g := NewFixture("packages.yml", t).WithFlags(DEBIAN).Build()
	g.Expect(cfg).To(ContainPackage("docker-ce"))
	aptGets := cfg.FindCmd("apt-get")
	g.Expect(aptGets).To(gomega.HaveCap(1))
	g.Expect(aptGets[0].Cmd).To(gomega.ContainSubstring("docker-ce"))
	g.Expect(aptGets[0].Cmd).To(gomega.ContainSubstring("socat"))
	g.Expect(aptGets[0].Cmd).To(gomega.ContainSubstring("netcat"))

}

func TestPackageMark(t *testing.T) {

}
