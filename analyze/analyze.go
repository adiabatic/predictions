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
	"fmt"
	"math"
	"os"

	"github.com/adiabatic/predictions/streams"
)

// Analysis is a dump of information all about a list of Streams.
type Analysis struct {
	Everything AnalyzedDocuments // title is ""

	EverythingByKey []AnalyzedDocuments // title is title + scope

	EverythingByTag []AnalyzedDocuments // title is tag

	EverythingByConfidence []AnalyzedDocuments
}

// An AnalyzedDocuments contains both an AnalysisUnit and a slice of PredictionDocument.
type AnalyzedDocuments struct {
	AnalysisUnit AnalysisUnit
	Documents    []streams.PredictionDocument
}

// AnalysisUnit provides information on a subset of an Analysis.
type AnalysisUnit struct {
	Title              string
	SquaredDifferences []float64 // (f_t-o_t)^2

	Called int //   predicted correctly
	Missed int //   predicted incorrectly

	Ongoing  int //   no “happened” value
	Excluded int //   has “cause for exclusion” key with value

	Unscorable int // lacks claim, lacks confidence, or both
}

// Total returns the sum of the scored items, the unscored items, and the unscorable items.
func (au *AnalysisUnit) Total() int { return au.Scored() + au.Unscored() + au.Unscorable }

// Scored returns the sum of the called items and the missed items.
func (au *AnalysisUnit) Scored() int { return au.Called + au.Missed }

// Unscored returns the sum of the ongoing and the excluded items.
func (au *AnalysisUnit) Unscored() int { return au.Ongoing + au.Excluded }

// OfTotalScored calculates the percent scored out of all the predictions made in this analysis unit.
func (au *AnalysisUnit) OfTotalScored() float64 {
	return 100.0 * float64(au.Scored()) / float64(au.Total())
}

// OfTotalCalled calcualtes the percent called (gotten right) out of all the predictions made in this analysis unit.
func (au *AnalysisUnit) OfTotalCalled() float64 {
	return 100.0 * float64(au.Called) / float64(au.Total())
}

// OfTotalMissed calculates the percent missed (gotten wrong) out of all the predictions made in this analysis unit.
func (au *AnalysisUnit) OfTotalMissed() float64 {
	return 100.0 * float64(au.Missed) / float64(au.Total())
}

// OfScoredCalled calculates the percent called (gotten right) out of all the scored predictions made in this analysis unit.
func (au *AnalysisUnit) OfScoredCalled() float64 {
	return 100.0 * float64(au.Called) / float64(au.Scored())
}

// OfScoredMissed calculates the percent missed (gotten wrong) out of all the scored predictions made in this analysis unit.
func (au *AnalysisUnit) OfScoredMissed() float64 {
	return 100.0 * float64(au.Missed) / float64(au.Scored())
}

// OfTotalUnscored calculates the percent unscored (whether because they’re open questions or excluded from consideration) out of all the predictions made in this analysis unit.
func (au *AnalysisUnit) OfTotalUnscored() float64 {
	return 100.0 * float64(au.Unscored()) / float64(au.Total())
}

// OfUnscoredOngoing calculates the percentage of ongoing (not answered yet) predictions out of all unscored predictions in this analysis unit.
func (au *AnalysisUnit) OfUnscoredOngoing() float64 {
	return 100.0 * float64(au.Ongoing) / float64(au.Unscored())
}

// OfUnscoredExcluded calculates the percentage of excluded (for cause) predictions out of all unscored predictions in this analysis unit.
func (au *AnalysisUnit) OfUnscoredExcluded() float64 {
	return 100.0 * float64(au.Excluded) / float64(au.Unscored())
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

// BrierScore calculates the Brier score of added squared differences.
//
// Returns NaN if no squared differences have been added.
func (au *AnalysisUnit) BrierScore() float64 {
	if len(au.SquaredDifferences) == 0 {
		return math.NaN()
	}
	var sum float64
	for _, d := range au.SquaredDifferences {
		sum += d
	}

	return sum / float64(len(au.SquaredDifferences))
}

// Confidence returns the confidence level of an AnalyzedDocuments if they all have the same confidence level, or nil otherwise.
func (ads *AnalyzedDocuments) Confidence() *float64 {
	const ε = 0.0001
	var ret *float64
	for _, d := range ads.Documents {
		if d.Confidence != nil {
			if ret == nil {
				ret = d.Confidence
			} else {
				if math.Abs(*ret-*d.Confidence) > ε {
					// preexisting saved value is too different from current value, so this bunch must be confidence-ly heterogenous
					return nil
				}
			}

		}
	}
	return ret
}

// Analyze calculates Brier scores for the given streams.
func Analyze(sts []streams.Stream) Analysis {
	ret := Analysis{}

	ret.Everything = Only(sts, streams.Everything)
	ret.Everything.AnalysisUnit.Title = "Everything"

	tagsUsed := streams.TagsUsed(sts)
	for _, tag := range tagsUsed {
		ds := Only(sts, streams.MatchingTag(tag))
		ds.AnalysisUnit.Title = fmt.Sprintf("Tag: %s", tag)
		ret.EverythingByTag = append(ret.EverythingByTag, ds)

	}

	keysUsed := streams.KeysUsed(sts)
	for _, key := range keysUsed {
		ds := Only(sts, streams.MatchingKey(key))
		ds.AnalysisUnit.Title = key
		ret.EverythingByKey = append(ret.EverythingByKey, ds)
	}

	confidencesUsed := streams.ConfidencesUsed(sts)
	for _, confidence := range confidencesUsed {
		ds := Only(sts, streams.MatchingConfidence(confidence))
		ds.AnalysisUnit.Title = fmt.Sprintf("At the %.0f%% confidence level", confidence)
		ret.EverythingByConfidence = append(ret.EverythingByConfidence, ds)
	}

	return ret
}

// Only analyzes only predictions in streams that pass a filter.
func Only(sts []streams.Stream, f streams.Filter) AnalyzedDocuments {
	ret := AnalyzedDocuments{}

	for _, st := range sts {
		for _, p := range st.Predictions {
			if !f(p) {
				continue
			}

			ret.Documents = append(ret.Documents, p)

			if p.Claim == "" || p.Confidence == nil {
				ret.AnalysisUnit.Unscorable++
				continue
			}

			if p.Happened == nil {
				if p.CauseForExclusion != "" {
					ret.AnalysisUnit.Excluded++
					continue
				}
				ret.AnalysisUnit.Ongoing++
				continue
			}

			if *(p.Happened) {
				ret.AnalysisUnit.Called++
			} else {
				ret.AnalysisUnit.Missed++
			}

			ret.AnalysisUnit.Add(float64(*(p.Confidence))/100.0, *(p.Happened))
		}
	}

	if len(ret.AnalysisUnit.SquaredDifferences) != ret.AnalysisUnit.Scored() {
		fmt.Fprintln(os.Stderr, "logic error: squared differences and scored items differ")
		os.Exit(4)
	}

	return ret
}
