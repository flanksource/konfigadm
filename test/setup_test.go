package test

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/onsi/gomega"
	. "github.com/onsi/gomega"
	"github.com/ory/dockertest"
	"gopkg.in/flanksource/yaml.v3"
)

// uses a sensible default on windows (tcp/http) and linux/osx (socket)
var docker, _ = dockertest.NewPool("")

//cwdVol contains the id of a volume in which the current working directory has been copied into
// this is required due to -v $PWD:$PWD not working on circleci
var cwdVol string
var cwd string
var env []string
var image string
var binary string

func init() {
	cwd, _ = os.Getwd()
	cwd = filepath.Dir(cwd)
	image = os.Getenv("IMAGE")
	cwdVol = os.Getenv("CWD_VOL")
	binary = cwd + "/.bin/konfigadm"
}

type Container struct {
	id        string
	container *dockertest.Resource
}

func (c *Container) Delete() {
	c.container.Close()
}

func (c *Container) Exec(args ...string) (string, error) {
	arg := []string{"exec", "-w", cwd, c.id}
	arg = append(arg, args...)
	cmd := exec.Command("docker", arg...)
	data, err := cmd.CombinedOutput()
	if err != nil {
		return string(data), err
	}
	cmd.Wait() // nolint: errcheck
	if !cmd.ProcessState.Success() {
		return string(data), errors.New("Failed")
	}
	return string(data), nil
}

//newContainer creates a new systemd based container and returns the container id
func newContainer() (*Container, error) {

	volumes := []string{
		fmt.Sprintf("%s:%s", cwdVol, cwd),
		"/sys/fs/cgroup:/sys/fs/cgroup",
	}

	if image == "" {
		image = "jrei/systemd-ubuntu:18.04"
	}
	opts := dockertest.RunOptions{
		Privileged: true,
		Env:        env,
		Repository: strings.Split(image, ":")[0],
		Tag:        strings.Split(image, ":")[1],
		Entrypoint: []string{"/lib/systemd/systemd"},
		Mounts:     volumes,
	}

	container, err := docker.RunWithOptions(&opts)
	if err != nil {
		return nil, err
	}
	return &Container{
		id:        container.Container.ID,
		container: container,
	}, nil
}

func setup(t *testing.T) (*gomega.WithT, *Container) {
	g := gomega.NewGomegaWithT(t)
	container, err := newContainer()
	if err != nil {
		t.Fatal(err)
	}
	return g, container
}

func TestVersion(t *testing.T) {
	g, container := setup(t)
	defer container.Delete()

	stdout, err := container.Exec(binary, "version")
	fmt.Println(stdout)
	if err != nil {
		t.Fatalf("%s: %s", err, stdout)
	}
	g.Expect(stdout).To(ContainSubstring("v"))
	g.Expect(strings.Split(stdout, "\n")).To(HaveLen(2))
}

var fixtures = []struct {
	in string
}{
	{"services.yml"},
	{"containers.yml"},
	{"docker.yml"},
	{"ansible.yml"},
	{"containerd.yml"},
	{"files.yml"},
	{"kubernetes.yml"},
	{"kernel.yml"},
	{"packages.yml"},
	{"trusted_ca.yml"},
	// {"sysctl.yml"},
}

func TestYamlRoundTrip(t *testing.T) {
	for _, f := range fixtures {
		t.Run(f.in, func(t *testing.T) {
			_, container := setup(t)
			defer container.Delete()
			stdout, err := container.Exec(binary, "minify", "-c", cwd+"/fixtures/"+f.in)
			if err != nil {
				t.Errorf("Minify failed %s:\n %s\n", f.in, stdout)
			}
			var data map[string]interface{}
			if err := yaml.Unmarshal([]byte(stdout), &data); err != nil {
				t.Errorf("Failed to unmarshall: %s\n%s", err, stdout)
			}
		})
	}
}

