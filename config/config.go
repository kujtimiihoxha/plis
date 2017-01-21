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

package config

import (
	"fmt"
	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigName("plis") // name of config file (without extension)
	viper.SetConfigType("json")
	viper.AddConfigPath("./") // adding home directory as first search path
	viper.AutomaticEnv()      // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		defaults()
	}
}
func defaults() {
	viper.Set("dir.base", "plis")
	viper.Set("dir.generators", "generators")
	viper.Set("dir.user", "user")
	viper.Set("dir.config", "config")
}
