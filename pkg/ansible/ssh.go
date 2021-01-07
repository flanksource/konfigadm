package ansible

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/flanksource/konfigadm/pkg/utils"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/knownhosts"
)

func CheckKnownHosts() ssh.HostKeyCallback {
	CreateKnownHosts()
	knownHosts, err := knownhosts.New(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil {
		log.Fatalf("Failed to check known hosts: %s", err)
	}
	return knownHosts
}

func CreateKnownHosts() {
	f, err := os.OpenFile(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"), os.O_CREATE, 0600)
	if err != nil {
		log.Fatalf("Failed to create known hosts file: %s", err)
	}
	f.Close()
}

func AddHostKey(host string, remote net.Addr, pubKey ssh.PublicKey) error {
	knownHosts := filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")
	f, err := os.OpenFile(knownHosts, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	fmt.Println(host)
	newHost := knownhosts.Normalize(host)
	newHost = knownhosts.HashHostname(newHost)
	_, fErr := f.WriteString(knownhosts.Line([]string{newHost}, pubKey) + "\n")
	return fErr
}

func VerifyHostKey(host string, remote net.Addr, pubKey ssh.PublicKey) error {
	knownHosts := CheckKnownHosts()
	err := knownHosts(host, remote, pubKey)
	var hostErr *knownhosts.KeyError
	if err != nil {
		hostErr = err.(*knownhosts.KeyError)
	}
	if hostErr != nil && len(hostErr.Want) > 0 { // There are existing keys for this host
		log.Warnf("key is not valid for %s, either a MiTM attack or %s has reconfigured the host pub key.", host, host)
		return hostErr
	} else if hostErr != nil && len(hostErr.Want) == 0 { // New host
		log.Warnf("%s is not trusted, adding key to known_hosts file", host)
		result := AddHostKey(host, remote, pubKey)
		return result
	}
	log.Debugf("Public key exists for %s", host)
	return nil
}

func SSHAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	key, err := ioutil.ReadFile(filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa"))
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}
	return ssh.PublicKeys(signer)
}

func ExecuteRemoteCommand(cmd string, host string) {
	log.Infof("Executing %s on %s\n", utils.LightGreenf(cmd), host)
	port := "22"
	user := "root"

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			SSHAgent(),
		},
		HostKeyCallback: VerifyHostKey,
		Timeout:         5 * time.Second,
	}

	client, err := ssh.Dial("tcp", host+":"+port, config)
	if err != nil {
		log.Fatalf("Failed to run on %s: %s", host, err)
	}
	defer client.Close()
	sess, err := client.NewSession()
	if err != nil {
		log.Fatalf("Failed to run on %s: %s", host, err)
	}
	defer sess.Close()

	// setup standard out and error
	// uses writer interface
	sess.Stdout = os.Stdout
	sess.Stderr = os.Stderr

	// run single command
	err = sess.Run(cmd)
	if err != nil {
		log.Fatalf("Failed to run on %s: %s", host, err)
	}
}
