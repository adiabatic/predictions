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
	"github.com/davecgh/go-spew/spew"
)

const (
	// ExcludedForCause is used to describe a prediction that has been deliberately excluded from consideration.
	ExcludedForCause = iota

	// Ongoing is used to describe a prediction that hasn’t been resolved yet.
	Ongoing

	// Resolved is used to describe a prediction that can be evaluated, but cannot be classified as “called” or “missed” because it predicts something will happen at a 50% confidence interval.
	Resolved

	//Called

	// CalledTruePositive is used to describe a prediction where something was predicted to happen, and that thing happened.
	CalledTruePositive

	// CalledTrueNegative is used to describe something that was predicted to not happen that did not happen.
	CalledTrueNegative

	// Missed

	// MissedFalsePositive is used to describe a prediction where something was predicted to happen, but didn’t.
	MissedFalsePositive

	// MissedFalseNegative is used to describe a prediction where something was predicted to not happen, but happened anyway.
	MissedFalseNegative
)

// Evaluate returns an int describing whether the prediction is excluded, ongoing, called, or missed.
func Evaluate(d streams.PredictionDocument) int {
	switch {
	case d.Happened == nil && d.CauseForExclusion != "": // must come before d.Happened checks
		return ExcludedForCause
	case d.Happened == nil:
		return Ongoing
	case (*d.Happened == true || *d.Happened == false) && d.Confidence != nil && *d.Confidence == 50:
		// Yes, everything in parentheses above is redundant but it communicates intent
		return Resolved
	case *d.Happened == true && d.Confidence != nil && *d.Confidence > 50:
		return CalledTruePositive
	case *d.Happened == false && d.Confidence != nil && *d.Confidence < 50:
		return CalledTrueNegative
	case *d.Happened == false && d.Confidence != nil && *d.Confidence > 50:
		return MissedFalsePositive
	case *d.Happened == true && d.Confidence != nil && *d.Confidence < 50:
		return MissedFalseNegative
	default:
		// Reduce noise in debug output
		d.Parent = nil
		scs := spew.NewDefaultConfig()
		scs.DisablePointerAddresses = true

		panic(fmt.Sprintf("logic error in formatters.Evaluate given document: %s", scs.Sdump(d)))
	}
}
