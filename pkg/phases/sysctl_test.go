package phases

import (
	"testing"

	"github.com/onsi/gomega"
)

func TestSysctl(t *testing.T) {
	cfg, g := NewFixture("sysctl.yml", t).Build()
	g.Expect(cfg.Files).To(gomega.HaveLen(1))
	g.Expect(cfg.PreCommands).To(gomega.HaveLen(2))
}
