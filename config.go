package udl

import (
	"errors"
	"fmt"
	"strings"
)

const (
	HTTPBasicAuthUsername = "username"
	HTTPBasicAuthPassword = "password"
	DataMode              = "dataMode"
)

var DataModeValues = []string{"TEST", "REAL", "SIMULATED", "EXERCISE"}

type Config struct {
	HTTPBasicAuthUsername string
	HTTPBasicAuthPassword string
	DataMode              string
}

func Parse(cfg map[string]string) (Config, error) {
	// validate supported data mode
	dm, ok := cfg[DataMode]
	if ok {
		if !supportStringValues(dm, DataModeValues) {
			return Config{}, errors.New(fmt.Sprintf("unsupported data mode (%s)", dm))
		}
	} else {
		dm = "TEST"
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
	}

	return parsed, nil
}

func supportStringValues(check string, supported []string) bool {
	for _, ds := range supported {
		if strings.ToLower(strings.TrimSpace(check)) == ds {
			return true
		}
	}
	return false
}
