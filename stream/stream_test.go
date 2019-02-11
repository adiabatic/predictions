package stream

import (
	"fmt"
	"strings"
	"testing"
)

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
	r := strings.NewReader(simpleStream)

	s, err := FromReader(r)
	if err != nil {
		t.Fatalf("unexpected error in FromReader: %v", err)
	}

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

const missingClaimsAndPredictions = `---
title: Comestibles prognostication
---
claim: I will eat a chocolate bar this week
confidence: 70
---
claim: I will eat a steak this week
# totally missing confidence
---
claim: I will eat a dinner salad this week
confidence: null
---
# missing both claim and confidence (“claims” doesn’t count)
claims: [I will eat ice cream this week, I will eat peanut-butter cups this week]
---
# missing everything, and its predecessor is missing a “claim”
`

func Example_missingClaimsAndConfidences() {
	s, err := FromReader(strings.NewReader(missingClaimsAndPredictions))
	if err != nil {
		panic("FromReader should at least work")
	}
	var sv Validator
	errs := sv.RunAll(s)
	for _, err := range errs {
		fmt.Println(err)
	}

	// Unordered output:
	// Prediction after “I will eat a dinner salad this week” has no claim in it
	// A prediction has no claim in it. Either it’s the first prediction or the prediction before it doesn’t have a claim in it, either
	// Prediction with claim “I will eat a steak this week” has no declared confidence
	// Prediction with claim “I will eat a dinner salad this week” has no declared confidence
	// Prediction after prediction with claim “I will eat a dinner salad this week” has no declared confidence
	// A prediction exists that lacks both a confidence and a claim, and its predecessor lacks a claim too
}
