package build

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/mitchellh/colorstring"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/flanksource/commons/files"
	"github.com/flanksource/konfigadm/pkg/types"
	"github.com/flanksource/konfigadm/pkg/utils"
)

type Qemu struct{}

func (q Qemu) Build(image string, config *types.Config) {
	var scratch Scratch
	if config.Context.CaptureLogs != "" {
		log.Infof("Using scratch directory / disk")
		scratch = NewScratch()
	}

	iso, err := createIso(config)
	if err != nil {
		log.Fatalf("Failed to build ISO %v", err)
	}
	if iso == "" {
		log.Fatalf("Empty ISO created")
	}
	cmdLine := qemuSystem(image, iso)
	if config.Context.CaptureLogs != "" {
		cmdLine += fmt.Sprintf(" -hdb %s", scratch.GetImg())
	}

	log.Infof("Executing %s\n", colorstring.Color("[light_green]"+cmdLine))
	if err := utils.Exec(cmdLine); err != nil {
		log.Fatalf("Failed to run: %s, %s", cmdLine, err)
	}
	if config.Context.CaptureLogs != "" {
		log.Infof("Coping captured logs to %s\n", config.Context.CaptureLogs)
		if err = scratch.UnwrapToDir(config.Context.CaptureLogs); err != nil {
			log.Fatalf("Failed to Unwrap: %s", err)
		}
	}
}

func (q Qemu) Test(image string, config *types.Config, privateKeyFile string, template string) error {
	var scratch Scratch
	if config.Context.CaptureLogs != "" {
		log.Infof("Using scratch directory / disk")
		scratch = NewScratch()
	}

	iso, err := createTestIso(config)
	if err != nil {
		return errors.Wrap(err, "failed to build ISO")
	}
	if iso == "" {
		return errors.New("Empty ISO created")
	}

	tempfile, err := ioutil.TempFile("", "test-image")
	if err != nil {
		return errors.Wrap(err, "failed to create temporary file for test image")
	}
	if err := files.Copy(image, tempfile.Name()); err != nil {
		return errors.Wrapf(err, "failed to copy %s to %s", image, tempfile.Name())
	}

	defer os.Remove(tempfile.Name())

	cmdLine := qemuSystem(tempfile.Name(), iso)
	if config.Context.CaptureLogs != "" {
		cmdLine += fmt.Sprintf(" -hdb %s", scratch.GetImg())
	}

	go func() {
		log.Infof("Executing %s\n", colorstring.Color("[light_green]"+cmdLine))
		if err := utils.ExecNoOutput(cmdLine); err != nil {
			log.Errorf("failed to run %s: %v", cmdLine, err)
			return
		}
		if config.Context.CaptureLogs != "" {
			log.Infof("Coping captured logs to %s\n", config.Context.CaptureLogs)
			if err = scratch.UnwrapToDir(config.Context.CaptureLogs); err != nil {
				log.Errorf("Failed to Unwrap to dit: %s", err)
			}
		}
	}()

	// Wait for SSH to be available
	timeout := time.Now().Add(5 * time.Minute)
	sshHost := "127.0.0.1:2022"
	sshUser := "root"
	fmt.Println("Waiting for ssh")
	for {
		if _, err := utils.RunSSHCommand(sshHost, sshUser, privateKeyFile, "true"); err != nil {
			log.Infof("failed to connect to ssh %s@%s: %v, retrying in 5 seconds", sshUser, sshHost, err)
			time.Sleep(5 * time.Second)
		} else {
			log.Infof("connection to %s@%s succeeded", sshUser, sshHost)
			break
		}
		if time.Now().After(timeout) {
			return errors.Errorf("failed to connect to ssh %s@%s after 5 minutes", sshUser, sshHost)
		}
	}

	defer func() {
		if _, err := utils.RunSSHCommand(sshHost, sshUser, privateKeyFile, "shutdown -h -t 5"); err != nil {
			log.Fatalf("failed to shutdown qemu VM: %v", err)
		}
	}()

	scriptTemplate := `
wget https://github.com/aelsabbahy/goss/releases/download/v0.3.13/goss-linux-amd64 -O /usr/bin/goss
chmod +x /usr/bin/goss
cat <<EOF > goss.yaml
%s
EOF
`
	script := fmt.Sprintf(scriptTemplate, template)

	output, err := utils.RunSSHScript(sshHost, sshUser, privateKeyFile, script)
	if output != nil {
		log.Infof("Output:\n%s", output)
	}
	if err != nil {
		return errors.Wrap(err, "failed to run script")
	}

	output, err = utils.RunSSHCommand(sshHost, sshUser, privateKeyFile, "goss validate")
	if output != nil {
		log.Infof("Output:\n%s", output)
	}
	if err != nil {
		return errors.Wrap(err, "failed to run goss validate")
	}

	return nil
}

func qemuSystem(disk, iso string) string {
	return fmt.Sprintf(`qemu-system-x86_64 \
		-nodefaults \
		-display none \
		-machine accel=kvm:hvf \
		-cpu host -smp cpus=2 \
		-m 1024 \
		-hda %s \
		-cdrom %s \
		-device virtio-serial-pci \
		-serial stdio \
		-net nic -net user,hostfwd=tcp:127.0.0.1:2022-:22`, disk, iso)
}
