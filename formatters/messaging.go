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

	"github.com/adiabatic/predictions/streams"
)

// Constants for describing the various states of a prediction.
const (
	ExcludedForCause = iota
	Ongoing
	Called
	Missed
)

// Evaluate returns an int describing whether the prediction is excluded, ongoing, called, or missed.
func Evaluate(d streams.PredictionDocument) int {
	switch {
	case d.Happened == nil && d.CauseForExclusion != "": // must come before d.Happened checks
		return ExcludedForCause
	case d.Happened == nil:
		return Ongoing
	case (*d.Happened == true && d.Confidence != nil && *d.Confidence > 50) ||
		(*d.Happened == false && d.Confidence != nil && *d.Confidence < 50):
		return Called
	case (*d.Happened == false && d.Confidence != nil && *d.Confidence > 50) ||
		(*d.Happened == true && d.Confidence != nil && *d.Confidence < 50):
		return Missed
	default:
		panic(fmt.Sprintf("logic error in formatters.Evaluate given document: %#v", d))
	}
}
