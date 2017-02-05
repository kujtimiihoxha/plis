package runtime

import (
	"github.com/kujtimiihoxha/plis/config"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type RTime interface {
	Run() error
	Initialize(c config.GeneratorConfig, gFs afero.Fs) RTime
}

func AddRuntime(cmd *cobra.Command, c config.GeneratorConfig, r RTime, gFs afero.Fs) {
	rt := r.Initialize(c, gFs)
	cmd.RunE = func(cd *cobra.Command, args []string) error {
		if cd.Name() == viper.GetString("plis.generator_project_name") {
			viper.Set("plis.is_generator_project",true)
		}
		return rt.Run()
	}
}
