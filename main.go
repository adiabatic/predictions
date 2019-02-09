package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/davecgh/go-spew/spew"
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

func main() {
	// no flags to parse yet, but we need to do this to make flag.Args() work

	flag.Parse()

	streams := make([]Stream, 0, 1)

	for _, fn := range flag.Args() {
		f, err := os.Open(fn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "couldn’t open file named “%v”: %v\n", fn, err)
			os.Exit(1)
		}
		defer f.Close()

		dec := yaml.NewDecoder(f)
		var s Stream
		var md MetadataDocument
		var pds []PredictionDocument

		err = dec.Decode(&md)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error while decoding metadata document of “%v”: %v\n", fn, err)
			os.Exit(2)
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
			if len(pds) >= 1 {
				knownGoodPrediction := pds[len(pds)-1]
				fmt.Fprintf(os.Stderr,
					"error reading the prediction after the one with the claim “%v”: %v\n",
					knownGoodPrediction.Claim,
					err,
				)
			} else {
				fmt.Fprintln(os.Stderr, "error reading the first prediction: ", err)
			}
			os.Exit(3)
		}

		s.Predictions = pds
		streams = append(streams, s)
	}

	fmt.Println("the streams:")
	spew.Dump(streams)
}
