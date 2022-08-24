package connector

import (
	"context"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/meroxa/conduit-connector-udl/udl"
)

type Destination struct {
	sdk.UnimplementedDestination
	config Config
	client *udl.Client
}

func NewDestination() sdk.Destination {
	return sdk.DestinationWithMiddleware(&Destination{}, sdk.DefaultDestinationMiddleware()...)
}

func (d *Destination) Parameters() map[string]sdk.Parameter {
	return map[string]sdk.Parameter{
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
		Endpoint: {
			Default:     "",
			Required:    true,
			Description: "The target UDL endpoint.",
		},
	}
}

func (d *Destination) Configure(ctx context.Context, cfg map[string]string) error {
	parsedCfg, err := d.ParseDestinationConfig(cfg)
	if err != nil {
		return err
	}
	d.config = parsedCfg
	return nil
}

func (d *Destination) Open(ctx context.Context) error {
	authProvider, err := generateBasicAuth(d.config.HTTPBasicAuthUsername, d.config.HTTPBasicAuthPassword)
	if err != nil {
		return err
	}
	c, err := udl.NewClient(d.config.BaseURL, udl.WithRequestEditorFn(authProvider.Intercept))
	if err != nil {
		return err
	}
	d.client = c
	return nil
}

func (d *Destination) Write(ctx context.Context, records []sdk.Record) (int, error) {
	var elsets []udl.ElsetIngest
	for _, r := range records {
		elset, err := toUDLElset(r.Payload.After.Bytes())
		if err != nil {
			return 0, err
		}
		elsets = append(elsets, elset)
	}
	d.client.FiledropUdlElsetPostId(ctx, elsets)
	return 0, nil
}

func generateBasicAuth(username, password string) (*securityprovider.SecurityProviderBasicAuth, error) {
	return securityprovider.NewSecurityProviderBasicAuth(username, password)
}
