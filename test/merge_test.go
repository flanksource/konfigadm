package test

import (
	"strings"
	"testing"

	_ "github.com/moshloop/konfigadm/pkg"
	"github.com/moshloop/konfigadm/pkg/types"

	"github.com/onsi/gomega"
)

func init() {

}

type Fixture struct {
	files []string
	vars  []string
	flags []types.Flag
	t     *testing.T
	g     *gomega.WithT
}

func (f *Fixture) WithVars(vars ...string) *Fixture {
	f.vars = vars
	return f
}

func (f *Fixture) WithFlags(flags ...types.Flag) *Fixture {
	f.flags = flags
	return f
}

func (f *Fixture) Build() (*types.Config, *gomega.WithT) {
	cfg, err := types.NewConfig(f.files...).
		WithFlags(f.flags...).
		WithVars(f.vars...).
		Build()
	if err != nil {
		f.t.Error(err)
	}
	return cfg, f.g
}

func NewFixture(t *testing.T, files ...string) *Fixture {
	return &Fixture{
		files: files,
		t:     t,
		g:     gomega.NewWithT(t),
	}
}

var PATH = "../fixtures/"

func TestImportKubernetesIntoPackages(t *testing.T) {
	cfg, g := NewFixture(t, PATH+"packages.yml", PATH+"kubernetes.yml").
		WithFlags(types.DEBIAN).
		Build()
	g.Expect(*cfg.Packages).NotTo(gomega.BeEmpty())
}

func TestImportKubernetesWithCommand(t *testing.T) {
	cfg, g := NewFixture(t, PATH+"kubernetes.yml").
		WithFlags(types.DEBIAN, types.DEBIAN_LIKE).
		Build()
	cfg.Extra.FileEncoding = ""
	count := strings.Count(cfg.ToCloudInit().String(), "echo 1.14.1 > /etc/kubernetes_version")
	g.Expect(count).To(gomega.BeEquivalentTo(1))
	g.Expect(*cfg.Packages).NotTo(gomega.BeEmpty())

}

func TestImportThree(t *testing.T) {
	cfg, g := NewFixture(t, PATH+"packages.yml", PATH+"kubernetes.yml", PATH+"ssh.yml").
		WithFlags(types.DEBIAN, types.DEBIAN_LIKE).
		Build()
	cfg.Extra.FileEncoding = ""
	cloudinit := cfg.ToCloudInit().String()
	g.Expect(cloudinit).To(gomega.ContainSubstring("kubeadm"))
}
