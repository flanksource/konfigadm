package cloudinit

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/flanksource/konfigadm/pkg/utils"
)

//CreateISO creates a new ISO with the user/meta data and returns a path to the iso
func CreateISO(hostname string, userData string) (string, error) {
	dir, err := ioutil.TempDir("", "cloudinit")
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd) //nolint: errcheck
	if err := os.Chdir(dir); err != nil {
		return "", fmt.Errorf("Failed to chdir %v", err)
	}
	if err != nil {
		return "", fmt.Errorf("Failed to create temp dir %s", err)
	}
	// userData = base64.StdEncoding.EncodeToString([]byte(userData))
	if err := ioutil.WriteFile(path.Join(dir, "user-data"), []byte(userData), 0644); err != nil {
		return "", fmt.Errorf("Failed to save user-data %s", err)
	}

	isoFilename, err := ioutil.TempFile("", "user-data*.iso")
	if err != nil {
		return "", fmt.Errorf("Failed to create temp iso %s", err)
	}

	metadata := fmt.Sprintf("instance-id: \nlocal-hostname: %s", hostname)
	if err := ioutil.WriteFile(path.Join(dir, "meta-data"), []byte(metadata), 0644); err != nil {
		return "", fmt.Errorf("Failed to write metadata %v", err)
	}

	var out string
	var ok bool
	if which("genisoimage") {
		out, ok = utils.SafeExec("genisoimage -output %s -volid cidata -joliet -rock user-data meta-data 2>&1", isoFilename.Name())
	} else if which("mkisofs") {
		out, ok = utils.SafeExec("mkisofs -output %s -volid cidata -joliet -rock user-data meta-data 2>&1", isoFilename.Name())
	} else {
		return "", fmt.Errorf("genisoimage or mkisofs not found")
	}
	if !ok && strings.Trim(out, " \n") != "" {
		return "", fmt.Errorf("Failed to create ISO %s", out)
	}
	info, _ := isoFilename.Stat()
	if info.Size() == 0 {
		return "", fmt.Errorf("Empty iso created")
	}
	return isoFilename.Name(), nil
}

func which(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
