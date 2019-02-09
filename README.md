# Predictions

`predictions` helps you figure out how how good you are at predicting things. You feed it a [YAML][] file full of predictions and the outcomes and it calculates how good you are at predicting things.

Not sure what this whole ‘predictions’ thing is? Have a look at [Scott Alexander’s 2018 predictions retrospective][ssc2018] to get an idea of what’s involved.

[yaml]: https://en.wikipedia.org/wiki/YAML "YAML Ain’t Markup Language"
[ssc2018]: https://slatestarcodex.com/2019/01/22/2018-predictions-calibration-results/

## Usage

`predictions [options] inputfile`

`predictions` always outputs to standard output, with errors and warnings printed to standard error. To get the output into a file, redirect standard output to a file with `>`, as in `predictions 2019.yaml > 2019.markdown`.

## The YAML file format

`predictions` takes in a multidocument YAML file. The first document specifies metadata. Each subsequent document contains one top-level mapping that specifies a prediction.

An example file:

```yaml
---
title: Predictions for 2019
scope: for 2019
salt: kkjskvjsdwolvkjsjv
totally private notes: >
  `predictions` doesn’t complain if you
  put in mapping keys that it doesn’t expect.
  Because of this, you can write notes
  about your predictions and these notes won’t
  make it into `predictions`’ output.
---
claim: I will not buy any socks
confidence: 80
# “happened” isn’t listed here because I might
# buy socks between now and the end of the year.
tags: [commerce]
---
claim: Dennis will not change jobs
confidence: 90
tags: [friends]
happened: no
hash: yes
notes: Left Acme on January 14. Go figure.
```

### Metadata-document mapping keys

#### `title`

The title of the file full of predictions.

#### `salt` (per-document) (not yet implemented)

A per-file salt used for hashing sensitive predictions.

See “Salt-and-hash rationale” below for why you might want to do this.

#### `scope` (not yet implemented)

This restricts the scope of predictions to a particular domain. The idea of scopes is that they’re one-per-file so one can combine, say, 2018 predictions, 2019 predictions, and 2020 predictions in an invocation of `predictions` and see combined results with each year’s prediction labeled as such.

#### `tag order` (not yet implemented)

A list showing the order you want tags displayed in.

Use this if you want to force an order in your output. Suppose you have three tags of predictions: one about U.S. politics, another for international politics, and a third about bananas. If you want to ensure that your output displays these three topics in that order (as opposed to boring your readers at the beginning with your heady pronouncements about bananas), then you should have `tag order: [U.S. politics, International politics, Bananas]` in your metadata document.

### Prediction-document mapping keys

#### `claim` (required)

A prediction that something will, or won’t, happen.

#### `confidence` (warns if missing)

How confident you are that this will happen, expressed as a percentage without the percent sign.

#### `tags`

A list of tags that you want to associate with a claim.

Tags are used for grouping purposes. For example, you might be very well-calibrated when it comes to predicting national affairs and international happenings, but wildly off-base when predicting what your friends will do. If you tag each type of prediction with its type, you’ll be able to find out if you’re better in some areas than others.

If some, but not all, of your predictions have tags, the untagged ones will be tagged “Untagged”.

#### `happened`

Either a boolean (yes/no/true/false) or null (~). Use yes/true for something that did happen, no/false for something that definitely didn’t happen, and ~ for weird ambiguous results that you want to throw out of any analysis.

#### `hash` (not yet implemented)

A boolean. If true, then this entry is hashed before going to a publicly-displayed output.

Uses a per-prediction salt if it exists. If it doesn’t, then it’ll fall back to the whole-file salt specified in the metadata document. If that doesn’t exist either, then the claim will be unsalted when published.

#### `salt` (per-prediction) (not yet implemented)

A per-prediction salt used for hashing sensitive predictions. If `salt` is specified in a prediction, then `hash` is implied.

See “salt-and-hash rationale” below for why you might want to do this.

#### `notes`

The contents of `notes` may be put into the output. The literal

## Command-line flags

### `-ot` (output type) (unimplemented)

Sets the output type. Options:

- `markdownsnippet`, useful for pasting into a blog post [like Scott Alexander’s][ssc2018]. Currently the default.
- `html`, useful for seeing your predictions, nicely formatted, in a standalone HTML page.

### `-Wall` (all warnings)

Emits all warnings such as “your confidence level for \_\_\_\_\_\_ was 100%” (see [Infinite Certainty][ic] for why this is a bad thing).

[ic]: https://www.readthesequences.com/Infinite-Certainty

## Salt-and-hash rationale

Sometimes you want to make predictions in public and display your score in public, but don’t want to publicly display what your prediction is. By salting and hashing your predictions, you make it nigh impossible for people to figure out what you’re predicting. Furthermore, if a given prediction has its own salt, then later, if you so choose, you’ll be able to reveal that particular claim at a later date without letting people guess what your other predictions are.

Have a look at the example file above. There’s a whole-file salt in the first document and the “Dennis will not change jobs” claim has `hash: yes` in it. When outputting to an output type like `markdownsnippet`, the claim won’t read “Dennis will not change jobs”, it’ll read “06b940c01111111849d11d0b90726f24a95e0ac2c5817ad1a49a0f298561adfb” (`printf "Dennis will not change jobskkjskvjsdwolvkjsjv" | shasum -a 256` on macOS).

If, after you’ve published your salted-and-hashed prediction, you want to reveal to the world (or Dennis) that you were very sure of your incorrect prediction, you need to reveal both the exact text of your claim as well as the salt, if any, that was used. This will let people run `printf` and `shasum` on their own computers to verify that 06b940c01111111849d11d0b90726f24a95e0ac2c5817ad1a49a0f298561adfb that you published at the beginning of the year was a prediction that you mis-guessed.

Of course, you could never reveal the claim text and the hash if you so choose. That way, Dennis will never know that you were pretty sure he was going to stay put at his then-current job.

Why salt? Because if Dennis or his friends have a lot of free time, they could type a bunch of guesses for your predictions into their computers and generate hashes for them. If Dennis manages to guess your exact wording for any of your hashed secret claims, they won’t be secret anymore. Guessing your exact wording plus a bunch of random letters and/or numbers and/or punctuation in a salt? That’s pretty much close to impossible.

## Unresolved questions

Should I encourage users to use YAML 1.1-only features like yes/no/on/off/~? in 1.2, only true/false/null are allowed.

## Build instructions

- `go get gopkg.in/yaml.v2`
- `go get github.com/davecgh/go-spew/spew`
