package phases_test

import (
	"testing"

	"github.com/onsi/gomega"
)

func TestSysctl(t *testing.T) {
	cfg, g := NewFixture("sysctl.yml", t).Build()
	files, commands, _ := cfg.ApplyPhases()
	g.Expect(files).To(gomega.HaveLen(1))
	g.Expect(commands).To(gomega.ContainSubstring("sysctl -w net.ipv6.conf.all.disable_ipv6 1"))
}
