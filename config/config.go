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

package config

const (
	HTTPBasicAuthUsername = "httpBasicAuthUsername"
	HTTPBasicAuthPassword = "httpBasicAuthPassword"
	DataMode              = "dataMode"
	DataType              = "dataType"
	BaseURL               = "baseURL"
	ClassificationMarking = "classificationMarking"
)

type Config struct {
	// The HTTP Basic Auth Username to use when accessing the UDL.
	HTTPBasicAuthUsername string `validate:"required"`
	// The HTTP Basic Auth Password to use when accessing the UDL.
	HTTPBasicAuthPassword string `validate:"required"`
	// The Data Mode to use when submitting requests to the UDL. Acceptable values are REAL, TEST, SIMULATED and EXERCISE.
	DataMode string `validate:"inclusion=REAL|TEST|SIMULATED|EXERCISE" default:"TEST"`
	// The Data Type that is being submitted to the UDL. Acceptable values are AIS and ELSET.
	DataType string `validate:"inclusion=AIS|ELSET|EPHEMERIS" default:"AIS"`
	// The Base URL to use to access the UDL. The default is https://unifieddatalibrary.com.
	BaseURL string `default:"https://unifieddatalibrary.com"`
	// Classification marking of the data in IC/CAPCO Portion-marked format. The default is U
	ClassificationMarking string `default:"U"`
}
