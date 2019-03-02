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

	"github.com/adiabatic/predictions/formatters"

	"github.com/adiabatic/predictions/analyze"
	"github.com/adiabatic/predictions/streams"
	"github.com/spf13/cobra"
)

func init() {
	publishCommand.AddCommand(publishHTMLCommand)
}

type payload struct {
	Title     string
	Scope     string
	PageTitle string
	Streams   []streams.Stream
	Analysis  analyze.Analysis
}

var publishHTMLCommand = &cobra.Command{
	Use:                   "html FILE …",
	Aliases:               []string{"h"},
	Short:                 "Formats your predictions as an HTML file",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		sts, err := streams.FromFiles(args)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		v := streams.Validator{}

		for _, st := range sts {
			errs := v.RunAll(st)
			for _, err := range errs {
				cmd.Println(err)
			}
		}

		err = formatters.HTMLFromStreams(os.Stdout, sts)
		if err != nil {
			cmd.Println("error when executing template: ", err)
			os.Exit(2)
		}
	},
}
