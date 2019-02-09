package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
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

// StreamsFromFiles generates a slice of Stream from the filenames specified.
func StreamsFromFiles(filenames []string) ([]Stream, error) {
	streams := make([]Stream, 0, 1)

	for _, fn := range filenames {
		f, err := os.Open(fn)
		if err != nil {
			return nil, errors.WithMessagef(err, "couldn’t open file named “%v”", fn)
		}
		defer f.Close()

		dec := yaml.NewDecoder(f)
		var s Stream
		var md MetadataDocument
		var pds []PredictionDocument

		err = dec.Decode(&md)
		if err != nil {
			return nil, errors.WithMessagef(err, "error while decoding metadata document in file “%v”", fn)
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
				return nil, errors.WithMessagef(err, "error reading the first prediction")
			}

			knownGoodPrediction := pds[len(pds)-1]
			return nil, errors.WithMessagef(err,
				"error reading the prediction after the one with the following claim: “%v”",
				knownGoodPrediction.Claim,
			)
		}

		s.Predictions = pds
		streams = append(streams, s)
	}

	return streams, nil
}

func main() {
	// no flags to parse yet, but we need to do this to make flag.Args() work

	flag.Parse()

	streams, err := StreamsFromFiles(flag.Args())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("the streams:")
	spew.Dump(streams)
}
