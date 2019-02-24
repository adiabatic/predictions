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

func DocumentAsMarkdown(d streams.PredictionDocument) string {
	meat := fmt.Sprintf("%v: %v%%", d.Claim, *(d.Confidence))
	withToppings := ""
	if d.Happened == nil {
		withToppings = fmt.Sprintf("- <i>%v</i>", meat)
	} else if *(d.Happened) {
		withToppings = fmt.Sprintf("- %v", meat)
	} else {
		withToppings = fmt.Sprintf("- <s>%v</s>", meat)
	}

	return withToppings + "\n"
}

func StreamAsMarkdown(st streams.Stream) string {
	var buf strings.Builder

	for _, d := range st.Predictions {
		if d.ShouldExclude() {
			continue
		}

		buf.WriteString(DocumentAsMarkdown(d))
	}
	buf.WriteString("\n")
	return buf.String()
}

func AsMarkdown(sts []streams.Stream) string {
	var buf strings.Builder

	// TODO: first by title/scope, then by each individual tag…
	for _, st := range sts {
		buf.WriteString(StreamAsMarkdown(st))
	}

	buf.WriteString("\n")
	return buf.String()
}
