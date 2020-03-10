package types_test

import (
	"testing"

	"github.com/flanksource/konfigadm/pkg/types"
	"github.com/onsi/gomega"
)

func TestAddPostCommand(t *testing.T) {
	t.Run("should add a post command to the config", func(t *testing.T) {
		g := gomega.NewWithT(t)
		config := &types.Config{}
		flag := types.DEBIAN
		anotherFlag := types.DEBIAN_LIKE
		expectedPostCommand := types.Command{
			Cmd:   "ls",
			Flags: []types.Flag{flag, anotherFlag},
		}

		modifiedConfig := config.AddPostCommand("ls", &flag, &anotherFlag)

		g.Expect(modifiedConfig).NotTo(gomega.BeNil())
		g.Expect(modifiedConfig.PostCommands).To(gomega.ContainElement(expectedPostCommand))
	})
}
