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

package connector

import (
	"errors"
	"fmt"
	"net/url"
)

const (
	HTTPBasicAuthUsername = "httpBasicAuthUsername"
	HTTPBasicAuthPassword = "httpBasicAuthPassword"
	DataMode              = "dataMode"
	DataType              = "dataType"
	BaseURL               = "baseURL"
)

type Config struct {
	// The HTTP Basic Auth Username to use when accessing the UDL.
	HTTPBasicAuthUsername string `validate:"required"`
	// The HTTP Basic Auth Password to use when accessing the UDL.
	HTTPBasicAuthPassword string `validate:"required"`
	// The Data Mode to use when submitting requests to the UDL. Acceptable values are REAL, TEST, SIMULATED and EXERCISE.
	DataMode string `default:"TEST"`
	// The Data Type that is being submitted to the UDL. Acceptable values are AIS and ELSET.
	DataType string `default:"AIS"`
	// The Base URL to use to access the UDL. The default is https://unifieddatalibrary.com.
	BaseURL string `default:"https://unifieddatalibrary.com"`
}

func (c *Destination) ParseDestinationConfig(cfg map[string]string) (Config, error) {
	// validate supported data mode
	dm, ok := cfg[DataMode]
	if ok {
		if !SupportedStringValues(dm, DataModeValues) {
			return Config{}, errors.New(fmt.Sprintf("unsupported data mode (%s)", dm))
		}
	} else {
		dm = c.Parameters()[DataMode].Default
	}

	// validate supported data type
	dt, ok := cfg[DataType]
	if ok {
		if !SupportedStringValues(dt, DataTypeValues) {
			return Config{}, errors.New(fmt.Sprintf("unsupported data type (%s)", dt))
		}
	} else {
		dt = c.Parameters()[DataType].Default
	}

	// validate/parse base URL
	u, ok := cfg[BaseURL]
	if ok {
		_, err := url.Parse(u)
		if err != nil {
			return Config{}, errors.New(fmt.Sprintf("invalid base URL (%s); err: %s", u, err))
		}
	} else {
		u = c.Parameters()[BaseURL].Default
	}

	// validate HTTP Basic Auth credentials
	if u, ok := cfg[HTTPBasicAuthUsername]; u == "" || !ok {
		return Config{}, errors.New("missing or invalid credentials")
	}

	if p, ok := cfg[HTTPBasicAuthPassword]; p == "" || !ok {
		return Config{}, errors.New("missing or invalid credentials")
	}

	parsed := Config{
		HTTPBasicAuthUsername: cfg[HTTPBasicAuthUsername],
		HTTPBasicAuthPassword: cfg[HTTPBasicAuthPassword],
		DataMode:              dm,
		DataType:              dt,
		BaseURL:               u,
	}

	return parsed, nil
}
