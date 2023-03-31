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
	BaseURL               = "baseURL"
	Endpoint              = "endpoint"
)

var DataModeValues = []string{"TEST", "REAL", "SIMULATED", "EXERCISE"}

type Config struct {
	// The HTTP Basic Auth Username to use when accessing the UDL.
	HTTPBasicAuthUsername string `validate:"required"`
	// The HTTP Basic Auth Password to use when accessing the UDL.
	HTTPBasicAuthPassword string `validate:"required"`
	// The Data Mode to use when submitting requests to the UDL. Acceptable values are REAL, TEST, SIMULATED and EXERCISE.
	DataMode string `default:"TEST"`
	// The Base URL to use to access the UDL. The default is https://unifieddatalibrary.com.
	BaseURL string `default:"https://unifieddatalibrary.com"`
	// The target UDL endpoint.
	Endpoint UDLEndpoint `validate:"required"`
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
		fmt.Printf("username: %s", u)
		fmt.Println("")
		fmt.Printf("ok: %t", ok)
		fmt.Println("")
		return Config{}, errors.New("missing or invalid credentials")
	}

	if p, ok := cfg[HTTPBasicAuthPassword]; p == "" || !ok {
		fmt.Printf("password: %s", p)
		fmt.Println("")
		fmt.Printf("ok: %t", ok)
		fmt.Println("")
		return Config{}, errors.New("missing or invalid credentials")
	}

	parsed := Config{
		HTTPBasicAuthUsername: cfg[HTTPBasicAuthUsername],
		HTTPBasicAuthPassword: cfg[HTTPBasicAuthPassword],
		DataMode:              dm,
		BaseURL:               u,
	}

	return parsed, nil
}

func supportedStringValues(check string, supported []string) bool {
	for _, ds := range supported {
		if strings.ToUpper(strings.TrimSpace(check)) == ds {
			s := fmt.Sprintf("%s is equal to %s", check, ds)
			fmt.Println(s)
			return true
		}
		s := fmt.Sprintf("%s is not equal to %s", check, ds)
		fmt.Println(s)
	}
	return false
}
