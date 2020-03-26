package apps_test

import (
	"testing"

	_ "github.com/flanksource/konfigadm/pkg"
	. "github.com/flanksource/konfigadm/pkg/types"
	. "github.com/onsi/gomega"
)

func TestMarkPackages(t *testing.T) {
	cfg, g := NewFixture("kubernetes.yml", t).WithFlags(DEBIAN_LIKE, UBUNTU).Build()
	_, commands, _ := cfg.ApplyPhases()
	g.Expect((*cfg.Packages)[0].Mark).To(BeTrue())
	g.Expect(commands).To(MatchCommand("apt-mark hold kubelet kubeadm kubectl"))
}
