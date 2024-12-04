package main

import (
	sdk "github.com/conduitio/conduit-connector-sdk"
	udl "github.com/meroxa/conduit-connector-udl-public"
)

func main() {
	sdk.Serve(udl.Connector)
}
