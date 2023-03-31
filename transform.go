package connector

import (
	"encoding/json"

	"fmt"

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

func toUDLAis(raw []byte) (udl.AISIngest, error) {
	originalMap := make(map[string]interface{})
	json.Unmarshal(raw, &originalMap)

	flattenedMap := make(map[string]interface{})
	flattenMap(originalMap, flattenedMap, "")
	fmt.Println(fmt.Sprintf("flattenedMap: %s", flattenedMap))
	fmt.Println(fmt.Sprintf("flattenedMap static data: %s", flattenedMap["staticData"]))
	fmt.Println(fmt.Sprintf("flattenedMap static data name: %s", flattenedMap["staticData.name"]))

	jsonData, _ := json.Marshal(flattenedMap)

	var ais udl.AISIngest
	// ais.Lat = flattenedMap["lastPositionUpdate.longitude"]
	// ais.Lat = originalMap["lastPositionUpdate"]
	// ais.ShipName = originalMap["staticData"]["name"]
	// shipName := flattenedMap["staticData.name"].(string)
	// ais.ShipName = &shipName
	// err := json.Unmarshal(raw, &ais)

	// TODO: Make these as part of config
	error := json.Unmarshal(jsonData, &ais)
	ais.DataMode = "TEST"
	ais.Source = "Spire"
	ais.ClassificationMarking = "U"
	return ais, error
}

func toUDLElset(raw []byte) (udl.ElsetIngest, error) {
	var elset udl.ElsetIngest
	err := json.Unmarshal(raw, &elset)
	return elset, err
}
