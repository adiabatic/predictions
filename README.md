# Predictions

`predictions` helps you figure out how how good you are at predicting things. You feed it a [YAML][] file full of predictions and the outcomes and it calculates how good you are at predicting things.

Not sure what this whole ‘predictions’ thing is? Have a look at [Scott Alexander’s 2018 predictions retrospective][ssc2018] to get an idea of what’s involved.

[yaml]: https://en.wikipedia.org/wiki/YAML "YAML Ain’t Markup Language"
[ssc2018]: https://slatestarcodex.com/2019/01/22/2018-predictions-calibration-results/

## Usage

`predictions [options] inputfile`

`predictions` always outputs to standard output, with errors and warnings printed to standard error. To get the output into a file, redirect standard output to a file with `>`, as in `predictions 2019.yaml > 2019.markdown`.

## The YAML file format

In order to use `predictions` you’ll need to write one or more YAML files in a text editor. See `README.5.md` to learn what `predictions` expects in those files.

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

## File-format forward compatibility issues

`gopkg.in/yaml.v2` supports YAML 1.1-only booleans and nulls like yes/no/on/off/~, even in documents that are assumed to be YAML 1.2. Version 3 of `go-yaml` will correct this someday. For maximum forward compatibility, you should use only “true” and “false” for boolean values and use only “null” for null values.

## Build instructions

```sh
go get gopkg.in/yaml.v2
go get github.com/davecgh/go-spew/spew
go get github.com/stretchr/testify
go get -u github.com/spf13/cobra/cobra
go get -u gopkg.in/russross/blackfriday.v2
```
