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

	"github.com/adiabatic/predictions/streams"
	"github.com/spf13/cobra"
)

func init() {
	publishCommand.AddCommand(publishMarkdownCommand)
}

var publishMarkdownCommand = &cobra.Command{
	Use:                   "markdown",
	Aliases:               []string{"m"},
	Args:                  cobra.MinimumNArgs(1),
	Short:                 "Formats your predictions in Markdown lists",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		sts, err := streams.StreamsFromFiles(args)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		for _, st := range sts {
			header := combineTitleAndScope(st.Metadata.Title, st.Metadata.Scope)

			if header != "" {
				fmt.Println("# " + header)
				fmt.Println()
			}

			fmt.Println(formatters.MarkdownFromStream(st))
		}
	},
}
