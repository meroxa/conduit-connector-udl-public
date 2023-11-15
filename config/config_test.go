package config

import (
	"testing"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/matryer/is"
)

var exampleConfig = map[string]string{
	"baseURL":               "https://example.com",
	"httpBasicAuthUsername": "user",
	"httpBasicAuthPassword": "pass",
	"dataType":              "AIS",
	"dataMode":              "TEST",
	"classificationMarking": "U",
}

func TestParseConfig(t *testing.T) {
	is := is.New(t)
	var got Config
	err := sdk.Util.ParseConfig(exampleConfig, &got)
	want := Config{
		BaseURL:               "https://example.com",
		HTTPBasicAuthUsername: "user",
		HTTPBasicAuthPassword: "pass",
		DataMode:              "TEST",
		DataType:              "AIS",
		ClassificationMarking: "U",
	}
	is.NoErr(err)
	is.Equal(want, got)
}

var exampleConfigWithSpireClassificationMarking = map[string]string{
	"baseURL":               "https://example.com",
	"httpBasicAuthUsername": "user",
	"httpBasicAuthPassword": "pass",
	"dataType":              "AIS",
	"dataMode":              "TEST",
	"classificationMarking": "U//PR-SPIRE-AIS",
}

func TestParseConfigWithSpireClassificationMarking(t *testing.T) {
	is := is.New(t)
	var got Config
	err := sdk.Util.ParseConfig(exampleConfigWithSpireClassificationMarking, &got)
	want := Config{
		BaseURL:               "https://example.com",
		HTTPBasicAuthUsername: "user",
		HTTPBasicAuthPassword: "pass",
		DataMode:              "TEST",
		DataType:              "AIS",
		ClassificationMarking: "U//PR-SPIRE-AIS",
	}
	is.NoErr(err)
	is.Equal(want, got)
}
