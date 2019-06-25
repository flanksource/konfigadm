package cloudinit

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/moshloop/konfigadm/pkg/utils"
)

//CreateISO creates a new ISO with the user/meta data and returns a path to the iso
func CreateISO(hostname string, userData string) (string, error) {

	dir, err := ioutil.TempDir("", "cloudinit")
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

	if which("genisoimage") {
		utils.Exec("genisoimage -output %s -volid cidata -joliet -rock user-data meta-data", isoFilename.Name())
	} else if which("mkisofs") {
		utils.Exec("mkisofs -output %s -volid cidata -joliet -rock user-data meta-data", isoFilename.Name())
	} else {
		return "", fmt.Errorf("genisoimage or mkisofs not found")
	}
	return isoFilename.Name(), nil
}

func which(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
