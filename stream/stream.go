package stream

import (
	"io"
	"os"
	"strings"

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
	Metadata    MetadataDocument
	Predictions []PredictionDocument
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
	Claim      string
	Confidence float64
	Tags       []string
	Happened   *bool
	Hash       bool
	Salt       string
	Notes      string
}

// FromReader decodes into a Stream from an io.Reader.
func FromReader(r io.Reader) (Stream, error) {
	dec := yaml.NewDecoder(r)
	var s Stream
	var md MetadataDocument
	var pds []PredictionDocument

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

// StreamsFromFiles generates a slice of Stream from the filenames specified.
func StreamsFromFiles(filenames []string) ([]Stream, error) {
	streams := make([]Stream, 0, 1)

	for _, fn := range filenames {
		f, err := os.Open(fn)
		if err != nil {
			return nil, errors.WithMessagef(err, "couldn’t open file named “%v”", fn)
		}
		defer f.Close()

		s, err := FromReader(f)
		if err != nil {
			return nil, errors.WithMessagef(err, "couldn’t make stream from file “%v”", fn)
		}

		streams = append(streams, s)
	}

	return streams, nil
}

// A ValidationFunction ensures that a Stream passes a sanity check.
type ValidationFunction func(Stream) error

// A Validator contains data useful for ValidationFunctions.
type Validator struct{}

// HasTitleOrScopeInMetadataBlock ensures that a stream has either title key or a scope key in the metadata block (or both). At least one of those keys’ values must be something other than the empty string.
func (sv *Validator) HasTitleOrScopeInMetadataBlock(s Stream) error {
	if s.Metadata.Title == "" || s.Metadata.Scope == "" {
		return NeitherTitleNorScopeInMetadataBlock
	}
	return nil
}

// A NoClaimError is returned when a prediction has no claim in it.
type NoClaimError struct {
	PreviousClaim string
}

func (e NoClaimError) Error() string {
	if e.PreviousClaim != "" {
		return "Prediction after “" + e.PreviousClaim + "” has no claim in it"
	}
	return "A prediction has no claim in it (possibly the first)"
}

// A NoClaimErrors is a slice of NoClaimError. Some streams will miss a bunch of claims.
type NoClaimErrors []NoClaimError

func (es NoClaimErrors) Error() string {
	sb := strings.Builder{}
	sb.WriteString("No-claim errors: ")

	for _, e := range es {
		sb.WriteString(e.Error())
		sb.WriteString(", ")
	}

	return sb.String()
}

// AllPredictionsHaveClaims ensures that all predictions in a stream have one claim in each.
func (sv *Validator) AllPredictionsHaveClaims(s Stream) error {
	errs := make([]NoClaimError, 0)

	for i, prediction := range s.Predictions {
		if prediction.Claim == "" {
			if i == 0 {
				errs = append(errs, NoClaimError{})
			} else {
				errs = append(errs, NoClaimError{s.Predictions[i-1].Claim})
			}
		}
	}

	return NoClaimErrors(errs)
}

// A NoConfidenceError is returned when one or more predictions in a stream doesn’t have an associated confidence level.
type NoConfidenceError struct {
	Claims []string
}

func (e NoConfidenceError) Error() string { return "Some claims lack confidence" }

// AllPredictionsHaveConfidences ensures that all predictions have a confidence key and a value of some sort.
func (sv *Validator) AllPredictionsHaveConfidences(s Stream) error {
	return nil
}
