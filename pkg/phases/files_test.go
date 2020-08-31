package phases_test

import (
	"io/ioutil"
	"net/http"
	"testing"

	. "github.com/flanksource/konfigadm/pkg/types"
	"github.com/onsi/gomega"
)

func TestCopy(t *testing.T) {
	cfg, g := NewFixture("files.yml", t).Build()
	fs, _, _ := cfg.ApplyPhases()
	data, _ := ioutil.ReadFile("../../fixtures/files.yml")
	g.Expect(fs).To(gomega.HaveKey("/etc/test"))
	g.Expect(fs["/etc/test"].Content).To(gomega.Equal(string(data)))
}
func TestCopyUrl(t *testing.T) {
	cfg, g := NewFixture("files.yml", t).Build()
	fs, _, _ := cfg.ApplyPhases()
	g.Expect(fs).To(gomega.HaveKey("/etc/testurl"))
	resp, err := http.Get("https://www.geotrust.com/resources/root_certificates/certificates/GeoTrust_Primary_CA.pem")
	g.Expect(err).Should(gomega.BeNil())
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	g.Expect(fs["/etc/testurl"].Content).To(gomega.Equal(string(data)))
}
func TestCopyDir(t *testing.T) {
	cfg, g := NewFixture("files.yml", t).Build()
	fs, _, _ := cfg.ApplyPhases()
	data, _ := ioutil.ReadFile("../../fixtures/files/en.yml")
	g.Expect(fs).To(gomega.HaveKey("/etc/testdir/tscope-master/config/locales/en.yml"))
	g.Expect(fs["/etc/testdir/tscope-master/config/locales/en.yml"].Content).To(gomega.Equal(string(data)))
}

func TestLookup(t *testing.T) {
}

func TestLookupFromContext(t *testing.T) {
}
