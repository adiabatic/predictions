# `predictions`’ errors and warnings

<!-- markdownlint-disable MD038 -->

## cannot unmarshal !!str `…` into float64

Whatever you wrote, probably a confidence level, isn’t recognized as such. Confidence levels need to be written as a number between 0 and 100, without the percent sign.

## [error.metadata.unexpected-claim]

TODO: reduce jargon

The first document in a stream is supposed to contain metadata about the predictions that follow. It is not supposed to contain things only predictions have, like a claim. This error is displayed when `claim: ` occurs in a document that `predictions` expects to contain only metadata, like `title: ` and `scope: `.

## [error.metadata-unexpected-confidence]

TODO: reduce jargon

The first document in a stream is supposed to contain metadata about the predictions that follow. It is not supposed to contain things only predictions have, like a confidence level. This error is displayed when `confidence: ` occurs in a document that `predictions` expects to contain only metadata, like `title: ` and `scope: `.

## [error.claim.missing]

A prediction doesn’t have a claim in it. Claims start with `claim: `.

## [error.confidence.missing]

A prediction doesn’t have a confidence level in it. Confidence levels start with `confidence: `.

## [error.confidence.impossible]

Confidence levels need to be written as a number between 0 and 100, corresponding to confidence levels of 0% and 100%. While fractional confidence levels are permissible (if unwise), negative numbers and numbers over 100 make no sense.

## [warn.confidence.zero]

It’s a bad idea to claim that something has a 0% chance of happening. [Infinite Certainty][ic] explains why.

## [warn.confidence.unity]

It’s a bad idea to claim that something has a 100% chance of happening. [Infinite Certainty][ic] explains why.

[ic]: https://www.readthesequences.com/Infinite-Certainty
