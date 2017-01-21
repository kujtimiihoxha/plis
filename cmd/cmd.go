// Copyright © 2016 Kujtim Hoxha <kujtimii.h@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"github.com/kujtimiihoxha/plis/fs"
	"github.com/kujtimiihoxha/plis/generators"
	"github.com/kujtimiihoxha/plis/helpers"
	"github.com/kujtimiihoxha/plis/scripts"
	"github.com/mattn/anko/vm"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"os"
)

func InitGenerators() {
	if ex, _ := afero.Exists(fs.WorkingDirFs(), helpers.GeneratorsPath()); !ex {
		return
	}
	d, _ := afero.ReadDir(fs.WorkingDirFs(), helpers.GeneratorsPath())
	for _, v := range d {
		CreateGenerator(helpers.RootGeneratorConfig(v.Name()), helpers.RootGeneratorScript(v.Name()), RootGenerator)
	}
}
func CreateGenerator(generatorConfigPath string, generatorScriptPath string, parent *generators.PlisGenerator) {
	config := generators.ReadConfig(fs.WorkingDirFs(), generatorConfigPath)
	cmd := createCmd(config)
	generator := generators.NewPlisGenerator(cmd, config, parent)
	addFlags(generator)
	createRunFunction(generator, generatorScriptPath)
	if generator.Config.SubCommands != nil {
		for _, v := range *generator.Config.SubCommands {
			CreateGenerator(
				helpers.ChildGeneratorConfig(generator.GetRootParent().Config.Name, v),
				helpers.ChildGeneratorScript(generator.GetRootParent().Config.Name, v),
				generator)
		}
	}

	//Todo add children.
	parent.Cmd.AddCommand(generator.Cmd)
}
func createRunFunction(gen *generators.PlisGenerator, generatorScriptPath string) {
	gen.Cmd.Run = func(c *cobra.Command, args []string) {
		gen.ValidateArguments(args)
		gen.ValidateFlags(c)
		var env = vm.NewEnv()
		scripts.Build(env, gen, args)
		if gen.GetRootParent().Config.Modules != nil {
			for _, v := range *gen.GetRootParent().Config.Modules {
				data, err := afero.ReadFile(fs.WorkingDirFs(), helpers.GeneratorModulesFile(gen.GetRootParent().Config.Name, v))
				if err != nil {
					fmt.Println(fmt.Sprintf("Could not read module '%s' ", v))
					os.Exit(-1)
				}
				_, err = env.Execute(string(data))
				if err != nil {
					fmt.Println(fmt.Sprintf("Error while executing module '%s' ", v), err)
					os.Exit(-1)
				}
			}
		}
		data, err := afero.ReadFile(fs.WorkingDirFs(), generatorScriptPath)
		if err != nil {
			fmt.Println("Could not read generator script")
			os.Exit(-1)
		}
		_, err = env.Execute(string(data))
		if err != nil {
			fmt.Println("Error while executing generator script")
			fmt.Println(err)
			os.Exit(-1)
		}
	}
}
func addFlags(gen *generators.PlisGenerator) {
	if gen.Config.Flags == nil {
		return
	}
	for _, v := range *gen.Config.Flags {
		switch v.Type {
		case "string":
			var df interface{}
			if v.Default != nil {
				df = *v.Default
				_, ok := df.(string)
				if !ok {
					fmt.Println(
						fmt.Sprintf(
							`The default value of type 'string' should be string in command '%s' flag '%s'`,
							gen.Config.Name,
							v.Long))
					os.Exit(-1)
				}
			} else {
				df = ""
			}
			if v.Persistent {
				gen.Cmd.PersistentFlags().StringP(v.Long, v.Short, df.(string), v.Description)
			} else {
				gen.Cmd.Flags().StringP(v.Long, v.Short, df.(string), v.Description)
			}
		case "int":
			var df interface{}
			if v.Default != nil {
				df = *v.Default
				_, ok := df.(float64)
				if !ok {
					fmt.Println(
						fmt.Sprintf(
							`The default value of type 'int' should be int in command '%s' flag '%s'`,
							gen.Config.Name,
							v.Long))
					os.Exit(-1)
				}
			} else {
				df = 0.0
			}
			if v.Persistent {
				gen.Cmd.PersistentFlags().Int64P(v.Long, v.Short, int64(df.(float64)), v.Description)
			} else {
				gen.Cmd.Flags().Int64P(v.Long, v.Short, int64(df.(float64)), v.Description)
			}
		case "float":
			var df interface{}
			if v.Default != nil {
				df = *v.Default
				_, ok := df.(float64)
				if !ok {
					fmt.Println(
						fmt.Sprintf(
							`The default value of type 'float' should be float in command '%s' flag '%s'`,
							gen.Config.Name,
							v.Long))
					os.Exit(-1)
				}
			} else {
				df = 0
			}
			if v.Persistent {
				gen.Cmd.PersistentFlags().Float64P(v.Long, v.Short, df.(float64), v.Description)
			} else {
				gen.Cmd.Flags().Float64P(v.Long, v.Short, df.(float64), v.Description)
			}
		case "bool":
			var df interface{}
			if v.Default != nil {
				df = *v.Default
				_, ok := df.(bool)
				if !ok {
					fmt.Println(
						fmt.Sprintf(
							`The default value of type 'bool' should be bool in command '%s' flag '%s'`,
							gen.Config.Name,
							v.Long))
					os.Exit(-1)
				}
			} else {
				df = false
			}
			if v.Persistent {
				gen.Cmd.PersistentFlags().BoolP(v.Long, v.Short, df.(bool), v.Description)
			} else {
				gen.Cmd.Flags().BoolP(v.Long, v.Short, df.(bool), v.Description)
			}
		default:
			if v.Persistent {
				gen.Cmd.PersistentFlags().StringP(v.Long, v.Short, "", v.Description)
			} else {
				gen.Cmd.Flags().StringP(v.Long, v.Short, "", v.Description)
			}
		}
	}
}
func createCmd(config *generators.GeneratorConfig) *cobra.Command {
	cmd := &cobra.Command{}
	cmd.Use = config.Name
	cmd.Short = config.Description
	if config.DescriptionL != nil {
		cmd.Long = config.LongDescription()
	}
	if config.Aliases != nil {
		cmd.Aliases = *config.Aliases
	}
	return cmd
}
