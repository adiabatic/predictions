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

// Package brier calculates Brier scores for slices of stream.Stream.
//
// Brier scores go from 0 to 1. 0 is the best score achievable, while 1 is the worst score achievable. If you put an event at the 50% confidence interval, your Brier score will be ¼ regardless of what happens. To aggregate Brier scores, add them all up and divide by how many you have, just like an arithmetic mean.
//
// Here’s how to calculate the Brier score of something when only two different outcomes are possible:
//
// "BS" = 1/N  ∑_(t=1)^N▒(f_t-o_t)^2
//
// …where f_t is the probability that was forecast (0–1, inclusive)
//
// …and o_t is the outcome (0 if it didn’t happen, 1 if it did)
//
// Prefer to read about this on Wikipedia? https://en.wikipedia.org/wiki/Brier_score#Definition
package brier

import (
	"math"

	"github.com/adiabatic/predictions/streams"
)

// A Filter is a function that filters out predictions if the filter returns false.
type Filter func(streams.PredictionDocument) bool

// Everything is a Filter that always returns true.
func Everything(_ streams.PredictionDocument) bool { return true }

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

// ForOnly calculates the Brier score for only predictions where f returns true.
//
// Returns NaN if there were no predictions scored.
func ForOnly(ss []streams.Stream, f Filter) float64 {
	eligiblePredictions := 0
	sum := 0.0

	for _, s := range ss {
		for _, p := range s.Predictions {
			if p.Claim == "" || p.Happened == nil || p.Confidence == nil {
				continue
			}

			if !f(p) {
				continue
			}

			confidence := *(p.Confidence) / 100.0
			outcome := 0.0
			if *(p.Happened) == true {
				outcome = 1.0
			}

			sum += math.Pow(confidence-outcome, 2.0)
			eligiblePredictions++
		}
	}

	if eligiblePredictions == 0 {
		return math.NaN()
	}

	return sum / float64(eligiblePredictions)

}
