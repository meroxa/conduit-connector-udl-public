package udl

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

const (
	HTTPBasicAuthUsername = "username"
	HTTPBasicAuthPassword = "password"
	DataMode              = "dataMode"
	BaseURL               = "baseurl"
)

var DataModeValues = []string{"TEST", "REAL", "SIMULATED", "EXERCISE"}

type Config struct {
	HTTPBasicAuthUsername string
	HTTPBasicAuthPassword string
	DataMode              string
	BaseURL               string
}

func Parse(cfg map[string]string) (Config, error) {
	// validate supported data mode
	dm, ok := cfg[DataMode]
	if ok {
		if !supportedStringValues(dm, DataModeValues) {
			return Config{}, errors.New(fmt.Sprintf("unsupported data mode (%s)", dm))
		}
	} else {
		dm = Specification().DestinationParams[DataMode].Default
	}

	// validate/parse base URL
	u, ok := cfg[BaseURL]
	if ok {
		_, err := url.Parse(u)
		if err != nil {
			return Config{}, errors.New(fmt.Sprintf("invalid base URL (%s); err: %s", u, err))
		}
	} else {
		u = Specification().DestinationParams[BaseURL].Default
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
		BaseURL:               u,
	}

	return parsed, nil
}

func supportedStringValues(check string, supported []string) bool {
	for _, ds := range supported {
		if strings.ToLower(strings.TrimSpace(check)) == ds {
			return true
		}
	}
	return false
}
