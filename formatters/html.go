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
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"

	"github.com/adiabatic/predictions/analyze"
	"github.com/adiabatic/predictions/streams"
	"github.com/gobuffalo/packr/v2"
	"github.com/russross/blackfriday/v2"
)

type payload struct {
	Title     string
	Scope     string
	PageTitle string
	Streams   []streams.Stream
	Analysis  analyze.Analysis

	ChartJS     template.JS
	PerfectData []Point
	GuessData   []Point
}

// A Point struct contains an x and y point. Used for Chart.js.
type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// MarshalJSON suppresses excess precision in float64’s default marshaling behavior.
//
// Do you want to see tooltips for points on a chart that say (0.7000000000000001, 0.7000000000000001)? Me neither.
func (p Point) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf(`{"x":%.2f, "y":%.2f}`, p.X, p.Y)
	return []byte(s), nil
}

func documentResult(d streams.PredictionDocument) (class, message string) {
	switch Evaluate(d) {
	case ExcludedForCause:
		return "excluded", "excluded"
	case Ongoing:
		return "ongoing", "ongoing"
	case Resolved:
		return "resolved", "resolved"
	case CalledTruePositive:
		return "true-positive", "called it"
	case CalledTrueNegative:
		return "true-negative", "called it"
	case MissedFalsePositive:
		return "false-positive", "missed it"
	case MissedFalseNegative:
		return "false-negative", "missed it"
	}
	return "logic-error", "logic error"
}

// HTMLFromStreams generates HTML output of streams and writes it to w.
func HTMLFromStreams(w io.Writer, sts []streams.Stream) error {
	markdownifyNotes(sts)

	var p payload

	box := packr.New("everything", "../templates")

	bs, err := box.Find("Chart.min.js")
	if err != nil {
		panic("could not load Chart.min.js")
	}

	p.ChartJS = template.JS(bs)

	if len(sts) == 1 {
		p.PageTitle = combineTitleAndScope(sts[0].Metadata.Title, sts[0].Metadata.Scope)
	}

	p.Streams = sts

	funcs := template.FuncMap{
		"explainResult": func(d streams.PredictionDocument) string {

			because := ""
			switch Evaluate(d) {
			case ExcludedForCause:
				because = "this prediction was deliberately excluded from consideration"
			case Ongoing:
				because = "it’s too soon to say whether this has happened or not"
			case Resolved:
				because = "whatever was predicted has come to pass, but if your prediction was that it had a 50/50 chance of happening, nobody can really say if your prediction was correct"
			case CalledTruePositive:
				because = "you said this would happen, and it did"
			case CalledTrueNegative:
				because = "you said this wouldn’t happen, and it didn’t"
			case MissedFalsePositive:
				because = "you said this would happen, but it didn’t"
			case MissedFalseNegative:
				because = "you said this wouldn’t happen, but it did anyway"
			}

			_, message := documentResult(d)
			ret := fmt.Sprintf("“%s” is here because %s.", message, because)

			return ret
		},
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

	addPerfectData(&p)
	addGuessData(&p)

	templateQuaString, err := box.FindString("template.html")
	if err != nil {
		panic("could not load template.html")
	}

	t := template.Must(template.New("").Funcs(funcs).Parse(templateQuaString))
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

func addPerfectData(pp *payload) {
	numbers := make([]float64, 0)

	var num float64
	for ; num <= 1.03; num += 0.05 {
		numbers = append(numbers, num)
	}
	numbers[0] = 0.01              // 0 is a bad idea
	numbers[len(numbers)-1] = 0.99 // 1 is a bad idea

	for _, n := range numbers {
		p := Point{n, n}
		pp.PerfectData = append(pp.PerfectData, p)
	}
}

func addGuessData(pp *payload) {
	confidenceGroupings := pp.Analysis.EverythingByConfidence
	for _, grouping := range confidenceGroupings {
		conf := grouping.Confidence()
		if conf == nil {
			fmt.Fprintln(os.Stderr, "Odd. I didn’t expect a nil confidence in addGuessData().")
			continue
		}

		p := Point{
			*conf / 100,
			grouping.AnalysisUnit.OfTotalCalled() / 100,
		}

		pp.GuessData = append(pp.GuessData, p)
	}
}
