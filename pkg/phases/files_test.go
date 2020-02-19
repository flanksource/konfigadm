package phases_test

import (
	. "github.com/flanksource/konfigadm/pkg/types"
	"github.com/onsi/gomega"
	"io/ioutil"
	"testing"
)

func TestCopy(t *testing.T) {
	cfg, g := NewFixture("files.yml", t).Build()
	fs, _, _ := cfg.ApplyPhases()
	data, _ := ioutil.ReadFile("../../fixtures/files.yml")
	g.Expect(fs).To(gomega.HaveKey("/etc/test"))
	g.Expect(fs["/etc/test"].Content).To(gomega.Equal(string(data)))
}

func TestLookup(t *testing.T) {
}

func TestLookupFromContext(t *testing.T) {
}
