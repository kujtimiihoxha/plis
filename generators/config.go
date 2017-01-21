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

package generators

import (
	"encoding/json"
	"fmt"
	"github.com/kujtimiihoxha/plis/generators/cmd"
	"github.com/spf13/afero"
	"os"
)

type GeneratorConfig struct {
	Name         string              `json:"name"`
	Description  string              `json:"description"`
	DescriptionL *[]string           `json:"description_l"`
	Aliases      *[]string           `json:"aliases"`
	Flags        *[]cmd.PlisFlag     `json:"flags"`
	Modules      *[]string           `json:"modules"`
	Arguments    *[]cmd.PlisArgument `json:"args"`
	SubCommands  *[]string           `json:"sub_commands,omitempty"`
}

func (gc *GeneratorConfig) LongDescription() (description string) {
	if gc.DescriptionL == nil {
		return
	}
	for _, v := range *gc.DescriptionL {
		description += description + v + "\n"
	}
	return
}

func ReadConfig(fs afero.Fs, pth string) *GeneratorConfig {
	data, err := afero.ReadFile(fs, pth)
	if err != nil {
		fmt.Println(fmt.Sprintf("Could not read config from `%s`", pth))
		os.Exit(-1)
	}
	config := &GeneratorConfig{}
	if err := json.Unmarshal(data, &config); err != nil {
		fmt.Println(fmt.Sprintf("Could not read json from `%s`", pth))
		os.Exit(-1)
	}
	if config.Flags != nil {
		for _, v := range *config.Flags {
			if v.Long == "" {
				fmt.Println("Flag needs a Long name")
				os.Exit(-1)
			}
		}
	}
	if config.Arguments != nil {
		for _, v := range *config.Arguments {
			if v.Name == "" {
				fmt.Println("Argument needs a Name.")
				os.Exit(-1)
			}
		}
	}
	return config
}
