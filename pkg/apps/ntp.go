package apps

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/flanksource/konfigadm/pkg/phases"
	. "github.com/flanksource/konfigadm/pkg/types"
)

var ntpPhase ntp = ntp{}

// NTP represents a phase which gives information
// on how to install NTP on the system
var NTP Phase = ntpPhase

type ntp struct{}

func (n ntp) ApplyPhase(sys *Config, ctx *SystemContext) ([]Command, Filesystem, error) {
	commands := []Command{}
	files := Filesystem{}

	if sys.NTP == nil {
		return commands, files, nil
	}

	sys.AddPackage("ntp", &DEBIAN_LIKE)

	configFileContent := createNTPConfigFileContent(sys.NTP)
	files["/tmp/konfigadm-ntp.conf"] = File{Content: configFileContent}
	sys.AddPostCommand("mv /tmp/konfigadm-ntp.conf /etc/ntp.conf", &DEBIAN_LIKE)
	sys.AddPostCommand("systemctl restart ntp", &DEBIAN_LIKE)

	return commands, files, nil
}

func createNTPConfigFileContent(ntpServers []string) string {
	var configFileContent strings.Builder
	for _, ntpServer := range ntpServers {
		configFileContent.WriteString(fmt.Sprintf("server %s\n", ntpServer))
	}
	return configFileContent.String()
}

func (n ntp) Verify(sys *Config, results *VerifyResults, flags ...Flag) bool {
	if sys.NTP == nil {
		return true
	}

	isNtpServiceRunning := phases.VerifyService("ntp", results)

	if !isNtpServiceRunning {
		return false
	}

	ntpConfigFilePath := "/etc/ntp.conf"
	expectedConfigFileContent := createNTPConfigFileContent(sys.NTP)
	actualConfigFileContentBytes, err := ioutil.ReadFile(ntpConfigFilePath)
	if err != nil {
		results.Fail("error while reading %s: %s", ntpConfigFilePath, err.Error())
		return false
	}

	if expectedConfigFileContent != string(actualConfigFileContentBytes) {
		results.Fail("ntp config file %s does not contain the expected content", ntpConfigFilePath)
		return false
	}

	results.Pass("ntp config file %s contains the expected content", ntpConfigFilePath)

	return true
}
