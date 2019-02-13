package brier

import (
	"math"
	"strings"
	"testing"

	"github.com/adiabatic/predictions/stream"
)

func mustStreams(s string) []stream.Stream {
	aStream, err := stream.FromReader(strings.NewReader(s))
	if err != nil {
		panic("FromReader should work")
	}
	ret := make([]stream.Stream, 0, 1)
	ret = append(ret, aStream)

	return ret
}

const ε = 0.000001

func closeEnough(l, r float64) bool {
	return math.Abs(l-r) < ε
}

const nailedIt = `
title: Everything I say will happen, happens
---
claim: I will wake up this morning
confidence: 100
happened: true
---
claim: I will breathe today
confidence: 100
happened: true
`

func TestNailedIt(t *testing.T) {
	ss := mustStreams(nailedIt)

	const expectedScore = 0

	actualScore := ForOnly(ss, Everything)

	if !closeEnough(actualScore, expectedScore) {
		t.Fatalf("expected score: %v; got: %v", expectedScore, actualScore)
	}

}

const wrongAboutEverything = `
title: I am maximally wrong about everything
---
claim: I will fly to Pluto and back today
confidence: 100
happened: false
---
claim: I will jump 300 feet in the air today
confidence: 100
happened: false
`

func TestWrongAboutEverything(t *testing.T) {
	ss := mustStreams(wrongAboutEverything)

	const expectedScore = 1

	actualScore := ForOnly(ss, Everything)

	if !closeEnough(actualScore, expectedScore) {
		t.Fatalf("expected score: %v; got: %v", expectedScore, actualScore)
	}

}

const kindaThoughtItWouldRainAndItDid = `
scope: today
---
claim: it will rain
confidence: 70
happened: true
`

func TestKindaThoughtItWouldRainAndItDid(t *testing.T) {
	ss := mustStreams(kindaThoughtItWouldRainAndItDid)
	const expectedScore = 0.09

	actualScore := ForOnly(ss, Everything)

	if !closeEnough(actualScore, expectedScore) {
		t.Fatalf("expected score: %v; got: %v", expectedScore, actualScore)
	}
}

const kindaThoughtItWouldNotRainButItDidAnyway = `
scope: today
---
claim: it will rain
confidence: 30
happened: true
`

func TestKindaThoughtItWouldNotRainButItDidAnyway(t *testing.T) {
	ss := mustStreams(kindaThoughtItWouldNotRainButItDidAnyway)

	const expectedScore = 0.49

	actualScore := ForOnly(ss, Everything)

	if !closeEnough(actualScore, expectedScore) {
		t.Fatalf("expected score: %v; got: %v", expectedScore, actualScore)
	}
}
