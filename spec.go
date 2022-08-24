package connector

import sdk "github.com/conduitio/conduit-connector-sdk"

func Specification() sdk.Specification {
	return sdk.Specification{
		Name:        "udl",
		Summary:     "A UDL (Unified Data Library) connector for Conduit, written in Go.",
		Description: "A UDL (Unified Data Library) connector for Conduit, written in Go.",
		Version:     "v0.1.0",
		Author:      "Meroxa, Inc.",
	}
}
