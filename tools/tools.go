package tools

import (
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/kujtimiihoxha/plis/cmd"
	"github.com/kujtimiihoxha/plis/config"
	"github.com/kujtimiihoxha/plis/fs"
	"github.com/kujtimiihoxha/plis/helpers"
	"github.com/kujtimiihoxha/plis/logger"
	"github.com/kujtimiihoxha/plis/runtime"
	"github.com/kujtimiihoxha/plis/runtime/js"
	"github.com/kujtimiihoxha/plis/runtime/lua"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strconv"
	"strings"
)

func find() (globalTools []string, projectTools []string) {
	dirs, err := afero.ReadDir(fs.GetPlisRootFs(), "tools")
	if err != nil {
		logger.GetLogger().Fatal(err)
	}
	for _, f := range dirs {
		if !strings.HasPrefix(f.Name(), ".") && strings.HasPrefix(f.Name(), "plis-") {
			globalTools = append(globalTools, strings.TrimPrefix(f.Name(), "plis-"))
		}
	}
	dirs, err = afero.ReadDir(fs.GetCurrentFs(), "plis"+afero.FilePathSeparator+"tools")
	if err != nil {
		if os.IsNotExist(err) {
			logger.GetLogger().Info("No project tools.")
		} else {
			logger.GetLogger().Fatal(err)
		}
	}
	for _, f := range dirs {
		if !strings.HasPrefix(f.Name(), ".") && strings.HasPrefix(f.Name(), "plis-") {
			projectTools = append(projectTools, strings.TrimPrefix(f.Name(), "plis-"))
		}
	}
	return
}
func Initialize() {
	globalTools, projectTools := find()
	for _, v := range globalTools {
		gFs := afero.NewBasePathFs(
			fs.GetPlisRootFs(),
			fmt.Sprintf("tools%splis-%s", afero.FilePathSeparator, v),
		)
		viper.Set(
			fmt.Sprintf("plis.tools.%s.root", v),
			fmt.Sprintf(
				"%s%stools%splis-%s",
				viper.GetString("plis.dir.root"),
				afero.FilePathSeparator,
				afero.FilePathSeparator,
				v,
			),
		)
		createToolCmd(gFs, cmd.RootCmd, v)
	}
	for _, v := range projectTools {
		if v == "get" || v == "install" {
			logger.GetLogger().Warnf("The commandd '%s' is a build in command so it can not be used", v)
			continue
		}
		toRemove := []*cobra.Command{}
		for _, gcmd := range cmd.RootCmd.Commands() {
			if gcmd.Name() == v {
				logger.GetLogger().Warnf(
					"The commandd '%s' exists as a global command it will be replaced by the project level command",
					v,
				)
				toRemove = append(toRemove, gcmd)
			}
		}
		cmd.RootCmd.RemoveCommand(toRemove...)
		gFs := afero.NewBasePathFs(fs.GetCurrentFs(), fmt.Sprintf("plis%stools%splis-%s", afero.FilePathSeparator, afero.FilePathSeparator, v))
		dr, _ := os.Getwd()
		viper.Set(
			fmt.Sprintf("plis.tools.%s.root", v),
			fmt.Sprintf(
				"%s%splis%stools%splis-%s",
				dr,
				afero.FilePathSeparator,
				afero.FilePathSeparator,
				afero.FilePathSeparator,
				v,
			),
		)
		createToolCmd(gFs, cmd.RootCmd, v)
	}
	checkIfToolProject()
}
func createToolCmd(fs afero.Fs, cmd *cobra.Command, tool string) {
	d, err := afero.ReadFile(fs, "config.json")
	if err != nil {
		logger.GetLogger().Errorf(
			"Tool `%s` has no config file and it will be ignored",
			tool,
		)
		return
	}
	c := config.ToolConfig{}
	err = json.Unmarshal(d, &c)
	if err != nil {
		logger.GetLogger().Errorf(
			"Could not read the config file of `%s` tool, this tool will be ignored",
			tool,
		)
		return
	}
	if c.Name == "" {
		c.Name = tool
	}
	addCmd(cmd, c, fs)
}
func addCmd(cmd *cobra.Command, c config.ToolConfig, gFs afero.Fs) {
	logger.GetLogger().Infof("Validating tool `%s`...", c.Name)
	if c.Validate() {
		logger.GetLogger().Info("Validation Ok")
		logger.GetLogger().Info("Validating flags...")
		flagsToKeep := []config.ToolFlag{}
		for _, v := range c.Flags {
			if v.Validate() {
				logger.GetLogger().Infof("Flag `%s` OK", v.Name)
				flagsToKeep = append(flagsToKeep, v)
			} else {
				logger.GetLogger().Warn("This flag will be ignored")
			}
		}
		c.Flags = flagsToKeep
		logger.GetLogger().Info("Validating args...")
		argsToKeep := []config.ToolArgs{}
		for _, v := range c.Args {
			if v.Validate() {
				logger.GetLogger().Infof("Argument `%s` OK", v.Name)
				argsToKeep = append(argsToKeep, v)
			} else {
				logger.GetLogger().Warn("This argument will be ignored")
			}
		}
		c.Args = argsToKeep

	} else {
		logger.GetLogger().Warn("This comand will be ignored")
		return
	}
	newC := createCommand(c, gFs)
	cmd.AddCommand(newC)
	for _, v := range c.SubCommands {
		_gFs := afero.NewBasePathFs(gFs, v)
		createToolCmd(_gFs, newC, v)
	}
	// update viper base.
}
func createCommand(c config.ToolConfig, gFs afero.Fs) *cobra.Command {
	genCmd := &cobra.Command{
		Use:     c.Name,
		Short:   c.Description,
		Long:    helpers.FromStringArrayToString(c.LongDescription),
		Aliases: c.Aliases,
	}
	addFlags(genCmd, c)
	genCmd.SetHelpTemplate(getUsageTemplate())
	switch c.ScriptType {
	case "lua":
		runtime.AddRuntime(genCmd, c, lua.NewLuaRuntime(gFs))
	case "js":
		runtime.AddRuntime(genCmd, c, js.NewJsRuntime(gFs))
	default:
		runtime.AddRuntime(genCmd, c, lua.NewLuaRuntime(gFs))
	}
	return genCmd
}
func addFlags(command *cobra.Command, c config.ToolConfig) {
	for _, v := range c.Flags {
		if v.Persistent {
			switch v.Type {
			case "string":
				command.PersistentFlags().StringP(v.Name, v.Short, v.Default, v.Description)
			case "int":
				f, _ := strconv.ParseFloat(v.Default, 32)
				iv := int(f)
				command.PersistentFlags().IntP(v.Name, v.Short, iv, v.Description)
			case "float":
				f, _ := strconv.ParseFloat(v.Default, 32)
				command.PersistentFlags().Float64P(v.Name, v.Short, f, v.Description)
			case "bool":
				b, _ := strconv.ParseBool(v.Default)
				command.PersistentFlags().BoolP(v.Name, v.Short, b, v.Description)
			}
		} else {
			switch v.Type {
			case "string":
				command.Flags().StringP(v.Name, v.Short, v.Default, v.Description)
			case "int":
				f, _ := strconv.ParseFloat(v.Default, 32)
				iv := int(f)
				command.Flags().IntP(v.Name, v.Short, iv, v.Description)
			case "float":
				f, _ := strconv.ParseFloat(v.Default, 32)
				command.Flags().Float64P(v.Name, v.Short, f, v.Description)
			case "bool":
				b, _ := strconv.ParseBool(v.Default)
				command.Flags().BoolP(v.Name, v.Short, b, v.Description)
			}
		}
	}
}
func getUsageTemplate() string {
	return `Usage:{{if .Runnable}}
  {{if .HasAvailableFlags}}{{appendIfNotPresent .UseLine "[flags]"}}{{else}}{{.UseLine}}{{end}}{{end}}{{if .HasAvailableSubCommands}}
  {{ .CommandPath}} [command]{{end}}{{if gt .Aliases 0}}
Aliases:
  {{.NameAndAliases}}
{{end}}{{if .HasExample}}
Examples:
{{ .Example }}{{end}}{{ if .HasAvailableSubCommands}}
Available Commands:{{range .Commands}}{{if .IsAvailableCommand}}
  {{rpad .NameAndAliases .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{ if .HasAvailableLocalFlags}}
Flags:
{{.LocalFlags.FlagUsages | trimRightSpace}}{{end}}{{ if .HasAvailableInheritedFlags}}
Global Flags:
{{.InheritedFlags.FlagUsages | trimRightSpace}}{{end}}{{if .HasHelpSubCommands}}
Additional help topics:{{range .Commands}}{{if .IsHelpCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{ if .HasAvailableSubCommands }}
Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
}
func checkIfToolProject() {
	logger.SetLevel(logrus.InfoLevel)
	d, err := afero.ReadFile(fs.GetCurrentFs(), ".plis-tool.json")
	if err != nil {
		return
	}
	c := config.ToolProjectConfig{}
	err = json.Unmarshal(d, &c)
	if err != nil {
		logger.GetLogger().Error(
			"Could not read the tool project config file",
			err,
		)
		return
	}
	v, _ := govalidator.ValidateStruct(c)
	if !v {
		logger.GetLogger().Error(
			"Could not calidate the tool project config file, make sure you specified all the required fields",
		)
		return
	}
	currentFs := fs.GetCurrentFs()
	fs.SetToolTestFs(afero.NewBasePathFs(currentFs, c.TestDir))
	viper.Set("plis.tool_project_name", c.ToolName)
	dr, _ := os.Getwd()
	viper.Set(
		fmt.Sprintf("plis.tools.%s.root", c.ToolName),
		fmt.Sprintf(
			dr,
		),
	)
	toRemove := []*cobra.Command{}
	for _, gcmd := range cmd.RootCmd.Commands() {
		if gcmd.Name() == c.ToolName {
			logger.GetLogger().Warnf(
				"The commandd '%s' exists as a global command it will be replaced by the project level command",
				c.ToolName,
			)
			toRemove = append(toRemove, gcmd)
		}
	}
	cmd.RootCmd.RemoveCommand(toRemove...)
	createToolCmd(currentFs, cmd.RootCmd, c.ToolName)

}
