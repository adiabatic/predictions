// Package stream provides a way to access predictions in a YAML stream.
//
// YAML streams contain one or more documents in them, separated by three hyphen-minuses (---). In streams that we’re using, the first document functions as a metadata header while all subsequent documents are predictions. Each prediction, with few exceptions, has both a claim and a confidence.
package stream

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// A ValidationError is returned when something about the stream isn’t right.
type ValidationError string

func (e ValidationError) Error() string { return string(e) }

// Various errors returned when decoding documents out of a stream and into structs.
const (
	NeitherTitleNorScopeInMetadataBlock = ValidationError("Neither title nor scope in metadata block")
)

// A Stream (YAML stream) contains a metadata document and predictions documents.
//
// In YAML, documents are prefixed by “---”.
type Stream struct {
	FromFilename string
	Metadata     MetadataDocument
	Predictions  []PredictionDocument
}

// A MetadataDocument contains information about the predictions in its Stream.
type MetadataDocument struct {
	Title string
	Scope string
	Salt  string
	Notes string
}

// A PredictionDocument contains a claim, the claim’s confidence, and so on.
type PredictionDocument struct {
	Claim             string
	Confidence        float64
	Tags              []string
	Happened          *bool
	CauseForExclusion string `yaml:"cause for exclusion"`
	Hash              bool
	Salt              string
	Notes             string
}

// ShouldExclude returns true if the receiver should be excluded from consideration.
func (d *PredictionDocument) ShouldExclude() bool {
	if d.Happened == nil || d.CauseForExclusion != "" {
		return true
	}

	return false
}

func fromReaderWithFilename(r io.Reader, filename string) (Stream, error) {
	dec := yaml.NewDecoder(r)
	var s Stream
	var md MetadataDocument
	var pds []PredictionDocument

	s.FromFilename = filename

	err := dec.Decode(&md)
	if err != nil {
		if err == io.EOF {
			return Stream{}, NeitherTitleNorScopeInMetadataBlock
		}
		return Stream{}, errors.WithMessage(err, "error while decoding metadata document")
	}

	s.Metadata = md

	for {
		var pd PredictionDocument
		err = dec.Decode(&pd)
		if err != nil {
			break
		}
		pds = append(pds, pd)
	}
	if err != io.EOF {
		if len(pds) == 0 {
			return Stream{}, errors.WithMessagef(err, "error reading the first prediction")
		}

		knownGoodPrediction := pds[len(pds)-1]
		return Stream{}, errors.WithMessagef(err,
			"error reading the prediction after the one with the following claim: “%v”",
			knownGoodPrediction.Claim,
		)
	}

	s.Predictions = pds

	return s, nil
}

// FromReader decodes into a Stream from an io.Reader.
func FromReader(r io.Reader) (Stream, error) {
	return fromReaderWithFilename(r, "")
}

// StreamsFromFiles generates a slice of Stream from the filenames specified.
func StreamsFromFiles(filenames []string) ([]Stream, error) {
	streams := make([]Stream, 0, 1)

	for _, fn := range filenames {
		f, err := os.Open(fn)
		if err != nil {
			return nil, errors.WithMessagef(err, "couldn’t open file named “%v”", fn)
		}
		defer f.Close()

		s, err := fromReaderWithFilename(f, fn)
		if err != nil {
			return nil, errors.WithMessagef(err, "couldn’t make stream from file “%v”", fn)
		}

		streams = append(streams, s)
	}

	return streams, nil
}

// A ValidationFunction ensures that a Stream passes a sanity check.
//
// Because many things can go wrong in a Stream that a user would want to know about all at once, ValidationFunction returns a slice of error.
type ValidationFunction func(Stream) []error

// A Validator contains data useful for ValidationFunctions.
type Validator struct{}

// RunValidationFunctions runs a bunch of validation functions on a stream.
//
// It returns a slice of error of everything that didn’t pass muster.
func (sv *Validator) RunValidationFunctions(s Stream, vfs ...ValidationFunction) []error {
	errs := make([]error, 0)
	for _, f := range vfs {
		ferrs := f(s)
		errs = append(errs, ferrs...)
	}

	return errs
}

