package generators

import (
	"encoding/json"
	"fmt"
	"github.com/kujtimiihoxha/plis/cmd"
	"github.com/kujtimiihoxha/plis/config"
	"github.com/kujtimiihoxha/plis/fs"
	"github.com/kujtimiihoxha/plis/helpers"
	"github.com/kujtimiihoxha/plis/logger"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
	"github.com/yuin/gopher-lua"
)

func find() (globalGenerators []string, projectGenerators []string) {
	dirs, err := afero.ReadDir(fs.GetPlisRootFs(), "generators")
	if err != nil {
		logger.GetLogger().Fatal(err)
	}
	for _, f := range dirs {
		if !strings.HasPrefix(f.Name(), ".") {
			globalGenerators = append(globalGenerators, f.Name())
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
		if !strings.HasPrefix(f.Name(), ".") {
			projectGenerators = append(projectGenerators, f.Name())
		}
	}
	return
}

func Initialize() {
	globalGenerators, projectGenerators := find()
	if len(globalGenerators) > 0 {
		createGeneratorCmd(
			fs.GetPlisRootFs(),
			globalGenerators,
			"",
		)
	}
	if len(projectGenerators) > 0 {
		createGeneratorCmd(
			afero.NewBasePathFs(
				fs.GetCurrentFs(),
				afero.FilePathSeparator+"plis"),
			projectGenerators,
			"",
		)
	}
}
func createGeneratorCmd(fs afero.Fs, generators []string, subPath string) {
	if subPath != "" {
		subPath = afero.FilePathSeparator + subPath + afero.FilePathSeparator
	} else {
		subPath = afero.FilePathSeparator
	}
	for _, v := range generators {
		d, err := afero.ReadFile(
			fs,
			"generators"+afero.FilePathSeparator+v+subPath+"config.json",
		)
		if err != nil {
			logger.GetLogger().Errorf(
				"Generator `%s` has no config file and it will be ignored",
				v,
			)
			continue
		}
		c := config.GeneratorConfig{}
		err = json.Unmarshal(d, &c)
		if err != nil {
			logger.GetLogger().Errorf(
				"Could not read the config file of `%s` generator, this generator will be ignored",
				v,
			)
			continue
		}
		if c.Name == "" {
			c.Name = v
		}
		addCmd(cmd.RootCmd, c, cmd.RootCmd.Name())
	}
}
func addCmd(cmd *cobra.Command, c config.GeneratorConfig, viperBase string) {
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
	newC := createCommand(c, viperBase)
	cmd.AddCommand(newC)
	// update viper base.
}
func createCommand(c config.GeneratorConfig, viperBase string) *cobra.Command {
	genCmd := &cobra.Command{
		Use:   c.Name,
		Short: c.Description,
		Long:  helpers.FromStringArrayToString(c.LongDescription),
		// Uncomment the following line if your bare application
		// has an action associated with it:
		Run: func(cmd *cobra.Command, args []string) {
			L := lua.NewState()
			defer L.Close()
			a := L.NewTable()
			L.NewUserData()
			a.RawSet(lua.LString("test"),lua.LNumber(123))
			L.SetGlobal("flags",a)
			b,_:=afero.ReadFile(fs.GetPlisRootFs(),"generators/test/run.lua")
			if err := L.DoString(string(b)); err != nil {
				panic(err)
			}
		},
	}
	addFlags(genCmd, c, viperBase)
	return genCmd
}
func addFlags(command *cobra.Command, c config.GeneratorConfig, viperBase string) {
	for _, v := range c.Flags {
		if v.Persistent {
			switch v.Type {
			case "string":
				command.PersistentFlags().StringP(v.Name, v.Short, v.Default.(string), v.Description)
			case "int":
				f := v.Default.(float64)
				iv := int(f)
				command.PersistentFlags().IntP(v.Name, v.Short,  iv, v.Description)
			case "float":
				f := v.Default.(float64)
				command.PersistentFlags().Float64P(v.Name,  v.Short, f, v.Description)
			case "bool":
				b := v.Default.(bool)
				command.PersistentFlags().BoolP(v.Name,  v.Short, b, v.Description)
			}
			n := fmt.Sprintf("%s.%s.flags.%s", viperBase, c.Name, v.Name)
			cmd.PersistentFlags = append(cmd.PersistentFlags,n)
			viper.BindPFlag(n, command.PersistentFlags().Lookup(v.Name))
		} else {
			switch v.Type {
			case "string":
				command.Flags().StringP(v.Name,  v.Short, v.Default.(string), v.Description)
				n := fmt.Sprintf("%s.%s.flags.%s", viperBase, c.Name, v.Name)
				viper.BindPFlag(n, command.PersistentFlags().Lookup(v.Name))
			case "int":
				f := v.Default.(float64)
				iv := int(f)
				command.Flags().IntP(v.Name, v.Short,  iv, v.Description)
			case "float":
				f := v.Default.(float64)
				command.Flags().Float64P(v.Name, v.Short,  f, v.Description)
			case "bool":
				b := v.Default.(bool)
				command.Flags().BoolP(v.Name, v.Short,  b, v.Description)
			}
			n := fmt.Sprintf("%s.%s.flags.%s", viperBase, c.Name, v.Name)
			viper.BindPFlag(n, command.Flags().Lookup(v.Name))
		}
	}
}
