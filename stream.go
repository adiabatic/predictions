package main

import (
	"io"
	"os"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// A StreamValidationError is returned when something about the stream isn’t right.
type StreamValidationError string

func (e StreamValidationError) Error() string { return string(e) }

// Various errors returned when decoding documents out of a stream and into structs.
const (
	NeitherTitleNorScopeInMetadataBlock = StreamValidationError("Neither title nor scope in metadata block")
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

// StreamFromReader decodes into a Stream from an io.Reader.
func StreamFromReader(r io.Reader) (Stream, error) {
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

		s, err := StreamFromReader(f)
		if err != nil {
			return nil, errors.WithMessagef(err, "couldn’t make stream from file “%v”", fn)
		}

		streams = append(streams, s)
	}

	return streams, nil
}

// A StreamValidationFunction ensures that a Stream passes a sanity check.
type StreamValidationFunction func(Stream) error

// A StreamValidator contains data useful for StreamValidationFunctions.
type StreamValidator struct{}

// HasTitleOrScopeInMetadataBlock ensures that a stream has either title key or a scope key in the metadata block (or both). At least one of those keys’ values must be something other than the empty string.
func (sv *StreamValidator) HasTitleOrScopeInMetadataBlock(s Stream) error {
	if s.Metadata.Title == "" || s.Metadata.Scope == "" {
		return NeitherTitleNorScopeInMetadataBlock
	}
	return nil
}
