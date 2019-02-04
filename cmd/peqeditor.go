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

// peqeditorCmd represents the peqeditor command
var peqeditorCmd = &cobra.Command{
	Use:   "peqeditor",
	Short: "PEQ Editor service on http://localhost:81",
	Long:  `PEQ Editor service on http://localhost:81`,
	Run: func(cmd *cobra.Command, args []string) {
		s, err := script.New(cmd)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = s.PEQEditorRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	startCmd.AddCommand(peqeditorCmd)
	stopCmd.AddCommand(peqeditorCmd)
	logsCmd.AddCommand(peqeditorCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// peqeditorCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// peqeditorCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
