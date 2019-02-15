package stream

import (
	"errors"
	"fmt"
)

type PredictionErrorMaker func([]PredictionDocument, int) error

func NewInsensibleConfidenceError(ps []PredictionDocument, i int) error {
	f := g(
		"Prediction with claim “%v” has a confidence level outside (0%%, 100%%)",
		"Prediction after prediction with claim “%v” has a confidence level outside (0%%, 100%%)",
		"A prediction exists with a confidence level outside (0%%, 100%%). It doesn’t have a claim and its immediate predecessor, if any, lacks a claim too",
	)

	return f(ps, i)

}

func g(at, atPrev, huh string) PredictionErrorMaker {
	return func(predictions []PredictionDocument, i int) error {
		if claim := predictions[i].Claim; claim != "" {
			return fmt.Errorf(at, claim)
		} else if i > 0 {
			if previousClaim := predictions[i-1].Claim; previousClaim != "" {
				return fmt.Errorf(atPrev, previousClaim)
			}
		}
		return errors.New(huh)
	}
}
