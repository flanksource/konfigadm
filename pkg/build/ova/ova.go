package ova

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/flanksource/konfigadm/pkg/utils"
)

var (
	options = `
	{
    "DiskProvisioning": "thin",
    "IPAllocationPolicy": "dhcpPolicy",
    "IPProtocol": "IPv4",
    "NetworkMapping": [
        {
            "Name": "VM Network",
            "Network": "%s"
        }
    ],
    "MarkAsTemplate": false,
    "PowerOn": false,
    "InjectOvfEnv": false,
    "WaitForIP": false,
    "Name": null
}
`
	base = `
.encoding = "UTF-8"
displayName = "%s"
disk.enableUUID = 1
cleanShutdown = "TRUE"
config.version = "8"
cpuid.coresPerSocket = "2"
ethernet0.addressType = "generated"
ethernet0.generatedAddressOffset = "0"
ethernet0.networkName = "VM Network"
ethernet0.connectionType = "bridged"
ethernet0.pciSlotNumber = "160"
ethernet0.present = "TRUE"
ethernet0.uptCompatibility = "TRUE"
ethernet0.virtualDev = "vmxnet3"
ethernet0.wakeOnPcktRcv = "FALSE"

floppy0.autodetect = "TRUE"
floppy0.startConnected = "FALSE"

guestOS = "other3xlinux-64"
hpet0.present = "TRUE"
ide1:0.autodetect = "TRUE"
ide1:0.clientDevice = "FALSE"
ide1:0.deviceType = "atapi-cdrom"
ide1:0.present = "TRUE"
ide1:0.startConnected = "FALSE"
memSize = "2048"
numvcpus = "2"
pciBridge0.present = "TRUE"
pciBridge4.functions = "8"
pciBridge4.present = "TRUE"
pciBridge4.virtualDev = "pcieRootPort"
pciBridge5.functions = "8"
pciBridge5.present = "TRUE"
pciBridge5.virtualDev = "pcieRootPort"
pciBridge6.functions = "8"
pciBridge6.present = "TRUE"
pciBridge6.virtualDev = "pcieRootPort"
pciBridge7.functions = "8"
pciBridge7.present = "TRUE"
pciBridge7.virtualDev = "pcieRootPort"

serial0.present = "TRUE"
serial0.fileType = "file"
serial0.autodetect = "TRUE"
serial0.fileName = "serial.out"
answer.msg.serial.file.open = "Replace"
answer.msg.uuid.altered = "I copied it"

scsi0:0.allowguestconnectioncontrol = "false"
scsi0:0.deviceType = "disk"
scsi0:0.fileName = "%s"
scsi0:0.mode = "persistent"
scsi0:0.present = "TRUE"
scsi0.present = "TRUE"
scsi0.virtualDev = "lsilogic"
svga.present = "TRUE"
svga.vramSize = "134217728"
tools.syncTime = "FALSE"
toolScripts.afterPowerOn = "TRUE"
toolScripts.afterResume = "TRUE"
toolScripts.beforePowerOff = "TRUE"
toolScripts.beforeSuspend = "TRUE"
virtualHW.productCompatibility = "hosted"
virtualhw.version = "11"
vmci0.present = "TRUE"
vmci0.unrestricted = "false"
`
)

func Create(name, image string, properties map[string]string) (string, error) {
	dir, _ := os.Getwd()
	base := utils.GetBaseName(image)
	vmdk := path.Join(dir, base+".vmdk")
	ovf := path.Join(dir, base+".ova")
	vmx := path.Join(dir, base+".vmx")
	if err := ioutil.WriteFile(vmx, []byte(getVmx(name, path.Base(vmdk), properties)), 0644); err != nil {
		return "", err
	}

	if !strings.HasSuffix(image, ".vmdk") {
		if runtime.GOOS == "Darwin" {
			return "", errors.New("qcow to vmdk conversion on MacOSX is broken in qemu, see https://bugs.launchpad.net/qemu/+bug/1776920")
		}
		log.Infof("Converting image to vmdk")
		if err := utils.Exec("qemu-img convert -O vmdk -p %s %s", image, vmdk); err != nil {
			return "", err
		}
	}

	log.Infof("Creating OVA")
	if err := utils.Exec("ovftool %s %s", vmx, ovf); err != nil {
		return "", err
	}
	return ovf, nil
}

func Import(name, ova, network string) error {
	tmp, _ := ioutil.TempFile("", "options*.json")
	if _, err := tmp.WriteString(getOptions(network)); err != nil {
		return err
	}
	if !log.IsLevelEnabled(log.TraceLevel) {
		defer os.Remove(tmp.Name())
	}
	return utils.Exec("govc import.ova --name %s --options %s %s", name, tmp.Name(), ova)
}

func ImportContentLibrary(library, name, ova string) error {
	if name != "" {
		return utils.Exec("govc library.import -n %s %s %s", name, library, ova)
	}
	return utils.Exec("govc library.import %s %s", library, ova)
}

func getVmx(name, image string, properties map[string]string) string {
	vmx := fmt.Sprintf(base, name, image)

	for k, v := range properties {
		vmx += fmt.Sprintf("%s=%s\n", k, v)
	}
	return vmx
}

func getOptions(network string) string {
	return fmt.Sprintf(options, network)
}
