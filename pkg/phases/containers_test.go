package phases_test

import (
	"testing"

	"github.com/onsi/gomega"
)

func TestContainers(t *testing.T) {
	cfg, g := NewFixture("containers.yml", t).Build()
	files, commands, _ := cfg.ApplyPhases()
	g.Expect(commands).To(gomega.ContainSubstring("systemctl start consul"))
	g.Expect(commands).To(gomega.ContainSubstring("systemctl enable consul"))
	g.Expect(files).To(gomega.HaveKey("/etc/environment.consul"))
	g.Expect(files).To(gomega.HaveKey("/etc/systemd/system/consul.service"))
}
