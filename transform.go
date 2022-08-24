package connector

import (
	"encoding/json"
	"github.com/meroxa/conduit-connector-udl/udl"
)

func toUDLElset(raw []byte) (udl.ElsetIngest, error) {
	var elset udl.ElsetIngest
	err := json.Unmarshal(raw, &elset)
	return elset, err
}
