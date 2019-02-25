// © 2019 Nathan Galt
//
// Licensed under the Apache License, Version 2.0 (the “License”);
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an “AS IS” BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/adiabatic/predictions/streams"
	"github.com/spf13/cobra"
)

func init() {
	rootCommand.AddCommand(checkCommand)
}

var checkCommand = &cobra.Command{
	Use:                   "check [filenames…]",
	Aliases:               []string{"c"},
	Short:                 "Ensure predictions files are sensible",
	DisableFlagsInUseLine: true,
	Args:                  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ss, err := streams.StreamsFromFiles(args)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		v := streams.Validator{}

		for _, s := range ss {
			errs := v.RunAll(s)
			for _, err := range errs {
				cmd.Println(err)
			}
		}
	},
}
