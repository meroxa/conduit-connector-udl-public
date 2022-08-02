package main

import (
	sdk "github.com/conduitio/conduit-connector-sdk"
	udl "github.com/meroxa/conduit-connector-udl"
)

func main() {
	sdk.Serve(
		udl.Specification,
		nil,
		udl.NewDestination,
	)
}
