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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/adiabatic/predictions/streams"
)

func TestConfidencesUnder50Percent(t *testing.T) {
	// First, a basic sanity check.

	type row struct {
		d         streams.PredictionDocument
		containee string
	}

	newRow := func(claim string, confidence float64, happened bool, containee string) row {
		return row{
			d: streams.PredictionDocument{
				Claim:      claim,
				Confidence: &confidence,
				Happened:   &happened,
			},
			containee: containee,
		}
	}

	rows := []row{
		newRow("The mail carrier will come tomorrow", 85.0, true, "<b>"),
		newRow("Ben Franklin will come tomorrow", 5.0, false, "<b>"),
		newRow("Microsoft will abandon browser-engine development", .1, true, "<s>"),
	}

	for _, aRow := range rows {
		s := MarkdownFromDocument(aRow.d)
		assert.Contains(t, s, aRow.containee)
	}
}
