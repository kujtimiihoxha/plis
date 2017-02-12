package generators

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/asaskevich/govalidator"
	"github.com/kujtimiihoxha/plis/cmd"
	"github.com/kujtimiihoxha/plis/config"
	"github.com/kujtimiihoxha/plis/fs"
	"github.com/kujtimiihoxha/plis/helpers"
	"github.com/kujtimiihoxha/plis/logger"
	"github.com/kujtimiihoxha/plis/runtime"
	"github.com/kujtimiihoxha/plis/runtime/lua"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func find() (globalGenerators []string, projectGenerators []string) {
	dirs, err := afero.ReadDir(fs.GetPlisRootFs(), "generators")
	if err != nil {
		logger.GetLogger().Fatal(err)
	}
	for _, f := range dirs {
		if !strings.HasPrefix(f.Name(), ".") && strings.HasPrefix(f.Name(), "plis-") {
			globalGenerators = append(globalGenerators, strings.TrimPrefix(f.Name(), "plis-"))
		}
	}
	dirs, err = afero.ReadDir(fs.GetCurrentFs(), "plis"+afero.FilePathSeparator+"generators")
	if err != nil {
		if os.IsNotExist(err) {
			logger.GetLogger().Info("No project generators.")
		} else {
			logger.GetLogger().Fatal(err)
		}
	}
	for _, f := range dirs {
		if !strings.HasPrefix(f.Name(), ".") && strings.HasPrefix(f.Name(), "plis-") {
			projectGenerators = append(projectGenerators, strings.TrimPrefix(f.Name(), "plis-"))
		}
	}
	return
}
func Initialize() {
	globalGenerators, projectGenerators := find()
	for _, v := range globalGenerators {
		gFs := afero.NewBasePathFs(
			fs.GetPlisRootFs(),
			fmt.Sprintf("generators%splis-%s", afero.FilePathSeparator, v),
		)
		viper.Set(
			fmt.Sprintf("plis.generators.%s.root", v),
			fmt.Sprintf(
				"%s%sgenerators%splis-%s",
				viper.GetString("plis.dir.root"),
				afero.FilePathSeparator,
				afero.FilePathSeparator,
				v,
			),
		)
		createGeneratorCmd(gFs, cmd.RootCmd, v)
	}
	for _, v := range projectGenerators {
		gFs := afero.NewBasePathFs(fs.GetCurrentFs(), fmt.Sprintf("plis%sgenerators%splis-%s", afero.FilePathSeparator, afero.FilePathSeparator, v))
		dr, _ := os.Getwd()
		viper.Set(
			fmt.Sprintf("plis.generators.%s.root", v),
			fmt.Sprintf(
				"%s%splis%sgenerators%splis-%s",
				dr,
				afero.FilePathSeparator,
				afero.FilePathSeparator,
				afero.FilePathSeparator,
				v,
			),
		)
		createGeneratorCmd(gFs, cmd.RootCmd, v)
	}
	checkIfGeneratorProject()
}
func createGeneratorCmd(fs afero.Fs, cmd *cobra.Command, generator string) {
	d, err := afero.ReadFile(fs, "config.json")
	if err != nil {
		logger.GetLogger().Errorf(
			"Generator `%s` has no config file and it will be ignored",
			generator,
		)
		return
	}
	c := config.GeneratorConfig{}
	err = json.Unmarshal(d, &c)
	if err != nil {
		logger.GetLogger().Errorf(
			"Could not read the config file of `%s` generator, this generator will be ignored",
			generator,
		)
		return
	}
	if c.Name == "" {
		c.Name = generator
	}
	addCmd(cmd, c, fs)
}
func addCmd(cmd *cobra.Command, c config.GeneratorConfig, gFs afero.Fs) {
	logger.GetLogger().Infof("Validating generator `%s`...", c.Name)
	if c.Validate() {
		logger.GetLogger().Info("Validation Ok")
		logger.GetLogger().Info("Validating flags...")
		flagsToKeep := []config.GeneratorFlag{}
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
		argsToKeep := []config.GeneratorArgs{}
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
		createGeneratorCmd(_gFs, newC, v)
	}
	// update viper base.
}
func createCommand(c config.GeneratorConfig, gFs afero.Fs) *cobra.Command {
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
	default:
		runtime.AddRuntime(genCmd, c, lua.NewLuaRuntime(gFs))
	}
	return genCmd
}
func addFlags(command *cobra.Command, c config.GeneratorConfig) {
	for _, v := range c.Flags {
		if v.Persistent {
			switch v.Type {
			case "string":
				command.PersistentFlags().StringP(v.Name, v.Short, v.Default.(string), v.Description)
			case "int":
				f := v.Default.(float64)
				iv := int(f)
				command.PersistentFlags().IntP(v.Name, v.Short, iv, v.Description)
			case "float":
				f := v.Default.(float64)
				command.PersistentFlags().Float64P(v.Name, v.Short, f, v.Description)
			case "bool":
				b := v.Default.(bool)
				command.PersistentFlags().BoolP(v.Name, v.Short, b, v.Description)
			}
		} else {
			switch v.Type {
			case "string":
				command.Flags().StringP(v.Name, v.Short, v.Default.(string), v.Description)
			case "int":
				f := v.Default.(float64)
				iv := int(f)
				command.Flags().IntP(v.Name, v.Short, iv, v.Description)
			case "float":
				f := v.Default.(float64)
				command.Flags().Float64P(v.Name, v.Short, f, v.Description)
			case "bool":
				b := v.Default.(bool)
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
func checkIfGeneratorProject() {
	logger.SetLevel(logrus.InfoLevel)
	d, err := afero.ReadFile(fs.GetCurrentFs(), ".plis-generator.json")
	if err != nil {
		return
	}
	c := config.GeneratorProjectConfig{}
	err = json.Unmarshal(d, &c)
	if err != nil {
		logger.GetLogger().Error(
			"Could not read the generator project config file",
			err,
		)
		return
	}
	v, _ := govalidator.ValidateStruct(c)
	if !v {
		logger.GetLogger().Error(
			"Could not calidate the generator project config file, make sure you specified all the required fields",
		)
		return
	}
	currentFs := fs.GetCurrentFs()
	fs.SetGeneratorTestFs(afero.NewBasePathFs(currentFs, c.TestDir))
	viper.Set("plis.generator_project_name", c.GeneratorName)
	dr, _ := os.Getwd()
	viper.Set(
		fmt.Sprintf("plis.generators.%s.root", c.GeneratorName),
		fmt.Sprintf(
			dr,
		),
	)
	createGeneratorCmd(currentFs, cmd.RootCmd, c.GeneratorName)

}
