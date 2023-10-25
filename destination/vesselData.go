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

type Dimensions struct {
	A      float64 `json:"a, omitempty"`
	B      float64 `json:"b, omitempty"`
	C      float64 `json:"c, omitempty"`
	D      float64 `json:"d, omitempty"`
	Length float64 `json:"length, omitempty"`
	Width  float64 `json:"width, omitempty"`
}

type StaticData struct {
	AISClass        string     `json:"aisClass"`
	CallSign        string     `json:"callsign"`
	Dimensions      Dimensions `json:"dimensions"`
	Flag            string     `json:"flag"`
	IMO             int        `json:"imo"`
	MMSI            int64      `json:"mmsi"`
	Name            string     `json:"name"`
	ShipSubType     string     `json:"shipSubType"`
	ShipType        string     `json:"shipType"`
	Timestamp       string     `json:"timestamp"`
	UpdateTimestamp string     `json:"updateTimestamp"`
}

type LastPositionUpdate struct {
	Accuracy           string  `json:"accuracy"`
	CollectionType     string  `json:"collectionType"`
	Course             float64 `json:"course"`
	Heading            float64 `json:"heading"`
	Latitude           float64 `json:"latitude"`
	Longitude          float64 `json:"longitude"`
	Maneuver           string  `json:"maneuver"`
	NavigationalStatus string  `json:"navigationalStatus"`
	ROT                float64 `json:"rot"`
	Speed              float64 `json:"speed, omitempty"`
	Timestamp          string  `json:"timestamp"`
	UpdateTimestamp    string  `json:"updateTimestamp"`
}

type GeoPoint struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
type Port struct {
	CenterPoint GeoPoint `json:"centerPoint"`
	Name        string   `json:"name"`
	Unlocode    string   `json:"unlocode"`
}

type MatchedPort struct {
	MatchScore float64 `json:"matchScore"`
	Port       Port    `json:"port"`
}

type CurrentVoyage struct {
	Destination     string      `json:"destination"`
	Draught         float64     `json:"draught"`
	ETA             string      `json:"eta"`
	MatchedPort     MatchedPort `json:"matchedPort"`
	Timestamp       string      `json:"timestamp"`
	UpdateTimestamp string      `json:"updateTimestamp"`
}

type VesselData struct {
	ID                 string             `json:"id"`
	StaticData         StaticData         `json:"staticData"`
	CurrentVoyage      CurrentVoyage      `json:"currentVoyage"`
	LastPositionUpdate LastPositionUpdate `json:"lastPositionUpdate"`
	UpdateTimestamp    string             `json:"updateTimestamp"`
}
