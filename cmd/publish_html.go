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

	"github.com/adiabatic/predictions/streams"
	"github.com/spf13/cobra"
)

func init() {
	publishCommand.AddCommand(publishHTMLCommand)
}

const htmlTemplate = `<!DOCTYPE html>
<html lang='en'>
<head>
	<meta charset='UTF-8'>
	<title>{{ .PageTitle }}</title>
	<style>
		html {
		    font-family: system-ui, sans-serif;
			--radius: 0rem;
			--border-style-internal: thin dotted gray;
		}

		section {
			padding: 0 .5rem;
		}
		
		.document {
            display: grid;
			
			/* percent, claim, result, tags, notes */
			grid-template:
				'p c'
				'r t'
				'r n'
				/ 15ch 1fr;

			border: thin solid gray;
			border-radius: var(--radius);
		}

		.document > * {
			background: #eee;
			padding: .25rem;
		}

		.percent {
			grid-area: p;
			font-size: 300%;
			font-weight: 900;
			text-align: center;
			border-top-left-radius: var(--radius);

			border-bottom: var(--border-style-internal);
			border-right: var(--border-style-internal);
		}

		.percent > div {
			padding-bottom: .5rem;
		}

		.claim {
			grid-area: c;
			font-size: 200%;

			border-top-right-radius: var(--radius);
			border-bottom: var(--border-style-internal);
		}

		.result {
			grid-area: r;

			border-bottom-left-radius: var(--radius);
			border-right: var(--border-style-internal);
		}

		.result.true {
			background: hsl(120, 50%, 90%);
			color: hsl(120, 50%, 50%);
		}

		.tags {
			grid-area: t;
			color: #666;

			border-bottom: var(--border-style-internal);

		}

		.notes {
			grid-area: n;
			border-bottom-right-radius: var(--radius);
		}

		.center-child {
			display: flex;
			align-items: center;
			justify-content: center;
		}

		.center-child-vertically {
			display: flex;
			align-items: center;
			justify-content: start;
		}
	</style>
</head>
<body>
	{{ range .Streams }}
	<section>
		<hgroup>
			{{ with .Metadata.Title }}<h1>{{ . }}</h1>{{ end }}
			{{ with .Metadata.Scope }}<h2>{{ . }}</h2>{{ end }}
		</hgroup>

		{{ range .Predictions }}
		<div class='document'>
			<div class='percent center-child'><div>{{ .Confidence }}%</div></div>
			<div class='result {{ printf "%v" (troolToString .Happened) }} center-child'><div>
				{{- if eq (troolToString .Happened) "true" -}}
					called it
				{{- else if eq (troolToString .Happened) "false" -}}
					missed it
				{{- else -}}
					ongoing
				{{- end -}}
			</div></div>    
			<div class='claim center-child-vertically'><div>{{ .Claim }}</div></div>
			<div class='tags'>tags: {{ commaSeparate .Tags }}</div>
			<div class='notes'>{{ .Notes }}</div>
		</div>
		{{ end }}
	</section>
	{{ end }}
</body>
</html>

`

type payload struct {
	Title     string
	Scope     string
	PageTitle string
	Streams   []streams.Stream
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

		var p payload

		if len(sts) == 1 {
			p.PageTitle = combineTitleAndScope(sts[0].Metadata.Title, sts[0].Metadata.Scope)
		}

		p.Streams = sts

		funcs := template.FuncMap{
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
		}

		t := template.Must(template.New("whatever").Funcs(funcs).Parse(htmlTemplate))
		err = t.Execute(os.Stdout, p)
		if err != nil {
			cmd.Println("error when executing template: ", err)
			os.Exit(2)
		}

	},
}
