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

package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/asaskevich/govalidator"
	"github.com/kujtimiihoxha/plis/cmd"
	"github.com/kujtimiihoxha/plis/config"
	"github.com/kujtimiihoxha/plis/generators"
	"github.com/kujtimiihoxha/plis/helpers"
	"github.com/kujtimiihoxha/plis/logger"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

func main() {
	setupConfigValidator()
	configDefaults()
	initFirstRun()
	logger.SetLevel(logrus.WarnLevel)
	generators.Initialize()
	cmd.Execute()
}
func initFirstRun() {
	checkGit()
	cmd.RootCmd.PersistentFlags().BoolP("debug", "d", false, "Is plis debugging")
	cmd.RootCmd.PersistentFlags().String("debug_folder", "", "Root folder of the debug mode")
	viper.BindPFlag("plis.debug_folder", cmd.RootCmd.PersistentFlags().Lookup("debug_folder"))
	var err error
	if _, err = os.Stat(viper.GetString("plis.dir.root")); err == nil {
		return
	}
	if os.IsNotExist(err) {
		logger.GetLogger().Info("Plis root does not exist")
		logger.GetLogger().Info("Initializing first run...")
		logger.GetLogger().Info(fmt.Sprintf(
			"Creating plis root in `%s`...",
			viper.GetString("plis.dir.root")))
		err := os.MkdirAll(viper.GetString("plis.dir.root")+afero.FilePathSeparator+"generators", os.ModePerm)
		if err != nil {
			logger.GetLogger().Fatal(err)
		}
		return
	}
	logger.GetLogger().Fatal(err)
}
func checkGit() {
	cmd := exec.Command("git", "--version")
	err := cmd.Run()
	if err != nil {
		logger.GetLogger().Fatal("Plis needs git to be installed please install git and try again.")
	}
}
func configDefaults() {
	usr, err := user.Current()
	if err != nil {
		logger.GetLogger().Fatal(err)
	}
	viper.Set("plis.dir.user", usr.HomeDir)
	viper.Set("plis.dir.root", usr.HomeDir+string(filepath.Separator)+".plis")
	viper.Set("plis.dir.generators", usr.HomeDir+string(filepath.Separator)+".plis"+
		string(filepath.Separator)+"generators")
}
func setupConfigValidator() {
	govalidator.CustomTypeTagMap.Set("inputType", govalidator.CustomTypeValidator(func(i interface{}, o interface{}) bool {
		return helpers.StringInSlice(i.(string), config.InputTypes)
	}))
	govalidator.CustomTypeTagMap.Set("scriptType", govalidator.CustomTypeValidator(func(i interface{}, o interface{}) bool {
		return helpers.StringInSlice(i.(string), config.ScriptTypes)
	}))
	govalidator.CustomTypeTagMap.Set("lenOne", govalidator.CustomTypeValidator(func(i interface{}, o interface{}) bool {
		return len(i.(string)) <= 1
	}))
}
