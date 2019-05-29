package phases

import (
	"github.com/moshloop/konfigadm/pkg/os"
	"github.com/moshloop/konfigadm/pkg/types"
)

func GetTags(os os.OS) []types.Flag {
	tags := []types.Flag{}
	for _, name := range os.GetTags() {
		tag := types.GetTag(name)
		if tag == nil {
			panic("Unknown tag: " + name)
		}
		tags = append(tags, *tag)
	}
	return tags
}
