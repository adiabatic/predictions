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
	"bytes"
	"fmt"
	"os"
	"sort"

	"github.com/adiabatic/predictions/formatters"
	"github.com/adiabatic/predictions/streams"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
)

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {

}

var rootCommand = &cobra.Command{
	Use:                   "predictions",
	Short:                 "predictions finds out how well-calibrated your predictions are",
	DisableFlagsInUseLine: true,
	Args:                  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ss, err := streams.StreamsFromFiles(args)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// fmt.Println("the streams:")
		// spew.Dump(streams)

		tags := streams.TagsUsed(ss)
		sort.Strings(tags)
		var v streams.Validator

		for _, s := range ss {
			errs := v.RunAll(s)
			if len(errs) > 0 {
				spew.Dump(errs)
			}
		}

		for _, tag := range tags {
			buf := &bytes.Buffer{}
			for _, s := range ss {

				for _, d := range s.Predictions {
					if d.ShouldExclude() {
						continue // only print out the ones with cause
					}

					if d.HasTag(tag) {
						fmt.Fprintln(buf, formatters.MarkdownFromDocument(d))
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
