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

package streams

import "math"

// A Filter removes predictions from consideration if the predicate returns false.
type Filter func(PredictionDocument) bool

// Everything is a Filter that filters nothing out.
func Everything(_ PredictionDocument) bool { return true }

// MatchingTag returns a Filter that returns true if the prediction’s tag matches the given tag.
func MatchingTag(tag string) Filter {
	return func(d PredictionDocument) bool {
		for _, predictionTag := range d.Tags {
			if tag == predictionTag {
				return true
			}
		}
		return false
	}
}

// MatchingKey returns a Filter that returns true if the prediction’s key matches the given key.
func MatchingKey(key string) Filter {
	return func(d PredictionDocument) bool {
		if d.Parent == nil {
			if key == "" {
				return true
			}
		}

		// this combining equality check exists elsewhere in the codebase. I should have a common function for both so they don’t get out of sync.
		return key == d.Parent.Metadata.Title+" "+d.Parent.Metadata.Scope
	}
}

// MatchingConfidence returns a Filter that returns true if the prediction’s confidence level matches (to within a fudge factor) the given confidence level.
func MatchingConfidence(conf float64) Filter {
	const ε = 0.0001
	return func(d PredictionDocument) bool {
		if d.Confidence != nil && math.Abs(*d.Confidence-conf) < ε {
			return true
		}

		return false
	}
}

// DocumentsMatching returns PredictionDocuments from streams that pass the given Filter.
func DocumentsMatching(sts []Stream, f Filter) []PredictionDocument {
	ret := make([]PredictionDocument, 0)
	for _, st := range sts {
		for _, p := range st.Predictions {
			if !f(p) {
				continue
			}

			ret = append(ret, p)
		}
	}

	return ret
}
