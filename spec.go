package udl

import sdk "github.com/conduitio/conduit-connector-sdk"

func Specification() sdk.Specification {
	return sdk.Specification{
		Name:        "udl",
		Summary:     "A UDL (Unified Data Library) connector for Conduit, written in Go.",
		Description: "A UDL (Unified Data Library) connector for Conduit, written in Go.",
		Version:     "v0.1.0",
		Author:      "Meroxa, Inc.",
		DestinationParams: map[string]sdk.Parameter{
			HTTPBasicAuthUsername: {
				Default:     "",
				Required:    true,
				Description: "The HTTP Basic Auth Username to use when accessing the UDL.",
			},
			HTTPBasicAuthPassword: {
				Default:     "",
				Required:    true,
				Description: "The HTTP Basic Auth Password to use when accessing the UDL.",
			},
			DataMode: {
				Default:     "TEST",
				Required:    false,
				Description: "The Data Mode to use when submitting requests to the UDL. Acceptable values are REAL, TEST, SIMULATED and EXERCISE.",
			},
			BaseURL: {
				Default:     "https://unifieddatalibrary.com",
				Required:    false,
				Description: "The Base URL to use to access the UDL. The default is https://unifieddatalibrary.com.",
			},
		},
	}
}
