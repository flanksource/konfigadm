package phases_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	. "github.com/flanksource/konfigadm/pkg/types"
	"github.com/onsi/gomega"
)

func TestTrustedCA(t *testing.T) {
	cfg, g := NewFixture("trusted_ca.yml", t).Build()
	files, commands, _ := cfg.ApplyPhases()
	g.Expect(files).To(gomega.HaveLen(4))

	g.Expect(files).To(gomega.HaveKey("/tmp/trustedCA/install_certs"))
	data, _ := ioutil.ReadFile("../../fixtures/trusted_ca.yml")

	file0 := files["/tmp/trustedCA/konfigadm-trusted-0.pem"]
	g.Expect(file0.Content).To(gomega.Equal(string(data)))
	file1 := files["/tmp/trustedCA/konfigadm-trusted-1.pem"]
	g.Expect(file1.Content).To(gomega.Equal(""))
	g.Expect(file1.ContentFromURL).To(gomega.Equal("https://certs.example.com/cert"))
	file2 := files["/tmp/trustedCA/konfigadm-trusted-2.pem"]
	fmt.Printf("File2\n%v\n", file2.Content)
	fmt.Printf("Equal\n%v\n", "-----BEGIN CERTIFICATE-----\nThird cert\n-----END CERTIFICATE-----")
	fmt.Printf("aaa\n")
	g.Expect(file2.Content).To(gomega.Equal("-----BEGIN CERTIFICATE-----\nThird cert\n-----END CERTIFICATE-----\n"))

	g.Expect(commands).To(gomega.HaveLen(4))
	g.Expect(commands[0].Cmd).To(gomega.Equal("/tmp/trustedCA/install_certs /tmp/trustedCA/konfigadm-trusted-0.pem"))
	g.Expect(commands[1].Cmd).To(gomega.Equal("/tmp/trustedCA/install_certs /tmp/trustedCA/konfigadm-trusted-1.pem"))
	g.Expect(commands[2].Cmd).To(gomega.Equal("/tmp/trustedCA/install_certs /tmp/trustedCA/konfigadm-trusted-2.pem"))
	g.Expect(commands[3].Cmd).To(gomega.Equal("rm -r /tmp/trustedCA/"))
}
