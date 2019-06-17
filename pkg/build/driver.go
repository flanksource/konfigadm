package build

import "os"

type Driver interface {
	Build(image string, config *os.File)
}
