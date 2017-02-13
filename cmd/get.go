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
	"github.com/kujtimiihoxha/plis/api"
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
	Short: "Get a generator from a git repository",
	Long:  `Get a generator from a git repository`,
	Run: func(cmd *cobra.Command, args []string) {
		getGenerator(args[0], viper.GetString("plis.get.branch"))
	},
}

func getGenerator(rep string, branch string) {
	dir := checkIfGeneratorFolderExists()
	repository := strings.Split(rep, "/")
	gen := repository[len(repository)-1]
	gen = strings.TrimSuffix(gen, ".git")
	b, _ := afero.Exists(fs.GetCurrentFs(), dir+afero.FilePathSeparator+gen)
	if b {
		logger.GetLogger().Warn("A generator with the same name already exists")
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
		fsAPI := api.NewFsAPI(fs.GetCurrentFs())
		b, err := fsAPI.Exists("plis.json")
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
			fsAPI.WriteFile("plis.json", string(data))
		} else {
			data, _ := fsAPI.ReadFile("plis.json")
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
				fsAPI.WriteFile("plis.json", string(d))
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
func checkIfGeneratorFolderExists() string {
	if viper.GetBool("plis.get.global") {
		fsAPI := api.NewFsAPI(fs.GetPlisRootFs())
		b, err := fsAPI.Exists("generators")
		if err != nil {
			logger.GetLogger().Fatal(err)
		}
		if !b {
			err = fsAPI.Mkdir("generators")
			if err != nil {
				logger.GetLogger().Fatal(err)
			}
		}

		return viper.GetString("plis.dir.generators")
	}
	fsAPI := api.NewFsAPI(fs.GetCurrentFs())
	b, err := fsAPI.Exists(fmt.Sprintf("plis%sgenerators", fsAPI.FileSeparator()))
	if err != nil {
		logger.GetLogger().Fatal(err)
	}
	if !b {
		err = fsAPI.MkdirAll(fmt.Sprintf("plis%sgenerators", fsAPI.FileSeparator()))
		if err != nil {
			logger.GetLogger().Fatal(err)
		}
	}
	return "plis" + fsAPI.FileSeparator() + "generators"
}
func init() {
	getCmd.Flags().BoolP("global", "g", false, "Use if the generator should be installed globally")
	getCmd.Flags().StringP("branch", "b", "", "Use if you want to get a specific branch of the generator")
	viper.BindPFlag("plis.get.global", getCmd.Flags().Lookup("global"))
	viper.BindPFlag("plis.get.branch", getCmd.Flags().Lookup("branch"))
	RootCmd.AddCommand(getCmd)
}
