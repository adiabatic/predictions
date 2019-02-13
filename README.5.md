# Predictions’ file format

This describes `predictions`’ file format. `predictions` uses [YAML][], so this document has a crash course in writing YAML.

[yaml]: https://yaml.org/

## Structure

- `predictions` takes in multidocument YAML streams.
- There is one YAML stream per file. A YAML stream (hereafter “file”) needs at least two documents in it to be interesting to `predictions`.
- YAML documents are separated by “---” on lines by themselves.
- The first document in a file is a mapping that specifies metadata.
- Mappings’ keys (the part before the “: ”) can be in any order.
- Each document after the metadata document contains one mapping.
- Each mapping after the metadata document is called a prediction.
- Each prediction has both a claim and a confidence unless the author forgot one (or both).

## Example

```yaml
---
# comments are great

# if you have a colon in a mapping value,
# you need to quote it somehow
title: 'Predictions: for 2019'
scope: for 2019
salt: kkjskvjsdwolvkjsjv
private notes: >
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
# the notes value is a good place to explain why
# “happened” has the value it does
notes: Left Acme on January 14. Go figure.
---
# use > to fold multiline strings into one line.
# newlines are turned into spaces.
# notice that we don’t need to add quotes for the colon.
claim: >
  I will finish
  The Legend of Zelda: Breath of the Wild
  on Master Mode
confidence: 95
tags: [games]
# use | to preserve line breaks.
# note that the “-” inside will be processed as Markdown.
notes: |
  This game will get significantly harder if I finish
  any more Divine Beasts. I should:

  - get all possible shrines before doing a second Divine Beast
  - avoid doing DLC shrines so I don't clutter my map
```

## Metadata-document mapping keys

### `title`

The title of the file full of predictions.

### `salt` (per-document) (not yet implemented)

A per-file salt used for hashing sensitive predictions.

See “Salt-and-hash rationale” below for why you might want to do this.

### `scope` (not yet implemented)

This restricts the scope of predictions to a particular domain. The idea of scopes is that they’re one-per-file so one can combine, say, 2018 predictions, 2019 predictions, and 2020 predictions in an invocation of `predictions` and see combined results with each year’s prediction labeled as such.

### `tag order` (not yet implemented)

A list showing the order you want tags displayed in.

Use this if you want to force an order in your output. Suppose you have three tags of predictions: one about U.S. politics, another for international politics, and a third about bananas. If you want to ensure that your output displays these three topics in that order (as opposed to boring your readers at the beginning with your heady pronouncements about bananas), then you should have `tag order: [U.S. politics, International politics, Bananas]` in your metadata document.

## Prediction-document mapping keys

### `claim` (required)

A prediction that something will, or won’t, happen.

### `confidence` (warns if missing)

How confident you are that this will happen, expressed as a percentage, without the percent sign.

### `tags`

A list of tags that you want to associate with a claim.

Tags are used for grouping purposes. For example, you might be very well-calibrated when it comes to predicting national and international affairs, but wildly off-base when predicting what your friends will do. If you tag each type of prediction with its type, you’ll be able to find out if you’re better in some areas than others.

If some, but not all, of your predictions have tags, the untagged ones will be tagged “Untagged”.

### `happened`

Either a boolean (yes/no/true/false) or null (~). Use yes/true for something that did happen, no/false for something that definitely didn’t happen, and ~ for weird ambiguous results that you want to throw out of any analysis.

### `cause for exclusion` (name subject to change)

A string. Put any explanation you want in it.

Some predictions don’t pan out the way you predict, and, of those, some just _go weird_. If a prediction has a `cause for exclusion` in it, that prediction won’t be included in calculations of your accuracy.

When would you use something like this? Suppose you have months of predictions for something like “I will park straight in a space when I get to work”, one per day. What would you write down if you had to parallel-park one day (it was super crowded) or you had to call in sick and didn’t drive there at all? A `cause for exclusion` will let you exclude a prediction from all the analysis yet still let you keep a record of having made the prediction.

### `hash` (not yet implemented)

A boolean. If true, then this entry is hashed before going to a publicly-displayed output.

Uses a per-prediction salt if it exists. If it doesn’t, then it’ll fall back to the whole-file salt specified in the metadata document. If that doesn’t exist either, then the claim will be unsalted when published.

### `salt` (per-prediction) (not yet implemented)

A per-prediction salt used for hashing sensitive predictions. If `salt` is specified in a prediction, then `hash` is implied.

See “salt-and-hash rationale” below for why you might want to do this.

### `notes` (Markdown processing not yet implemented)

`notes` is for you to write notes about the prediction. Frequently, it’s helpful to write down why `happened` has the value it does. For example, if your claim is “I will weigh less than 185 pounds”, then it might be nice to write down the date you first dropped below 185 pounds. Similarly, if you’re trying to predict world events, it’s handy to link to newspaper articles substantiating whether your claim happened (or not).

The contents of `notes` may be put into the output, whether HTML or Markdown. If you want to write notes that never will be put into the output, use `private notes` or similar mapping keys.

Because multiline notes are frequently easier to read, the literal-style indicator (“|”) can be helpful here.

## Forward-compatibility concerns

Currently, `predictions` uses [go-yaml][] version 2, which supports YAML 1.2 but still (erroneously) allows yes/no/on/off/~ from YAML 1.1. If you use any of these instead of true/false/null, be prepared to search-and-replace in your files in a few years’ time.

[go-yaml]: https://github.com/go-yaml/yaml
