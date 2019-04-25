package cloudmeta

import (
	"testing"

	"github.com/onsi/gomega"
)

func TestContainers(t *testing.T) {
	cfg, g := SetupFixture("containers.yml", t)
	g.Expect(cfg.PreCommands).To(gomega.ContainElement("systemctl start consul"))
	g.Expect(cfg.PreCommands).To(gomega.ContainElement("systemctl enable consul"))
	g.Expect(cfg.Files).To(gomega.HaveKey("/etc/environment.consul"))
	g.Expect(cfg.Files).To(gomega.HaveKey("/etc/systemd/system/consul.service"))
}
