package phases

import (
	"github.com/moshloop/konfigadm/pkg/os"
	"github.com/moshloop/konfigadm/pkg/types"
)

func GetTags(os os.OS) []types.Flag {
	tags := []types.Flag{}
	for _, tag := range os.GetTags() {
		tags = append(tags, *types.GetTag(tag))
	}
	return tags
}
