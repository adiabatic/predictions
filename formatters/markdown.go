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
	"strings"

	"github.com/adiabatic/predictions/streams"
)

// MarkdownFromDocument makes a Markdown-formatted document.
//
// Formatting:
//
// - indeterminate results are italicized
//
// - things that happened are plain
//
// - things that didn't happen are struck through
//
// Note that MarkdownFromDocument also uses HTML for the italics and the strikethrough. This may be a problem in some contexts that allow markdown but not HTML, like some forum software in some configurations.
func MarkdownFromDocument(d streams.PredictionDocument) string {
	meat := fmt.Sprintf("%v: %v%%", d.Claim, *(d.Confidence))
	withToppings := ""

	switch {
	case d.Happened == nil && d.CauseForExclusion != "":
		withToppings = fmt.Sprintf("- <i>%v</i>", meat)
	case d.Happened == nil:
		withToppings = fmt.Sprintf("- %v", meat)
	case *d.Happened == true:
		withToppings = fmt.Sprintf("- <b>%v</b>", meat)
	case *d.Happened == false:
		withToppings = fmt.Sprintf("- <s>%v</s>", meat)
	}

	return withToppings + "\n"
}

// MarkdownFromStream makes a markdown-formatted stream.
func MarkdownFromStream(st streams.Stream) string {
	var buf strings.Builder

	for _, d := range st.Predictions {
		if d.ShouldExclude() {
			continue
		}

		buf.WriteString(MarkdownFromDocument(d))
	}
	buf.WriteString("\n")
	return buf.String()
}

// MarkdownFromStreams makes a Markdown-formatted version of a slice of streams.
func MarkdownFromStreams(sts []streams.Stream) string {
	var buf strings.Builder

	// TODO: first by title/scope, then by each individual tag…
	for _, st := range sts {
		buf.WriteString(MarkdownFromStream(st))
	}

	buf.WriteString("\n")
	return buf.String()
}
