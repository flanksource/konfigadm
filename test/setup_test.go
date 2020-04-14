package test

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/flanksource/yaml"
	"github.com/onsi/gomega"
	. "github.com/onsi/gomega"
	"github.com/ory/dockertest"
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
	cmd.Wait()
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

func check(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
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
	{"packages.yml"},
	{"trusted_ca.yml"},
	// {"sysctl.yml"},
}

type konfigadm struct {
}

func (c konfigadm) Verify(config ...string) bool {
	return false
}

func (c konfigadm) Apply(config ...string) bool {
	return false
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
