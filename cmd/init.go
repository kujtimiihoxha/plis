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
	"github.com/Songmu/prompter"
	"github.com/kujtimiihoxha/plis/fs"
	"github.com/kujtimiihoxha/plis/helpers"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/alioygur/godash.v0"
	"gopkg.in/flosch/pongo2.v3"
	"os"
	"strings"
)

var plisFolder string
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initiate plis project",
	Long: `This generator is used to create necessary folders and files for plis to work.
Init wil generat the plis folder where you can store your generators and the plis config.`,
}

func init() {
	initCmd.Flags().StringVarP(&plisFolder, "folder", "f", "plis", "The base plis folder name.")
	initCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			fmt.Println("Please add the project name as an argument ex. `plis init my-project`")
			os.Exit(-1)
		}
		if a, _ := afero.Exists(fs.WorkingDirFs(), "plis.json"); a {
			if !prompter.YN("The plis config file already exists do you want to override it ?", false) {
				return nil
			}
		}

		viper.Set("dir.base", plisFolder)
		name := args[0]
		if a, _ := afero.Exists(fs.WorkingDirFs(), helpers.BasePath()); !a {
			fs.WorkingDirFs().MkdirAll(helpers.BasePath(), os.ModePerm)

		}
		if a, _ := afero.Exists(fs.WorkingDirFs(), helpers.BasePath()+"user/config"); !a {
			fs.WorkingDirFs().MkdirAll(helpers.BasePath()+"user/config", os.ModePerm)
		}
		config := `{
  "name":"{{ name }}",
  "dir":{
    "base":"{{ base }}",
    "generators":"generators",
    "user":"user"
  }
}`
		t, _ := pongo2.FromString(config)
		s, _ := t.Execute(map[string]interface{}{
			"name": strings.Replace(godash.ToSnakeCase(name), "_", "-", -1),
			"base": plisFolder,
		})

		return afero.WriteFile(fs.WorkingDirFs(), "plis.json", []byte(s), os.ModePerm)
	}
	RootCmd.AddCommand(initCmd)
}
