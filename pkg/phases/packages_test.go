package phases_test

import (
	"fmt"
	"testing"

	"github.com/flanksource/konfigadm/pkg/phases"
	. "github.com/flanksource/konfigadm/pkg/types"
	"github.com/onsi/gomega"
)

func init() {
	// log.SetLevel(log.TraceLevel)
}

func TestPackageRuntimeFlag(t *testing.T) {

}

func TestFlagsToOS(t *testing.T) {
	g := gomega.NewWithT(t)
	fixtures := map[string]string{
		"ubuntu": "ubuntu",
		"debian": "debian",
	}

	for tag, os := range fixtures {
		matchedTag := GetTag(tag)
		g.Expect(matchedTag.String()).To(gomega.Equal(tag))
		matchedOs, _ := phases.GetOSForTag(*GetTag(tag))
		g.Expect(fmt.Sprintf("%s", matchedOs)).To(gomega.Equal(os))
	}

}

func TestPackageDebian(t *testing.T) {
	cfg, g := NewFixture("packages.yml", t).WithFlags(DEBIAN, DEBIAN_LIKE).Build()
	g.Expect(cfg).To(ContainPackage("netcat-openbsd"))
	g.Expect(cfg).NotTo(ContainPackage("nmap-netcat"))
	g.Expect(cfg).NotTo(ContainPackage("nano"))

}

func TestPackagePhoton(t *testing.T) {
	cfg, g := NewFixture("packages.yml", t).WithFlags(PHOTON).Build()
	g.Expect(cfg).To(ContainPackage("lvm2"))
	g.Expect(cfg).To(ContainPackage("netcat"))
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
