package cloudmeta

import (
	"testing"

	"github.com/onsi/gomega"
)

func TestSysctl(t *testing.T) {
	cfg, g := SetupFixture("sysctl.yml", t)
	g.Expect(cfg.Files).To(gomega.HaveLen(1))
	g.Expect(cfg.PreCommands).To(gomega.HaveLen(2))
}
