package apps_test

import (
	"testing"

	_ "github.com/flanksource/konfigadm/pkg"
	"github.com/flanksource/konfigadm/pkg/types"
	. "github.com/flanksource/konfigadm/pkg/types"
	"github.com/onsi/gomega"
)

func TestNTPInstallation(t *testing.T) {
	t.Run("should install ntp when ntp config is present", func(t *testing.T) {
		cfg, g := NewFixture("ntp.yml", t).WithFlags(types.DEBIAN_LIKE, types.UBUNTU).Build()

		expectedPackage := Package{
			Name:  "ntp",
			Flags: []Flag{DEBIAN_LIKE},
		}
		expectedFileContent := `server 0.in.pool.ntp.org
server 1.in.pool.ntp.org
server 2.in.pool.ntp.org
`
		expectedSystemCtlCommand := Command{
			Cmd:   "systemctl restart ntp",
			Flags: []Flag{DEBIAN_LIKE},
		}
		expectedMoveConfigCommand := Command{
			Cmd:   "mv /tmp/konfigadm-ntp.conf /etc/ntp.conf",
			Flags: []Flag{DEBIAN_LIKE},
		}

		fileSystem, commands, err := cfg.ApplyPhases()
		g.Expect(err).To(gomega.BeNil())

		g.Expect(*cfg.Packages).To(gomega.ContainElement(expectedPackage))

		fileDetails, ok := fileSystem["/tmp/konfigadm-ntp.conf"]
		g.Expect(ok).To(gomega.BeTrue())
		g.Expect(fileDetails.Content).To(gomega.Equal(expectedFileContent))

		commandsLength := len(commands)
		g.Expect(commandsLength).Should(gomega.BeNumerically(">=", 2))
		g.Expect(commands[commandsLength-1]).To(gomega.Equal(expectedSystemCtlCommand))
		g.Expect(commands[commandsLength-2]).To(gomega.Equal(expectedMoveConfigCommand))
	})

	t.Run("should not install ntp when ntp config is not present", func(t *testing.T) {
		cfg, g := NewFixture("no-ntp.yml", t).WithFlags(types.DEBIAN_LIKE, types.UBUNTU).Build()

		expectedCommands := []Command{
			{
				Cmd: "echo 'dummy command'",
			},
		}

		fileSystem, commands, err := cfg.ApplyPhases()
		g.Expect(err).To(gomega.BeNil())

		g.Expect(*cfg.Packages).To(gomega.BeEmpty())
		g.Expect(fileSystem).To(gomega.BeEmpty())
		g.Expect(commands).To(gomega.Equal(expectedCommands))
	})
}
