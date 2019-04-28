package phases

import (
	"io/ioutil"
	"testing"

	"github.com/onsi/gomega"
)

func TestCopy(t *testing.T) {
	cfg, g := NewFixture("files.yml", t).Build()
	data, _ := ioutil.ReadFile("../../fixtures/files.yml")
	g.Expect(cfg.Files).To(gomega.HaveKey("/etc/test"))
	g.Expect(cfg.Files["/etc/test"]).To(gomega.Equal(string(data)))
}

func TestLookup(t *testing.T) {
}

func TestLookupFromContext(t *testing.T) {
}
