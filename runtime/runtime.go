package runtime

import "github.com/kujtimiihoxha/plis/config"

type RTime interface {
	Run(data []byte)
	Initialize(c config.GeneratorConfig, s string)
	Stop()
}
