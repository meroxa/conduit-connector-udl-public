package connector

type Dimensions struct {
	A      float64 `json:"a"`
	B      float64 `json:"b"`
	C      float64 `json:"c"`
	D      float64 `json:"d"`
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
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
	Speed              float64 `json:"speed"`
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
