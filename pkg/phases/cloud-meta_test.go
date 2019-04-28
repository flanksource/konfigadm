package phases

import (
	"fmt"
	"testing"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

func SetupFixtureWithArgs(name string, vars []string, t *testing.T) (*SystemConfig, *gomega.WithT) {
	cfg := NewSystemConfigFromFixture(name, vars, t)
	g := gomega.NewWithT(t)
	return cfg, g
}
func SetupFixture(name string, t *testing.T) (*SystemConfig, *gomega.WithT) {
	cfg := NewSystemConfigFromFixture(name, []string{}, t)
	g := gomega.NewWithT(t)
	return cfg, g
}

func NewSystemConfigFromFixture(name string, vars []string, t *testing.T) *SystemConfig {
	cfg, err := NewSystemConfig(
		vars,
		[]string{"../../fixtures/" + name},
	)

	if err != nil {
		t.Error(err)
	}
	return cfg
}

func MatchCommand(expected interface{}) types.GomegaMatcher {
	return &CommandMatcher{
		expected: expected,
	}
}

type CommandMatcher struct {
	expected interface{}
}

func (matcher *CommandMatcher) Match(actual interface{}) (success bool, err error) {
	sys, ok := actual.(*SystemConfig)
	if !ok {
		return false, fmt.Errorf("CommandMatcher matcher expects a SystemConfig")
	}
	for _, cmd := range sys.PreCommands {
		if cmd.Cmd == matcher.expected {
			return true, nil
		}
	}
	for _, cmd := range sys.Commands {
		if cmd.Cmd == matcher.expected {
			return true, nil
		}
	}
	return false, nil
}

func (matcher *CommandMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected \t%#v to contain command: \t%#v", actual, matcher.expected)
}

func (matcher *CommandMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected \t%#v to NOT contain command: \t%#v", actual, matcher.expected)
}

func TestMultipleConfigs(t *testing.T) {
}

func TestRemoteConfig(t *testing.T) {
}
