package connector

import sdk "github.com/conduitio/conduit-connector-sdk"

var Connector = sdk.Connector{
	NewSpecification: Specification,
	NewDestination:   NewDestination,
}