// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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
	"encoding/json"
	"fmt"
	"github.com/kujtimiihoxha/plis/config"
	"github.com/kujtimiihoxha/plis/fs"
	"github.com/kujtimiihoxha/plis/logger"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"strings"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a tool from a git repository",
	Long:  `Get a tool from a git repository`,
	Run: func(cmd *cobra.Command, args []string) {
		getTool(args[0], viper.GetString("plis.get.branch"))
	},
}

func getTool(rep string, branch string) {
	dir := checkIfToolFolderExists()
	repository := strings.Split(rep, "/")
	gen := repository[len(repository)-1]
	gen = strings.TrimSuffix(gen, ".git")
	b, _ := afero.Exists(fs.GetCurrentFs(), dir+afero.FilePathSeparator+gen)
	if b {
		logger.GetLogger().Warn("A tool with the same name already exists")
		return
	}
	dir += afero.FilePathSeparator + gen
	Args := []string{
		"clone",
	}
	if branch != "" {
		Args = append(Args, "-b", branch, "--single-branch")
	}
	Args = append(Args, rep, dir)
	_cmd := exec.Command("git", Args...)
	_cmd.Stdout = os.Stdout
	_cmd.Stdin = os.Stdin
	_cmd.Stderr = os.Stderr
	err := _cmd.Run()
	if err != nil {
		logger.GetLogger().Error(err)
	}
	if !viper.GetBool("plis.get.global") {
		//fsAPI := api.NewFsAPI(fs.GetCurrentFs())
		_fs := fs.GetCurrentFs()
		b, err := afero.Exists(_fs, "plis.json")
		if err != nil {
			logger.GetLogger().Fatal(err)
		}
		if branch == "" {
			branch = "master"
		}
		pd := config.PlisDependency{
			Repository: rep,
			Branch:     branch,
		}
		if !b {
			pc := config.PlisConfig{
				Dependencies: []config.PlisDependency{
					pd,
				},
			}
			data, _ := json.MarshalIndent(pc, "", "    ")
			afero.WriteFile(_fs, "plis.json", data, os.ModePerm)
		} else {
			data, _ := afero.ReadFile(_fs, "plis.json")
			pc := config.PlisConfig{}
			json.Unmarshal([]byte(data), &pc)
			exists := false
			for _, v := range pc.Dependencies {
				if v.Repository == pd.Repository {
					exists = true
				}
			}
			if !exists {
				pc.Dependencies = append(pc.Dependencies, pd)
				d, _ := json.MarshalIndent(pc, "", "    ")
				afero.WriteFile(_fs, "plis.json", d, os.ModePerm)
			}
		}
	}
	_fs := afero.NewBasePathFs(fs.GetCurrentFs(), dir)
	if viper.GetBool("plis.get.global") {
		_fs = afero.NewBasePathFs(afero.NewOsFs(), dir)
	}
	if b, _ := afero.Exists(_fs, "plis.json"); b {
		installDependencies(viper.GetBool("plis.get.global"), dir)
	}
}
func checkIfToolFolderExists() string {
	if viper.GetBool("plis.get.global") {
		//fsAPI := api.NewFsAPI(fs.GetPlisRootFs())
		_fs := fs.GetPlisRootFs()
		b, err := afero.Exists(_fs, "tools")
		if err != nil {
			logger.GetLogger().Fatal(err)
		}
		if !b {
			err = _fs.Mkdir("tools", os.ModePerm)
			if err != nil {
				logger.GetLogger().Fatal(err)
			}
		}

		return viper.GetString("plis.dir.tools")
	}
	//fsAPI := api.NewFsAPI(fs.GetCurrentFs())
	_fs := fs.GetCurrentFs()
	b, err := afero.Exists(_fs, fmt.Sprintf("plis%stools", afero.FilePathSeparator))
	if err != nil {
		logger.GetLogger().Fatal(err)
	}
	if !b {
		err = _fs.MkdirAll(fmt.Sprintf("plis%stools", afero.FilePathSeparator), os.ModePerm)
		if err != nil {
			logger.GetLogger().Fatal(err)
		}
	}
	return "plis" + afero.FilePathSeparator + "tools"
}
func init() {
	getCmd.Flags().BoolP("global", "g", false, "Use if the tool should be installed globally")
	getCmd.Flags().StringP("branch", "b", "", "Use if you want to get a specific branch of the tool")
	viper.BindPFlag("plis.get.global", getCmd.Flags().Lookup("global"))
	viper.BindPFlag("plis.get.branch", getCmd.Flags().Lookup("branch"))
	RootCmd.AddCommand(getCmd)
}
