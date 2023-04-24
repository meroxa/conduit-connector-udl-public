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
	"errors"
	"fmt"
	"net/url"
	"strings"
)

const (
	HTTPBasicAuthUsername = "httpBasicAuthUsername"
	HTTPBasicAuthPassword = "httpBasicAuthPassword"
	DataMode              = "dataMode"
	DataType              = "dataType"
	BaseURL               = "baseURL"
)

var DataModeValues = []string{"TEST", "REAL", "SIMULATED", "EXERCISE"}

var DataTypeValues = []string{"AIS", "ELSET", "EPHEMERIS"}

type Config struct {
	// The HTTP Basic Auth Username to use when accessing the UDL.
	HTTPBasicAuthUsername string `validate:"required"`
	// The HTTP Basic Auth Password to use when accessing the UDL.
	HTTPBasicAuthPassword string `validate:"required"`
	// The Data Mode to use when submitting requests to the UDL. Acceptable values are REAL, TEST, SIMULATED and EXERCISE.
	DataMode string `default:"TEST"`
	// The Data Type that is being submitted to the UDL. Acceptable values are AIS, ELSET, and EPHEMERIS.
	DataType string `default:"AIS"`
	// The Base URL to use to access the UDL. The default is https://unifieddatalibrary.com.
	BaseURL string `default:"https://unifieddatalibrary.com"`
}

func (c *Destination) ParseDestinationConfig(cfg map[string]string) (Config, error) {
	// validate supported data mode
	dm, ok := cfg[DataMode]
	fmt.Printf("dm: %s", dm)
	fmt.Println("")
	fmt.Printf("ok: %t", ok)
	fmt.Println("")
	if ok {
		if !supportedStringValues(dm, DataModeValues) {
			return Config{}, errors.New(fmt.Sprintf("unsupported data mode (%s)", dm))
		}
	} else {
		dm = c.Parameters()[DataMode].Default
	}

	// validate supported data type
	dt, ok := cfg[DataType]
	fmt.Printf("dt: %s", dt)
	fmt.Println("")
	fmt.Printf("ok: %t", ok)
	fmt.Println("")
	if ok {
		if !supportedStringValues(dt, DataTypeValues) {
			return Config{}, errors.New(fmt.Sprintf("unsupported data type (%s)", dt))
		}
	} else {
		dt = c.Parameters()[DataType].Default
	}

	// validate/parse base URL
	u, ok := cfg[BaseURL]
	fmt.Printf("baseURL: %s", u)
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

func supportedStringValues(check string, supported []string) bool {
	for _, ds := range supported {
		if strings.ToUpper(strings.TrimSpace(check)) == ds {
			return true
		}
	}
	return false
}