func TestFull(t *testing.T) {
	for _, f := range fixtures {
		if strings.Contains(image, "fedora") && f.in == "kernel.yml" {
			continue
		}
		t.Run(f.in, func(t *testing.T) {
			g, container := setup(t)
			defer container.Delete()
			stdout, err := container.Exec(binary, "verify", "-c", cwd+"/fixtures/"+f.in)
			if err == nil {
				t.Errorf("Verify should have failed %s:\n %s\n", f.in, stdout)
			}
			os.Stderr.WriteString(stdout)
			stdout, err = container.Exec("cat", "/etc/os-release")
			os.Stderr.WriteString(stdout)
			stdout, err = container.Exec(binary, "apply", "-c", cwd+"/fixtures/"+f.in)
			if err != nil {
				t.Errorf("Apply should succeed %s: %s\n", err, stdout)
			}

			os.Stderr.WriteString(stdout)

			g.Eventually(func() string {
				stdout, err = container.Exec(binary, "verify", "-c", cwd+"/fixtures/"+f.in)
				os.Stderr.WriteString(stdout + "\n")
				return stdout
			}, "30s", "3s").Should(ContainSubstring("0 failed"))

		})
	}
}

func TestSSH(t *testing.T) {
	f := fixtures[0]
	g, targetContainer := setup(t)
	defer targetContainer.Delete()

	sshContainer, err := newContainer()
	if err != nil {
		t.Fatal(err)
	}
	defer sshContainer.Delete()

	stdout, err := targetContainer.Exec(binary, "verify", "-c", cwd+"/fixtures/"+f.in)
	if err == nil {
		t.Errorf("Verify should have failed %s:\n %s\n", f.in, stdout)
	}
	os.Stderr.WriteString(stdout)
	stdout, err = targetContainer.Exec("cat", "/etc/os-release")
	os.Stderr.WriteString(stdout)

	stdout, err = SetupSSHContainers(sshContainer, targetContainer)
	if err != nil {
		t.Errorf("Unable to setup SSH container: %s: %s\n", err, stdout)
	}

	stdout, err = sshContainer.Exec(binary, "apply", "-c", cwd+"/fixtures/"+f.in,
		"--inventory", cwd+"/inventory")
	if err != nil {
		_, _ = sshContainer.Exec("rm", cwd+"/inventory")
		t.Errorf("Apply should succeed %s: %s\n", err, stdout)
	}

	os.Stderr.WriteString(stdout)
	_, _ = sshContainer.Exec("rm", cwd+"/inventory")

	g.Eventually(func() string {
		stdout, err = targetContainer.Exec(binary, "verify", "-c", cwd+"/fixtures/"+f.in)
		os.Stderr.WriteString(stdout + "\n")
		if err != nil {
			t.Fatal(err)
		}
		return stdout
	}, "30s", "3s").Should(ContainSubstring("0 failed"))
}

func SetupSSHContainers(sshContainer *Container, targetContainer *Container) (string, error) {
	stdout, err := sshContainer.Exec("sh", "-c",
		"echo "+targetContainer.container.Container.NetworkSettings.IPAddress+" > "+cwd+"/inventory")
	if err != nil {
		return stdout, err
	}
	stdout, err = sshContainer.Exec("mkdir", "-p", "/root/.ssh")
	if err != nil {
		return stdout, err
	}
	stdout, err = sshContainer.Exec("chmod", "700", "/root/.ssh")
	if err != nil {
		return stdout, err
	}
	stdout, err = targetContainer.Exec("mkdir", "-p", "/root/.ssh")
	if err != nil {
		return stdout, err
	}
	stdout, err = targetContainer.Exec("chmod", "700", "/root/.ssh")
	if err != nil {
		return stdout, err
	}
	// TODO: Need to set up SSH daemon on the target container with SSH keys available on the SSH container
	// SSH keys need adding to SSH agent or ClientConfig needs to support ~/.ssh/id_* directly.
	stdout, err = sshContainer.Exec("ssh-keygen", "-t", "rsa", "-b", "4096", "-N", "", "-f", "/root/.ssh/id_rsa")
	if err != nil {
		return stdout, err
	}
	sshKey, err := sshContainer.Exec("cat", "/root/.ssh/id_rsa.pub")
	if err != nil {
		return sshKey, err
	}
	stdout, err = targetContainer.Exec("sh", "-c",
		"echo \""+sshKey+"\" > /root/.ssh/authorized_keys")
	if err != nil {
		return stdout, err
	}
	stdout, err = targetContainer.Exec("chmod", "600", "/root/.ssh/authorized_keys")
	if err != nil {
		return stdout, err
	}
	return "", nil
}
