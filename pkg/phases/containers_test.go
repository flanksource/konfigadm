package phases_test

import (
	_ "github.com/flanksource/konfigadm/pkg"
	. "github.com/flanksource/konfigadm/pkg/types"
	"github.com/onsi/gomega"
	"testing"
)

func TestContainers(t *testing.T) {
	cfg, g := NewFixture("containers.yml", t).Build()
	files, commands, _ := cfg.ApplyPhases()
	g.Expect(commands).To(MatchCommand("systemctl start consul"))
	g.Expect(commands).To(MatchCommand("systemctl enable consul"))
	g.Expect(files).To(gomega.HaveKey("/etc/environment.consul"))
	g.Expect(files).To(gomega.HaveKey("/etc/systemd/system/consul.service"))
}
