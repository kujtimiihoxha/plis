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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/kujtimiihoxha/plis/fs"
	"github.com/spf13/afero"
	"github.com/kujtimiihoxha/plis/logger"
	"os"
	"os/exec"
	"strings"
	"github.com/kujtimiihoxha/plis/api"
	"github.com/kujtimiihoxha/plis/config"
	"encoding/json"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a generator from a git repository",
	Long: `Get a generator from a git repository`,
	Run: func(cmd *cobra.Command, args []string) {
		getGenerator(args[0],viper.GetString("get.branch"))
	},
}

func getGenerator(rep string, branch string) {
	dir := checkIfGeneratorFolderExists()
	repository:= strings.Split(rep,"/")
	gen:= repository[len(repository)-1]
	gen = strings.TrimSuffix(gen,".git")
	dir += afero.FilePathSeparator + gen
	cmdArgs :=[]string{
		"clone",
	}
	if branch != ""{
		cmdArgs = append(cmdArgs,"-b",branch,"--single-branch")
	}
	cmdArgs = append(cmdArgs,rep,dir)
	cmd := exec.Command("git",cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		logger.GetLogger().Error(err)
	}
	if !viper.GetBool("get.global"){
		fsAPI := api.NewFsAPI(fs.GetCurrentFs())
		b,err := fsAPI.Exists("plis.json")
		if err != nil {
			logger.GetLogger().Fatal(err)
		}
		if branch == ""{
			branch = "master"
		}
		pd := config.PlisDependency{
			Repository:rep,
			Branch:branch,
		}
		if !b {
			pc := config.PlisConfig{
				Dependencies:[]config.PlisDependency{
					pd,
				},
			}
			data,_:=json.MarshalIndent(pc,"", "    ")
			fsAPI.WriteFile("plis.json",string(data))
			return
		}
		data,err := fsAPI.ReadFile("plis.json")
		pc := config.PlisConfig{}
		json.Unmarshal([]byte(data),&pc)
		pc.Dependencies = append(pc.Dependencies,pd)
		d,_ := json.MarshalIndent(pc,"", "    ")
		fsAPI.WriteFile("plis.json",string(d))
	}

}
func checkIfGeneratorFolderExists() string{
	if viper.GetBool("get.global"){
		fsAPI := api.NewFsAPI(fs.GetPlisRootFs())
		b,err:=fsAPI.Exists("generators")
		if err!= nil {
			logger.GetLogger().Fatal(err)
		}
		if !b {
			err = fsAPI.Mkdir("generators")
			if err!= nil {
				logger.GetLogger().Fatal(err)
			}
		}

		return  viper.GetString("plis.dir.generators")
	}
	fsAPI := api.NewFsAPI(fs.GetCurrentFs())
	b,err:=fsAPI.Exists("plis/generators")
	if err!= nil {
		logger.GetLogger().Fatal(err)
	}
	if !b {
		err = fsAPI.MkdirAll("plis/generators")
		if err!= nil {
			logger.GetLogger().Fatal(err)
		}
	}
	return "plis" + fsAPI.FileSeparator() + "generators"
}
func init() {
	getCmd.Flags().BoolP("global","g",false,"Use if the generator should be installed globally")
	getCmd.Flags().StringP("branch","b","","Use if you want to get a specific branch of the generator")
	viper.BindPFlag("get.global", getCmd.Flags().Lookup("global"))
	viper.BindPFlag("get.branch", getCmd.Flags().Lookup("branch"))
	RootCmd.AddCommand(getCmd)
}
