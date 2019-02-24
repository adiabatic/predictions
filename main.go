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

package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/adiabatic/predictions/stream"
	"github.com/davecgh/go-spew/spew"
)

func deduplicateTags(ss []string) []string {
	seen := make(map[string]struct{}, len(ss))
	j := 0
	for _, v := range ss {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		ss[j] = v
		j++
	}
	return ss[:j]
}

// TagsUsed returns a list of all tags used in the given slice of stream.Stream.
func TagsUsed(ss []stream.Stream) []string {
	ret := make([]string, 0)
	for _, s := range ss {
		for _, p := range s.Predictions {
			for _, t := range p.Tags {
				ret = append(ret, t)
			}
		}
	}
	return deduplicateTags(ret)
}

func asMarkdown(d stream.PredictionDocument) string {
	meat := fmt.Sprintf("%v: %v%%", d.Claim, *(d.Confidence))
	withToppings := ""
	if d.Happened == nil {
		withToppings = fmt.Sprintf("- <i>%v</i>", meat)
	} else if *(d.Happened) {
		withToppings = fmt.Sprintf("- %v", meat)
	} else {
		withToppings = fmt.Sprintf("- <s>%v</s>", meat)
	}

	return withToppings
}

func hasTag(d stream.PredictionDocument, tag string) bool {
	for _, t := range d.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

func main() {
	// no flags to parse yet, but we need to do this to make flag.Args() work

	flag.Parse()

	streams, err := stream.StreamsFromFiles(flag.Args())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// fmt.Println("the streams:")
	// spew.Dump(streams)

	tags := TagsUsed(streams)
	sort.Strings(tags)
	var v stream.Validator

	for _, s := range streams {
		errs := v.RunAll(s)
		if len(errs) > 0 {
			spew.Dump(errs)
		}
	}

	for _, tag := range tags {
		buf := &bytes.Buffer{}
		for _, s := range streams {

			for _, d := range s.Predictions {
				if d.ShouldExclude() {
					continue // only print out the ones with cause
				}

				if hasTag(d, tag) {
					fmt.Fprintln(buf, asMarkdown(d))
				}
			}
		}

		if buf.Len() > 0 {
			fmt.Println("#", tag)
			fmt.Println()
			fmt.Println(buf.String())
		}

	}

}
