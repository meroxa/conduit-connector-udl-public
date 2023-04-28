# Conduit Connector for the Unified Data Library (UDL)

[Conduit](https://conduit.io) Destination Connector for the [Unified Data Library](https://unifieddatalibrary.com/).

## How to build?

Run `make build` to build the connector.

## Destination

The destination connector pushes data to the Unified Data Library (UDL). The connection supports various data types as specified by the UDL and pushes to those respective endpoints

### Configuration

A UDL username and password is required to use this connector

| name                    | description                                                                                                         | required | default value |
| ----------------------- | ------------------------------------------------------------------------------------------------------------------- | -------- | ------------- | ------------------------------ |
| `httpBasicAuthUsername` | The HTTP Basic Auth Username to use when accessing the UDL.                                                         | true     |               |
| `httpBasicAuthPassword` | The HTTP Basic Auth Password to use when accessing the UDL.                                                         | true     |               |
| `dataMode`              | The Data Mode to use when submitting requests to the UDL. Acceptable values are REAL, TEST, SIMULATED and EXERCISE. | false    | TEST          |
| `dataType`              | The Data Type that is being submitted to the UDL. Acceptable values are AIS, ELSET, and EPHEMERIS.                  | false    | AIS           |
| `baseURL`               | The Base URL to use to access the UDL. The default is https://unifieddatalibrary.com.                               | false    | AIS           | https://unifieddatalibrary.com |
