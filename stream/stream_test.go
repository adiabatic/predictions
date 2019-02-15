package stream

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mustStreamFromString(t *testing.T, s string) Stream {
	const msg = "FromReader should at least work"

	st, err := FromReader(strings.NewReader(s))
	if err != nil {
		if t == nil {
			panic(err)
		}
		t.Fatalf(err.Error())
	}
	return st
}

// YAML strings must be at the top level of indentation. goimports will indent raw-string blocks in functions, adding tabs to most lines inside the string that we cannot handle.

const simpleStream = `---
title: test stream
scope: in the year 2000
---
claim: I will park my car straight at the gym today
confidence: 60
`

// This is called table-driven testing, isn’t it?

type KeyValueTable []KeyValueRow

type KeyValueRow struct {
	Key string

	ExpectedValue string
	ActualValue   string
}

func TestMetadata(t *testing.T) {
	s := mustStreamFromString(t, simpleStream)

	table := []KeyValueRow{
		{
			Key:           "title",
			ExpectedValue: "test stream",
			ActualValue:   s.Metadata.Title,
		},
		{
			Key:           "scope",
			ExpectedValue: "in the year 2000",
			ActualValue:   s.Metadata.Scope,
		},
	}

	for _, row := range table {
		name := fmt.Sprintf("Ensure that %v is filled correctly", row.Key)
		t.Run(name, func(t *testing.T) {
			if row.ExpectedValue != row.ActualValue {
				t.Errorf("Metadata mismatch. Expected “%v”; got “%v”",
					row.ExpectedValue,
					row.ActualValue,
				)
			}
		})
	}
}

const nullStream = ``

func TestNullStream(t *testing.T) {
	expectedError := NeitherTitleNorScopeInMetadataBlock
	_, actualError := FromReader(strings.NewReader(nullStream))
	if actualError != expectedError {
		t.Fatalf("unexpected error from a null stream. expected: %v; got: %v", expectedError, actualError)
	}
}

const nonMapMetadata = `---
title yum
scope: in the year 2000`

func TestMetadataIsAMap(t *testing.T) {
	_, actualError := FromReader(strings.NewReader(nonMapMetadata))
	if actualError == nil {
		t.Fatal("a metadata document should fail if it’s not a proper mapping")
	}
}

const missingClaimsAndConfidences = `---
title: Comestibles prognostication
scope: this week
---
claim: I will eat a steak
---
claim: I will eat a chocolate bar
confidence: 70
---
claim: I will eat a dinner salad
confidence: null
---
# missing both claim and confidence (“claims” doesn’t count)
claims: [I will eat ice cream, I will eat peanut-butter cups]
---
# missing everything, and its predecessor is missing a “claim”
`

func Example_missingClaimsAndConfidences() {
	s := mustStreamFromString(nil, missingClaimsAndConfidences)
	var sv Validator
	errs := sv.RunValidationFunctions(s,
		sv.AllPredictionsHaveConfidences,
	)
	for _, err := range errs {
		fmt.Println(err)
	}

	// Unordered output:
	// not it
}

const missingClaims = `---
title: My taste buds
scope: 2019
notes: > 
  Only the second prediction has a “claim” key.
  None of the others do.
---
claims: [I will like fish, I will like shellfish]
confidence: 5
---
claim: I will like red meat
confidence: 99
---
claims: [I will like parsnips, I will like turnips]
confidence: 30
---
claims: [I will like hoppy beer, I will like lamb]
confidence: 20
claims: [I like water, I like food]
confidence: 32
`

func TestMissingClaims(t *testing.T) {
	s := mustStreamFromString(nil, missingClaims)
	var sv Validator
	errs := sv.RunAll(s)

	assert.Equal(t, len(errs), 3)
	expecteds := []string{
		"first prediction has no claim",
		"claim after “I will like red meat” has no claim",
		"prediction exists that has no claim, and neither does the one before it",
	}
	for i, expected := range expecteds {
		assert.Equal(t, expected, errs[i].Error())
	}

}

const questionableConfidences = `
title: I am the universe
---
claim: green is spiky
confidence: -20
---
claim: my left arm will turn into a tentacle
confidence: 0
---
claim: the sun will rise tomorrow
confidence: 100
---
claim: I will marry my middle-school crush
confidence: 245
`

func TestQuestionableConfidences(t *testing.T) {
	s := mustStreamFromString(t, questionableConfidences)
	var sv Validator
	errs := sv.RunValidationFunctions(s,
		sv.AllConfidencesBetweenZeroAndOneHundredExclusive,
	)

	assert.Len(t, errs, 4)
	expecteds := []string{
		"first prediction, with claim “green is spiky”, has a confidence level outside (0%, 100%)",
		"prediction with claim “my left arm will turn into a tentacle” has a confidence level outside (0%, 100%)",
		"prediction with claim “the sun will rise tomorrow” has a confidence level outside (0%, 100%)",
		"prediction with claim “I will marry my middle-school crush” has a confidence level outside (0%, 100%)",
	}
	for i, expected := range expecteds {
		assert.Equal(t, expected, errs[i].Error())
	}

}