// RunAll is a convenience function to run all known stream validators.
func (sv *Validator) RunAll(s Stream) []error {
	return sv.RunValidationFunctions(s,
		sv.HasTitleOrScopeInMetadataBlock,
		sv.AllPredictionsHaveClaims,
		sv.AllPredictionsHaveConfidences,
		sv.AllConfidencesBetweenZeroAndOneHundredExclusive,
	)
}

// HasTitleOrScopeInMetadataBlock ensures that a stream has either title key or a scope key in the metadata block (or both). At least one of those keys’ values must be something other than the empty string.
func (sv *Validator) HasTitleOrScopeInMetadataBlock(s Stream) []error {
	errs := make([]error, 0)
	if s.Metadata.Title == "" && s.Metadata.Scope == "" {
		return append(errs, NeitherTitleNorScopeInMetadataBlock)
	}
	return errs
}

// AllPredictionsHaveClaims ensures that all predictions in a stream have one claim in each.
func (sv *Validator) AllPredictionsHaveClaims(s Stream) []error {
	errs := make([]error, 0)

	for i, prediction := range s.Predictions {
		if prediction.Claim == "" {
			errs = append(errs, NewNoClaimError(s, i))
		}
	}

	return errs
}

// AllPredictionsHaveConfidences ensures that all predictions have a confidence key and a value of some sort.
func (sv *Validator) AllPredictionsHaveConfidences(s Stream) []error {
	errs := make([]error, 0)
	for i, pred := range s.Predictions {
		if pred.Confidence == 0.0 {
			errs = append(errs, NewNoConfidenceError(s, i))
		}
	}
	return errs
}

// A ConfidenceOutOfRange is returned when one or more predictions in a stream has a confidence that is too high or too low.
type ConfidenceOutOfRange struct {
	Claim         string
	PreviousClaim string
}

// NewConfidenceOutOfRange returns a reasonable error for the location it’s found in.
func NewConfidenceOutOfRange(predictions []PredictionDocument, i int) ConfidenceOutOfRange {
	if predictions[i].Claim != "" {
		return ConfidenceOutOfRange{
			Claim: predictions[i].Claim,
		}
	} else if i > 0 && predictions[i-1].Claim != "" {
		return ConfidenceOutOfRange{
			PreviousClaim: predictions[i-1].Claim,
		}
	}
	return ConfidenceOutOfRange{}
}

func (e ConfidenceOutOfRange) Error() string {
	if e.Claim != "" {
		return fmt.Sprintf("Prediction with claim “%v” has a too-weird confidence level", e.Claim)
	} else if e.PreviousClaim != "" {
		return fmt.Sprintf("Prediction after prediction with claim “%v” has a too-weird confidence level", e.PreviousClaim)
	}
	return "A prediction exists with a too-weird confidence level, it also doesn’t have a claim, and its predecessor lacks a claim too"
}

// AllConfidencesBetweenZeroAndOneHundredInclusive ensures all confidences are on [0, 100].
func (sv *Validator) AllConfidencesBetweenZeroAndOneHundredInclusive(s Stream) []error {
	errs := make([]error, 0)
	for i, pred := range s.Predictions {
		if pred.Confidence < 0.0 || pred.Confidence > 100.0 {
			errs = append(errs, NewConfidenceOutOfRange(s.Predictions, i))
		}
	}
	return errs
}

// AllConfidencesBetweenZeroAndOneHundredExclusive ensures all confidences are on (0, 100).
func (sv *Validator) AllConfidencesBetweenZeroAndOneHundredExclusive(s Stream) []error {
	errs := make([]error, 0)
	for i, pred := range s.Predictions {
		if pred.Confidence <= 0.0 || pred.Confidence >= 100.0 {
			errs = append(errs, NewInsensibleConfidenceError(s, i))
		}
	}
	return errs
}
