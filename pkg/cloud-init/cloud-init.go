package cloudinit

import (
	"gopkg.in/yaml.v2"
)

func (init CloudInit) String() string {
	data, err := yaml.Marshal(init)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func (init *CloudInit) AddCommand(cmd string) *CloudInit {
	init.Runcmd = append(init.Runcmd, []string{"once-per-instance", cmd})
	return init
}

func (init *CloudInit) AddFile(path string, contents string) *CloudInit {
	init.WriteFiles = append(init.WriteFiles, File{
		Path:    path,
		Content: contents,
	})
	return init
}
