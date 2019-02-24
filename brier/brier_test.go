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

package brier

import (
	"math"
	"strings"
	"testing"

	"github.com/adiabatic/predictions/streams"
)

func mustStreams(s string) []streams.Stream {
	aStream, err := streams.FromReader(strings.NewReader(s))
	if err != nil {
		panic("FromReader should work")
	}
	ret := make([]streams.Stream, 0, 1)
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
scope: today
---
claim: I will fly to Pluto and back
confidence: 100
happened: false
---
claim: I will jump 300 feet in the air
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

const howGoodAmIAtKnowingMyOwnCooking = `
scope: today
---
claim: I will like the dinner I just made
confidence: 70
happened: true
tags: [food]
---
claim: I will like the new Brian Eno album I just bought
confidence: 99
happened: false
tags: [music]
`

func TestMatchingTag(t *testing.T) {
	ss := mustStreams(howGoodAmIAtKnowingMyOwnCooking)

	const expectedScore = 0.09

	actualScore := ForOnly(ss, MatchingTag("food"))
	if !closeEnough(actualScore, expectedScore) {
		t.Fatalf("expected score: %v; got: %v", expectedScore, actualScore)
	}

}
