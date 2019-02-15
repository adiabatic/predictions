package stream

import (
	"errors"
	"fmt"
)

// A PredictionErrorMaker takes a Stream and an index and returns an error. The index is meant to be the index of the prediction, so the first prediction is referred to with a zero index.
type PredictionErrorMaker func(Stream, int) error

// NewInsensibleConfidenceError returns an error for a prediction at an index with an insensible confidence level (not between 0% and 100%, exclusive)
func NewInsensibleConfidenceError(s Stream, i int) error {
	return makePredictionErrorMaker(
		"has a confidence level outside (0%%, 100%%)",
	)(s, i)
}

func makePredictionErrorMaker(meme string) PredictionErrorMaker {
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

		first := "first prediction, with claim “%v”, " + meme
		at := "prediction with claim “%v” " + meme
		atPrev := "prediction after prediction with claim “%v” " + meme
		huh := "prediction exists that " + meme + "; neither it nor its predecessor have a claim"

		switch {
		case i == 0:
			return fmt.Errorf(prefix+first, claim)
		case claim != "":
			return fmt.Errorf(prefix+at, claim)
		case previousClaim != "":
			return fmt.Errorf(prefix+atPrev, previousClaim)
		default:
			return errors.New(prefix + huh)
		}
	}
}

// NewNoConfidenceError returns an error describing a prediction that lacks a confidence level.
func NewNoConfidenceError(s Stream, i int) error {
	return makePredictionErrorMaker(
		"has no confidence level specified",
	)(s, i)
}

// NewNoClaimError returns an error that describes the approximate location of a prediction that has no claim.
func NewNoClaimError(s Stream, i int) error {
	prefix := ""
	if s.FromFilename != "" {
		prefix = s.FromFilename + ": "
	}

	previousClaim := ""
	if i > 0 {
		previousClaim = s.Predictions[i-1].Claim
	}

	switch {
	case i == 0:
		return fmt.Errorf(prefix + "first prediction has no claim")
	case previousClaim != "":
		return fmt.Errorf(prefix+"claim after “%v” has no claim", previousClaim)
	default:
		return fmt.Errorf(prefix + "prediction exists that has no claim, and neither does the one before it")
	}
}
