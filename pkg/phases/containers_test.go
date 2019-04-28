package phases

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestContainers(t *testing.T) {
	cfg, g := NewFixture("containers.yml", t).Build()
	g.Expect(cfg).To(MatchCommand("systemctl start consul"))
	g.Expect(cfg).To(MatchCommand("systemctl enable consul"))
	g.Expect(cfg.Files).To(HaveKey("/etc/environment.consul"))
	g.Expect(cfg.Files).To(HaveKey("/etc/systemd/system/consul.service"))
}
