package connector

import (
	"context"

	"encoding/json"
	"fmt"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/deepmap/oapi-codegen/pkg/securityprovider"

	"github.com/meroxa/conduit-connector-udl/udl"
)

type Destination struct {
	sdk.UnimplementedDestination
	config Config
	client udl.ClientInterface
}

func NewDestination() sdk.Destination {
	return sdk.DestinationWithMiddleware(&Destination{}, sdk.DefaultDestinationMiddleware()...)
}

func (d *Destination) Parameters() map[string]sdk.Parameter {
	return d.config.Parameters()
	// return map[string]sdk.Parameter{
	// 	HTTPBasicAuthUsername: {
	// 		Default:     "",
	// 		Required:    true,
	// 		Description: "The HTTP Basic Auth Username to use when accessing the UDL.",
	// 	},
	// 	HTTPBasicAuthPassword: {
	// 		Default:     "",
	// 		Required:    true,
	// 		Description: "The HTTP Basic Auth Password to use when accessing the UDL.",
	// 	},
	// 	DataMode: {
	// 		Default:     "TEST",
	// 		Required:    false,
	// 		Description: "The Data Mode to use when submitting requests to the UDL. Acceptable values are REAL, TEST, SIMULATED and EXERCISE.",
	// 	},
	// 	BaseURL: {
	// 		Default:     "https://unifieddatalibrary.com",
	// 		Required:    false,
	// 		Description: "The Base URL to use to access the UDL. The default is https://unifieddatalibrary.com.",
	// 	},
	// 	Endpoint: {
	// 		Default:     "",
	// 		Required:    true,
	// 		Description: "The target UDL endpoint.",
	// 	},
	// }
}

func (d *Destination) Configure(ctx context.Context, cfg map[string]string) error {
	sdk.Logger(ctx).Debug().Msg("Configuring Destination connector...")
	sdk.Logger(ctx).Debug().Msgf("config: %w", cfg)
	parsedCfg, err := d.ParseDestinationConfig(cfg)
	if err != nil {
		return err
	}
	sdk.Logger(ctx).Debug().Msgf("parsedCfg: %w", parsedCfg)
	d.config = parsedCfg
	return nil
}

func (d *Destination) Open(ctx context.Context) error {
	authProvider, err := generateBasicAuth(d.config.HTTPBasicAuthUsername, d.config.HTTPBasicAuthPassword)
	if err != nil {
		return err
	}
	c, err := udl.NewClient(d.config.BaseURL, udl.WithRequestEditorFn(authProvider.Intercept))
	sdk.Logger(ctx).Debug().Msgf("client url: %w", c.Client)
	if err != nil {
		return err
	}
	d.client = c
	return nil
}

func (d *Destination) Write(ctx context.Context, records []sdk.Record) (int, error) {
	// var elsets []udl.ElsetIngest
	var aisData []udl.AISIngest
	for _, r := range records {
		sdk.Logger(ctx).Debug().Msgf("record: %w", r)
		// elset, err := toUDLElset(r.Payload.After.Bytes())
		// sdk.Logger(ctx).Debug().Msgf("elset: %w", elset)
		ais, err := toUDLAis(r.Payload.After.Bytes())
		sdk.Logger(ctx).Debug().Msgf("ais output: %+v", ais)
		sdk.Logger(ctx).Debug().Msgf("ais ShipName: %s", ais.ShipName)
		payload := make(sdk.StructuredData)
		if err := json.Unmarshal(r.Payload.After.Bytes(), &payload); err != nil {
			return 0, fmt.Errorf("unmarshal payload: %w", err)
		}
		sdk.Logger(ctx).Debug().Msgf("payload: %w", payload)
		aisData = append(aisData, ais)
		if err != nil {
			return 0, err
		}
		// elsets = append(elsets, elset)
	}

	resp, err := d.client.FiledropUdlAisPostId(ctx, aisData)
	sdk.Logger(ctx).Debug().Msgf("resp: %w", resp)
	if err != nil || resp.StatusCode > 300 {
		sdk.Logger(ctx).Debug().Msgf("err: %w", err)
		return 0, err
	}
	sdk.Logger(ctx).Debug().Msgf("resp: %w", resp)

	// resp, err := d.client.FiledropUdlElsetPostId(ctx, elsets)
	// if err != nil || resp.StatusCode > 300 {
	// 	return 0, err
	// }
	return len(aisData), nil
}

func generateBasicAuth(username, password string) (*securityprovider.SecurityProviderBasicAuth, error) {
	return securityprovider.NewSecurityProviderBasicAuth(username, password)
}
