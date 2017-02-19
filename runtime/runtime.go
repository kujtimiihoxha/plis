package runtime

import (
	"github.com/kujtimiihoxha/plis/config"
	"github.com/kujtimiihoxha/plis/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type RTime interface {
	Run() error
	Initialize(cmd *cobra.Command, args map[string]string, c config.ToolConfig)
}

func AddRuntime(cmd *cobra.Command, c config.ToolConfig, rt RTime) {
	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		if !validateArgs(c.Args, args) {
			logger.GetLogger().Fatal("Please add all the reqired arguments")
		}
		_cmd := cmd
		for _cmd != nil {
			if _cmd.Name() == viper.GetString("plis.tool_project_name") {
				viper.Set("plis.is_tool_project", true)
			}
			_cmd = _cmd.Parent()
		}
		switch c.ScriptType {

		}
		rt.Initialize(cmd, createArgumentMap(args, c.Args), c)
	}
	cmd.RunE = func(cd *cobra.Command, args []string) error {
		return rt.Run()
	}
}
func createArgumentMap(args []string, cnfArgs []config.ToolArgs) (m map[string]string) {
	i := 0
	m = make(map[string]string)
	for _, v := range cnfArgs {
		if v.Required {
			m[v.Name] = args[i]
			i++
		}
	}
	for _, v := range cnfArgs {
		if !v.Required {
			if i <= len(args)-1 {
				m[v.Name] = args[i]
				i++
			}
		}
	}
	for _,v :=range cnfArgs {
		if m[v.Name] == "" {
			switch v.Type {
			case "int":
				m[v.Name] = "0"
			case "float":
				m[v.Name] = "0.0"
			case "bool":
				m[v.Name] = "false"
			}
		}
	}
	return
}
func validateArgs(cnfArgs []config.ToolArgs, args []string) bool {
	i := 0
	for _, v := range cnfArgs {
		if v.Required {
			i++
		}
	}
	if len(args) >= i {
		return true
	}
	return false
}
