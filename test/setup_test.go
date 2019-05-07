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
)

var binary = "dist/configadm"

// uses a sensible default on windows (tcp/http) and linux/osx (socket)
var docker, _ = dockertest.NewPool("")

type Container struct {
	id        string
	container *dockertest.Resource
}

func (c *Container) Delete() {
	c.container.Close()
}

func (c *Container) Exec(args ...string) (string, error) {
	arg := []string{"exec", c.id}
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
	env := os.Environ()
	cwd, _ := os.Getwd()
	cwd = filepath.Dir(cwd)
	binary = cwd + "/" + binary
	volumes := []string{
		fmt.Sprintf("%s:%s", cwd, cwd),
		"/sys/fs/cgroup:/sys/fs/cgroup",
	}

	opts := dockertest.RunOptions{
		Privileged: true,
		Env:        env,
		Repository: "moshloop/docker-ubuntu1804-ansible",
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
	g.Expect(stdout).To(ContainSubstring("built"))
	g.Expect(strings.Split(stdout, "\n")).To(HaveLen(2))
}
