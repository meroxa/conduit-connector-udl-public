// Copyright © 2022 Meroxa, Inc.
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

package connector

import (
	"context"
	"net/http"
	"testing"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/matryer/is"

	"github.com/meroxa/udl-go"
)

type mockClient struct {
	udl.ClientInterface
}

func (c *mockClient) FiledropUdlAisPostId(ctx context.Context, body udl.FiledropUdlAisPostIdJSONRequestBody, reqEditors ...udl.RequestEditorFn) (*http.Response, error) {
	// Add your custom logic or validations here if needed
	// For example, you can print the request body and other information:
	// fmt.Println("request.body:")
	// fmt.Println(body[0])

	// Since this is a mock function, you can simply return a successful HTTP response without making an actual API call
	return &http.Response{
		StatusCode: http.StatusOK,
	}, nil
}

func TestParameters(t *testing.T) {
	is := is.New(t)
	d := Destination{}
	params := d.Parameters()
	is.Equal(len(params), 5) // Assumes there are 6 parameters in the config
}

func TestConfigure(t *testing.T) {
	is := is.New(t)
	dest := Destination{}
	ctx := context.Background()
	err := dest.Configure(ctx, map[string]string{
		"baseURL":               "https://example.com",
		"httpBasicAuthUsername": "user",
		"httpBasicAuthPassword": "pass",
		"dataType":              "AIS",
		"dataMode":              "TEST",
	})
	is.NoErr(err)
	is.Equal(dest.config.BaseURL, "https://example.com")
}

func TestConfigureWithInvalidConfig(t *testing.T) {
	is := is.New(t)
	dest := Destination{}
	ctx := context.Background()
	err := dest.Configure(ctx, map[string]string{
		"invalid_key": "invalid_value",
	})
	is.True(err != nil)
}

func TestOpen(t *testing.T) {
	is := is.New(t)
	dest := Destination{}
	ctx := context.Background()
	err := dest.Configure(ctx, map[string]string{
		"baseURL":               "https://example.com",
		"httpBasicAuthUsername": "user",
		"httpBasicAuthPassword": "pass",
		"dataType":              "AIS",
		"dataMode":              "TEST",
	})
	is.NoErr(err)
	err = dest.Open(ctx)
	is.NoErr(err)
	is.True(dest.client != nil)
}

func TestWrite(t *testing.T) {
	// Dummy vessel data generated with ChatGPT
	validJSON := []byte(`{
		"antennaRefDimensions": [10.0, 20.0, 5.0, 5.0],
		"avgSpeed": 15.5,
		"callSign": "ABCD123",
		"cargoType": "General Cargo",
		"classificationMarking": "UNCLASSIFIED",
		"course": 45.0,
		"createdAt": "2022-01-01T00:00:00Z",
		"createdBy": "user1",
		"currentPortGUID": "USNYC",
		"currentPortLOCODE": "USNYC",
		"dataMode": "REAL",
		"destination": "New York",
		"destinationETA": "2022-01-10T00:00:00Z",
		"distanceToGo": 500.0,
		"distanceTravelled": 1000.0,
		"currentVoyage.draught": 9.0,
		"engagedIn": "Fishing",
		"etaCalculated": "2022-01-09T00:00:00Z",
		"etaUpdated": "2022-01-05T00:00:00Z",
		"id": "abc123",
		"idTrack": "track1",
		"idVessel": "vessel1",
		"imon": 1234567,
		"lastPortGUID": "USMIA",
		"lastPortLOCODE": "USMIA",
		"lat": 40.7128,
		"length": 100.0,
		"lastPositionUpdate.longitude": -74.0060,
		"maxSpeed": 20.0,
		"staticData.mmsi": 987654321,
		"lastPositionUpdate.navigationalStatus": "Underway Using Engine",
		"nextPortGUID": "USLAX",
		"nextPortLOCODE": "USLAX",
		"origNetwork": "AIS",
		"origin": "AIS",
		"posDeviceType": "GPS",
		"posHiAccuracy": true,
		"posHiLatency": false,
		"rateOfTurn": 5.0,
		"shipDescription": "General Cargo Ship",
		"staticData.name": "MV_Sample_Ship",
		"staticData.shipType": "Cargo",
		"source": "AIS",
		"specialCraft": "Pilot Vessel",
		"specialManeuver": true,
		"speed": 12.0,
		"trueHeading": 45.0,
		"updateTimestamp": "2022-01-05T12:00:00Z",
		"vesselFlag": "US",
		"width": 30.0
	  }`)
	is := is.New(t)
	dest := Destination{}
	ctx := context.Background()
	dest.config.DataType = "AIS"
	dest.client = &mockClient{}
	records := []sdk.Record{{Payload: sdk.Change{After: sdk.RawData(validJSON)}}}
	num, err := dest.Write(ctx, records)
	is.NoErr(err)
	is.Equal(num, len(records))
}

func TestWriteWithInvalidDataType(t *testing.T) {
	is := is.New(t)
	dest := Destination{}
	ctx := context.Background()
	dest.config.DataType = "INVALID"
	dest.client = &mockClient{}
	records := []sdk.Record{}
	num, err := dest.Write(ctx, records)
	is.Equal(err.Error(), "invalid data type: INVALID; expecting the data type to be on of the following values: 'AIS', 'ELSET'")
	is.Equal(num, 0)
}
