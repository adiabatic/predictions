package cmd

import (
	"bytes"
	"fmt"
	"os"
	"sort"

	"github.com/adiabatic/predictions/stream"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
)

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {

}

var rootCmd = &cobra.Command{
	Use:   "predictions",
	Short: "predictions finds out how well-calibrated your predictions are",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		streams, err := stream.StreamsFromFiles(args)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// fmt.Println("the streams:")
		// spew.Dump(streams)

		tags := stream.TagsUsed(streams)
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

	},
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// extra stuff that belongs elsewhere

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
