// Copyright © 2022 Meroxa, Inc.
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

import (
	"context"
	"fmt"
	"log"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/deepmap/oapi-codegen/pkg/securityprovider"

	"github.com/meroxa/udl-go"
)

type Destination struct {
	sdk.UnimplementedDestination
	Config Config
	client udl.ClientInterface
}

func NewDestination() sdk.Destination {
	return sdk.DestinationWithMiddleware(&Destination{}, sdk.DefaultDestinationMiddleware()...)
}

func (d *Destination) Parameters() map[string]sdk.Parameter {
	return d.Config.Parameters()
}

func (d *Destination) Configure(ctx context.Context, cfg map[string]string) error {
	sdk.Logger(ctx).Debug().Msg("Configuring Destination connector...")
	err := sdk.Util.ParseConfig(cfg, &d.Config)
	if err != nil {
		sdk.Logger(context.Background()).Err(err).Msgf("invalid config")
		return err
	}
	return nil
}

func (d *Destination) Open(ctx context.Context) error {
	authProvider, err := generateBasicAuth(d.Config.HTTPBasicAuthUsername, d.Config.HTTPBasicAuthPassword)
	if err != nil {
		return err
	}
	c, err := udl.NewClient(d.Config.BaseURL, udl.WithRequestEditorFn(authProvider.Intercept))
	if err != nil {
		return err
	}
	d.client = c
	return nil
}

func (d *Destination) Write(ctx context.Context, records []sdk.Record) (int, error) {
	dataType := d.Config.DataType

	log.Printf("dataType selected: %s\n", d.Config.DataType)
	sdk.Logger(context.Background()).Info().Msgf("dataType selected: %s", d.Config.DataType)
	switch dataType {
	case "AIS":
		return d.writeAisToUDL(ctx, records)
	case "Elset":
		return d.writeElsetToUDL(ctx, records)
	case "EPHEMERIS":
		return d.writeEphemerisToUDL(ctx, records)
	default:
		return 0, fmt.Errorf("unsupported data type: %s;", dataType)
	}
}

func (d *Destination) Teardown(ctx context.Context) error {
	// Teardown signals to the plugin that there will be no more calls to any
	// other function. After Teardown returns, the plugin should be ready for a
	// graceful shutdown.
	return nil
}

func generateBasicAuth(username, password string) (*securityprovider.SecurityProviderBasicAuth, error) {
	return securityprovider.NewSecurityProviderBasicAuth(username, password)
}
