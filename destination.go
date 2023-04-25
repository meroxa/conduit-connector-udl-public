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

package connector

import (
	"context"
	"fmt"
	"strings"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/deepmap/oapi-codegen/pkg/securityprovider"

	"github.com/meroxa/udl-go"
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
}

func (d *Destination) Configure(ctx context.Context, cfg map[string]string) error {
	sdk.Logger(ctx).Debug().Msg("Configuring Destination connector...")
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
	dataType := d.config.DataType
	dataTypeValuesStr := strings.Join(DataTypeValues, "', '")

	// Check if dataType is one of the acceptable data types listed in DataTypeValues and return an error if it's not
	if !SupportedStringValues(dataType, DataTypeValues) {
		return 0, fmt.Errorf("invalid data type: %s; expecting the data type to be on of the following values: '%s'", dataType, dataTypeValuesStr)
	}

	switch dataType {
	case "AIS":
		return writeAisToUDL(ctx, records, d)
	case "Elset":
		return writeElsetToUDL(ctx, records, d)
	default:
		return 0, fmt.Errorf("unsupported data type: %s;", dataType)
	}
}

func generateBasicAuth(username, password string) (*securityprovider.SecurityProviderBasicAuth, error) {
	return securityprovider.NewSecurityProviderBasicAuth(username, password)
}
