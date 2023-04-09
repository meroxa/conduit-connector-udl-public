package connector

import (
	"encoding/json"

	"github.com/meroxa/conduit-connector-udl/udl"
)

func flattenMap(inputMap map[string]interface{}, flattenedMap map[string]interface{}, parentKey string) {
	for key, value := range inputMap {
		newKey := key
		if parentKey != "" {
			newKey = parentKey + "." + key
		}
		switch v := value.(type) {
		case map[string]interface{}:
			flattenMap(v, flattenedMap, newKey)
		default:
			flattenedMap[newKey] = value
		}
	}
}

func toUDLAis(raw []byte) udl.AISIngest {
	originalMap := make(map[string]interface{})
	json.Unmarshal(raw, &originalMap)

	var vesselData VesselData
	json.Unmarshal(raw, &vesselData)

	UDLClassification := "U"
	var ais udl.AISIngest

	ais.Id = &vesselData.ID
	ais.ClassificationMarking = UDLClassification
	ais.Ts = vesselData.UpdateTimestamp
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
	ais.DestinationETA = &vesselData.CurrentVoyage.ETA

	if port, ok := originalMap["currentVoyage"].(map[string]interface{})["matchedPort"]; ok {
		if unlocode, ok := port.(map[string]interface{})["port"].(map[string]interface{})["unlocode"].(string); ok {
			ais.CurrentPortLOCODE = &unlocode
		}
	}

	ais.DataMode = "TEST"
	ais.Source = "Spire"
	return ais
}

func toUDLElset(raw []byte) (udl.ElsetIngest, error) {
	var elset udl.ElsetIngest
	err := json.Unmarshal(raw, &elset)
	return elset, err
}
