package destination

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	// e.g. 2022  7  6 13 18 13.00000000
	sp3cTimeLayout = "2006  1  2 15  4  5.00000000"
)

type Report struct {
	SatelliteName string
	Entries       []Entry
}

type Entry struct {
	Timestamp time.Time
	Position  Position
	Velocity  Velocity
}

type Position struct {
	// Each satellite name is denoted by a 3-digit number representing the flight module (FM) number
	FlightModuleNumber int
	// X coords in km
	X string
	// Y coords in km
	Y string
	// Z coords in km
	Z string
	// Estimated receiver clock error from true GPS time in microseconds
	ClockError string
}

type Velocity struct {
	// Each satellite name is denoted by a 3-digit number representing the flight module (FM) number
	FlightModuleNumber int
	// X velocity in decimeters/s
	X string
	// Y velocity in decimeters/s
	Y string
	// Z velocity in decimeters/s
	Z string
	// Estimated receiver clock error rate of change in units of 10-4 microseconds/sec. This is typically not estimated
	// from the precise orbit determination software and is therefore set to 999999.999999
	ClockErrorRateOfChange string
}

func Parse(raw []byte) (Report, error) {
	satName, err := extractSatelliteName(raw)
	if err != nil {
		return Report{}, err
	}

	entries, err := splitEntries(raw)
	if err != nil {
		return Report{}, err
	}

	return Report{
		SatelliteName: satName,
		Entries:       entries,
	}, nil
}

// extractSatelliteName returns the Satellite Name
func extractSatelliteName(raw []byte) (string, error) {
	// setup line by line reader
	reader := bytes.NewReader(raw)
	scanner := bufio.NewScanner(reader)
	lineNum := 0

	// skip to Satellite Name
	for lineNum < 21 {
		if !scanner.Scan() {
			return "", errors.New("invalid input")
		}
		lineNum++
	}

	// extract Satellite Name
	fullLine := scanner.Text()
	satName := strings.Replace(fullLine, "/* SATELLITE NAME: ", "", 1)
	satName = strings.TrimSpace(satName)
	return satName, nil
}

func splitEntries(raw []byte) ([]Entry, error) {
	// setup line by line reader
	reader := bytes.NewReader(raw)
	scanner := bufio.NewScanner(reader)

	// skip past header
	lineNum := 0
	for lineNum < 22 {
		if !scanner.Scan() {
			return nil, errors.New("invalid input")
		}
		lineNum++
	}

	var entries []Entry

	// there are more entries
	for scanner.Scan() {
		lineNum++
		var entry Entry

		// extract epoch line
		epochLine := scanner.Text()
		if handleEOF(epochLine) {
			return entries, nil
		}
		timestamp, err := parseEpoch(epochLine)
		if err != nil {
			log.Println("error on line: ", lineNum)
			return nil, err
		}
		entry.Timestamp = timestamp

		// extract position
		scanner.Scan() // next line
		lineNum++
		positionLine := scanner.Text()
		position, err := parsePosition(positionLine)
		if err != nil {
			log.Println("error on line: ", lineNum)
			return nil, err
		}
		entry.Position = position

		// extract velocity
		scanner.Scan() // next line
		lineNum++
		velocityLine := scanner.Text()
		velocity, err := parseVelocity(velocityLine)
		if err != nil {
			log.Println("error on line: ", lineNum)
			return nil, err
		}
		entry.Velocity = velocity

		entries = append(entries, entry)
	}

	return entries, nil
}

func handleEOF(line string) bool {
	if strings.Contains(line, "EOF") {
		return true
	}
	return false
}

