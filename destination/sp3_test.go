package destination

import (
	"testing"
	"time"

	"github.com/matryer/is"
)

func sampleFile() []byte {
	return []byte(`
#cV2022 07 06 01 18 13.00000000    3546     d IGS08 FIT SPIR
## 2217 263893.00000000     1.00000000 59766 0.0543171297759
+    1   143  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0
+          0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0
+          0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0
+          0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0
+          0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0
++         0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0
++         0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0
++         0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0
++         0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0
++         0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0
%c G  cc GPS ccc cccc cccc cccc cccc ccccc ccccc ccccc ccccc
%c cc cc ccc ccc cccc cccc cccc cccc ccccc ccccc ccccc ccccc
%f  0.0000000  0.000000000  0.00000000000  0.000000000000000
%f  0.0000000  0.000000000  0.00000000000  0.000000000000000
%i    0    0    0    0      0      0      0      0         0
%i    0    0    0    0      0      0      0      0         0
/* NOTE: Spire sp3c satellite names are denoted by a        
/* 3-digit number representing flight module (FM) number.   
/* SATELLITE NAME: LEMUR-2-JOHN-TREIRES                     
/* CCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC
*  2022  7  6  1 18 13.00000000
P143  -6658.162753  -1527.302901   -971.376727  -3827.755483
V143   6844.820031  17028.031395 -74566.102286 999999.999999
*  2022  7  6  1 18 14.00000000
P143  -6657.474119  -1525.599224   -978.832746  -3827.858503
V143   6927.842084  17045.492627 -74554.222095 999999.999999
`)
}

func TestExtractSatelliteName(t *testing.T) {
	is := is.New(t)

	// Example raw input that includes a satellite name header
	rawInput := sampleFile()

	expectedSatName := "LEMUR-2-JOHN-TREIRES"
	satName, err := extractSatelliteName(rawInput)

	is.NoErr(err)
	is.Equal(satName, expectedSatName) // Satellite name should be extracted correctly
}

func TestParseEpoch(t *testing.T) {
	is := is.New(t)

	// A line that starts with "*  " followed by a date in the sp3c format
	line := "*  2022  7  6 13 18 13.00000000"

	expectedTimestamp, _ := time.Parse(sp3cTimeLayout, "2022  7  6 13 18 13.00000000")
	timestamp, err := parseEpoch(line)

	is.NoErr(err)
	is.Equal(timestamp, expectedTimestamp) // Timestamp should match the expected value
}

func TestCompact(t *testing.T) {
	is := is.New(t)

	compactInput := "This   is   a    test"
	expectedOutput := "This is a test"
	output := compact(compactInput)

	is.Equal(output, expectedOutput) // Compacted string should have no consecutive spaces
}

// Note: The following is an approach to test, assuming sample inputs for `parsePosition` and `parseVelocity`.
// You should create more detailed tests based on actual and varied inputs for better coverage.

func TestParsePosition(t *testing.T) {
	is := is.New(t)

	positionLine := samplePositionLine()
	expectedPosition := Position{
		FlightModuleNumber: 143,
		X:                  "-6657.474119",
		Y:                  "-1525.599224",
		Z:                  "-978.832746",
		ClockError:         "-3827.858503",
	}
	position, err := parsePosition(positionLine)

	is.NoErr(err)
	is.Equal(position, expectedPosition)
}

func TestParseVelocity(t *testing.T) {
	is := is.New(t)

	velocityLine := sampleVelocityLine()
	expectedVelocity := Velocity{
		FlightModuleNumber:     143,
		X:                      "6927.842084",
		Y:                      "17045.492627",
		Z:                      "-74554.222095",
		ClockErrorRateOfChange: "999999.999999",
	}
	velocity, err := parseVelocity(velocityLine)

	is.NoErr(err)
	is.Equal(velocity, expectedVelocity)
}

func samplePositionLine() string {
	return `P143  -6657.474119  -1525.599224   -978.832746  -3827.858503`
}

func sampleVelocityLine() string {
	return `V143   6927.842084  17045.492627 -74554.222095 999999.999999`
}
