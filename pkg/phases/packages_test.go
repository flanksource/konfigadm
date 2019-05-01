package phases_test

import (
	"testing"

	. "github.com/moshloop/configadm/pkg/types"

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
	_, commands, _ := cfg.ApplyPhases()
	g.Expect(commands).To(gomega.ContainSubstring("apt-get install -y"))
	g.Expect(commands).To(gomega.ContainSubstring("docker-ce"))
	g.Expect(commands).To(gomega.ContainSubstring("socat"))
	g.Expect(commands).To(gomega.ContainSubstring("netcat"))
	g.Expect(commands).NotTo(gomega.ContainSubstring("yum"))
}

func TestPackageMark(t *testing.T) {

}
