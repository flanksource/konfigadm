package phases_test

import (
	"testing"

	. "github.com/flanksource/konfigadm/pkg/types"
)

func init() {
	// log.SetLevel(log.TraceLevel)
}

func TestPackageRuntimeFlag(t *testing.T) {

}

func TestPackageDebian(t *testing.T) {
	cfg, g := NewFixture("packages.yml", t).WithFlags(DEBIAN, DEBIAN_LIKE).Build()
	g.Expect(cfg).To(ContainPackage("netcat-openbsd"))
	g.Expect(cfg).NotTo(ContainPackage("nmap-netcat"))
	g.Expect(cfg).NotTo(ContainPackage("nano"))

}

func TestPackageUbuntu(t *testing.T) {
	cfg, g := NewFixture("packages.yml", t).WithFlags(UBUNTU, DEBIAN_LIKE).Build()
	g.Expect(cfg).To(ContainPackage("netcat-openbsd"))
	g.Expect(cfg).NotTo(ContainPackage("nmap-netcat"))
	g.Expect(cfg).To(ContainPackage("nano"))

}

func TestPackageRedhat(t *testing.T) {
	cfg, g := NewFixture("packages.yml", t).WithFlags(REDHAT, REDHAT_LIKE).Build()
	g.Expect(cfg).NotTo(ContainPackage("netcat-openbsd"))
	g.Expect(cfg).To(ContainPackage("nmap-ncat"))
}

func TestPackageUninstall(t *testing.T) {

}

func TestPackageInstall(t *testing.T) {
	cfg, g := NewFixture("packages.yml", t).WithFlags(DEBIAN_LIKE).Build()
	_, commands, _ := cfg.ApplyPhases()
	g.Expect(commands).To(MatchCommand("apt-get install -y"))
	g.Expect(commands).To(MatchCommand("socat"))
	g.Expect(commands).To(MatchCommand("netcat"))
	g.Expect(commands).NotTo(MatchCommand("yum"))
}

func TestPackageMark(t *testing.T) {

}
