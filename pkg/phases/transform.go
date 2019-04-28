package phases

import (
	"fmt"
	"log"
	"reflect"
	"runtime"
)

var (
	transformers = make([]Transformer, 0)
)

func Register(fn Transformer) {
	log.Printf("Registering %+v\n", runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name())
	transformers = append(transformers, fn)
}

func (cfg *SystemConfig) Transform(ctx SystemContext) error {
	for _, t := range transformers {
		c, f, e := t(cfg, &ctx)
		if e != nil {
			return fmt.Errorf("%s", e)
		}

		cfg.PreCommands = append(cfg.PreCommands, c...)
		for k, v := range f {
			cfg.Files[k] = v.Content
		}
	}
	return nil
}
