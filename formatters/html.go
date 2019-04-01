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

package formatters

import (
	"html/template"
	"io"
	"io/ioutil"
	"strings"

	"github.com/adiabatic/predictions/analyze"
	"github.com/adiabatic/predictions/streams"
	"gopkg.in/russross/blackfriday.v2"
)

type payload struct {
	Title     string
	Scope     string
	PageTitle string
	Streams   []streams.Stream
	Analysis  analyze.Analysis

	ChartJS template.JS
}

func documentResult(d streams.PredictionDocument) (class, message string) {
	switch Evaluate(d) {
	case ExcludedForCause:
		return "excluded", "excluded"
	case Ongoing:
		return "ongoing", "ongoing"
	case Resolved:
		return "resolved", "resolved"
	case Called:
		return "true", "called it"
	case Missed:
		return "false", "missed it"
	}
	return "logic-error", "logic error"
}

// HTMLFromStreams generates HTML output of streams and writes it to w.
func HTMLFromStreams(w io.Writer, sts []streams.Stream) error {
	markdownifyNotes(sts)

	var p payload

	bs, err := ioutil.ReadFile("Chart.min.js")
	if err != nil {
		panic("could not load Chart.min.js")
	}

	p.ChartJS = template.JS(bs)

	if len(sts) == 1 {
		p.PageTitle = combineTitleAndScope(sts[0].Metadata.Title, sts[0].Metadata.Scope)
	}

	p.Streams = sts

	funcs := template.FuncMap{
		"resultClass": func(d streams.PredictionDocument) string {
			class, _ := documentResult(d)
			return class
		},
		"resultMessage": func(d streams.PredictionDocument) string {
			_, message := documentResult(d)
			return message
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
	return t.Execute(w, p)
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

func combineTitleAndScope(title, scope string) string {
	r := ""
	if title != "" {
		r += title
		if scope != "" {
			r += ": " + scope
		}
	} else {
		r += scope
	}
	return r
}
