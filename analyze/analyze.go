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

package analyze

import (
	"math"

	"github.com/adiabatic/predictions/streams"
)

// Analysis is a dump of information all about a list of Streams.
type Analysis struct {
	Everything AnalysisUnit // title is ""

	EverythingByKey []AnalysisUnit // title is title + scope

	EverythingByTag []AnalysisUnit // title is tag

}

// AnalysisUnit provides information on a subset of an Analysis.
type AnalysisUnit struct {
	Title              string
	SquaredDifferences []float64 // (f_t-o_t)^2
}

// Add adds a prediction to the unit.
//
// confidence must be on [0, 1].
func (au *AnalysisUnit) Add(confidence float64, happened bool) {
	if au.SquaredDifferences == nil {
		au.SquaredDifferences = make([]float64, 0)
	}

	outcome := 0.0
	if happened {
		outcome = 1.0
	}

	ret := math.Pow(confidence-outcome, 2.0)

	au.SquaredDifferences = append(au.SquaredDifferences, ret)
}

// Count returns the number of added predictions.
func (au *AnalysisUnit) Count() int {
	if au == nil {
		return 0
	}
	return len(au.SquaredDifferences)
}

// BrierScore calculates the Brier score of added squared differences.
//
// Returns NaN if no squared differences have been added.
func (au *AnalysisUnit) BrierScore() float64 {
	if au.Count() == 0 {
		return math.NaN()
	}
	var sum float64
	for _, d := range au.SquaredDifferences {
		sum += d
	}

	return sum
}

// Analyze calculates Brier scores for the given streams.
func Analyze(sts []streams.Stream) Analysis {
	ret := Analysis{}

	ret.Everything = AnalyzeOnly(sts, Everything)

	tagsUsed := streams.TagsUsed(sts)
	for _, tag := range tagsUsed {
		ret.EverythingByTag = append(ret.EverythingByTag,
			AnalyzeOnly(sts, MatchingTag(tag)),
		)
	}

	keysUsed := streams.KeysUsed(sts)
	for _, key := range keysUsed {
		ret.EverythingByKey = append(ret.EverythingByKey,
			AnalyzeOnly(sts, MatchingKey(key)),
		)
	}

	return ret
}

// A Filter removes predictions from consideration if the predicate returns false.
type Filter func(streams.PredictionDocument) bool

// Everything is a Filter that filters nothing out.
func Everything(_ streams.PredictionDocument) bool { return true }

// AnalyzeOnly analyzes only predictions in streams that pass a filter.
func AnalyzeOnly(sts []streams.Stream, f Filter) AnalysisUnit {
	ret := AnalysisUnit{}

	for _, st := range sts {
		for _, p := range st.Predictions {
			if p.Claim == "" || p.Happened == nil || p.Confidence == nil {
				continue
			}

			if !f(p) {
				continue
			}

			ret.Add(*(p.Confidence), *(p.Happened))
		}
	}

	return ret
}

// MatchingTag returns a Filter that returns true if the prediction’s tag matches the given tag.
func MatchingTag(tag string) Filter {
	return func(d streams.PredictionDocument) bool {
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
	return func(d streams.PredictionDocument) bool {
		if d.Parent == nil {
			if key == "" {
				return true
			}
		}
		return key == d.Parent.Metadata.Title+d.Parent.Metadata.Scope

	}
}
