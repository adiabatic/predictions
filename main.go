package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/adiabatic/predictions/stream"
	"github.com/davecgh/go-spew/spew"
)

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
}
