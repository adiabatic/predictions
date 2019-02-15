package stream

import (
	"errors"
	"fmt"
)

// A PredictionErrorMaker blah blah
type PredictionErrorMaker func(Stream, int) error

func NewInsensibleConfidenceError(s Stream, i int) error {
	f := g(
		"first prediction, with claim “%v”, has a confidence level outside (0%%, 100%%)",
		"Prediction with claim “%v” has a confidence level outside (0%%, 100%%)",
		"Prediction after prediction with claim “%v” has a confidence level outside (0%%, 100%%)",
		"A prediction exists with a confidence level outside (0%%, 100%%). It doesn’t have a claim and its immediate predecessor, if any, lacks a claim too",
	)

	return f(s, i)

}

func g(first, at, atPrev, huh string) PredictionErrorMaker {
	return func(s Stream, i int) error {
		prefix := ""
		if s.FromFilename != "" {
			prefix = s.FromFilename + ": "
		}

		claim := s.Predictions[i].Claim
		previousClaim := ""
		if i > 0 {
			previousClaim = s.Predictions[i-1].Claim
		}

		// now for testing

		if i == 0 {
			return fmt.Errorf(prefix+first, claim)
		}

		if claim != "" {
			return fmt.Errorf(prefix+at, claim)
		}

		if previousClaim != "" {
			return fmt.Errorf(prefix+atPrev, previousClaim)
		}

		return errors.New(prefix + huh)
	}
}

// NewNoClaimError returns an error that describes the approximate location of a prediction that has no claim.
func NewNoClaimError(s Stream, predictionIndex int) error {
	prefix := ""
	if s.FromFilename != "" {
		prefix = s.FromFilename + ": "
	}

	if predictionIndex > 0 {
		previousClaim := s.Predictions[predictionIndex-1].Claim
		if previousClaim != "" {
			return fmt.Errorf("%vclaim after “%v” has no claim", prefix, previousClaim)
		}
		return fmt.Errorf("%va prediction has no claim, and neither does the one before it", prefix)

	}
	return fmt.Errorf("%vthe first prediction has no claim", prefix)
}
