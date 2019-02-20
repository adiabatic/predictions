# Errors and warnings

<!-- markdownlint-disable MD038 -->

## cannot unmarshal !!str `…` into float64

Whatever you wrote as a confidence isn’t recognized as such. Confidence levels need to be written as a number between 0 and 100, without the percent sign.

## [error.claim.missing]

A prediction doesn’t have a claim in it. Claims start with `claim: `.

## [error.confidence.missing]

A prediction doesn’t have a confidence level in it. Confidence levels start with `confidence: `.

## [error.confidence.impossible]

…

## [warn.confidence.zero]

…

## [warn.confidence.unity]

…
(“unity” is a fancy name for “one”.)
