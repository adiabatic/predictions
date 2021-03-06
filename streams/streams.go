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

// Package streams provides a way to access predictions in a YAML stream.
//
// YAML streams contain one or more documents in them, separated by three hyphen-minuses (---). In streams that we’re using, the first document functions as a metadata header while all subsequent documents are predictions. Each prediction, with few exceptions, has both a claim and a confidence.
package streams

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/xtgo/set"
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

	// These are here to detect when a user accidentally omits a metadata document in a stream.
	MisplacedClaim      string `yaml:"claim"`
	MisplacedConfidence string `yaml:"confidence"`
}

// A PredictionDocument contains a claim, the claim’s confidence, and so on.
type PredictionDocument struct {
	Claim             string
	Confidence        *float64
	Tags              []string
	Happened          *bool
	CauseForExclusion string `yaml:"cause for exclusion"`
	Hash              bool
	Salt              string
	Notes             string

	Parent *Stream
}

// ShouldExclude returns true if the receiver should be excluded from stats calculation.
func (d *PredictionDocument) ShouldExclude() bool {
	if d == nil {
		return false
	}

	if d.Confidence == nil || d.Happened == nil || d.CauseForExclusion != "" {
		return true
	}

	return false
}

// HasTag returns true if the given tag is in the receiver’s tag list.
func (d *PredictionDocument) HasTag(tag string) bool {
	if d == nil {
		return false
	}

	for _, t := range d.Tags {
		if t == tag {
			return true
		}
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

	filenamePrefix := ""
	if filename != "" {
		filenamePrefix = filename + ": "
	}

	if md.MisplacedClaim != "" {
		return Stream{}, fmt.Errorf("%s[error.metadata.unexpected-claim]: claim of “%s” in first (metadata) document",
			filenamePrefix, md.MisplacedClaim)
	}

	if md.MisplacedConfidence != "" {
		return Stream{}, fmt.Errorf("%s[error.metadata.unexpected-confidence]: confidence of “%s” in first (metadata) document", filenamePrefix, md.MisplacedConfidence)
	}

	s.Metadata = md

	for {
		var pd PredictionDocument
		err = dec.Decode(&pd)
		if err != nil {
			break
		}
		pd.Parent = &s
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

// FromFiles generates a slice of Stream from the filenames specified.
func FromFiles(filenames []string) ([]Stream, error) {
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

// RunMinimal is a convenience function to run the bare minimum of sanity checks.
func (sv *Validator) RunMinimal(s Stream) []error {
	return sv.RunValidationFunctions(s,
		sv.AllPredictionsHaveClaims,
		sv.AllPredictionsHaveConfidences,
	)
}

// RunAll is a convenience function to run all known stream validators.
func (sv *Validator) RunAll(s Stream) []error {
	return sv.RunValidationFunctions(s,
		sv.HasTitleOrScopeInMetadataBlock,
		sv.AllPredictionsHaveClaims,
		sv.AllPredictionsHaveConfidences,
		sv.AllConfidencesSensible,
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
			errs = append(errs, NewErrorClaimMissing(s, i))
		}
	}

	return errs
}

// AllPredictionsHaveConfidences ensures that all predictions have a confidence key and a value of some sort.
func (sv *Validator) AllPredictionsHaveConfidences(s Stream) []error {
	errs := make([]error, 0)
	for i, pred := range s.Predictions {
		if pred.Confidence == nil {
			errs = append(errs, NewErrorConfidenceMissing(s, i))
		}
	}
	return errs
}

// AllConfidencesBetweenZeroAndOneHundredInclusive ensures all confidences are on [0, 100].
func (sv *Validator) AllConfidencesBetweenZeroAndOneHundredInclusive(s Stream) []error {
	errs := make([]error, 0)
	for i, pred := range s.Predictions {
		if pred.Confidence != nil &&
			*(pred.Confidence) < 0.0 || *(pred.Confidence) > 100.0 {
			errs = append(errs, NewErrorConfidenceImpossible(s, i))
		}
	}
	return errs
}

// AllConfidencesSensible ensures no confidence level is 0 or 100.
func (sv *Validator) AllConfidencesSensible(s Stream) []error {
	errs := make([]error, 0)
	for i, pred := range s.Predictions {
		if pred.Confidence != nil {
			if *(pred.Confidence) == 0.0 {
				errs = append(errs, NewErrorConfidenceZero(s, i))
			} else if *(pred.Confidence) == 100.0 {
				errs = append(errs, NewErrorConfidenceUnity(s, i))
			}
		}
	}
	return errs
}

// other stuff

func deduplicateStrings(ss []string) []string {
	seen := make(map[string]struct{}, len(ss))
	j := 0
	for _, v := range ss {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		ss[j] = v
		j++
	}
	return ss[:j]
}

// TagsUsed returns a list of all tags used in the given slice of stream.Stream.
func TagsUsed(sts []Stream) []string {
	ret := make([]string, 0)
	for _, s := range sts {
		for _, p := range s.Predictions {
			for _, t := range p.Tags {
				ret = append(ret, t)
			}
		}
	}
	return deduplicateStrings(ret)
}

// KeysUsed returns a list of all keys used in the given Streams.
//
// A “key”, here, is the title and scope of a given stream, with a space in between.
func KeysUsed(sts []Stream) []string {
	ret := make([]string, 0)
	for _, s := range sts {
		ret = append(ret, s.Metadata.Title+" "+s.Metadata.Scope)
	}
	return deduplicateStrings(ret)

}

// ConfidencesUsed returns a list of all confidence intervals used in the given Streams.
func ConfidencesUsed(sts []Stream) []float64 {
	ret := make([]float64, 0)
	for _, s := range sts {
		for _, pred := range s.Predictions {
			if !pred.ShouldExclude() {
				ret = append(ret, *pred.Confidence)
			}
		}
	}
	return set.Float64s(ret)
}
