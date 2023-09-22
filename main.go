package main

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/integrii/flaggy"
	"github.com/oklog/ulid/v2"
	"log"
	"math"
	"os"
	"strconv"
	"time"
)

func main() {

	var (
		rfc3339, rfc3339utc, unix, milli, float, ulid bool
	)
	flaggy.Bool(&rfc3339utc, "", "rfc3339-utc", "Format time as RFC 3339 UTC (default)")
	flaggy.Bool(&rfc3339, "r", "rfc3339", "Format time as RFC 3339")
	flaggy.Bool(&unix, "u", "unix", "Format time as unix time")
	flaggy.Bool(&milli, "m", "unix-milli", "Format time as unix milli")
	flaggy.Bool(&float, "f", "unix-float", "Format time as unix float")
	flaggy.Bool(&ulid, "id", "uuid", "Format time as UUID (ULID)")
	flaggy.SetDescription("A tiny convenience tool to convert time formats.\n\n  Example:\n    reformat-time -u -- \"$(date)\"")
	flaggy.Parse()

	// Set default
	if !(rfc3339 || rfc3339utc || unix || milli || float || ulid) {
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
	} else if p, err := parseUUID(input); err == nil {
		parsed = p
	} else {
		log.Fatalf("Failed to parse %s\n", input)
	}

	// Output
	if rfc3339utc {
		fmt.Println("FRC3339Nano UTC:", parsed.UTC().Format(time.RFC3339Nano))
	}
	if rfc3339 {
		fmt.Println("RFC3339Nano:", parsed.In(time.Local).Format(time.RFC3339Nano))
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
	if ulid {
		converted, err := convertToUUID(parsed)
		if err != nil {
			log.Fatalf("Failed to convert parsed timestamp to a ULID: %s\n", parsed)
		}
		fmt.Printf("ULID: %s\n", converted)
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

func parseUUID(input string) (time.Time, error) {
	id, err := uuid.Parse(input)
	if err != nil {
		return time.Time{}, err
	}
	return ulid.Time(uuid2ulid(id).Time()), nil
}

type zeroReader struct{}

func (zeroReader) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = 0
	}
	return len(p), nil
}

// convertToUUID generates a new ulid.ULID for a given timestamp with entropy bits set to zero and returns it as UUID
func convertToUUID(time time.Time) (uuid.UUID, error) {
	ul, err := ulid.New(ulid.Timestamp(time), zeroReader{})
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to create new ulid: %w", err)
	}
	return ulid2uuid(ul), nil

}

// ulid2uuid converts a ulid.ULID to a uuid.UUID
func ulid2uuid(ulid ulid.ULID) uuid.UUID {
	bt, err := ulid.MarshalBinary()
	if err != nil {
		panic(err)
	}
	uu, err := uuid.FromBytes(bt)
	if err != nil {
		panic(err)
	}
	return uu
}

// uuid2ulid converts a uuid.UUID to a ulid.ULID
func uuid2ulid(uuid uuid.UUID) (ulid ulid.ULID) {
	bt, err := uuid.MarshalBinary()
	if err != nil {
		panic(err)
	}
	err = ulid.UnmarshalBinary(bt)
	if err != nil {
		panic(err)
	}
	return ulid
}
