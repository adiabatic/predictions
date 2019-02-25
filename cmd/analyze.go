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
	"sort"
	"strings"

	"github.com/adiabatic/predictions/formatters"
	"github.com/adiabatic/predictions/streams"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
)

func init() {
	rootCommand.AddCommand(analyzeCommand)
}

var analyzeCommand = &cobra.Command{
	Use:                   "analyze [filenames…]",
	Aliases:               []string{"a", "analyse"},
	Short:                 "Runs analyses on your predictions",
	DisableFlagsInUseLine: true,
	Args:                  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sts, err := streams.StreamsFromFiles(args)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// fmt.Println("the streams:")
		// spew.Dump(streams)

		tags := streams.TagsUsed(sts)
		sort.Strings(tags)
		var v streams.Validator

		for _, st := range sts {
			errs := v.RunAll(st)
			if len(errs) > 0 {
				spew.Dump(errs)
			}
		}

		for _, tag := range tags {
			var buf strings.Builder
			for _, st := range sts {

				for _, d := range st.Predictions {
					if d.ShouldExclude() {
						continue // only print out the ones with cause
					}

					if d.HasTag(tag) {
						buf.WriteString(formatters.MarkdownFromDocument(d))
					}
				}
			}

			if buf.Len() > 0 {
				fmt.Println("#", tag)
				fmt.Println()
				fmt.Println(buf.String())
			}
		}
	},
}

// Execute runs the root command.
func Execute() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
