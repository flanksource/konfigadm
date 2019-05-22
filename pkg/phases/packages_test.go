package phases_test

import (
	"testing"

	. "github.com/moshloop/konfigadm/pkg/types"
)

func TestPackageRuntimeFlag(t *testing.T) {

}

func TestPackageDebian(t *testing.T) {
	cfg, g := NewFixture("packages.yml", t).WithFlags(DEBIAN).Build()
	g.Expect(cfg).To(ContainPackage("netcat-openbsd"))
	g.Expect(cfg).NotTo(ContainPackage("nmap-netcat"))

}

func TestPackageRedhat(t *testing.T) {
	cfg, g := NewFixture("packages.yml", t).WithFlags(REDHAT).Build()
	g.Expect(cfg).NotTo(ContainPackage("netcat-openbsd"))
	g.Expect(cfg).To(ContainPackage("nmap-netcat"))
}

func TestPackageUninstall(t *testing.T) {

}

func TestPackageInstall(t *testing.T) {
	cfg, g := NewFixture("packages.yml", t).WithFlags(DEBIAN).Build()
	_, commands, _ := cfg.ApplyPhases()
	g.Expect(commands).To(MatchCommand("apt-get install -y"))
	g.Expect(commands).To(MatchCommand("socat"))
	g.Expect(commands).To(MatchCommand("netcat"))
	g.Expect(commands).NotTo(MatchCommand("yum"))
}

func TestPackageMark(t *testing.T) {

}
