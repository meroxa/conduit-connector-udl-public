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

package destination

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"time"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/meroxa/udl-go"
)

func ToUDLAis(raw []byte, dataMode udl.AISIngestDataMode, classificationMarking string) (udl.AISIngest, error) {
	var vesselData VesselData
	err := json.Unmarshal(raw, &vesselData)
	if err != nil {
		return udl.AISIngest{}, err
	}

	// Replace underscores with spaces in vesselData strings
	replaceUnderscoresInStruct(&vesselData)

	var ais udl.AISIngest

	sdk.Logger(context.Background()).Debug().Msgf("vesselData: %+v", vesselData)

	layout := "2006-01-02T15:04:05.999Z"

	ts, tsErr := time.Parse(layout, vesselData.UpdateTimestamp)
	if tsErr != nil {
		sdk.Logger(context.Background()).Err(err).Msgf("Error parsing timestamp")
		return udl.AISIngest{}, err
	}
	ais.Ts = ts

	if vesselData.ID != "" {
		ais.Id = &vesselData.ID
	}
	ais.ClassificationMarking = classificationMarking
	if vesselData.StaticData.MMSI != 0 {
		ais.Mmsi = &vesselData.StaticData.MMSI
	}
	if vesselData.StaticData.Name != "" {
		ais.ShipName = &vesselData.StaticData.Name
	}
	if vesselData.StaticData.ShipType != "" {
		ais.ShipType = &vesselData.StaticData.ShipType
	}
	if vesselData.StaticData.CallSign != "" {
		ais.CallSign = &vesselData.StaticData.CallSign
	}

	if vesselData.StaticData.Flag != "" {
		ais.VesselFlag = &vesselData.StaticData.Flag
	}
	if vesselData.StaticData.Flag != "" {
		ais.VesselFlag = &vesselData.StaticData.Flag
	}

	if vesselData.LastPositionUpdate.Latitude != 0.0 {
		ais.Lat = &vesselData.LastPositionUpdate.Latitude
	}

	if vesselData.LastPositionUpdate.Longitude != 0.0 {
		ais.Lon = &vesselData.LastPositionUpdate.Longitude
	}
	hiAccuracy := vesselData.LastPositionUpdate.Accuracy == "HIGH"
	ais.PosHiAccuracy = &hiAccuracy

	if vesselData.LastPositionUpdate.Heading != 0.0 {
		ais.TrueHeading = &vesselData.LastPositionUpdate.Heading
	}

	if vesselData.LastPositionUpdate.Course != 0.0 {
		ais.Course = &vesselData.LastPositionUpdate.Course
	}

	if vesselData.LastPositionUpdate.NavigationalStatus != "" {
		ais.NavStatus = &vesselData.LastPositionUpdate.NavigationalStatus
	}

	dimensionsSlice := []float64{
		vesselData.StaticData.Dimensions.A,
		vesselData.StaticData.Dimensions.B,
		vesselData.StaticData.Dimensions.C,
		vesselData.StaticData.Dimensions.D,
	}

	if !allZero(dimensionsSlice) {
		ais.AntennaRefDimensions = &dimensionsSlice
	}

	if vesselData.StaticData.Dimensions.Length != 0.0 {
		ais.Length = &vesselData.StaticData.Dimensions.Length
	}

	if vesselData.StaticData.Dimensions.Width != 0.0 {
		ais.Width = &vesselData.StaticData.Dimensions.Width
	}

	if vesselData.CurrentVoyage.Draught != 0.0 {
		ais.Draught = &vesselData.CurrentVoyage.Draught
	}

	// Not every vessel has vesselData.CurrentVoyage.ETA; so those vessels will result in a DestinationETA being set to nil
	if vesselData.CurrentVoyage.ETA != "" {
		eta, etaErr := time.Parse(layout, vesselData.CurrentVoyage.ETA)
		// The UDL endpoint doesn't require DestinationETA to be populated to accept a valid payload; so we set it to nil instead of erroring out
		if etaErr != nil {
			ais.DestinationETA = nil
		}
		ais.DestinationETA = &eta
	}

	if vesselData.CurrentVoyage.MatchedPort.Port.Unlocode != "" {
		ais.CurrentPortLOCODE = &vesselData.CurrentVoyage.MatchedPort.Port.Unlocode
	}

	ais.DataMode = dataMode
	ais.Source = "Spire"

	return ais, err
}

func ToUDLElset(raw []byte) (udl.ElsetIngest, error) {
	var elset udl.ElsetIngest
	err := json.Unmarshal(raw, &elset)
	return elset, err
}

func ToUDLEphemeris(raw []byte, dataMode udl.EphemerisIngestDataMode, classificationMarking string) (UDLReport, error) {
	// parse raw lines to sp3 report
	sp3Report, err := Parse(raw)
	if err != nil {
		log.Printf("error parsing decoded bytes: %s", err)
	}

	log.Printf("name: %s Timestamp: %s  FlightModuleNumber: %d", sp3Report.SatelliteName, sp3Report.Entries[0].Timestamp, sp3Report.Entries[0].Position.FlightModuleNumber)

	// convert to UDL Report
	//ur, err := SP3cToUDL(sp3Report)
	ur, err := SP3cToUDL(sp3Report)
	if err != nil {
		log.Printf("error converting to udl report: %s", err)
	}

	return ur, err

}

func KeyPayload(rawKey string) (string, error) {
	// get the key payload
	keyMap := make(map[string]interface{})

	err := json.Unmarshal([]byte(rawKey), &keyMap)
	if err != nil {
		return "", err
	}

	// decode key
	decoded, err := base64.StdEncoding.DecodeString(keyMap["payload"].(string))
	if err != nil {
		log.Printf("error: %s; raw key payload: %+v", err, keyMap)
	}

	return string(decoded), nil
}
