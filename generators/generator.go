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
	"fmt"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"os"
)

type PlisGenerator struct {
	Cmd    *cobra.Command
	Config *GeneratorConfig
	Parent *PlisGenerator
}

func NewPlisGenerator(cmd *cobra.Command, config *GeneratorConfig, parent *PlisGenerator) *PlisGenerator {
	return &PlisGenerator{
		Cmd:    cmd,
		Config: config,
		Parent: parent,
	}
}
func (pg *PlisGenerator) GetRootParent() *PlisGenerator {
	res := pg
	for res.Parent.Config.Name != "plis" {
		res = res.Parent
	}
	return res
}
func (pg *PlisGenerator) ValidateArguments(args []string) {
	required := 0
	if pg.Config.Arguments == nil {
		return
	}
	for _, v := range *pg.Config.Arguments {
		if v.Required {
			required++
		}
	}
	if len(args) < required {
		fmt.Println("Please add all requred arguments.")
		os.Exit(-1)
	}
}
func (pg *PlisGenerator) ValidateFlags(c *cobra.Command) {
	if pg.Config.Flags == nil {
		return
	}
	for _, v := range *pg.Config.Flags {
		if !flagChanged(c.Flags(), v.Long) && v.Required {
			fmt.Println(fmt.Sprintf("Please add required flag , `--%s` is required", v.Long))
			os.Exit(-1)
		}
	}
}
func flagChanged(flags *flag.FlagSet, key string) bool {
	flag := flags.Lookup(key)
	if flag == nil {
		return false
	}
	return flag.Changed
}
