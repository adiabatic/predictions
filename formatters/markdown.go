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
// - ongoing predictions are plain
//
// - predictions that were called correctly are bold
//
// - predictions that were mis-called are struck through
//
// - excluded-for-cause predictions are italicized
//
// Note that MarkdownFromDocument also uses HTML for the italics and the strikethrough. This may be a problem in some contexts that allow markdown but not HTML, like some forum software in some configurations.
func MarkdownFromDocument(d streams.PredictionDocument) string {
	meat := fmt.Sprintf("%v: %v%%", d.Claim, *(d.Confidence))
	withToppings := ""

	switch Evaluate(d) {
	case ExcludedForCause:
		withToppings = fmt.Sprintf("- <i>%v</i>", meat)
	case Ongoing:
		withToppings = fmt.Sprintf("- %v", meat)
	case Called:
		withToppings = fmt.Sprintf("- <b>%v</b>", meat)
	case Missed:
		withToppings = fmt.Sprintf("- <s>%v</s>", meat)
	default:
		panic(fmt.Sprintf("logic error in MarkdownFromDocument given document: %#v", d))
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
	return buf.String()
}

// MarkdownFromStreams makes a Markdown-formatted version of a slice of streams.
func MarkdownFromStreams(sts []streams.Stream, options ...Option) string {
	var buf strings.Builder

	o := formattingOptions{}
	for _, f := range options {
		f(&o)
	}

	// TODO: first by title/scope, then by each individual tag…
	buf.WriteString("# Everything\n\n")
	for _, st := range sts {
		buf.WriteString(MarkdownFromStream(st))
	}

	buf.WriteString("\n")

	keysUsed := streams.KeysUsed(sts)
	if len(keysUsed) > 1 {
		for _, key := range keysUsed {
			fmt.Fprintf(&buf, "# %s\n\n", key)

			for _, d := range streams.DocumentsMatching(sts, streams.MatchingKey(key)) {
				buf.WriteString(MarkdownFromDocument(d))
			}

			buf.WriteString("\n")
		}
	}

	tagsUsed := streams.TagsUsed(sts)
	if len(tagsUsed) > 0 {
		for _, tag := range tagsUsed {
			fmt.Fprintf(&buf, "# %s\n\n", tag)

			for _, d := range streams.DocumentsMatching(sts, streams.MatchingTag(tag)) {
				buf.WriteString(MarkdownFromDocument(d))
			}

			buf.WriteString("\n")
		}
	}

	return buf.String()
}
