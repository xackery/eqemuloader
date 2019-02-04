// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	"os"

	"git.rebuildeq.com/eq/loader/script"
	"github.com/spf13/cobra"
)

// mariadbCmd represents the mariadb command
var mariadbCmd = &cobra.Command{
	Use:   "mariadb",
	Short: "Database",
	Long:  `Database that stores all information related to Everquest`,
	Run: func(cmd *cobra.Command, args []string) {
		s, err := script.New(cmd)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = s.MariaDBRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	startCmd.AddCommand(mariadbCmd)
	stopCmd.AddCommand(mariadbCmd)
	logsCmd.AddCommand(mariadbCmd)
	injectCmd.AddCommand(mariadbCmd)
	dumpCmd.AddCommand(mariadbCmd)
	pullCmd.AddCommand(mariadbCmd)
}