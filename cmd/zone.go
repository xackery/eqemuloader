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

// zoneCmd represents the zone command
var zoneCmd = &cobra.Command{
	Use:   "zone",
	Short: "An Everquest zone service",
	Long:  `An everquest zone service is used to spin up a zone. Multiple zones can be requested`,
	Run: func(cmd *cobra.Command, args []string) {
		s, err := script.New(cmd)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = s.ZoneRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	startCmd.AddCommand(zoneCmd)
	stopCmd.AddCommand(zoneCmd)
	logsCmd.AddCommand(zoneCmd)
	crashCmd.AddCommand(zoneCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// zoneCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// zoneCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	zoneCmd.Flags().IntP("count", "c", 1, "number of zones to start")
}
