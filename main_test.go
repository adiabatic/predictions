package main

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

	s, err := StreamFromReader(r)
	if err != nil {
		t.Fatalf("unexpected error in StreamFromReader: %v", err)
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
	_, actualError := StreamFromReader(strings.NewReader(nullStream))
	if actualError != expectedError {
		t.Fatalf("unexpected error from a null stream. expected: %v; got: %v", expectedError, actualError)
	}

}
