package types

import (
	"testing"

	"github.com/onsi/gomega"
)

type Fixture struct {
	name  string
	vars  []string
	flags []Flag
	t     *testing.T
	g     *gomega.WithT
}

func (f *Fixture) WithVars(vars ...string) *Fixture {
	f.vars = vars
	return f
}

func (f *Fixture) WithFlags(flags ...Flag) *Fixture {
	f.flags = flags
	return f
}

func (f *Fixture) Build() (*Config, *gomega.WithT) {
	cfg, err := NewConfig("../../fixtures/" + f.name).
		WithFlags(f.flags...).
		WithVars(f.vars...).
		Build()
	if err != nil {
		f.t.Error(err)
	}
	return cfg, f.g
}

func NewFixture(name string, t *testing.T) *Fixture {
	return &Fixture{
		name: name,
		t:    t,
		g:    gomega.NewWithT(t),
	}
}
