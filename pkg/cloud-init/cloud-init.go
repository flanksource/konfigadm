package cloudinit

import (
	"encoding/base64"

	"gopkg.in/yaml.v3"
)

func (init CloudInit) String() string {
	data, err := yaml.Marshal(init)
	if err != nil {
		panic(err)
	}
	return "#cloud-config\n" + string(data)
}

func (init *CloudInit) AddCommand(cmd string) *CloudInit {
	init.Runcmd = append(init.Runcmd, []string{"sh", cmd})
	return init
}

func (init *CloudInit) AddFile(path string, contents string) *CloudInit {
	if init.FileEncoding == "base64" {
		init.WriteFiles = append(init.WriteFiles, File{
			Path:     path,
			Content:  base64.StdEncoding.EncodeToString([]byte(contents)),
			Encoding: "base64",
		})
	} else {
		init.WriteFiles = append(init.WriteFiles, File{
			Path:    path,
			Content: contents,
		})
	}

	return init
}
