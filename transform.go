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

package connector

import (
	"encoding/json"
	"log"
	"time"

	"github.com/meroxa/conduit-connector-udl/udl"
)

func ToUDLAis(raw []byte) (udl.AISIngest, error) {
	var vesselData VesselData
	err := json.Unmarshal(raw, &vesselData)

	UDLClassification := "U"
	var ais udl.AISIngest

	layout := "2006-01-02T15:04:05.999Z"

	ts, tsErr := time.Parse(layout, vesselData.UpdateTimestamp)
	if tsErr != nil {
		log.Fatalf("Error parsing timestamp: %v", tsErr)
	}
	ais.Ts = ts

	ais.Id = &vesselData.ID
	ais.ClassificationMarking = UDLClassification
	ais.Mmsi = &vesselData.StaticData.MMSI
	ais.ShipName = vesselData.StaticData.Name
	ais.ShipType = &vesselData.StaticData.ShipType
	ais.CallSign = &vesselData.StaticData.CallSign
	ais.VesselFlag = &vesselData.StaticData.Flag
	ais.Lat = &vesselData.LastPositionUpdate.Latitude
	ais.Lon = &vesselData.LastPositionUpdate.Longitude
	hiAccuracy := vesselData.LastPositionUpdate.Accuracy == "HIGH"
	ais.PosHiAccuracy = &hiAccuracy
	ais.TrueHeading = &vesselData.LastPositionUpdate.Heading
	ais.Course = &vesselData.LastPositionUpdate.Course
	ais.NavStatus = &vesselData.LastPositionUpdate.NavigationalStatus
	dimensionsSlice := []float64{
		vesselData.StaticData.Dimensions.A,
		vesselData.StaticData.Dimensions.B,
		vesselData.StaticData.Dimensions.C,
		vesselData.StaticData.Dimensions.D,
	}

	ais.AntennaRefDimensions = &dimensionsSlice
	ais.Length = &vesselData.StaticData.Dimensions.Length
	ais.Width = &vesselData.StaticData.Dimensions.Width
	ais.Draught = &vesselData.CurrentVoyage.Draught

	eta, etaErr := time.Parse(layout, vesselData.CurrentVoyage.ETA)
	if etaErr != nil {
		ais.DestinationETA = nil
	} else {
		ais.DestinationETA = &eta
	}

	if &vesselData.CurrentVoyage.MatchedPort.Port.Unlocode != nil {
		ais.CurrentPortLOCODE = &vesselData.CurrentVoyage.MatchedPort.Port.Unlocode
	}

	ais.DataMode = "TEST"
	ais.Source = "Spire"
	return ais, err
}

func ToUDLElset(raw []byte) (udl.ElsetIngest, error) {
	var elset udl.ElsetIngest
	err := json.Unmarshal(raw, &elset)
	return elset, err
}
