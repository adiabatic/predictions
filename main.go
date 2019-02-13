package main

import (
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

func main() {
	// no flags to parse yet, but we need to do this to make flag.Args() work

	flag.Parse()

	streams, err := stream.StreamsFromFiles(flag.Args())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("the streams:")
	spew.Dump(streams)

	tags := TagsUsed(streams)
	sort.Strings(tags)
	fmt.Println("tags used: ", tags)
}
