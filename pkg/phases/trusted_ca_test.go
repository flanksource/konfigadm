package phases_test

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/flanksource/konfigadm/pkg/types"
	"github.com/onsi/gomega"
)

func TestTrustedCA(t *testing.T) {
	cfg, g := NewFixture("trusted_ca.yml", t).Build()
	wd, _ := os.Getwd()
	err := os.Chdir("../..")
	g.Expect(err).ToNot(gomega.HaveOccurred())
	defer os.Chdir(wd) // nolint: errcheck
	files, commands, err := cfg.ApplyPhases()
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(files).To(gomega.HaveLen(4))

	g.Expect(files).To(gomega.HaveKey("/tmp/install_certs"))
	data, _ := ioutil.ReadFile("fixtures/files/example-k8s-ca.pem")

	file0 := files["/tmp/konfigadm-trusted-0.pem"]
	g.Expect(file0.Content).To(gomega.Equal(string(data)))
	file1 := files["/tmp/konfigadm-trusted-1.pem"]
	g.Expect(file1.Content).To(gomega.Equal(""))
	g.Expect(file1.ContentFromURL).To(gomega.Equal("https://dl.cacerts.digicert.com/CybertrustGlobalRoot.crt.pem"))
	file2 := files["/tmp/konfigadm-trusted-2.pem"]
	g.Expect(file2.Content).To(gomega.Equal(string(cfg.TrustedCA[2])))

	g.Expect(commands).To(gomega.HaveLen(5))
	g.Expect(commands[0].Cmd).To(gomega.Equal("/tmp/install_certs /tmp/konfigadm-trusted-0.pem"))
	g.Expect(commands[1].Cmd).To(gomega.Equal("/tmp/install_certs /tmp/konfigadm-trusted-1.pem"))
	g.Expect(commands[2].Cmd).To(gomega.Equal("/tmp/install_certs /tmp/konfigadm-trusted-2.pem"))
	g.Expect(commands[3].Cmd).To(gomega.Equal("rm -r /tmp/konfigadm-trusted-*.pem"))
	g.Expect(commands[4].Cmd).To(gomega.Equal("rm -r /tmp/install_certs"))
}
