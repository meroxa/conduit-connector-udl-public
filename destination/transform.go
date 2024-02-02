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
	"encoding/json"
	"slices"
	"time"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/meroxa/udl-go"
)

func allZero(slice []float64) bool {
	for _, v := range slice {
		if v != 0 {
			return false
		}
	}
	return true
}

var spireToNavcenShipTypeMapping = map[string]string{
	"ANTI POLLUTION":      "Other",
	"CAR CARRIER":         "Cargo",
	"COMBINATION CARRIER": "Combination Carrier",
	"CONTAINER":           "Container",
	"DIVE VESSEL":         "Other",
	"DREDGER":             "Other",
	"DRY BULK":            "Cargo",
	"FISHING":             "Fishing",
	"GAS CARRIER":         "Gas Carrier",
	"GENERAL CARGO":       "Cargo",
	"GENERAL TANKER":      "General Tanker",
	"HIGH SPEED CRAFT":    "Other",
	"LAW ENFORCEMENT":     "Special Craft",
	"LIVESTOCK":           "Livestock Carrier",
	"LNG CARRIER":         "LNG Carrier",
	"MEDICAL TRANS":       "Special Craft",
	"MILITARY OPS":        "Other",
	"OFFSHORE":            "Offshore Vessel",
	"OTHER":               "Other",
	"PASSENGER":           "Passenger",
	"PILOT VESSEL":        "Special Craft",
	"PLEASURE CRAFT":      "Other",
	"PORT TENDER":         "Special Craft",
	"REEFER":              "Cargo",
	"ROLL ON ROLL OFF":    "Cargo",
	"SAILING":             "Other",
	"SEARCH AND RESCUE":   "Special Craft",
	"SPECIAL CRAFT":       "Special Craft",
	"TANKER CHEMICALS":    "Tanker Chemicals",
	"TANKER CRUDE":        "Tanker Crude",
	"TANKER PRODUCT":      "Tanker Product",
	"TUG":                 "Tug",
	"VEHICLE PASSENGER":   "Vehicle Passenger",
}

var cargoTypeMapping = map[string]string{
	"CAR CARRIER":      "Carrier",
	"DRY BULK":         "Dry Bulk",
	"GENERAL CARGO":    "General Cargo",
	"REEFER":           "Reefer",
	"ROLL ON ROLL OFF": "Roll on Roll off Cargo",
}

var engagedIn = []string{"ANTI POLLUTION", "DIVE VESSEL", "DREDGER", "HIGH SPEED CRAFT", "MILITARY OPS", "OTHER", "PLEASURE CRAFT", "SAILING"}
var specialCraft = []string{"LAW ENFORCEMENT", "MEDICAL TRANS", "PILOT VESSEL", "PORT TENDER", "SEARCH AND RESCUE", "SPECIAL CRAFT"}
var cargoType = []string{"CAR CARRIER", "DRY BULK", "GENERAL CARGO", "REEFER", "ROLL ON ROLL OFF"}

func setAISShipSubType(vesselData VesselData, ais *udl.AISIngest) {
	if slices.Contains(engagedIn, vesselData.StaticData.ShipType) {
		sdk.Logger(context.Background()).Debug().Msgf("%s is in the engagedIn\n", vesselData.StaticData.ShipType)
		if vesselData.StaticData.ShipSubType != "" {
			sdk.Logger(context.Background()).Debug().Msgf("ShipSubType: %s", vesselData.StaticData.ShipSubType)
			ais.EngagedIn = &vesselData.StaticData.ShipSubType
		}
	}
	if slices.Contains(specialCraft, vesselData.StaticData.ShipType) {
		sdk.Logger(context.Background()).Debug().Msgf("%s is in the specialCraft\n", vesselData.StaticData.ShipType)
		specialCraftValue := toTitleCase(vesselData.StaticData.ShipType)
		sdk.Logger(context.Background()).Debug().Msgf("specialCraftValue - %s\n", specialCraftValue)
		ais.SpecialCraft = &specialCraftValue
	}
	if slices.Contains(cargoType, vesselData.StaticData.ShipType) {
		sdk.Logger(context.Background()).Debug().Msgf("%s is in the cargoType\n", vesselData.StaticData.ShipType)
		cargoTypeValue := cargoTypeMapping[vesselData.StaticData.ShipType]
		sdk.Logger(context.Background()).Debug().Msgf("cargoTypeValue - %s\n", cargoTypeValue)
		ais.CargoType = &cargoTypeValue
	}
}

func ToUDLAis(raw []byte, dataMode udl.AISIngestDataMode, classificationMarking string) (udl.AISIngest, error) {
	var vesselData VesselData
	err := json.Unmarshal(raw, &vesselData)

	var ais udl.AISIngest

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
		navCenShipType := spireToNavcenShipTypeMapping[vesselData.StaticData.ShipType]
		ais.ShipType = &navCenShipType

		sdk.Logger(context.Background()).Info().Msgf("ais shipType: %s", *ais.ShipType)
		sdk.Logger(context.Background()).Info().Msgf("ShipType: %s", vesselData.StaticData.ShipType)
		setAISShipSubType(vesselData, &ais)

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