func parseEpoch(line string) (time.Time, error) {
	// string line prefix
	line = strings.TrimPrefix(line, "*  ")
	// sp3c timestamp format yyyy mm dd hh ss.ssssssss
	t, err := time.Parse(sp3cTimeLayout, line)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func parsePosition(line string) (Position, error) {
	// strip P prefix
	line = strings.TrimPrefix(line, "P")

	// compact
	line = compact(line)

	// split 5 components FM number, X, Y, Z and clock error
	components := strings.Split(line, " ")
	if len(components) != 5 {
		log.Printf("components(%d): %+v", len(components), components)
		return Position{}, errors.New("invalid position line")
	}
	fmNumber, err := strconv.Atoi(components[0])
	if err != nil {
		return Position{}, err
	}
	return Position{
		FlightModuleNumber: fmNumber,
		X:                  components[1],
		Y:                  components[2],
		Z:                  components[3],
		ClockError:         components[4],
	}, nil
}

func parseVelocity(line string) (Velocity, error) {
	// strip V prefix
	line = strings.TrimPrefix(line, "V")

	// compact
	line = compact(line)

	// split 5 components FM number, X, Y, Z and clock error
	components := strings.Split(line, " ")
	if len(components) != 5 {
		log.Printf("components(%d): %+v", len(components), components)
		return Velocity{}, errors.New("invalid velocity line")
	}
	fmNumber, err := strconv.Atoi(components[0])
	if err != nil {
		return Velocity{}, err
	}
	return Velocity{
		FlightModuleNumber:     fmNumber,
		X:                      components[1],
		Y:                      components[2],
		Z:                      components[3],
		ClockErrorRateOfChange: components[4],
	}, nil
	return Velocity{}, nil
}

// compact takes a string and removes duplicated (padded) spaces
func compact(in string) string {
	space := regexp.MustCompile(`\s+`)
	return space.ReplaceAllString(in, " ")
}

func SP3cToUDL(report Report) (UDLReport, error) {
	var uReport UDLReport

	// the idOnOrbit is the FlightModuleNumber, so we take the first one and
	// error if it changes while going through the entries
	fm := report.Entries[0].Position.FlightModuleNumber
	for _, e := range report.Entries {
		// if FlightModuleNumber has changed mid-report, return error
		if e.Position.FlightModuleNumber != fm || e.Velocity.FlightModuleNumber != fm {
			return UDLReport{}, errors.New("report contains multiple flight modules")
		}
		var udlEntry UDLEntry

		// reformat date
		udlEntry.Timestamp = sp3cTimestampToUDL(e.Timestamp)

		// reformat position
		udlEntry.Position = sp3cPositionToUDL(e.Position)

		// reformat velocity
		var err error
		udlEntry.Velocity, err = sp3cVelocityToUDL(e.Velocity)
		if err != nil {
			return UDLReport{}, err
		}
		uReport.Entries = append(uReport.Entries, udlEntry)
	}

	// map Spire Flight Module number to NORAD ID (for use in idOnOrbit)
	nID, ok := fmMap()[fm]
	if !ok {
		return UDLReport{}, fmt.Errorf("no norad mapping for flight ID %d", fm)
	}
	uReport.ID = strconv.Itoa(nID)

	return uReport, nil
}

func sp3cTimestampToUDL(t time.Time) string {
	return t.Format(udlTimeLayout)
}

func sp3cPositionToUDL(p Position) UDLPosition {
	return UDLPosition{
		X: p.X,
		Y: p.Y,
		Z: p.Z,
	}
}

func sp3cVelocityToUDL(v Velocity) (UDLVelocity, error) {
	fX, err := strconv.ParseFloat(v.X, 64)
	if err != nil {
		return UDLVelocity{}, err
	}

	fY, err := strconv.ParseFloat(v.Y, 64)
	if err != nil {
		return UDLVelocity{}, err
	}

	fZ, err := strconv.ParseFloat(v.Z, 64)
	if err != nil {
		return UDLVelocity{}, err
	}

	return UDLVelocity{
		X: fixedWidthFloat(fX/10000, 11),
		Y: fixedWidthFloat(fY/10000, 11),
		Z: fixedWidthFloat(fZ/10000, 11),
	}, nil
}

// fixedWidthFloat returns a fixed width float64
func fixedWidthFloat(num float64, width int) string {
	whole := fmt.Sprintf("%.0f", num)
	lenWhole := len(strings.TrimPrefix(whole, "-"))
	prec := width - lenWhole
	return strconv.FormatFloat(num, 'f', prec, 64)
}
