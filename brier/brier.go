package brier

import (
	"math"

	"github.com/adiabatic/predictions/stream"
)

// A Filter is a function that filters out predictions if the filter returns false.
type Filter func(stream.PredictionDocument) bool

// Everything is a brierFilter that always returns true.
func Everything(_ stream.PredictionDocument) bool { return true }

// MatchingTag returns a brierFilter that returns true if the prediction’s tag matches the given tag.
func MatchingTag(tag string) Filter {
	return func(d stream.PredictionDocument) bool {
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
func ForOnly(ss []stream.Stream, f Filter) float64 {
	eligiblePredictions := 0
	sum := 0.0

	for _, s := range ss {
		for _, p := range s.Predictions {
			if p.Claim == "" || p.Happened == nil {
				continue
			}

			if !f(p) {
				continue
			}

			confidence := p.Confidence / 100.0
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