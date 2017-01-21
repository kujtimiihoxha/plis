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
	"github.com/kujtimiihoxha/plis/config"
	"github.com/kujtimiihoxha/plis/fs"
	"github.com/kujtimiihoxha/plis/generators"
	"github.com/spf13/cobra"
	"os"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "plis",
	Short: "Plis the most simple code generator framework",
	Long: `Plis is a framework to create code generators for all types of projects.

Plis is very easy to use and can be used for very simple tasks to very complicated generators.You can use other open
sourced generators as plugins for plis.

Plis is created by @kujtimiihoxha and is written in golang.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
var RootGenerator *generators.PlisGenerator

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootGenerator = generators.NewPlisGenerator(RootCmd, RootConfig(), nil)
	config.Init()
	fs.Init()
	InitGenerators()
}

func RootConfig() *generators.GeneratorConfig {
	return &generators.GeneratorConfig{
		Name:        "plis",
		Description: "Plis the most simple code generator framework",
		DescriptionL: &[]string{
			"Plis is a framework to create code generators for all types of projects.",
			"",
			"Plis is very easy to use and can be used for very simple tasks to very complicated generators.You can use other open",
			"sourced generators as plugins for plis.",
			"",
			"Plis is created by @kujtimiihoxha and is written in golang.",
		},
	}
}
