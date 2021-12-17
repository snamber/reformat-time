package main

import (
	"errors"
	"fmt"
	"github.com/integrii/flaggy"
	"log"
	"math"
	"os"
	"strconv"
	"time"
)

func main() {

	var (
		rfc3339, rfc3339utc, unix, milli, float bool
	)
	flaggy.Bool(&rfc3339utc, "", "rfc3339-utc", "Format time as RFC 3339 UTC (default)")
	flaggy.Bool(&rfc3339, "r", "rfc3339", "Format time as RFC 3339")
	flaggy.Bool(&unix, "u", "unix", "Format time as unix time")
	flaggy.Bool(&milli, "m", "unix-milli", "Format time as unix milli")
	flaggy.Bool(&float, "f", "unix-float", "Format time as unix float")
	flaggy.SetDescription("A tiny convenience tool to convert time formats.\n\n  Example:\n    reformat-time -u -- \"$(date)\"")
	flaggy.Parse()

	// Set default
	if !(rfc3339 || rfc3339utc || unix || milli || float) {
		rfc3339utc = true
	}

	if len(flaggy.TrailingArguments) != 1 {
		flaggy.ShowHelp(fmt.Sprintf("time-format needs exactly one command line argument. got %v\n", flaggy.TrailingArguments))
		os.Exit(1)
	}
	input := flaggy.TrailingArguments[0]

	// Parse
	parsed := time.Time{}
	if p, err := time.Parse(time.RFC3339Nano, input); err == nil {
		parsed = p
	} else if p, err := time.Parse(time.UnixDate, input); err == nil {
		parsed = p
	} else if p, err := parseUnix(input); err == nil {
		parsed = p
	} else if p, err := parseUnixMilli(input); err == nil {
		parsed = p
	} else if p, err := parseFloat(input); err == nil {
		parsed = p
	} else {
		log.Fatalf("Failed to parse %s\n", input)
	}

	// Output
	if rfc3339utc {
		fmt.Println("FRC3339Nano UTC:", parsed.UTC().Format(time.RFC3339Nano))
	}
	if rfc3339 {
		fmt.Println("RFC3339Nano:", parsed.Format(time.RFC3339Nano))
	}
	if unix {
		fmt.Println("Unix:", parsed.Unix())
	}
	if milli {
		fmt.Println("UnixMilli:", parsed.UnixMilli())
	}
	if float {
		fmt.Printf("Float: %f\n", float64(parsed.UnixNano())/1e9)
	}
}

func parseUnix(input string) (time.Time, error) {
	i, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	t := time.Unix(i, 0)
	if t.Year() > 10_000 {
		return time.Time{}, errors.New("year > 10000")
	}
	return t, nil
}

func parseUnixMilli(input string) (time.Time, error) {
	i, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	t := time.UnixMilli(i)
	return t, nil
}

func parseFloat(input string) (time.Time, error) {
	f, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return time.Time{}, err
	}
	floor := math.Floor(f)
	remainder := f - floor
	nano := remainder * 1e9
	return time.Unix(int64(floor), int64(nano)), nil
}
