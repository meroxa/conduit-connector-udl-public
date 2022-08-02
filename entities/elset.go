package entities

import "time"

type Elset struct {
	ElsetID          string   `json:"idElset"`
	OnOrbitID        string   `json:"idOnOrbit"`
	OriginObjectID   string   `json:"origOrbitId"`
	Descriptor       string   `json:"descriptor"`
	RawFileURI       string   `json:"rawFileURI"`
	Origin           string   `json:"origin"`
	Source           string   `json:"source"`
	Tags             []string `json:"tags"`
	Algorithm        string   `json:"algorithm"`
	SourcedData      []string `json:"sourcedData"`
	SourcedDataTypes []string `json:"sourcedDataTypes"`
	TransactionID    string   `json:"transactionId"`

	SatelliteNo                 int       `json:"satNo"`
	Classification              string    `json:"classificationMarking"`
	Epoch                       time.Time `json:"epoch"`
	MeanMotion                  float64   `json:"meanMotion"`
	Inclination                 float64   `json:"inclination"`
	RightAscensionAscendingNode float64   `json:"raan"`
	Eccentricity                float64   `json:"eccentricity"`
	MeanAnomaly                 float64   `json:"meanAnomaly"`
	MeanMotionDot               float64   `json:"meanMotionDot"`
	MeanMotionDDot              float64   `json:"meanMotionDDot"`
	BSTARDrag                   float64   `json:"bStar"`
	ArgumentOfPerigree          float64   `json:"argOfPerigree"`
	RevNumber                   float64   `json:"revNo"`
	EphemType                   int       `json:"ephemType"`
	Line1                       string    `json:"line1,omitempty"`
	Line2                       string    `json:"line2,omitempty"`
}
