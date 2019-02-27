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
	"html/template"
	"os"
	"strings"

	"gopkg.in/russross/blackfriday.v2"

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
	Use:                   "html",
	Aliases:               []string{"h"},
	Short:                 "Formats your predictions as an HTML file",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		sts, err := streams.StreamsFromFiles(args)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		markdownifyNotes(sts)

		var p payload

		if len(sts) == 1 {
			p.PageTitle = combineTitleAndScope(sts[0].Metadata.Title, sts[0].Metadata.Scope)
		}

		p.Streams = sts

		funcs := template.FuncMap{
			// a three-valued bool should be a “trool”, right?
			"troolToString": func(v *bool) string {
				switch {
				case v == nil:
					return "null"
				case *v == true:
					return "true"
				case *v == false:
					return "false"
				}

				return "unexpected"
			},
			"commaSeparate": func(ss []string) string {
				return strings.Join(ss, ", ")
			},
			"safeHTML": func(s string) template.HTML {
				return template.HTML(s)
			},
			"refAU": func(au analyze.AnalysisUnit) *analyze.AnalysisUnit {
				return &au
			},
		}

		p.Analysis = analyze.Analyze(sts)

		t := template.Must(template.New("template").Funcs(funcs).ParseFiles("template.html"))
		err = t.Execute(os.Stdout, p)
		if err != nil {
			cmd.Println("error when executing template: ", err)
			os.Exit(2)
		}

	},
}

func markdownifyNotes(sts []streams.Stream) {
	for _, st := range sts {
		for i, d := range st.Predictions {
			if d.Notes != "" {
				st.Predictions[i].Notes = string(blackfriday.Run([]byte(d.Notes)))
			}
		}
	}
}
