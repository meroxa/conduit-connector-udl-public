// Copyright © 2023 Meroxa, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package destination

import (
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestToUDLAis(t *testing.T) {
	is := is.New(t)

	// Define a valid raw JSON input for the test
	raw := []byte(`{
		"ID": "1",
		"UpdateTimestamp": "2022-01-01T00:00:00.000Z",
		"StaticData": {
			"MMSI": 123456789,
			"Name": "Vessel Name",
			"ShipType": "Cargo",
			"CallSign": "CALLSIGN",
			"Flag": "US",
			"Dimensions": {
				"A": 10,
				"B": 20,
				"C": 30,
				"D": 40,
				"Length": 50,
				"Width": 60
			}
		},
		"LastPositionUpdate": {
			"Latitude": 12.34,
			"Longitude": 56.78,
			"Accuracy": "HIGH",
			"Heading": 90,
			"Course": 100,
			"NavigationalStatus": "Underway"
		},
		"CurrentVoyage": {
			"Draught": 5.5,
			"ETA": "2022-01-10T00:00:00.000Z",
			"MatchedPort": {
				"Port": {
					"Unlocode": "USNYC"
				}
			}
		}
	}`)

	// Call the toUDLAis function with the test input
	ais, err := ToUDLAis(raw, "TEST")
	is.NoErr(err) // Check for no errors

	// Verify the output fields
	is.Equal(*ais.Id, "1")
	is.Equal(ais.ClassificationMarking, "U")
	is.Equal(*ais.Mmsi, int64(123456789))
	is.Equal(*ais.ShipName, "Vessel Name")
	is.Equal(*ais.ShipType, "Cargo")
	is.Equal(*ais.CallSign, "CALLSIGN")
	is.Equal(*ais.VesselFlag, "US")
	is.Equal(*ais.Lat, 12.34)
	is.Equal(*ais.Lon, 56.78)
	is.Equal(*ais.PosHiAccuracy, true)
	is.Equal(*ais.TrueHeading, float64(90))
	is.Equal(*ais.Course, float64(100))
	is.Equal(*ais.NavStatus, "Underway")
	expectedDimensions := []float64{10, 20, 30, 40}
	is.Equal(*ais.AntennaRefDimensions, expectedDimensions)
	is.Equal(*ais.Length, 50.0)
	is.Equal(*ais.Width, 60.0)
	is.Equal(*ais.Draught, 5.5)
	expectedEta, _ := time.Parse("2006-01-02T15:04:05.999Z", "2022-01-10T00:00:00.000Z")
	is.Equal(*ais.DestinationETA, expectedEta)
	is.Equal(*ais.CurrentPortLOCODE, "USNYC")
	is.Equal(string(ais.DataMode), "TEST")
	is.Equal(ais.Source, "Spire")
}

func TestToUDLElset(t *testing.T) {
	is := is.New(t)

	// Define a valid raw JSON input for the test
	raw := []byte(`{
		"idOnOrbit": "1",
		"epoch": "2022-01-01T00:00:00.000Z"
	}`)

	// Call the toUDLElset function with the test input
	elset, err := ToUDLElset(raw)
	is.NoErr(err) // Check for no errors

	// Verify the output fields
	is.Equal(*elset.IdOnOrbit, "1")
	expectedTimestamp, _ := time.Parse("2006-01-02T15:04:05.999Z", "2022-01-01T00:00:00.000Z")
	is.Equal(elset.Epoch, expectedTimestamp)
}
