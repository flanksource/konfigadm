package phases

import (
	"github.com/moshloop/konfigadm/pkg/os"
	"github.com/moshloop/konfigadm/pkg/types"
)

func GetTags(os os.OS) []types.Flag {
	return []types.Flag{*types.GetTag(os.GetTag())}
}
