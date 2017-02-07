package runtime

import (
	"github.com/kujtimiihoxha/plis/config"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/kujtimiihoxha/plis/logger"
)

type RTime interface {
	Run() error
	Initialize(cmd *cobra.Command, args map[string]string, c config.GeneratorConfig, gFs afero.Fs) RTime
}

func AddRuntime(cmd *cobra.Command, c config.GeneratorConfig, r RTime, gFs afero.Fs) {
	var rt RTime
	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		if !validateArgs(c.Args, args) {
			logger.GetLogger().Fatal("Please add all the reqired arguments")
		}
		if cmd.Name() == viper.GetString("plis.generator_project_name") {
			viper.Set("plis.is_generator_project", true)
		}
		rt = r.Initialize(cmd, createFlagMap(args, c.Args), c, gFs)
	}
	cmd.RunE = func(cd *cobra.Command, args []string) error {
		return rt.Run()
	}
}
func createFlagMap(args []string, cnfArgs []config.GeneratorArgs) (m map[string]string) {
	i := 0
	m = make(map[string]string)
	for _,v :=range cnfArgs {
		if v.Required {
			m[v.Name] = args[i]
			i++
		}
	}
	for _,v :=range cnfArgs {
		if !v.Required {
			if i <= len(args) - 1 {
				m[v.Name] = args[i]
				i++
			}
		}
	}
	return
}
func validateArgs(cnfArgs []config.GeneratorArgs, args []string) bool {
	i:=0
	for _,v := range cnfArgs {
		if v.Required {
			i++
		}
	}
	if len(args) >= i {
		return true
	}
	return false
}