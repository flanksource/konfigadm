package phases_test

import (
	"testing"
	_ "github.com/moshloop/konfigadm/pkg"
	. "github.com/moshloop/konfigadm/pkg/types"
	"github.com/onsi/gomega"
)

func TestContainers(t *testing.T) {
	cfg, g := NewFixture("containers.yml", t).Build()
	files, commands, _ := cfg.ApplyPhases()
	g.Expect(commands).To(MatchCommand("systemctl start consul"))
	g.Expect(commands).To(MatchCommand("systemctl enable consul"))
	g.Expect(files).To(gomega.HaveKey("/etc/environment.consul"))
	g.Expect(files).To(gomega.HaveKey("/etc/systemd/system/consul.service"))
}
